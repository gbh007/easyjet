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
        <v-btn prepend-icon="mdi-play" @click="runProject(project.id)">
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

function editProject(id: number) {
  router.push(`/projects/${id}/edit`);
}

function runProject(id: number) {
  console.log(id);
}

const project = ref<Project>();

function load() {
  axios
    .get(`/api/v1/projects/${route.params.id}`)
    .then((v) => {
      project.value = v.data;
    })
    .catch((err) => {
      console.log(err);
    });
}

onMounted(() => {
  load();
});
</script>
