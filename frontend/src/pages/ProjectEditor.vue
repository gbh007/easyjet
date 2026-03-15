<template>
  <v-container class="d-flex ga-4 flex-column" v-if="form">
    <v-sheet elevation="2" class="d-flex pa-4 flex-column">
      <h2>
        {{ isEdit ? "Редактирование" : "Создание" }} проекта{{
          isEdit ? " #" + form.id : ""
        }}
      </h2>

      <v-text-field
        v-model="form.name"
        label="Название"
        class="mt-4"
        required
      />

      <v-text-field v-model="form.dir" label="Директория" class="mt-2" />

      <v-text-field v-model="form.git_url" label="Git URL" class="mt-2" />

      <v-text-field v-model="form.git_branch" label="Git branch" class="mt-2" />

      <div class="d-flex flex-row justify-space-between align-center mt-4">
        <h3>Этапы</h3>
        <v-btn prepend-icon="mdi-plus" @click="addStage"> Добавить этап </v-btn>
      </div>

      <v-sheet
        v-for="(stage, index) in form.stages"
        :key="index"
        elevation="1"
        class="d-flex pa-4 flex-column mt-2"
      >
        <div class="d-flex flex-row justify-space-between align-center">
          <span class="text-h6">Этап {{ index + 1 }}</span>
          <v-btn
            icon="mdi-delete"
            size="small"
            variant="text"
            color="error"
            @click="removeStage(index)"
          />
        </div>
        <v-textarea
          v-model="stage.script"
          label="Скрипт"
          rows="4"
          class="mt-2"
          monospace
        />
      </v-sheet>

      <div class="d-flex flex-row ga-2 mt-4">
        <v-btn color="success" @click="saveProject"> Сохранить </v-btn>
        <v-btn variant="outlined" @click="cancel"> Отмена </v-btn>
      </div>
    </v-sheet>
  </v-container>
</template>

<script setup lang="ts">
import axios from "axios";
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const router = useRouter();
const route = useRoute();

interface Stage {
  number: number;
  script: string;
}

interface ProjectForm {
  id: number;
  name: string;
  dir?: string;
  git_url?: string;
  git_branch?: string;
  stages: Stage[];
}

const form = ref<ProjectForm | null>(null);

const isEdit = computed(() => {
  return form.value?.id !== undefined && form.value.id !== 0;
});

function addStage() {
  if (!form.value) return;
  const nextNumber =
    form.value.stages.length > 0
      ? Math.max(...form.value.stages.map((s) => s.number)) + 1
      : 1;
  form.value.stages.push({ number: nextNumber, script: "" });
}

function removeStage(index: number) {
  if (!form.value) return;
  form.value.stages.splice(index, 1);
}

function saveProject() {
  if (!form.value) return;
  if (isEdit.value) {
    axios
      .put(`/api/v1/projects/${form.value.id}`, form.value)
      .then(() => {
        router.push(`/projects/${form.value!.id}`);
      })
      .catch((err) => {
        console.error(err);
      });
  } else {
    axios
      .post(`/api/v1/projects`, form.value)
      .then((response) => {
        router.push(`/projects/${response.data.id}`);
      })
      .catch((err) => {
        console.error(err);
      });
  }
}

function cancel() {
  if (isEdit.value && form.value) {
    router.push(`/projects/${form.value.id}`);
  } else {
    router.push("/");
  }
}

function load() {
  const id = route.params.id;
  if (!id || id === "new") {
    form.value = {
      id: 0,
      name: "",
      dir: "",
      git_url: "",
      git_branch: "",
      stages: [],
    };
    return;
  }

  axios
    .get(`/api/v1/projects/${id}`)
    .then((v) => {
      const data = v.data;
      form.value = {
        id: data.id,
        name: data.name,
        dir: data.dir || "",
        git_url: data.git_url || "",
        git_branch: data.git_branch || "",
        stages: data.stages?.map((s: Stage) => ({ ...s })) || [],
      };
    })
    .catch((err) => {
      console.error(err);
    });
}

onMounted(() => {
  load();
});
</script>
