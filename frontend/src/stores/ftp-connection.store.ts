import {defineStore} from "pinia";
import {ref} from "vue";
import {Connect} from "../../wailsjs/go/main/FtpService";

export type ConnectionId = string

export class FtpConnection {
    constructor(
        public readonly connectionId: ConnectionId,
        public readonly url: string,
        public readonly username: string,
        public readonly alias: string,
    ) {
    }
}

export const useFtpConnectionStore = defineStore('ftp-connection', () => {

    const connections = ref<Map<ConnectionId, FtpConnection>>(new Map())

    async function connect(url: string, username: string, password: string, alias: string): Promise<string> {

        alert('getcha')
        const connectionId = await Connect(url, username, password)
        alert(connectionId)

        alias = alias === undefined || alias.trim().length === 0
            ? `${username}@${url}`
            : alias

        const connection = new FtpConnection(connectionId, url, username, alias)

        connections.value.set(connectionId, connection)

        return connectionId
    }

    return {
        connections,
        connect,
    }
})