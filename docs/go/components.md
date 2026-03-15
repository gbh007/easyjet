# Компоненты

Возможно эта документация будет удалена в дальнейшем.

## Ядро (Core)

### Сущности (Entity)

| Сущность            | Описание                                                         |
| ------------------- | ---------------------------------------------------------------- |
| `Project`           | Проект: имя, директория, git URL/branch, этапы выполнения        |
| `ProjectStage`      | Этап выполнения проекта: номер, скрипт                           |
| `ProjectRun`        | Запуск проекта: статусы (success/pending/processing), лог ошибок |
| `ProjectRunStage`   | Результат выполнения этапа: номер, успех, лог                    |
| `ProjectRunGitLogs` | Git коммиты для запуска: хэш, subject                            |
| `Commit`            | Git коммит: хэш, subject                                         |

### Порты (Port)

| Порт         | Тип       | Описание                                                   |
| ------------ | --------- | ---------------------------------------------------------- |
| `Exec`       | Исходящий | Выполнение команд в shell                                  |
| `FileSystem` | Исходящий | Работа с файловой системой (создание директорий, скриптов) |
| `Git`        | Исходящий | Git операции (init, pull, diff)                            |
| `Database`   | Исходящий | Репозиторий для персистентности                            |
| `Service`    | Входящий  | Публичный интерфейс сервиса                                |

### Сервисы (Service)

| Метод                     | Описание                                      |
| ------------------------- | --------------------------------------------- |
| `Project(ctx, id)`        | Получить проект                               |
| `Projects(ctx)`           | Список проектов                               |
| `CreateProject(ctx, p)`   | Создать проект                                |
| `UpdateProject(ctx, p)`   | Обновить проект                               |
| `RunProject(ctx, id)`     | Создать запуск проекта (pending)              |
| `HandleRun(ctx, runID)`   | Обработать запуск (pull git, выполнить этапы) |
| `PendingProjectRuns(ctx)` | Получить ожидающие запуски                    |
| `ProjectRun(ctx, runID)`  | Получить запуск проекта                       |
| `ProjectRuns(ctx, id)`    | Список запусков проекта                       |

## Адаптеры (Adapter)

### Входящие (Driving)

| Адаптер              | Описание                                                  |
| -------------------- | --------------------------------------------------------- |
| `httpapi.Controller` | HTTP API на Echo: REST endpoints для проектов и запусков  |
| `worker.Controller`  | Фоновый worker: опрашивает pending запуски каждую секунду |

**HTTP Endpoints:**

- `POST /api/v1/projects` - создать проект
- `GET /api/v1/projects` - список проектов
- `PUT /api/v1/projects/:project_id` - обновить проект
- `GET /api/v1/projects/:project_id` - получить проект
- `POST /api/v1/projects/:project_id/runs` - запустить проект
- `GET /api/v1/projects/:project_id/runs` - список запусков
- `GET /api/v1/projects/:project_id/runs/:run_id` - детали запуска

### Исходящие (Driven)

| Адаптер              | Порт         | Описание                                       |
| -------------------- | ------------ | ---------------------------------------------- |
| `shellexec.Adapter`  | `Exec`       | Выполнение shell команд                        |
| `filesystem.Adapter` | `FileSystem` | Работа с ФС: создание директорий, .sh скриптов |
| `shellgit.Adapter`   | `Git`        | Git через shell: init, pull, diff              |
| `gorm.Repository`    | `Database`   | SQLite/PostgreSQL через GORM                   |
