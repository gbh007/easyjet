<template>
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
            <tr v-for="envVar in envVariables" :key="envVar.name">
              <td>
                <code class="text-primary">{{ envVar.name }}</code>
              </td>
              <td>{{ envVar.description }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script setup lang="ts">
interface EnvVariable {
  name: string;
  description: string;
}

const envVariables: EnvVariable[] = [
  {
    name: 'BUILD_NUMBER',
    description: 'Текущий номер сборки, например "153"',
  },
  {
    name: 'BUILD_ID',
    description: 'Текущий ID сборки в формате "BUILD_NUMBER-STAGE_NUMBER"',
  },
  {
    name: 'BUILD_URL',
    description: 'Внешний web URL, где можно найти результаты этой сборки',
  },
  {
    name: 'NODE_NAME',
    description: 'Имя узла, на котором запущена текущая сборка. Равно "master" для главного узла',
  },
  {
    name: 'JOB_NAME',
    description: 'Имя проекта (задания), которое было указано при создании',
  },
  {
    name: 'BUILD_TAG',
    description:
      'Строка формата "JOB_NAME-BUILD_NUMBER", название приводится к нижнему регистру без пробелов',
  },
  {
    name: 'JENKINS_URL',
    description: 'URL EasyJet-сервера, который запускает сборку (для совместимости с Jenkins)',
  },
  {
    name: 'EASYJET_URL',
    description: 'URL EasyJet-сервера, который запускает сборку',
  },
  {
    name: 'EXECUTOR_NUMBER',
    description:
      'Уникальный номер исполнителя (среди исполнителей той же машины), выполняющего эту сборку. Начинается с 0',
  },
  {
    name: 'WORKSPACE',
    description: 'Абсолютный путь к рабочей директории проекта',
  },
  {
    name: 'GIT_COMMIT',
    description: 'Для проектов на Git: хэш коммита, выбранного для сборки (требуется плагин git)',
  },
  {
    name: 'GIT_URL',
    description: 'Для проектов на Git: URL Git-репозитория (требуется плагин git)',
  },
  {
    name: 'GIT_BRANCH',
    description: 'Для проектов на Git: имя ветки Git, выбранной для сборки (требуется плагин git)',
  },
];
</script>
