<template>
  <v-app-bar elevation="2">
    <v-app-bar-title>
      <router-link class="text-decoration-none text-inherit" to="/">
        EasyJet
      </router-link>
    </v-app-bar-title>

    <v-btn icon @click="toggleTheme">
      <v-icon>{{ isDark ? 'mdi-white-balance-sunny' : 'mdi-weather-night' }}</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<script setup lang="ts">
  import { useThemeStore } from '@/stores/theme'

  const { isDark, toggle } = useThemeStore()

  function toggleTheme () {
    toggle()
  }
</script>
