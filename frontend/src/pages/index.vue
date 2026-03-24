<template>
  <v-container class="d-flex ga-4 flex-column">
    <div class="d-flex flex-row justify-space-between align-center">
      <h1>Проекты</h1>
      <v-btn prepend-icon="mdi-plus" @click="createProject"> Создать проект </v-btn>
    </div>
    <v-card v-for="project in projects" :key="project.id">
      <v-card-text class="d-flex flex-column ga-2">
        <div class="d-flex flex-row justify-space-between align-center">
          <b>#{{ project.id }} {{ project.name }}</b>
          <span class="d-flex flex-row ga-2">
            <v-btn prepend-icon="mdi-eye" @click="openProject(project.id)"> Посмотреть </v-btn>
            <v-btn prepend-icon="mdi-pencil" @click="editProject(project.id)">
              Редактировать
            </v-btn>
          </span>
        </div>
        <div class="d-flex flex-row ga-4 text-caption text-medium-emphasis">
          <span v-if="project.cron_enabled">
            <v-icon class="mr-1" icon="mdi-clock-outline" size="small" />
            Cron включен
          </span>
          <span v-if="project.last_run">
            <v-icon
              class="mr-1"
              :color="getLastRunColor(project.last_run)"
              :icon="getLastRunIcon(project.last_run)"
              size="small"
            />
            Последний запуск: {{ formatLastRun(project.last_run) }}
            <span
              v-if="
                project.last_run.duration &&
                !project.last_run.processing &&
                !project.last_run.pending
              "
              class="text-medium-emphasis ml-1"
            >
              ({{ formatDuration(project.last_run.duration) }})
            </span>
          </span>
          <span v-if="project.last_successful_run_at">
            <v-icon class="mr-1" color="success" icon="mdi-check-circle-outline" size="small" />
            Последний успешный: {{ formatDateTime(project.last_successful_run_at) }}
          </span>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import type { ProjectLastRun, ProjectListItem } from '@/api/generated.schemas';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { getEasyJetAPI } from '@/api/generated';
import { formatDuration } from '@/utils/formatDuration';

const router = useRouter();
const api = getEasyJetAPI();

const projects = ref<Array<ProjectListItem>>(new Array<ProjectListItem>());

function load() {
  api
    .getProjects()
    .then((v) => {
      const projectsData = v.data.projects;
      if (projectsData) {
        projects.value = projectsData;
      }
    })
    .catch((error) => {
      console.log(error);
    });
}

function openProject(id: number) {
  router.push(`/projects/${id}`);
}

function editProject(id: number) {
  router.push(`/projects/${id}/edit`);
}

function createProject() {
  router.push('/projects/new');
}

function getLastRunIcon(lastRun: ProjectLastRun): string {
  if (lastRun.processing) {
    return 'mdi-loading';
  }
  if (lastRun.pending) {
    return 'mdi-clock-outline';
  }
  if (lastRun.success) {
    return 'mdi-check-circle-outline';
  }
  return 'mdi-alert-circle-outline';
}

function getLastRunColor(lastRun: ProjectLastRun): string {
  if (lastRun.processing) {
    return 'info';
  }
  if (lastRun.pending) {
    return 'warning';
  }
  if (lastRun.success) {
    return 'success';
  }
  return 'error';
}

function formatLastRun(lastRun: ProjectLastRun): string {
  if (lastRun.processing) {
    return 'Выполняется...';
  }
  if (lastRun.pending) {
    return 'В очереди...';
  }
  if (lastRun.success) {
    return `Успешно ${formatDateTime(lastRun.created_at)}`;
  }
  return `Ошибка ${formatDateTime(lastRun.created_at)}`;
}

function formatDateTime(dateString?: string): string {
  if (!dateString) {
    return '';
  }
  const date = new Date(dateString);
  return date.toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

onMounted(() => {
  load();
});
</script>
