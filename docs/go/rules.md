# Правила при написании кода на Go

## ⚠️ КРИТИЧЕСКИ ВАЖНО

**ВСЕГДА** используй команды через `task` — **НИКОГДА** не запускай их напрямую или альтернативными способами.

Это обязательное требование, а не рекомендация.

## Обязательные команды

После **любой** модификации Go кода **ОБЯЗАТЕЛЬНО** выполни по порядку:

1. `task go:tidy` — нормализация зависимостей
2. `task go:format` — форматирование кода (gofumpt + goimports)
3. `task go:lint` — проверка линтером
   - При ошибках: `task go:lint:fix`
4. `task go:test` — запуск тестов
5. `task build:server` — сборка сервера

## Запрещено

- ❌ **НЕЛЬЗЯ** запускать `go mod tidy` напрямую — только `task go:tidy`
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
