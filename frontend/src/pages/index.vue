<template>
  <v-container class="d-flex ga-4 flex-column">
    <div class="d-flex flex-row justify-space-between align-center">
      <h1>Проекты</h1>
      <v-btn prepend-icon="mdi-plus" @click="createProject">
        Создать проект
      </v-btn>
    </div>
    <v-card v-for="project in projects">
      <v-card-text class="d-flex flex-row justify-space-between align-center">
        <b>#{{ project.id }} {{ project.name }}</b>
        <span class="d-flex flex-row ga-2">
          <v-btn prepend-icon="mdi-eye" @click="openProject(project.id)">
            Посмотреть
          </v-btn>
          <v-btn prepend-icon="mdi-pencil" @click="editProject(project.id)">
            Редактировать
          </v-btn>
        </span>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
  import axios from 'axios'
  import { onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'

  const router = useRouter()

  interface Project {
    id: number
    name: string
  }

  const projects = ref<Array<Project>>(new Array<Project>())

  function load () {
    axios
      .get('/api/v1/projects')
      .then(v => {
        projects.value = v.data.projects
      })
      .catch(error => {
        console.log(error)
      })
  }

  function openProject (id: number) {
    router.push(`/projects/${id}`)
  }

  function editProject (id: number) {
    router.push(`/projects/${id}/edit`)
  }

  function createProject () {
    router.push('/projects/new')
  }

  onMounted(() => {
    load()
  })
</script>
