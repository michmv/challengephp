# API

## Создание одного события

POST: http://localhost:12001/events

```bash
curl --location 'http://localhost:12001/events' \
--header 'Content-Type: application/json' \
--data '{
  "user_id": 1234,
  "event_type": "click",
  "timestamp": "2025-08-28T12:34:56Z",
  "metadata": {
    "page": "/home"
  }
}'
```

## Получить события с пагинацией и сортировкой по времени

GET: http://localhost:12001/events?page=1&limit=100

```bash
curl --location 'http://localhost:12001/events?page=1&limit=100'
```

## Последние 1000 событий конкретного пользователя

GET: http://localhost:12001/users/1/events

```bash
curl --location 'http://localhost:12001/users/1/events'
```

## Агрегированная статистика

GET: http://localhost:12001/stats?from=2025-04-01+00:00:00&to=2025-09-01+00:00:00&type=click

```bash
curl --location 'http://localhost:12001/stats?from=2025-04-01+00%3A00%3A00&to=2025-09-01+00%3A00%3A00&type=click'
```
