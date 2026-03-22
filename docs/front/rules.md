# Правила при написании кода на TypeScript/Vue

## ⚠️ КРИТИЧЕСКИ ВАЖНО

**ВСЕГДА** используй команды через `task` — **НИКОГДА** не запускай их напрямую или альтернативными способами.

Это обязательное требование, а не рекомендация.

## Обязательные команды

После **любой** модификации TypeScript/Vue кода **ОБЯЗАТЕЛЬНО** выполни по порядку:

1. `task ts:format` — форматирование кода (Prettier)
2. `task ts:lint` — проверка линтером
   - При ошибках: `task ts:lint:fix`
3. `task build:front` — сборка фронтенда

## Запрещено

- ❌ **НЕЛЬЗЯ** запускать `prettier` напрямую — только `task ts:format`
- ❌ **НЕЛЬЗЯ** запускать `eslint` напрямую — только `task ts:lint`
- ❌ **НЕЛЬЗЯ** запускать `npm run build` напрямую — только `task build:front`

## Почему это важно

Команды через `task`:

- Гарантируют единообразие выполнения
- Используют правильные флаги и настройки
- Легче поддерживать и изменять
- Позволяют избежать человеческих ошибок

**Нарушение этих правил приводит к техническому долгу и несогласованности кода.**

## API спецификация и Orval

### Генерация кода

**ВСЕГДА** используй `task ts:generate` после изменения `openapi.yaml`:

```bash
task ts:generate
```

Это запустит генерацию TypeScript-клиента из OpenAPI спецификации с помощью Orval.

### Структура генерации

- Генерация производится в директорию `frontend/src/api/`
- Исходный файл для генерации: `openapi.yaml`
- Генерируемые файлы:
  - `frontend/src/api/generated.ts` — API функции (axios client)
  - `frontend/src/api/generated.schemas.ts` — TypeScript типы и интерфейсы

### Правила работы с API клиентом

1. **Не изменяй сгенерированные файлы** — все изменения будут потеряны при следующей генерации
2. **Используй сгенерированный клиент** — все HTTP-запросы к API должны выполняться через функции из `generated.ts`
3. **Импортируй из `@/api`** — используй баррел-экспорт для импорта API функций и типов
4. **Обрабатывай undefined поля** — все поля в сгенерированных типах могут быть `undefined`, используй `??` для значений по умолчанию

### Пример использования API клиента

```typescript
// ✅ ПРАВИЛЬНО: использование сгенерированного клиента
import { getEasyJetAPI, type ProjectCreate } from "@/api";

const api = getEasyJetAPI();

// Получить список проектов
const { data } = await api.getProjects();

// Создать проект с правильной обработкой undefined полей
const payload: ProjectCreate = {
  name: form.value.name,
  dir: form.value.dir || undefined,
  stages: form.value.stages,
};

const { data } = await api.createProject(payload);

// Обработка ответа с undefined полями
const projectData = response.data;
const id = projectData.id ?? 0;
const name = projectData.name ?? "";
```

```typescript
// ❌ НЕПРАВИЛЬНО: прямые вызовы axios
import axios from "axios";

// Прямой HTTP-вызов вместо сгенерированного клиента
const response = await axios.get("/api/v1/projects");
```

### Структура файлов API модуля

```plain
frontend/src/api/
├── axios.ts              # Настройка axios instance (baseURL, headers)
├── generated.ts          # Сгенерированные API функции (Orval) ⚠️ не изменять
├── generated.schemas.ts  # Сгенерированные TypeScript типы (Orval) ⚠️ не изменять
└── index.ts              # Баррел-экспорт для удобного импорта
```

### Конвертация сгенерированных типов в локальные

Если тебе нужны более строгие типы (без `undefined`), создай локальные интерфейсы и сконвертируй данные:

```typescript
// Локальный интерфейс с обязательными полями
interface Project {
  id: number;
  name: string;
  created_at: string;
}

// Конвертация из сгенерированного типа
const projectData = response.data;
const project: Project = {
  id: projectData.id ?? 0,
  name: projectData.name ?? "",
  created_at: projectData.created_at ?? "",
};
```

**Преимущества использования Orval:**

- ✅ Типобезопасность — все типы генерируются автоматически из OpenAPI spec
- ✅ Единый стиль — все API вызовы выполняются одинаково
- ✅ Легкость рефакторинга — при изменении спецификации типы обновляются автоматически
- ✅ Меньше бойлерплейта — не нужно вручную описывать типы и HTTP-клиенты
