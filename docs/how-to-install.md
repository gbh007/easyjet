# Установка EasyJet как linux daemon

1. Создаем нового пользователя

   ```shell
   sudo useradd -m easyjet -s /bin/sh
   ```

2. Добавляем пользователя в нужные группы (пример для docker)

   ```shell
   sudo usermod -aG docker easyjet
   ```

3. Логинимся под новым пользователем

   ```shell
   sudo su - easyjet
   ```

4. Генерируем SSH ключ для пользователя (если необходимо)

   ```shell
   ssh-keygen
   cat .ssh/id_rsa.pub
   ```

5. Прописываем ключ на машинах куда нужен доступ в `.ssh/authorized_keys` (если необходимо)
6. Вносим авторизационные данные в `.netrc` (если необходимо)

   ```plain
   machine <HOST>
   login <LOGIN>
   password <PASSWORD/TOKEN>
   ```

7. Собираем приложение и размещаем его в `/home/easyjet/`  
   Можно воспользоваться `build.sh`  
   Права на файлы должны быть переданы пользователю `easyjet`

   ```shell
   sudo chown -R easyjet:easyjet /home/easyjet/easyjet /home/easyjet/web
   ```

8. При необходимости создаем базу в PostgreSQL (либо используем sqlite)

   ```sql
   CREATE ROLE easyjetuser WITH LOGIN PASSWORD 'easyjetpass';
   CREATE DATABASE easyjet OWNER easyjetuser;
   ```

9. Настраиваем конфигурацию приложения `config.toml`, пример под PostgreSQL

   ```toml
   [app]
   project_dir = "/home/easyjet/projects"

   [log]
   level = "info"
   format = "json"

   [database]
   type = "postgres"
   dns = "postgres://easyjetuser:easyjetpass@localhost:5432/easyjet?sslmode=disable"

   [server]
   addr = ":8080"
   user = "jet"
   pass = "easy"
   static_files_path = "/home/easyjet/web"
   ```

10. Возвращаемся под основного пользователя
11. Создаем настройку для linux daemon в `/etc/systemd/system/easyjet.service`

    ```ini
    [Unit]
    Description=EasyJet CD
    Requires=network.target
    After=network.target

    [Service]
    ExecStart=/home/easyjet/easyjet
    Restart=always
    User=easyjet
    Group=easyjet
    WorkingDirectory=/home/easyjet

    [Install]
    WantedBy=multi-user.target
    ```

12. Обновляем конфигурацию демонов

    ```shell
    sudo systemctl daemon-reload
    ```

13. Запускаем демона

    ```shell
    sudo systemctl enable --now easyjet.service
    ```

14. Проверяем демона и его логи

    ```shell
    systemctl status easyjet.service
    journalctl -u easyjet.service
    ```
