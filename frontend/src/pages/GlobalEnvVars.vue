<template>
  <v-container class="d-flex ga-4 flex-column">
    <v-sheet class="d-flex pa-4 flex-column" elevation="2">
      <h2>Глобальные переменные окружения</h2>

      <p class="text-body-2 mt-2">
        Глобальные переменные окружения доступны всем проектам в системе. Переменные проекта имеют
        более высокий приоритет и могут переопределять глобальные переменные.
      </p>

      <v-sheet class="d-flex pa-4 flex-column">
        <EnvVarsInfo />
      </v-sheet>

      <div class="d-flex flex-row ga-2 mt-4 mb-2">
        <v-btn color="primary" @click="addEnvVar">
          <v-icon icon="mdi-plus" />
          Добавить переменную
        </v-btn>
      </div>

      <v-table v-if="envVars.length > 0" density="comfortable" variant="outlined">
        <thead>
          <tr>
            <th class="text-left">Имя</th>
            <th class="text-left">Значение</th>
            <th class="text-center">Использует другие</th>
            <th class="text-right">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(envVar, index) in envVars" :key="envVar.id || `new-${index}`">
            <td>
              <v-text-field
                v-if="!envVar.id || editingIndex === index"
                v-model="envVar.name"
                density="compact"
                hide-details
                label="Имя"
                placeholder="VAR_NAME"
                variant="outlined"
              />
              <span v-else class="font-weight-medium">{{ envVar.name }}</span>
            </td>
            <td>
              <v-text-field
                v-if="!envVar.id || editingIndex === index"
                v-model="envVar.value"
                density="compact"
                hide-details
                label="Значение"
                placeholder="value"
                variant="outlined"
              />
              <span v-else>{{ envVar.value }}</span>
            </td>
            <td class="text-center">
              <v-checkbox
                v-model="envVar.uses_other_variables"
                class="d-flex justify-center"
                density="compact"
                :disabled="!!envVar.id && editingIndex !== index"
                hide-details
              />
            </td>
            <td class="text-right">
              <v-btn
                v-if="!envVar.id"
                class="mr-2"
                color="success"
                icon="mdi-content-save"
                size="small"
                variant="text"
                @click="saveNewEnvVar(index)"
              />
              <v-btn
                v-if="envVar.id && editingIndex === index"
                class="mr-2"
                color="success"
                icon="mdi-content-save"
                size="small"
                variant="text"
                @click="saveEdit(index)"
              />
              <v-btn
                v-if="envVar.id && editingIndex === index"
                class="mr-2"
                color="grey"
                icon="mdi-close"
                size="small"
                variant="text"
                @click="cancelEdit(index)"
              />
              <v-btn
                v-if="envVar.id && editingIndex !== index"
                class="mr-2"
                color="primary"
                icon="mdi-pencil"
                size="small"
                variant="text"
                @click="startEdit(index)"
              />
              <v-btn
                color="error"
                icon="mdi-delete"
                size="small"
                variant="text"
                @click="deleteEnvVar(index)"
              />
            </td>
          </tr>
        </tbody>
      </v-table>

      <v-alert v-else class="mt-4" text="Нет глобальных переменных окружения." type="info" />
    </v-sheet>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { getEasyJetAPI } from '@/api/generated';
import EnvVarsInfo from '@/components/EnvVarsInfo.vue';

const api = getEasyJetAPI();

interface EnvVar {
  id?: number;
  name: string;
  value: string;
  uses_other_variables: boolean;
}

const envVars = ref<EnvVar[]>([]);
const editingIndex = ref<number | null>(null);

function addEnvVar() {
  envVars.value.push({ name: '', value: '', uses_other_variables: false });
}

function saveNewEnvVar(index: number) {
  const envVar = envVars.value[index];
  const payload = {
    name: envVar.name,
    value: envVar.value,
    uses_other_variables: envVar.uses_other_variables,
  };

  api
    .createGlobalEnvVar(payload)
    .then((response) => {
      envVar.id = response.data.id;
    })
    .catch((error) => {
      console.error(error);
    });
}

function startEdit(index: number) {
  editingIndex.value = index;
}

function saveEdit(index: number) {
  const envVar = envVars.value[index];
  if (!envVar.id) return;

  const payload = {
    name: envVar.name,
    value: envVar.value,
    uses_other_variables: envVar.uses_other_variables,
  };

  api
    .updateGlobalEnvVar(envVar.id, payload)
    .then(() => {
      editingIndex.value = null;
    })
    .catch((error) => {
      console.error(error);
    });
}

function cancelEdit(index: number) {
  // Reload to restore original values
  load();
}

function deleteEnvVar(index: number) {
  const envVar = envVars.value[index];
  if (envVar.id) {
    api
      .deleteGlobalEnvVar(envVar.id)
      .then(() => {
        envVars.value.splice(index, 1);
      })
      .catch((error) => {
        console.error(error);
      });
  } else {
    envVars.value.splice(index, 1);
  }
}

function load() {
  api
    .getGlobalEnvVars()
    .then((v) => {
      const data = v.data.env_vars || [];
      envVars.value = data.map((ev) => ({
        id: ev.id,
        name: ev.name,
        value: ev.value,
        uses_other_variables: ev.uses_other_variables ?? false,
      }));
    })
    .catch((error) => {
      console.error(error);
    });
}

onMounted(() => {
  load();
});
</script>
