# Список используемых библиотек для языка Go

## Основные зависимости

- **Конфигурация** - [`github.com/BurntSushi/toml`](https://github.com/BurntSushi/toml) - парсинг TOML-конфигурации
- **Логирование** - `log/slog` (стандартная библиотека Go 1.21+)
  - [`github.com/lmittmann/tint`](https://github.com/lmittmann/tint) - форматирование логов для локальной разработки
  - [`github.com/golang-cz/devslog`](https://github.com/golang-cz/devslog) - цветное форматирование для отладки
- **База данных** - [`gorm.io/gorm`](https://gorm.io/) - ORM для работы с БД
  - [`gorm.io/driver/postgres`](https://gorm.io/driver/postgres) - драйвер для PostgreSQL
  - [`github.com/glebarez/sqlite`](https://github.com/glebarez/sqlite) - драйвер для SQLite (pure Go)
  - [`github.com/orandin/slog-gorm`](https://github.com/orandin/slog-gorm) - логирование запросов GORM через slog
- **Миграции БД** - [`github.com/pressly/goose/v3`](https://github.com/pressly/goose) - управление миграциями базы данных
- **Веб-сервер** - [`github.com/labstack/echo/v4`](https://github.com/labstack/echo) - HTTP-фреймворк
  - [`github.com/go-playground/validator/v10`](https://github.com/go-playground/validator) - валидация данных
  - [`github.com/samber/slog-echo`](https://github.com/samber/slog-echo) - логирование запросов Echo через slog
- **Планировщик задач** - [`github.com/go-co-op/gocron/v2`](https://github.com/go-co-op/gocron) - выполнение задач по расписанию (cron)
- **Утилиты** - [`github.com/samber/lo`](https://github.com/samber/lo) - функциональные утилиты (Lodash-style для Go)
- **Парсинг SQL** - [`github.com/Masterminds/squirrel`](https://github.com/Masterminds/squirrel) - конструктор SQL-запросов

## Внутренние зависимости

- [`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx) - PostgreSQL драйвер для Go (используется напрямую в repository/postgres)
- [`golang.org/x/sync`](https://pkg.go.dev/golang.org/x/sync) - примитивы синхронизации (context group и др.)
