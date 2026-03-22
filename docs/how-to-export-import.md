# Утилита configsync

**configsync** — утилита командной строки для экспорта и импорта конфигураций проектов EasyJet.

## Назначение

- **Экспорт** — выгрузка всех проектов из EasyJet в JSON-файл для резервного копирования или переноса
- **Импорт** — загрузка проектов из JSON-файла в EasyJet

## Сборка

```sh
go build -o configsync cmd/configsync/main.go
```

После сборки бинарный файл `configsync` будет доступен в текущей директории.

## Команды

### export — Экспорт проектов

Экспортирует все проекты из EasyJet в JSON-файл.

#### Синтаксис

```sh
./configsync export [опции]
```

**Полный пример с опциями:**

```sh
./configsync export \
  -url http://easyjet.example.com:8080 \
  -username admin \
  -password secret123 \
  -file /backup/easyjet-backup.json \
  -timeout 2m
```

#### Опции

| Опция       | По умолчанию            | Описание                                 |
| ----------- | ----------------------- | ---------------------------------------- |
| `-url`      | `http://localhost:8080` | URL сервера EasyJet                      |
| `-username` | (пусто)                 | Имя пользователя для аутентификации      |
| `-password` | (пусто)                 | Пароль для аутентификации                |
| `-file`     | `dump.json`             | Путь к файлу для экспорта                |
| `-timeout`  | `1m`                    | Таймаут запросов (например, `30s`, `2m`) |

### import — Импорт проектов

Импортирует проекты из JSON-файла в EasyJet.

#### Синтаксис

```sh
./configsync import [опции]
```

**Полный пример с опциями:**

```sh
./configsync import \
  -url http://easyjet.example.com:8080 \
  -username admin \
  -password secret123 \
  -file /backup/easyjet-backup.json \
  -timeout 2m
```

#### Опции

| Опция       | По умолчанию            | Описание                                 |
| ----------- | ----------------------- | ---------------------------------------- |
| `-url`      | `http://localhost:8080` | URL сервера EasyJet                      |
| `-username` | (пусто)                 | Имя пользователя для аутентификации      |
| `-password` | (пусто)                 | Пароль для аутентификации                |
| `-file`     | `dump.json`             | Путь к файлу для импорта                 |
| `-timeout`  | `1m`                    | Таймаут запросов (например, `30s`, `2m`) |
