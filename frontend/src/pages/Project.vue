<template>
  <v-container v-if="project" class="d-flex ga-4 flex-column">
    <v-sheet
      class="d-flex pa-4 flex-row justify-space-between align-center"
      elevation="2"
    >
      <v-sheet class="d-flex pa-4 flex-column">
        <v-sheet class="d-flex flex-row justify-space-between align-center">
          <h2>#{{ project.id }} {{ project.name }}</h2>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Создан</b>
          <span>{{ new Date(project.created_at).toLocaleString() }}</span>
        </v-sheet>
        <v-sheet v-if="project.dir" class="d-flex flex-row ga-2">
          <b>Директория</b>
          <span>{{ project.dir }}</span>
        </v-sheet>
        <v-sheet v-if="project.git_url" class="d-flex flex-row ga-2">
          <b>GIT</b>
          <span>{{ project.git_url }}</span>
        </v-sheet>
        <v-sheet v-if="project.git_branch" class="d-flex flex-row ga-2">
          <b>GIT branch</b>
          <span>{{ project.git_branch }}</span>
        </v-sheet>
      </v-sheet>
      <v-sheet class="d-flex ga-2 flex-column">
        <v-btn prepend-icon="mdi-pencil" @click="editProject(project.id)">
          Редактировать
        </v-btn>
        <v-btn
          :disabled="running"
          :loading="running"
          prepend-icon="mdi-play"
          @click="runProject(project.id)"
        >
          Запустить
        </v-btn>
      </v-sheet>
    </v-sheet>

    <v-sheet
      v-for="stage in project.stages"
      class="d-flex pa-4 flex-column"
      elevation="2"
    >
      <h3>Этап {{ stage.number }}</h3>
      <pre>{{ stage.script }}</pre>
    </v-sheet>

    <v-sheet class="d-flex pa-4 flex-column" elevation="2">
      <div class="d-flex flex-row justify-space-between align-center">
        <h3>История запусков</h3>
        <v-btn
          icon="mdi-refresh"
          size="small"
          variant="text"
          @click="loadRuns"
        />
      </div>
      <v-data-table
        :headers="runsHeaders"
        hover
        item-value="id"
        :items="runs"
        :loading="runsLoading"
        @click:row="handleRowClick"
      >
        <template #item.success="{ item }">
          <v-chip :color="getStatusColor(item)" size="small" variant="tonal">
            {{ getStatusText(item) }}
          </v-chip>
        </template>
        <template #item.created_at="{ item }">
          {{ new Date(item.created_at).toLocaleString() }}
        </template>
      </v-data-table>
    </v-sheet>
  </v-container>
</template>

<script setup lang="ts">
  import axios from 'axios'
  import { onMounted, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'

  const router = useRouter()
  const route = useRoute()

  interface Project {
    id: number
    name: string
    created_at: string
    dir?: string
    git_url?: string
    git_branch?: string
    stages?: Array<{
      number: number
      script: string
    }>
  }

  interface ProjectRun {
    id: number
    created_at: string
    project_id: number
    success: boolean
    pending: boolean
    processing: boolean
    fail_log: string
  }

  const runsHeaders = [
    { title: 'ID', key: 'id', sortable: true },
    { title: 'Статус', key: 'success', sortable: true },
    { title: 'Дата', key: 'created_at', sortable: true },
  ]

  function editProject (id: number) {
    router.push(`/projects/${id}/edit`)
  }

  function handleRowClick (_event: Event, { item }: { item: ProjectRun }) {
    router.push(`/projects/${route.params.id}/runs/${item.id}`)
  }

  const running = ref(false)
  const runs = ref<ProjectRun[]>([])
  const runsLoading = ref(false)

  function runProject (id: number) {
    running.value = true
    axios
      .post(`/api/v1/projects/${id}/runs`)
      .catch(error => {
        console.error(error)
      })
      .finally(() => {
        loadRuns()
        running.value = false
      })
  }

  function loadRuns () {
    runsLoading.value = true
    axios
      .get(`/api/v1/projects/${route.params.id}/runs`)
      .then(response => {
        runs.value = response.data.runs
        runs.value.sort((a: ProjectRun, b: ProjectRun) => {
          return b.id - a.id
        })
      })
      .catch(error => {
        console.error(error)
      })
      .finally(() => {
        runsLoading.value = false
      })
  }

  const project = ref<Project>()

  function load () {
    axios
      .get(`/api/v1/projects/${route.params.id}`)
      .then(v => {
        project.value = v.data
      })
      .catch(error => {
        console.error(error)
      })
  }

  function getStatusColor (item: ProjectRun): string {
    if (item.pending) return 'warning'
    if (item.processing) return 'info'
    return item.success ? 'success' : 'error'
  }

  function getStatusText (item: ProjectRun): string {
    if (item.pending) return 'Ожидание'
    if (item.processing) return 'Выполняется'
    return item.success ? 'Успешно' : 'Ошибка'
  }

  onMounted(() => {
    load()
    loadRuns()
  })
</script>
