# geo-alerts-system
Ядро системы геооповещений на Go с интеграцией через вебхуки для новостного портала (Django).

## Возможности
- CRUD инцидентов для оператора (API-key).
- Проверка координат с возвратом ближайших опасных зон.
- Асинхронные вебхуки через Redis-очередь + retry.
- Кэш активных инцидентов в Redis.
- Статистика по зонам за окно времени.
- Health-check эндпоинт.

## Стек
- Go 1.24+
- PostgreSQL 15
- Redis
- Gin

## Быстрый старт (Docker)
1) Запуск сервисов:
```
docker compose up --build
```

2) Миграции:
```
docker compose exec -T db psql -U geoalerts -d geoalerts_db < migrations/001_init.sql
docker compose exec -T db psql -U geoalerts -d geoalerts_db < migrations/002_indexes.sql
```

3) Сервис доступен на `http://localhost:8080`.

## Запуск локально
1) Поднимите Postgres и Redis.
2) Создайте `.env` из `.env.example` и отредактируйте значения.
3) Примените миграции:
```
psql -h localhost -U geoalerts -d geoalerts_db < migrations/001_init.sql
psql -h localhost -U geoalerts -d geoalerts_db < migrations/002_indexes.sql
```
4) Запустите сервис:
```
go run ./cmd/server
```

## Тесты
```
go test ./...
```
Unit-тесты лежат в `tests/unit`.

Интеграционные тесты (PostgreSQL/Redis):
```
export TEST_DB_DSN="postgres://geoalerts:password@localhost:5432/geoalerts_db?sslmode=disable"
export TEST_REDIS_ADDR="localhost:6379"
export TEST_REDIS_DB="1"
go test -tags=integration ./tests/integration
```
Если у вас локальный Postgres слушает `localhost:5432`, тесты могут подключиться к нему вместо контейнера. В этом случае остановите локальный Postgres, смените порт в `docker-compose.yml` или укажите IP машины вместо `localhost` в `TEST_DB_DSN` (например, `192.168.1.10`).

## Webhook mock + ngrok
1) Запустите заглушку:
```
go run ./cmd/webhook-mock
```
2) Поднимите туннель:
```
ngrok http 9090
```
3) Установите `WEBHOOK_URL` на публичный URL ngrok (например, `https://xxxxx.ngrok-free.app/webhook`).

## Переменные окружения
Смотри `.env.example`.

Ключевые:
- `API_KEY` — для защищенных эндпоинтов (header `X-API-Key`).
- `WEBHOOK_URL` — URL вебхука.
- `STATS_TIME_WINDOW_MINUTES` — окно статистики.
- `CACHE_TTL_SECONDS` — TTL кэша активных инцидентов.

Дополнительно:
- `DB_MAX_CONNS`, `DB_MIN_CONNS`, `DB_MAX_CONN_LIFETIME_SECONDS`, `DB_MAX_CONN_IDLE_SECONDS`.
- `WEBHOOK_TIMEOUT_SECONDS`, `HTTP_READ_TIMEOUT_SECONDS`, `HTTP_WRITE_TIMEOUT_SECONDS`, `HTTP_IDLE_TIMEOUT_SECONDS`, `SHUTDOWN_TIMEOUT_SECONDS`.
- `HEALTH_TIMEOUT_SECONDS`.

## API
### Health-check
`GET /api/v1/system/health`
Возвращает статус сервиса и зависимостей (PostgreSQL/Redis). При деградации — HTTP 503.

### Инциденты (требуется `X-API-Key`)
`POST /api/v1/incidents`
```
curl -X POST http://localhost:8080/api/v1/incidents \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev_api_key_12345" \
  -d '{
    "title": "Пожар в районе",
    "description": "Сильное задымление",
    "severity": "high",
    "latitude": 55.751244,
    "longitude": 37.618423,
    "radius_meters": 1200
  }'
```

`GET /api/v1/incidents?page=1&page_size=20`
```
curl -H "X-API-Key: dev_api_key_12345" \
  "http://localhost:8080/api/v1/incidents?page=1&page_size=20"
```

`GET /api/v1/incidents/{id}`
```
curl -H "X-API-Key: dev_api_key_12345" \
  http://localhost:8080/api/v1/incidents/{id}
```

`PUT /api/v1/incidents/{id}`
```
curl -X PUT http://localhost:8080/api/v1/incidents/{id} \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev_api_key_12345" \
  -d '{"radius_meters": 1500}'
```

`DELETE /api/v1/incidents/{id}` (деактивация)
```
curl -X DELETE http://localhost:8080/api/v1/incidents/{id} \
  -H "X-API-Key: dev_api_key_12345"
```

### Статистика по зонам (требуется `X-API-Key`)
`GET /api/v1/incidents/stats`
```
curl -H "X-API-Key: dev_api_key_12345" \
  http://localhost:8080/api/v1/incidents/stats
```

Ответ:
```
{
  "window_minutes": 60,
  "stats": [
    {
      "incident_id": "uuid",
      "title": "Пожар в районе",
      "user_count": 12
    }
  ]
}
```

### Проверка координат (публичный)
`POST /api/v1/location/check`
```
curl -X POST http://localhost:8080/api/v1/location/check \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "latitude": 55.751244,
    "longitude": 37.618423
  }'
```

Ответ:
```
{
  "check_id": "uuid",
  "is_in_danger_zone": true,
  "checked_at": "2025-01-01T12:00:00Z",
  "incidents": [
    {
      "id": "uuid",
      "title": "Пожар в районе",
      "severity": "high",
      "latitude": 55.751244,
      "longitude": 37.618423,
      "radius_meters": 1200,
      "distance_meters": 350.2
    }
  ]
}
```

### Формат вебхука
```
{
  "check_id": "uuid",
  "user_id": "user-123",
  "latitude": 55.751244,
  "longitude": 37.618423,
  "is_in_danger_zone": true,
  "checked_at": "2025-01-01T12:00:00Z",
  "incidents": [
    {
      "id": "uuid",
      "title": "Пожар в районе",
      "severity": "high",
      "latitude": 55.751244,
      "longitude": 37.618423,
      "radius_meters": 1200,
      "distance_meters": 350.2
    }
  ]
}
```
