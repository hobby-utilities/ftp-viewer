<script setup lang="ts">
import {GetList} from "../../wailsjs/go/main/FtpService";
import {computed, onMounted, ref} from "vue";
import DataTable, {DataTableRowSelectEvent} from "primevue/datatable";
import Column from "primevue/column";

const PreviousDirectory = '../'

// FtpViewer를 호출 시 입력되는 파라미터
const props = defineProps<{ connectionId: string }>()

// FTP에서 조회할 경로를 관리하기 위한 폴더의 이름으록 구성된 배열
const currentPathTokens = ref<string[]>([''])

// `currentPathTokens`에 입력된 값에 따라 FTP 조회 시 사용할 실제 경로를 생성하는 반응형 변수
const currentPath = computed(() => {
  const p = currentPathTokens.value.join('/')
  return p.length === 0 ? '/' : p
})
const _clickedRow = ref();

// FTP에서 정의된 파일 유형
enum FileType {File, Directory}

// FTP에서 반환되는 파일 데이터 구조체 정의
type FileInfo = { Name: string; Size: number; Type: FileType; Time: string; }

// 현재 경로에 존재하는 파일 또는 폴더들에 대한 배열
const directoryFiles = ref<FileInfo[]>([])

/**
 * FtpViewer 호출 시 1회 실행되는 함수
 */
onMounted(async () => {
  await getFiles(currentPath.value)
})

/**
 * DataTable의 열 선택 시 실행되는 이벤트 함수
 * @param event
 */
function onRowSelect(event: DataTableRowSelectEvent) {

  const data: FileInfo = event.data;
  if (data.Type === FileType.File)
    return;

  // ../ 선택 시
  if (data.Name === PreviousDirectory) {
    // 현재 경로가 root(/)일 경우 동작하지 않음.
    if (currentPathTokens.value.length === 1)
      return;
    // 마지막 경로를 제거한다. (e.g. /src/frontend 를 의미하는 ['src', 'frontend'] 에서 'frontend' 제거)
    currentPathTokens.value.pop()
  } else {
    // ../ 이외의 정상 폴더의 경우 경로 토큰에 폴더 이름을 추가한다.
    currentPathTokens.value.push(data.Name)
  }

  // FTP에서 변경된 경로에 대한 파일 목록을 요청한다.
  getFiles(currentPath.value)
}

async function getFiles(path: string): Promise<void> {
  const files = await GetList(props.connectionId, path)
  const l: FileInfo[] = [{Name: PreviousDirectory, Size: 0, Type: FileType.Directory, Time: ''}]
  l.push(...files)
  directoryFiles.value = l
}

</script>

<template>
  <div style="margin-bottom: 8px">
    <div>Current path: {{ currentPath }}</div>
  </div>

  <!-- 폴더 내에 파일이 없을 경우 Empty를 출력하도록 설정 (Warnings 때문에 추가됨. 현재는 ../ 때문에 출력되지 않음) -->
  <template v-if="directoryFiles.length === 0">
    Empty
  </template>

  <!-- 폴더 내 파일 정보들을 DataTable 형태로 표현 -->
  <template v-else>
    <DataTable :value="directoryFiles" data-key="Name"
               v-model:selection="_clickedRow"
               selection-mode="single"
               table-style="min-width: 50rem" @row-select="onRowSelect">
      <Column field="Type" header="type">
        <template #body="slotProps">
          <template v-if="slotProps.data.Type === FileType.File">FILE</template>
          <template v-else>DIR</template>
        </template>
      </Column>
      <Column field="Name" header="name"/>
      <Column field="Size" header="size"/>
      <Column field="Time" header="time"/>
    </DataTable>
  </template>
</template>

<style scoped>

</style>