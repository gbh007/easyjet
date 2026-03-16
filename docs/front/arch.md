# Архитектура фронтенда на Vue.js

Приложение использует **Vue 3** с **Composition API** и **Pinia** для управления состоянием.

## Структура проекта

```plain
frontend/
├── src/
│   ├── assets/                 # Статические ресурсы (шрифты, изображения)
│   ├── components/             # Переиспользуемые Vue компоненты
│   ├── pages/                  # Страничные компоненты (route views)
│   ├── plugins/                # Плагины (Vuetify и др.)
│   ├── router/                 # Конфигурация Vue Router
│   ├── stores/                 # Pinia stores (глобальное состояние)
│   ├── styles/                 # Глобальные стили (SCSS)
│   ├── App.vue                 # Корневой компонент
│   └── main.ts                 # Точка входа, инициализация приложения
├── public/                     # Публичные файлы (копируются в dist)
├── index.html                  # HTML шаблон
├── package.json                # Зависимости и скрипты
├── tsconfig.json               # Конфигурация TypeScript
├── vite.config.mts             # Конфигурация Vite
└── .prettierrc.json            # Конфигурация Prettier
```

## Архитектурные принципы

### Composition API

Все компоненты используют `<script setup>` синтаксис с Composition API.

### Pinia для состояния

Глобальное состояние управляется через Pinia stores (`src/stores/`).

### Vue Router для навигации

Маршрутизация настроена через Vue Router (`src/router/`).

### Vuetify для UI

Компоненты Material Design из библиотеки Vuetify 4.
