<template>
  <v-container class="d-flex ga-4 flex-column" v-if="project">
    <v-sheet
      elevation="2"
      class="d-flex pa-4 flex-row justify-space-between align-center"
    >
      <v-sheet class="d-flex pa-4 flex-column">
        <v-sheet class="d-flex flex-row justify-space-between align-center">
          <h2>#{{ project.id }} {{ project.name }}</h2>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Создан</b>
          <span>{{ new Date(project.created_at).toLocaleString() }}</span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2" v-if="project.dir">
          <b>Директория</b>
          <span>{{ project.dir }}</span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2" v-if="project.git_url">
          <b>GIT</b>
          <span>{{ project.git_url }}</span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2" v-if="project.git_branch">
          <b>GIT branch</b>
          <span>{{ project.git_branch }}</span>
        </v-sheet>
      </v-sheet>
      <v-sheet class="d-flex ga-2 flex-column">
        <v-btn prepend-icon="mdi-pencil" @click="editProject(project.id)">
          Редактировать
        </v-btn>
        <v-btn
          prepend-icon="mdi-play"
          :loading="running"
          :disabled="running"
          @click="runProject(project.id)"
        >
          Запустить
        </v-btn>
      </v-sheet>
    </v-sheet>

    <v-sheet
      elevation="2"
      class="d-flex pa-4 flex-column"
      v-for="stage in project.stages"
    >
      <h3>Этап {{ stage.number }}</h3>
      <pre>{{ stage.script }}</pre>
    </v-sheet>

    <v-sheet elevation="2" class="d-flex pa-4 flex-column">
      <div class="d-flex flex-row justify-space-between align-center">
        <h3>История запусков</h3>
        <v-btn
          icon="mdi-refresh"
          variant="text"
          size="small"
          @click="loadRuns"
        />
      </div>
      <v-data-table
        :headers="runsHeaders"
        :items="runs"
        item-value="id"
        :loading="runsLoading"
        hover
        @click:row="handleRowClick"
      >
        <template v-slot:item.success="{ item }">
          <v-chip
            :color="item.success ? 'success' : 'error'"
            size="small"
            variant="tonal"
          >
            {{ item.success ? "Успешно" : "Ошибка" }}
          </v-chip>
        </template>
        <template v-slot:item.created_at="{ item }">
          {{ new Date(item.created_at).toLocaleString() }}
        </template>
      </v-data-table>
    </v-sheet>
  </v-container>
</template>

<script setup lang="ts">
import axios from "axios";
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const router = useRouter();
const route = useRoute();

interface Project {
  id: number;
  name: string;
  created_at: string;
  dir?: string;
  git_url?: string;
  git_branch?: string;
  stages?: Array<{
    number: number;
    script: string;
  }>;
}

interface ProjectRun {
  id: number;
  created_at: string;
  project_id: number;
  success: boolean;
  fail_log: string;
}

const runsHeaders = [
  { title: "ID", key: "id", sortable: true },
  { title: "Статус", key: "success", sortable: true },
  { title: "Дата", key: "created_at", sortable: true },
];

function editProject(id: number) {
  router.push(`/projects/${id}/edit`);
}

function handleRowClick(_event: Event, { item }: { item: ProjectRun }) {
  router.push(`/projects/${route.params.id}/runs/${item.id}`);
}

const running = ref(false);
const runs = ref<ProjectRun[]>([]);
const runsLoading = ref(false);

function runProject(id: number) {
  running.value = true;
  axios
    .post(`/api/v1/projects/${id}/runs`)
    .catch((err) => {
      console.error(err);
    })
    .finally(() => {
      loadRuns();
      running.value = false;
    });
}

function loadRuns() {
  runsLoading.value = true;
  axios
    .get(`/api/v1/projects/${route.params.id}/runs`)
    .then((response) => {
      runs.value = response.data.runs;
      runs.value.sort((a: ProjectRun, b: ProjectRun) => {
        return b.id - a.id;
      });
    })
    .catch((err) => {
      console.error(err);
    })
    .finally(() => {
      runsLoading.value = false;
    });
}

const project = ref<Project>();

function load() {
  axios
    .get(`/api/v1/projects/${route.params.id}`)
    .then((v) => {
      project.value = v.data;
    })
    .catch((err) => {
      console.error(err);
    });
}

onMounted(() => {
  load();
  loadRuns();
});
</script>
