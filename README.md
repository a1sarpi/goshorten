# GoShorten - Сервис сокращения URL

GoShorten - это сервис сокращения URL, написанный на Go. Он поддерживает как in-memory хранилище для быстрого тестирования, так и PostgreSQL для production использования.

## Возможности

- Сокращение длинных URL в короткие
- Поддержка TTL (время жизни) для ссылок
- Два типа хранилища: in-memory и PostgreSQL
- Swagger документация
- CORS настройки для безопасности
- Graceful shutdown

## Требования

- Go 1.21 или выше
- PostgreSQL 13 или выше (если используется PostgreSQL хранилище)
- Docker (опционально, для запуска PostgreSQL)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/a1sarpi/goshorten.git
cd goshorten
```

2. Установите зависимости:
```bash
go mod download
```

3. Установите Swagger CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. Сгенерируйте Swagger документацию:
```bash
swag init -g api/router.go
```

## Запуск

### С In-Memory хранилищем

```bash
go run cmd/server/main.go
```

### С PostgreSQL

1. Запустите PostgreSQL с помощью Docker Compose:
```bash
docker-compose up -d db
```

2. Запустите приложение:
```bash
go run cmd/server/main.go
```

## API

### Создание короткой ссылки

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}' \
  http://localhost:8080/shorten
```

Ответ:
```json
{
  "short_url": "http://localhost:8080/{shortcode}",
  "original_url": "https://example.com",
  "expires_at": 1234567890
}
```

### Переход по короткой ссылке

```bash
curl -L http://localhost:8080/{shortcode}
```

## Swagger документация

Swagger UI доступен по адресу: http://localhost:8080/docs/

## Тестирование

Запуск всех тестов:
```bash
go test ./...
```

## Docker

Для запуска всего приложения с PostgreSQL в Docker:

```bash
docker-compose up -d
```

Приложение будет доступно по адресу: http://localhost:8080
PostgreSQL будет доступен на порту 5432
