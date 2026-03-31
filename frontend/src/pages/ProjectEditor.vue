<template>
  <v-container v-if="form" class="d-flex ga-4 flex-column">
    <v-sheet class="d-flex pa-4 flex-column" elevation="2">
      <h2>
        {{ isEdit ? 'Редактирование' : 'Создание' }} проекта{{ isEdit ? ' #' + form.id : '' }}
      </h2>

      <v-alert
        class="mt-4"
        density="compact"
        text="Выберите шаблон для быстрого заполнения настроек проекта"
        type="info"
        variant="tonal"
      />

      <div class="d-flex flex-row ga-2 align-center mt-2">
        <v-select
          v-model="selectedTemplateId"
          density="compact"
          hide-details
          item-title="name"
          item-value="id"
          :items="templates"
          label="Применить шаблон"
          style="max-width: 400px"
          variant="outlined"
        />
        <v-btn
          color="primary"
          :disabled="!selectedTemplateId"
          variant="tonal"
          @click="applyTemplate"
        >
          Применить
        </v-btn>
      </div>

      <v-divider class="mt-4 mb-2" />

      <v-text-field v-model="form.name" class="mt-4" label="Название" required />

      <v-text-field v-model="form.dir" class="mt-2" label="Директория" />

      <v-text-field v-model="form.git_url" class="mt-2" label="Git URL" />

      <v-text-field v-model="form.git_branch" class="mt-2" label="Git branch" />

      <div class="d-flex flex-row align-center mt-4">
        <v-switch
          v-model="form.cron_enabled"
          class="mr-2"
          color="primary"
          hide-details
          label="Включить расписание"
        />
        <v-text-field
          v-model="form.cron_schedule"
          class="flex-grow-1"
          :disabled="!form.cron_enabled"
          hint="Формат: минута час день месяц день_недели (например, 0 5 * * * = ежедневно в 5:00)"
          label="Cron расписание"
          persistent-hint
          placeholder="0 5 * * *"
          :rules="cronRules"
        />
      </div>

      <v-switch
        v-model="form.restart_after"
        class="mt-4"
        color="primary"
        hide-details
        label="Перезапускать приложение после выполнения задачи"
      />

      <v-switch
        v-model="form.with_root_env"
        class="mt-4"
        color="primary"
        hide-details
        label="Использовать переменные окружения хоста (withRootEnv)"
      />

      <v-switch
        v-model="form.is_template"
        class="mt-4"
        color="purple"
        hide-details
        label="Это шаблон (не может быть запущен)"
      />

      <v-text-field
        v-model.number="form.retention_count"
        class="flex-grow-1"
        hint="Количество последних запусков для хранения (0 - отключить)"
        label="Сохранять запусков"
        min="0"
        persistent-hint
        type="number"
      />

      <v-divider class="mt-4 mb-2" />

      <div class="d-flex flex-row justify-space-between align-center mt-4">
        <h3>Этапы</h3>
        <v-btn prepend-icon="mdi-plus" @click="addStage"> Добавить этап </v-btn>
      </div>

      <v-sheet class="d-flex pa-4 flex-column">
        <v-expansion-panels variant="accordion">
          <v-expansion-panel>
            <v-expansion-panel-title class="d-flex align-center">
              <v-icon class="mr-2" color="info" icon="mdi-information-outline" size="small" />
              <span class="text-subtitle-1 font-weight-bold">Переменные окружения</span>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <div class="text-body-2 mt-2">
                В скриптах этапов доступны следующие переменные окружения:
              </div>
              <v-table class="mt-2" density="compact" variant="outlined">
                <thead>
                  <tr>
                    <th class="text-left">Переменная</th>
                    <th class="text-left">Описание</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><code class="text-primary">WORKSPACE</code></td>
                    <td>Рабочая директория проекта (путь к папке на сервере)</td>
                  </tr>
                </tbody>
              </v-table>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-sheet>

      <v-sheet class="d-flex pa-4 flex-column mt-2" color="grey-lighten-4" variant="outlined">
        <div class="d-flex flex-row justify-space-between align-center mb-2">
          <h4 class="text-h6 mb-0">Переменные окружения проекта</h4>
          <v-btn color="primary" size="small" variant="tonal" @click="addEnvVar">
            <v-icon icon="mdi-plus" size="small" />
            Добавить переменную
          </v-btn>
        </div>

        <v-alert
          v-if="form.env_vars.length === 0"
          class="mb-2"
          density="compact"
          text="Нет переменных окружения. Добавьте переменные для использования в скриптах."
          type="info"
          variant="tonal"
        />

        <v-sheet
          v-for="(envVar, index) in form.env_vars"
          :key="envVar.id || `new-${index}`"
          class="d-flex flex-row ga-2 align-center mb-2"
          variant="flat"
        >
          <v-text-field
            v-model="envVar.name"
            class="flex-grow-1"
            density="compact"
            hide-details
            label="Имя"
            placeholder="VAR_NAME"
            style="max-width: 200px"
          />
          <v-text-field
            v-model="envVar.value"
            class="flex-grow-1"
            density="compact"
            hide-details
            label="Значение"
            placeholder="value"
            type="text"
          />
          <v-checkbox
            v-model="envVar.uses_other_variables"
            class="flex-grow-0"
            density="compact"
            hide-details
            label="Использует другие переменные"
            style="max-width: 220px"
          />
          <v-btn
            class="align-self-center"
            color="error"
            icon="mdi-delete"
            size="small"
            variant="text"
            @click="removeEnvVar(index)"
          />
        </v-sheet>
      </v-sheet>

      <v-sheet
        v-for="(stage, index) in form.stages"
        :key="index"
        class="d-flex pa-4 flex-column mt-2"
        elevation="1"
      >
        <div class="d-flex flex-row justify-space-between align-center">
          <span class="text-h6">Этап {{ index + 1 }}</span>
          <v-btn
            color="error"
            icon="mdi-delete"
            size="small"
            variant="text"
            @click="removeStage(index)"
          />
        </div>
        <v-textarea v-model="stage.script" class="mt-2" label="Скрипт" monospace rows="4" />
      </v-sheet>

      <div class="d-flex flex-row ga-2 mt-4">
        <v-btn color="success" @click="saveProject"> Сохранить </v-btn>
        <v-btn variant="outlined" @click="cancel"> Отмена </v-btn>
      </div>
    </v-sheet>
  </v-container>
