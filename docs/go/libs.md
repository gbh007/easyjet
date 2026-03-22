# Список используемых библиотек для языка Go

## Основные зависимости

- **Конфигурация** - [`github.com/BurntSushi/toml`](https://github.com/BurntSushi/toml) - парсинг TOML-конфигурации
- **Логирование** - `log/slog` (стандартная библиотека Go 1.21+)
  - [`github.com/lmittmann/tint`](https://github.com/lmittmann/tint) - форматирование логов для локальной разработки
  - [`github.com/golang-cz/devslog`](https://github.com/golang-cz/devslog) - цветное форматирование для отладки
- **База данных** - [`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx) - PostgreSQL driver с connection pool
  - [`github.com/glebarez/go-sqlite`](https://github.com/glebarez/go-sqlite) - драйвер для SQLite (pure Go)
  - [`github.com/Masterminds/squirrel`](https://github.com/Masterminds/squirrel) - конструктор SQL-запросов
- **Миграции БД** - [`github.com/pressly/goose/v3`](https://github.com/pressly/goose) - управление миграциями базы данных
- **Веб-сервер** - [`github.com/ogen-go/ogen`](https://github.com/ogen-go/ogen) - генерация HTTP-сервера из OpenAPI спецификации
- **Планировщик задач** - [`github.com/go-co-op/gocron/v2`](https://github.com/go-co-op/gocron) - выполнение задач по расписанию (cron)
- **Утилиты** - [`github.com/samber/lo`](https://github.com/samber/lo) - функциональные утилиты (Lodash-style для Go)
- **Метрики и трассировка** - OpenTelemetry
  - [`go.opentelemetry.io/otel`](https://github.com/open-telemetry/opentelemetry-go) - трассировка и метрики
  - [`go.opentelemetry.io/otel/metric`](https://github.com/open-telemetry/opentelemetry-go) - API метрик
  - [`go.opentelemetry.io/otel/trace`](https://github.com/open-telemetry/opentelemetry-go) - API трассировки

## Внутренние зависимости

- [`github.com/go-faster/errors`](https://github.com/go-faster/errors) - утилиты для обработки ошибок (используется ogen)
- [`github.com/go-faster/jx`](https://github.com/go-faster/jx) - быстрое JSON кодирование/декодирование (используется ogen)
- [`golang.org/x/sync`](https://pkg.go.dev/golang.org/x/sync) - примитивы синхронизации (errgroup, context group и др.)
