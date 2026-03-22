<template>
  <v-container v-if="project" class="d-flex ga-4 flex-column">
    <v-sheet class="d-flex pa-4 flex-row justify-space-between align-center" elevation="2">
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
        <v-sheet class="d-flex flex-row ga-2">
          <b>Cron</b>
          <span>
            <v-chip
              v-if="project.cron_enabled && project.cron_schedule"
              class="mr-2"
              color="success"
              size="small"
              variant="tonal"
            >
              Активно
            </v-chip>
            <v-chip
              v-else-if="!project.cron_enabled && project.cron_schedule"
              class="mr-2"
              color="warning"
              size="small"
              variant="tonal"
            >
              Пауза
            </v-chip>
            <span v-if="project.cron_schedule">{{ project.cron_schedule }}</span>
            <span v-else>Не настроено</span>
          </span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Ротация запусков</b>
          <span>
            <v-chip
              v-if="project.retention_count && project.retention_count > 0"
              class="mr-2"
              color="success"
              size="small"
              variant="tonal"
            >
              Активно
            </v-chip>
            <v-chip v-else class="mr-2" color="secondary" size="small" variant="tonal">
              Отключена
            </v-chip>
            <span v-if="project.retention_count && project.retention_count > 0">
              Сохранять последних {{ project.retention_count }} запусков
            </span>
          </span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Перезапуск после задачи</b>
          <span>
            <v-chip
              v-if="project.restart_after"
              class="mr-2"
              color="success"
              size="small"
              variant="tonal"
            >
              Включено
            </v-chip>
            <v-chip v-else class="mr-2" color="secondary" size="small" variant="tonal">
              Отключено
            </v-chip>
          </span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Использовать окружение хоста</b>
          <span>
            <v-chip
              v-if="project.with_root_env"
              class="mr-2"
              color="success"
              size="small"
              variant="tonal"
            >
              Включено
            </v-chip>
            <v-chip v-else class="mr-2" color="secondary" size="small" variant="tonal">
              Отключено
            </v-chip>
          </span>
        </v-sheet>
      </v-sheet>
      <v-sheet class="d-flex ga-2 flex-column">
        <v-btn prepend-icon="mdi-pencil" @click="editProject(project.id)"> Редактировать </v-btn>
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
      :key="stage.number"
      class="d-flex pa-4 flex-column"
      elevation="2"
    >
      <h3>Этап {{ stage.number }}</h3>
      <pre>{{ stage.script }}</pre>
    </v-sheet>

    <v-sheet class="d-flex pa-4 flex-column" elevation="2">
      <div class="d-flex flex-row justify-space-between align-center">
        <h3>История запусков</h3>
        <v-btn icon="mdi-refresh" size="small" variant="text" @click="loadRuns" />
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
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getEasyJetAPI } from '@/api';

const router = useRouter();
const route = useRoute();
const api = getEasyJetAPI();

interface Project {
  id: number;
  name: string;
  created_at: string;
  dir?: string;
  git_url?: string;
  git_branch?: string;
  cron_enabled: boolean;
  cron_schedule: string;
  restart_after: boolean;
  with_root_env: boolean;
  retention_count: number;
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
  pending: boolean;
  processing: boolean;
  fail_log: string;
}

const runsHeaders = [
  { title: 'ID', key: 'id', sortable: true },
  { title: 'Статус', key: 'success', sortable: true },
  { title: 'Дата', key: 'created_at', sortable: true },
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
  api
    .createProjectRun(id)
    .catch((error) => {
      console.error(error);
    })
    .finally(() => {
      loadRuns();
      running.value = false;
    });
}

function loadRuns() {
  runsLoading.value = true;
  api
    .getProjectRuns(Number(route.params.id))
    .then((response) => {
      const runsData = response.data.runs;
      if (runsData) {
        runs.value = runsData.map((r) => ({
          id: r.id ?? 0,
          created_at: r.created_at ?? '',
          project_id: r.project_id ?? 0,
          success: r.success ?? false,
          pending: r.pending ?? false,
          processing: r.processing ?? false,
          fail_log: r.fail_log ?? '',
        }));
      }
      runs.value.sort((a: ProjectRun, b: ProjectRun) => {
        return b.id - a.id;
      });
    })
    .catch((error) => {
      console.error(error);
    })
    .finally(() => {
      runsLoading.value = false;
    });
}

const project = ref<Project>();

function load() {
  api
    .getProject(Number(route.params.id))
    .then((v) => {
      const data = v.data;
      project.value = {
        id: data.id ?? 0,
        name: data.name ?? '',
        created_at: data.created_at ?? '',
        dir: data.dir,
        git_url: data.git_url,
        git_branch: data.git_branch,
        cron_enabled: data.cron_enabled ?? false,
        cron_schedule: data.cron_schedule ?? '',
        restart_after: data.restart_after ?? false,
        with_root_env: data.with_root_env ?? false,
        retention_count: data.retention_count ?? 0,
        stages: data.stages?.map((s) => ({
          number: s.number ?? 0,
          script: s.script ?? '',
        })),
      };
    })
    .catch((error) => {
      console.error(error);
    });
}

function getStatusColor(item: ProjectRun): string {
  if (item.pending) return 'warning';
  if (item.processing) return 'info';
  return item.success ? 'success' : 'error';
}

function getStatusText(item: ProjectRun): string {
  if (item.pending) return 'Ожидание';
  if (item.processing) return 'Выполняется';
  return item.success ? 'Успешно' : 'Ошибка';
}

onMounted(() => {
  load();
  loadRuns();
});
</script>
