<template>
  <v-container v-if="project" class="d-flex ga-4 flex-column">
    <v-sheet class="d-flex pa-4 flex-row justify-space-between align-center" elevation="2">
      <v-sheet class="d-flex pa-4 flex-column">
        <v-sheet class="d-flex flex-row justify-space-between align-center">
          <div class="d-flex flex-row align-center ga-2">
            <h2>#{{ project.id }} {{ project.name }}</h2>
            <v-chip v-if="project.is_template" color="purple" size="small" variant="tonal">
              <v-icon icon="mdi-content-save-cog-outline" size="small" />
              Шаблон
            </v-chip>
          </div>
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
          :disabled="running || project.is_template"
          :loading="running"
          prepend-icon="mdi-play"
          @click="runProject(project.id)"
        >
          Запустить
        </v-btn>
      </v-sheet>
    </v-sheet>

    <v-sheet
      v-if="project.env_vars && project.env_vars.length > 0"
      class="d-flex pa-4 flex-column"
      elevation="2"
    >
      <h3>Переменные окружения</h3>
      <v-table>
        <thead>
          <tr>
            <th>Имя</th>
            <th>Значение</th>
            <th>Использует другие переменные</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="envVar in project.env_vars" :key="envVar.id">
            <td>{{ envVar.name }}</td>
            <td>
              <code>{{ envVar.value }}</code>
            </td>
            <td>
              <v-chip
                v-if="envVar.uses_other_variables"
                class="mr-2"
                color="success"
                size="small"
                variant="tonal"
              >
                Да
              </v-chip>
              <v-chip v-else class="mr-2" color="secondary" size="small" variant="tonal">
                Нет
              </v-chip>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-sheet>

    <v-sheet class="d-flex pa-4 flex-column" elevation="2">
      <h3>Этапы</h3>
      <v-expansion-panels>
        <v-expansion-panel v-for="stage in project.stages" :key="stage.number">
          <v-expansion-panel-title> Этап {{ stage.number }} </v-expansion-panel-title>
          <v-expansion-panel-text>
            <pre class="pa-2 rounded code-block">{{ stage.script }}</pre>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
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
        <template #item.status="{ item }">
          <v-chip :color="getStatusColor(item)" size="small" variant="tonal">
            {{ getStatusText(item) }}
          </v-chip>
        </template>
        <template #item.duration="{ item }">
          {{ item.duration ? formatDuration(item.duration) : '—' }}
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
import { getEasyJetAPI } from '@/api/generated';
import { formatDuration } from '@/utils/formatDuration';

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
  is_template: boolean;
  stages?: Array<{
    number: number;
    script: string;
  }>;
  env_vars?: Array<{
    id?: number;
    name: string;
    value: string;
    uses_other_variables?: boolean;
  }>;
}

interface ProjectRun {
  id: number;
  created_at: string;
  project_id: number;
  status: string;
  fail_log: string;
  duration?: number;
}

const runsHeaders = [
  { title: 'ID', key: 'id', sortable: true },
  { title: 'Статус', key: 'status', sortable: true },
  { title: 'Длительность', key: 'duration', sortable: true },
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
          status: r.status ?? 'failed',
          fail_log: r.fail_log ?? '',
          duration: r.duration ?? 0,
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
        is_template: data.is_template ?? false,
        stages: data.stages?.map((s) => ({
          number: s.number ?? 0,
          script: s.script ?? '',
        })),
        env_vars: data.env_vars?.map((ev) => ({
          id: ev.id,
          name: ev.name ?? '',
          value: ev.value ?? '',
          uses_other_variables: ev.uses_other_variables ?? false,
        })),
      };
    })
    .catch((error) => {
      console.error(error);
    });
}

function getStatusColor(item: ProjectRun): string {
  if (item.status === 'pending') return 'warning';
  if (item.status === 'processing') return 'info';
  if (item.status === 'success') return 'success';
  return 'error';
}

function getStatusText(item: ProjectRun): string {
  if (item.status === 'pending') return 'Ожидание';
  if (item.status === 'processing') return 'Выполняется';
  if (item.status === 'success') return 'Успешно';
  return 'Ошибка';
}

onMounted(() => {
  load();
  loadRuns();
});
</script>

<style scoped>
.code-block {
  max-height: 400px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
