# Установка EasyJet как linux daemon

1. Создаем нового пользователя

   ```sh
   sudo useradd -m easyjet -s /bin/sh
   ```

2. Добавляем пользователя в нужные группы (пример для docker)

   ```sh
   sudo usermod -aG docker easyjet
   ```

3. Логинимся под новым пользователем

   ```sh
   sudo su - easyjet
   ```

4. Генерируем SSH ключ для пользователя (если необходимо)

   ```sh
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

   ```sh
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

10. При необходимости добавьте нужные переменные окружения для пользователя в `.profile`, на примере Go

    ```sh
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    ```

11. Возвращаемся под основного пользователя
12. Создаем настройку для linux daemon в `/etc/systemd/system/easyjet.service`

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

13. Обновляем конфигурацию демонов

    ```sh
    sudo systemctl daemon-reload
    ```

14. Запускаем демона

    ```sh
    sudo systemctl enable --now easyjet.service
    ```

15. Проверяем демона и его логи

    ```sh
    sudo systemctl status easyjet.service
    sudo journalctl -u easyjet.service
    ```
