<script setup lang="ts">
import {GetList} from "../../wailsjs/go/main/FtpService";
import {computed, onMounted, ref} from "vue";
import DataTable, {DataTableRowSelectEvent} from "primevue/datatable";
import Column from "primevue/column";

const PreviousDirectory = '../'

const props = defineProps<{ connectionId: string }>()

const currentPathTokens = ref<string[]>([''])
const currentPath = computed(() => {
  const p = currentPathTokens.value.join('/')
  return p.length === 0 ? '/' : p
})
const _clickedRow = ref();

enum FileType {File, Directory}
type FileInfo = { Name: string; Size: number; Type: FileType; Time: string; }
const directoryFiles = ref<FileInfo[]>([])

onMounted(async () => {
  await getFiles(currentPath.value)
})

function onRowSelect(event: DataTableRowSelectEvent) {

  const data: FileInfo = event.data;
  if (data.Type === FileType.File)
    return;

  if (data.Name === PreviousDirectory) {
    if (currentPathTokens.value.length === 1)
      return;
    currentPathTokens.value.pop()
  } else {
    currentPathTokens.value.push(data.Name)
  }
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
  <template v-if="directoryFiles.length === 0">
    Empty
  </template>
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