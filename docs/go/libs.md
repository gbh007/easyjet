# Список используемых библиотек для языка Go

- Для работы с конфигурацией - <github.com/BurntSushi/toml>
- Для работы с логированием - `log/slog`
  - Для форматирования логов в production - `slog.JSONHandler`
  - Для формирования логов при локальном запуске - <github.com/lmittmann/tint>
  - Для формирования логов для отладки человеком локально - <github.com/golang-cz/devslog>
- Для работы с базой данных - <gorm.io/gorm>
  - Драйвер для PostgreSQL - <gorm.io/driver/postgres>
  - Драйвер для Sqlite - <github.com/glebarez/sqlite>
  - Библиотека для логирования - <github.com/orandin/slog-gorm>
- Для работы с веб сервером - <github.com/labstack/echo/v4>
  - Библиотека для валидации данных - <github.com/go-playground/validator/v10>
  - Библиотека для логирования - <github.com/samber/slog-echo>
- Для трансформации данных функциональными методами - <github.com/samber/lo>
