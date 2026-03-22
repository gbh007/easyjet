# Правила при написании кода на Go

## ⚠️ КРИТИЧЕСКИ ВАЖНО

**ВСЕГДА** используй команды через `task` — **НИКОГДА** не запускай их напрямую или альтернативными способами.

Это обязательное требование, а не рекомендация.

## Обязательные команды

После **любой** модификации Go кода **ОБЯЗАТЕЛЬНО** выполни по порядку:

1. `task go:tidy` — нормализация зависимостей
2. `task go:generate` — генерация кода из OpenAPI spec (если изменялась спецификация)
3. `task go:format` — форматирование кода (gofumpt + goimports)
4. `task go:lint` — проверка линтером
   - При ошибках: `task go:lint:fix`
5. `task go:test` — запуск тестов
6. `task build:server` — сборка сервера

## Запрещено

- ❌ **НЕЛЬЗЯ** запускать `go mod tidy` напрямую — только `task go:tidy`
- ❌ **НЕЛЬЗЯ** запускать `go generate` напрямую — только `task go:generate`
- ❌ **НЕЛЬЗЯ** запускать `gofumpt`, `goimports` напрямую — только `task go:format`
- ❌ **НЕЛЬЗЯ** запускать `golangci-lint` напрямую — только `task go:lint`
- ❌ **НЕЛЬЗЯ** запускать `go test` напрямую — только `task go:test`
- ❌ **НЕЛЬЗЯ** запускать `go build` напрямую — только `task build:server`

## Почему это важно

Команды через `task`:

- Гарантируют единообразие выполнения
- Используют правильные флаги и настройки
- Легче поддерживать и изменять
- Позволяют избежать человеческих ошибок

**Нарушение этих правил приводит к техническому долгу и несогласованности кода.**

## Правила написания SQL-запросов (squirrel)

### Используй `SetMap` для INSERT и UPDATE

**ВСЕГДА** используй `.SetMap(map[string]any{...})` вместо:

- ❌ `.Columns(...).Values(...)`
- ❌ `.Set("column", value)` для нескольких полей

**Пример правильного использования:**

```go
// INSERT
insertQuery, insertArgs, err := repo.psql.
    Insert("runs").
    SetMap(map[string]any{
        "project_id": run.ProjectID,
        "success":    run.Success,
        "pending":    run.Pending,
    }).
    Suffix("RETURNING id").
    ToSql()

// UPDATE
updateQuery, updateArgs, err := repo.psql.
    Update("runs").
    SetMap(map[string]any{
        "updated_at": run.UpdatedAt,
        "success":    run.Success,
        "pending":    run.Pending,
    }).
    Where(squirrel.Eq{"id": run.ID}).
    ToSql()
```

**Преимущества `SetMap`:**

- ✅ Чище и читаемее — видно соответствие колонка=значение
- ✅ Легче добавлять/удалять поля
- ✅ Меньше шансов ошибиться с порядком значений
- ✅ Единый стиль для INSERT и UPDATE

## API спецификация и ogen

### Генерация кода

**ВСЕГДА** используй `task go:generate` после изменения `openapi.yaml`:

```bash
task go:generate
```

Это запустит генерацию кода из OpenAPI спецификации с помощью ogen.

### Структура генерации

- Генерация производится в директорию `internal/adapter/handler/httpapiogen/`
- Исходный файл для генерации: `openapi.yaml`
- Директория генерации: `internal/adapter/handler/httpapiogen/`

### Правила работы с handlers

1. **Не изменяй сгенерированные файлы** — все изменения будут потеряны при следующей генерации
2. **Реализуй интерфейсы в отдельных файлах** — создавай файлы с реализацией хендлеров отдельно от сгенерированного кода
3. **Используй валидацию** — ogen автоматически генерирует валидацию на основе OpenAPI spec
4. **Следуй бизнес-логике** — вызывай сервисный слой через порт `port.Service`

### Структура файлов handlers

Каждый endpoint должен быть реализован в отдельном файле с именованием по названию метода:

- **Файл обработчика** — именуется по названию метода в snake_case (например, `create_project.go` для метода `CreateProject`)
- **Файл конвертера** — если для endpoint есть конвертер, он размещается в том же файле
- **Файл общей модели** — если конвертер переиспользуется между несколькими endpoint'ами, выноси его в отдельный файл с префиксом `model_convert_` и названием функции в snake_case (например, `model_convert_project_to_ogen.go` для функции `convertProjectToOgen`)

**Пример структуры:**

```plain
internal/adapter/handler/httpapi/
├── handler.go                      # Базовая структура Handler и NewError
├── create_project.go               # Обработчик: CreateProject
├── get_project.go                  # Обработчик: GetProject
├── update_project.go               # Обработчик: UpdateProject
├── create_project_run.go           # Обработчик: CreateProjectRun
├── get_project_run.go              # Обработчик: GetProjectRun
├── get_projects.go                 # Обработчик: GetProjects
├── get_project_runs.go             # Обработчик: GetProjectRuns
├── model_convert_project_create.go # Конвертер: convertProjectCreate
├── model_convert_project_update.go # Конвертер: convertProjectUpdate
├── model_convert_project_to_ogen.go    # Конвертер: convertProjectToOgen
├── model_convert_projects_to_ogen.go   # Конвертер: convertProjectsToOgen
├── model_convert_project_run_to_ogen.go  # Конвертер: convertProjectRunToOgen
└── model_convert_project_runs_to_ogen.go # Конвертер: convertProjectRunsToOgen
```

**Правила именования:**

- Файлы обработчиков: `<метод в snake_case>.go` (например, `create_project.go` для `CreateProject`)
- Файлы конвертеров: `model_convert_<функция в snake_case>.go` (например, `model_convert_project_to_ogen.go` для `convertProjectToOgen`)
- Конвертеры, используемые только в одном обработчике, могут быть размещены в том же файле что и обработчик

### Пример реализации хендлера

```go
// internal/adapter/handler/httpapiogen/projects.go
package httpapiogen

import (
    "context"
    "net/http"

    "github.com/gbh007/easyjet/internal/adapter/handler/httpapiogen/ogenapi"
    "github.com/gbh007/easyjet/internal/core/port"
)

type Handler struct {
    service port.Service
}

func NewHandler(service port.Service) *Handler {
    return &Handler{service: service}
}

// GetProjects implements ogenapi.Handler.
func (h *Handler) GetProjects(ctx context.Context, req *ogenapi.GetProjectsReq) (*ogenapi.GetProjectsRes, error) {
    projects, err := h.service.Projects(ctx)
    if err != nil {
        return nil, err
    }

    return &ogenapi.GetProjectsRes{
        Projects: projects,
    }, nil
}
```
