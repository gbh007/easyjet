# Установка EasyJet через Docker

⚠️ **ВАЖНО:** инструкция в процессе обкатки

⚠️ **Важные особенности при запуске в docker:**

- Тома монтируемые по пути используют внешние пути а не в контейнере
- Файл docker compose используется по пути внутри контейнера
- Работа с docker compose проверенна не полностью и возможны прочие особенности не описанные здесь

## Общие рекомендации

Для того чтобы не включать передачу параметров окружения системы, рекомендуется настроить глобально вот эти переменные:

- `PATH`, пример `/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin`
- `HOME`, пример `/var/lib/easyjet`

## Перенос с linux daemon

Процесс описан для случая когда настройка была произведена по [этой инструкции](./how-to-install.md)  
⚠️ **ВАЖНО:** процесс описан под настройку под root, если требуется более защищенная настройка, то расширьте самостоятельно

1. Создайте директорию под проект и настройте ее под соответствующего пользователя

   ```sh
   sudo mkdir /srv/easyjet
   sudo chown -R root:root /srv/easyjet
   ```

2. Скопируйте конфиг

   ```sh
   cp /home/easyjet/config.toml /srv/easyjet/config.toml
   ```

3. Скопируйте данные и доступы

   ```sh
   sudo mkdir -p /var/lib/easyjet

   sudo cp -R /home/easyjet/projects /var/lib/easyjet/projects

   # Данные авторизаций если есть
   sudo cp -R /home/easyjet/.ssh /home/easyjet/.netrc /var/lib/easyjet

   # Исправляем владельца
   sudo chown -R root:root /var/lib/easyjet
   ```

4. Перенастраиваем новый конфиг под изменения (указаны только те параметры на которые надо обратить внимание)

   ```toml
   [app]
   # Меняем путь
   project_dir = "/var/lib/easyjet/projects"

   [database]
   # Меняем хост (если не используется сеть хоста)
   dns = "postgres://easyjetuser:easyjetpass@postgres:5432/easyjet?sslmode=disable"

   [server]
   # Меняем т.к. в образе уже есть настроенная отдача
   static_files_path = "/app/static"
   ```

5. Собираем образ (выполнять в склонированом проекте на хосте)

   ```sh
   docker build -f Dockerfile -t easyjet:latest .
   ```

6. Создаем `/srv/easyjet/docker-compose.yml`

   ```yaml
   services:
     easyjet:
       image: easyjet:latest
       restart: always
       volumes:
         - /var/lib/easyjet:/var/lib/easyjet
         - ./config.toml:/app/config.toml
         - /var/run/docker.sock:/var/run/docker.sock
         # Если нужно работать с корневой системой
         # - /:/host:ro
         # Если используешь SSH
         # - /var/lib/easyjet/.ssh:/root/.ssh
       environment:
         HOME: /var/lib/easyjet
         # Настрой нужную таймзону для крона
         TZ: Asia/Krasnoyarsk
       ports:
         - 8080:8080
       logging:
         driver: json-file
            options:
               max-size: 10m
               max-file: 2
   ```

7. Отключаем linux-daemon

   ```sh
   sudo systemctl disable --now easyjet.service
   ```

## Настройка само-обновления

Для самообновления можно использовать подобный скрипт (при настройке по инструкции)

```sh
docker build -f Dockerfile -t easyjet:latest .

docker run --rm -d --name easyjet-updater \
-v /srv/easyjet:/srv/easyjet \
-v /var/run/docker.sock:/var/run/docker.sock \
docker:29 \
/bin/sh -c "
sleep 10
docker compose -f /srv/easyjet/docker-compose.yml up -d --remove-orphans
"
```