</template>

<script setup lang="ts">
import type { ProjectCreate, ProjectUpdate } from '@/api/generated.schemas';
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getEasyJetAPI } from '@/api/generated';

const router = useRouter();
const route = useRoute();
const api = getEasyJetAPI();

interface Stage {
  number: number;
  script: string;
}

interface EnvVar {
  id?: number;
  name: string;
  value: string;
  uses_other_variables: boolean;
}

interface ProjectForm {
  id: number;
  name: string;
  dir?: string;
  git_url?: string;
  git_branch?: string;
  cron_enabled: boolean;
  cron_schedule: string;
  restart_after: boolean;
  with_root_env: boolean;
  is_template: boolean;
  retention_count: number;
  stages: Stage[];
  env_vars: EnvVar[];
}

const form = ref<ProjectForm | null>(null);
const templates = ref<Array<{ id: number; name: string }>>([]);
const selectedTemplateId = ref<number | null>(null);

const cronRules = [
  (v: string) => {
    if (!v || v.trim() === '') return true; // Empty is valid (no schedule)
    // Basic cron validation (5 fields)
    const parts = v.trim().split(/\s+/);
    if (parts.length !== 5) return 'Cron выражение должно содержать 5 полей';
    return true;
  },
];

const isEdit = computed(() => {
  return form.value?.id !== undefined && form.value.id !== 0;
});

function addStage() {
  if (!form.value) return;
  const nextNumber =
    form.value.stages.length > 0 ? Math.max(...form.value.stages.map((s) => s.number)) + 1 : 1;
  form.value.stages.push({ number: nextNumber, script: '' });
}

function removeStage(index: number) {
  if (!form.value) return;
  form.value.stages.splice(index, 1);
}

