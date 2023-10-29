<script lang="ts" setup>
import {reactive} from 'vue'
import {useFtpConnectionStore} from "../stores/ftp-connection.store";
import TabPanel from "primevue/tabpanel";
import TabView from "primevue/tabview";
import FtpViewer from "./FtpViewer.vue";

const ftpConnectionStore = useFtpConnectionStore()

/** 사용자가 입력한 FTP 연결에 필요한 입력값들 */
const connectionInputs = reactive({
  url:'',
  username: '',
  password: '',
  alias: '',
  modifiedAliasByUser: false,
})

/**
 * ID 입력 시 `connectionInputs.alias`를 필요에 따라 자동 입력하는 이벤트
 */
function onChangeUsername() {
  _updateAlias()
}

/**
 * FTP URL 입력 시 `connectionInputs.alias`를 필요에 따라 자동 입력하는 이벤트
 */
function onChangeUrl() {
  _updateAlias()
}

/**
 * alias 입력 시 username 또는 url 입력에 따라 자동으로 alias가 변하지 않도록 설정
 */
function onChangeAlias() {
  connectionInputs.modifiedAliasByUser = true
}

/**
 * 연결(connect) 클릭 시 FTP에 연결하고 connection ID를 store에 등록한다.
 * 등록될 경우 새로운 tab이 생성된다.
 */
function onClickConnect() {
  const reg = /^.+:\d/g

  ftpConnectionStore.connect(
      reg.exec(connectionInputs.url) === null
          ? `${connectionInputs.url}:21`
          : connectionInputs.url,
      connectionInputs.username,
      connectionInputs.password,
      connectionInputs.alias)
}

function _updateAlias() {
  if (connectionInputs.modifiedAliasByUser) {
    if (connectionInputs.alias.trim().length === 0)
      connectionInputs.modifiedAliasByUser = false;
  }

  if (!connectionInputs.modifiedAliasByUser)
    connectionInputs.alias = `${connectionInputs.username}@${connectionInputs.url}`
}

</script>

<template>
  <main>
    <div>
      <div>FTP 로그인 정보</div>
      <input id="" class="connectionInput" v-model="connectionInputs.url" placeholder="ftp.*.com" @change="onChangeUrl"/>
      <input id="" class="connectionInput" v-model="connectionInputs.username" placeholder="ID" @change="onChangeUsername"/>
      <input id="" class="connectionInput" v-model="connectionInputs.password" placeholder="password"/>
      <input id="" class="connectionInput" v-model="connectionInputs.alias" placeholder="alias"/>
      <button @click=onClickConnect>connect</button>
    </div>


      <template v-if="ftpConnectionStore.connections.size === 0">
        No connections
      </template>
      <template v-else>
        <TabView :scrollable=true>
          <template v-for="connection in ftpConnectionStore.connections.values()" :key="connection.connectionId">
            <TabPanel :header="connection.alias">
              <FtpViewer :connection-id="connection.connectionId"/>
            </TabPanel>
          </template>
        </TabView>
      </template>
  </main>
</template>

<style scoped>
.connectionInput {
  margin-left: 4px;
  margin-right: 4px;
}

</style>
