package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jlaffaye/ftp"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type FtpConnection struct {
	id         string
	connection *ftp.ServerConn
}

type FtpService struct {
	ctx         context.Context
	connections map[string]*ftp.ServerConn
}

func NewFtpService() *FtpService {
	return &FtpService{
		connections: map[string]*ftp.ServerConn{},
	}
}

func (service *FtpService) Connect(url string, username string, password string) (connId *string, err error) {
	connection, err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = connection.Login(username, password)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	id := uuid.NewString()
	connId = &id

	service.connections[id] = connection

	return connId, nil
}

func (service *FtpService) GetList(connId string, directoryPath string) ([]*FileInfo, error) {
	if service.connections[connId] == nil {
		return nil, errors.New("connection ID is wrong")
	}
	entries, err := service.connections[connId].List(directoryPath)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return FromEntries(entries), nil
}

func (service *FtpService) Disconnect(connId string) {
	if service.connections[connId] == nil {
	}

	err := service.connections[connId].Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func addressString(address string) string {
	var newAddress = address

	if strings.HasPrefix(newAddress, "ftp://") {
		newAddress, _ = strings.CutPrefix(newAddress, "ftp://")
	}
	if strings.HasPrefix(newAddress, "ftps://") {
		newAddress, _ = strings.CutPrefix(newAddress, "ftps://")
	}

	regex, _ := regexp.Compile(".+:\\d")

	if strings.HasSuffix(newAddress, "/") {
		newAddress, _ = strings.CutSuffix(newAddress, "/")
	}

	if !regex.MatchString(newAddress) {
		newAddress = newAddress + ":21"
	}

	return newAddress
}

func connect(address string, username string, password string) (*ftp.ServerConn, error) {
	connection, err := ftp.Dial(address, ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = connection.Login(username, password)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return connection, nil
}

var eucKrDecoder = korean.EUCKR.NewDecoder()

func FromEntries(l []*ftp.Entry) []*FileInfo {
	fiList := make([]*FileInfo, 0)

	for _, element := range l {
		fi, err := FromEntry(element)
		if err != nil {
			return nil
		}
		fiList = append(fiList, fi)
	}

	return fiList
}

func FromEntry(entry *ftp.Entry) (*FileInfo, error) {
	return NewFtpFileInfo(entry.Name, entry.Size, entry.Type, entry.Time)
}

func NewFtpFileInfo(name string, size uint64, typ ftp.EntryType, time time.Time) (*FileInfo, error) {

	newName := getUtf8Name(name)

	if newName == nil {
		return nil, errors.New("string `name` is not UTF8 or EUC-KR")
	}

	return &FileInfo{
		Name: *newName,
		Size: size,
		Type: typ,
		Time: time,
	}, nil
}

type FileInfo struct {
	Name string
	Size uint64
	Type ftp.EntryType
	Time time.Time
}

func getUtf8Name(name string) *string {
	var validatedName *string = nil

	if utf8.ValidString(name) {
		validatedName = &name
	} else {
		utf8Name, _, err := transform.String(eucKrDecoder, name)

		if err != nil {
			log.Fatal(err)
			return nil
		}
		validatedName = &utf8Name
	}

	return validatedName
}