function saveProject() {
  if (!form.value) return;
  if (isEdit.value) {
    const payload: ProjectUpdate = {
      id: form.value.id,
      name: form.value.name,
      dir: form.value.dir || undefined,
      git_url: form.value.git_url || undefined,
      git_branch: form.value.git_branch || undefined,
      cron_enabled: form.value.cron_enabled,
      cron_schedule: form.value.cron_schedule || undefined,
      restart_after: form.value.restart_after,
      retention_count: form.value.retention_count,
      with_root_env: form.value.with_root_env,
      is_template: form.value.is_template,
      stages: form.value.stages,
      env_vars: form.value.env_vars.map((ev) => ({
        name: ev.name,
        value: ev.value,
        uses_other_variables: ev.uses_other_variables,
      })),
    };
    api
      .updateProject(form.value.id, payload)
      .then(() => {
        router.push(`/projects/${form.value!.id}`);
      })
      .catch((error) => {
        console.error(error);
      });
  } else {
    const payload: ProjectCreate = {
      name: form.value.name,
      dir: form.value.dir || undefined,
      git_url: form.value.git_url || undefined,
      git_branch: form.value.git_branch || undefined,
      cron_enabled: form.value.cron_enabled,
      cron_schedule: form.value.cron_schedule || undefined,
      restart_after: form.value.restart_after,
      retention_count: form.value.retention_count,
      with_root_env: form.value.with_root_env,
      is_template: form.value.is_template,
      stages: form.value.stages,
      env_vars: form.value.env_vars.map((ev) => ({
        name: ev.name,
        value: ev.value,
        uses_other_variables: ev.uses_other_variables,
      })),
    };
    api
      .createProject(payload)
      .then((response) => {
        router.push(`/projects/${response.data.id}`);
      })
      .catch((error) => {
        console.error(error);
      });
  }
}

function cancel() {
  if (isEdit.value && form.value) {
    router.push(`/projects/${form.value.id}`);
  } else {
    router.push('/');
  }
}

function addEnvVar() {
  if (!form.value) return;
  form.value.env_vars.push({ name: '', value: '', uses_other_variables: false });
}

function removeEnvVar(index: number) {
  if (!form.value) return;
  form.value.env_vars.splice(index, 1);
}

function applyTemplate() {
  if (!form.value || !selectedTemplateId.value) return;

  api
    .getProject(selectedTemplateId.value)
    .then((v) => {
      const template = v.data;
      form.value!.name = template.name ?? '';
      form.value!.dir = template.dir ?? '';
      form.value!.git_url = template.git_url ?? '';
      form.value!.git_branch = template.git_branch ?? '';
      form.value!.cron_enabled = template.cron_enabled ?? false;
      form.value!.cron_schedule = template.cron_schedule ?? '';
      form.value!.restart_after = template.restart_after ?? false;
      form.value!.with_root_env = template.with_root_env ?? false;
      form.value!.retention_count = template.retention_count ?? 0;
      form.value!.stages = template.stages?.map((s: Stage) => ({ ...s })) || [];
      form.value!.env_vars =
        template.env_vars?.map((ev) => ({
          id: ev.id,
          name: ev.name,
          value: ev.value,
          uses_other_variables: ev.uses_other_variables ?? false,
        })) || [];
    })
    .catch((error) => {
      console.error(error);
    });
}

function load() {
  const id = route.params.id;
  if (!id || id === 'new') {
    form.value = {
      id: 0,
      name: '',
      dir: '',
      git_url: '',
      git_branch: '',
      cron_enabled: false,
      cron_schedule: '',
      restart_after: false,
      with_root_env: false,
      is_template: false,
      retention_count: 0,
      stages: [],
      env_vars: [],
    };
  } else {
    api
      .getProject(Number(id))
      .then((v) => {
        const data = v.data;
        form.value = {
          id: data.id ?? 0,
          name: data.name ?? '',
          dir: data.dir || '',
          git_url: data.git_url || '',
          git_branch: data.git_branch || '',
          cron_enabled: data.cron_enabled || false,
          cron_schedule: data.cron_schedule || '',
          restart_after: data.restart_after || false,
          with_root_env: data.with_root_env || false,
          is_template: data.is_template || false,
          retention_count: data.retention_count || 0,
          stages: data.stages?.map((s: Stage) => ({ ...s })) || [],
          env_vars:
            data.env_vars?.map((ev) => ({
              id: ev.id,
              name: ev.name,
              value: ev.value,
              uses_other_variables: ev.uses_other_variables ?? false,
            })) || [],
        };
      })
      .catch((error) => {
        console.error(error);
      });
  }

  api
    .getProjects({ type: 'template' })
    .then((v) => {
      const projectsData = v.data.projects;
      if (projectsData) {
        templates.value = projectsData.map((p) => ({ id: p.id, name: p.name }));
      }
    })
    .catch((error) => {
      console.error(error);
    });
}

onMounted(() => {
  load();
});
</script>
