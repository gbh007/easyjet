<template>
  <v-container v-if="run" class="d-flex ga-4 flex-column">
    <v-sheet class="d-flex pa-4 flex-row justify-space-between align-center" elevation="2">
      <v-sheet class="d-flex pa-4 flex-column">
        <v-sheet class="d-flex flex-row justify-space-between align-center">
          <h2>Запуск #{{ run.id }}</h2>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Проект</b>
          <span>#{{ project?.id }} {{ project?.name }}</span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Дата</b>
          <span>{{ new Date(run.created_at).toLocaleString() }}</span>
        </v-sheet>
        <v-sheet class="d-flex flex-row ga-2">
          <b>Статус</b>
          <v-chip :color="getStatusColor(run)" size="small" variant="tonal">
            {{ getStatusText(run) }}
          </v-chip>
        </v-sheet>
        <v-sheet
          v-if="run.duration && !run.processing && !run.pending"
          class="d-flex flex-row ga-2"
        >
          <b>Длительность</b>
          <span>{{ formatDuration(run.duration) }}</span>
        </v-sheet>
      </v-sheet>
      <v-sheet class="d-flex ga-2 flex-column">
        <v-btn variant="outlined" @click="goBack"> Назад </v-btn>
      </v-sheet>
    </v-sheet>

    <v-alert v-if="!run.success && run.fail_log" class="mb-2" title="Ошибка" type="error">
      {{ run.fail_log }}
    </v-alert>

    <v-sheet
      v-if="run.git_commits && run.git_commits.length > 0"
      class="d-flex pa-4 flex-column"
      elevation="2"
    >
      <h3>Git изменения</h3>
      <v-list>
        <v-list-item v-for="commit in run.git_commits" :key="commit.number">
          <template #prepend>
            <v-icon icon="mdi-git" size="small" />
          </template>
          <v-list-item-title>{{ commit.subject }}</v-list-item-title>
          <v-list-item-subtitle>{{ commit.hash }}</v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-sheet>

    <v-sheet
      v-if="run.stages && run.stages.length > 0"
      class="d-flex pa-4 flex-column"
      elevation="2"
    >
      <h3>Результаты выполнения этапов</h3>
      <v-expansion-panels>
        <v-expansion-panel v-for="stage in run.stages" :key="stage.stage_number">
          <v-expansion-panel-title>
            <div class="d-flex align-center ga-2">
              <v-icon
                :color="stage.success ? 'success' : 'error'"
                :icon="stage.success ? 'mdi-check-circle' : 'mdi-alert-circle'"
              />
              Этап {{ stage.stage_number }}
              <v-chip
                v-if="stage.duration"
                class="ml-2"
                color="primary"
                size="x-small"
                variant="tonal"
              >
                {{ formatDuration(stage.duration) }}
              </v-chip>
            </div>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <pre class="pa-2 rounded code-block">{{ stage.log }}</pre>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
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
}

interface ProjectRunStage {
  stage_number: number;
  success: boolean;
  log: string;
  duration?: number;
}

interface ProjectRunGitCommit {
  number: number;
  hash: string;
  subject: string;
}

interface ProjectRun {
  id: number;
  created_at: string;
  project_id: number;
  success: boolean;
  pending: boolean;
  processing: boolean;
  fail_log: string;
  duration?: number;
  stages?: ProjectRunStage[];
  git_commits?: ProjectRunGitCommit[];
}

const run = ref<ProjectRun | null>(null);
const project = ref<Project | null>(null);
const loading = ref(false);

function goBack() {
  router.push(`/projects/${route.params.project_id}`);
}

function load() {
  loading.value = true;
  api
    .getProjectRun(Number(route.params.project_id), Number(route.params.run_id))
    .then((response) => {
      const data = response.data;
      run.value = {
        id: data.id ?? 0,
        created_at: data.created_at ?? '',
        project_id: data.project_id ?? 0,
        success: data.success ?? false,
        pending: data.pending ?? false,
        processing: data.processing ?? false,
        fail_log: data.fail_log ?? '',
        duration: data.duration ?? 0,
        stages: data.stages?.map((s) => ({
          stage_number: s.stage_number ?? 0,
          success: s.success ?? false,
          log: s.log ?? '',
          duration: s.duration ?? 0,
        })),
        git_commits: data.git_commits?.map((c) => ({
          number: c.number ?? 0,
          hash: c.hash ?? '',
          subject: c.subject ?? '',
        })),
      };
      return api.getProject(Number(route.params.project_id));
    })
    .then((response) => {
      const data = response.data;
      project.value = {
        id: data.id ?? 0,
        name: data.name ?? '',
      };
    })
    .catch((error) => {
      console.error(error);
    })
    .finally(() => {
      loading.value = false;
    });
}

function getStatusColor(run: ProjectRun): string {
  if (run.pending) return 'warning';
  if (run.processing) return 'info';
  return run.success ? 'success' : 'error';
}

function getStatusText(run: ProjectRun): string {
  if (run.pending) return 'Ожидание';
  if (run.processing) return 'Выполняется';
  return run.success ? 'Успешно' : 'Ошибка';
}

onMounted(() => {
  load();
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
