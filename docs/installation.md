# Docker

```bash
cd server
echo "DB_PASSWORD=<password>" > .env
docker compose up -d --build # команду надо повторить, не всегда оба контейнера поднимаются успешно 
```

После этой команды будет:
- поднят контейнер с базой данных
- поднят контейнер в котором можно запускать приложение
- в корне проетка будет создана папка `.storage`, в которой будут лежать файл postgres и golang

# Сборка

```bash
cd server
docker compose exec site bash

# Сборка сервера
go build -o challengephp cmd/server/main.go

# Сборка утилиты для базы данных
go build -o database cmd/database/main.go
```

## Запуск приложения

Сначала нужно подправить конфиг, указав в нем правильный пароль для доступа к базе данных, тот пароль что указывался при создании контейнеров.

```bash
cd application
cp config_copy.yml config.yml
vim config.yml
```

```bash
cd server
docker compose exec site bash
./database init # пересоздают таблицы в базе данных
./database seed # заполняет таблицу десятью миллионом случайных событий
./challengephp server # запуск сервера
```

Тест что сервер работает - http://localhost:12001/ping
