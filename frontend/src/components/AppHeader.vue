<template>
  <v-app-bar elevation="2">
    <v-app-bar-title>
      <router-link class="text-decoration-none text-inherit" to="/"> EasyJet </router-link>
    </v-app-bar-title>

    <v-menu>
      <template #activator="{ props }">
        <v-btn icon v-bind="props">
          <v-icon>mdi-menu</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item to="/">
          <v-list-item-title>Проекты</v-list-item-title>
        </v-list-item>
        <v-list-item to="/global-env-vars">
          <v-list-item-title>Глобальные переменные</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>

    <v-btn icon @click="toggleTheme">
      <v-icon>{{ isDark ? 'mdi-white-balance-sunny' : 'mdi-weather-night' }}</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme';

const { isDark, toggle } = useThemeStore();

function toggleTheme() {
  toggle();
}
</script>
