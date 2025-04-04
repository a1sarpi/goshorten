# GoShorten - Сервис сокращения URL

GoShorten - это сервис сокращения URL, написанный на Go. Он поддерживает как in-memory хранилище для быстрого тестирования, так и PostgreSQL для production использования.

## Возможности

- Сокращение длинных URL в короткие
- Поддержка TTL (время жизни) для ссылок
- Два типа хранилища: in-memory и PostgreSQL
- Swagger документация
- Rate limiting для защиты от перегрузок
- CORS настройки для безопасности
- Метрики для мониторинга
- Graceful shutdown
- Конфигурация через YAML и переменные окружения

## Требования

- Go 1.21 или выше
- PostgreSQL 12 или выше (если используется PostgreSQL хранилище)
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

3. Скопируйте пример конфигурации:
```bash
cp config/config.example.yaml config/config.yaml
```

4. Настройте конфигурацию в `config/config.yaml`

## Запуск

### С In-Memory хранилищем

```bash
go run cmd/server/main.go
```

### С PostgreSQL

1. Запустите PostgreSQL (с помощью Docker):
```bash
docker run --name goshorten-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=goshorten \
  -p 5432:5432 \
  -d postgres:latest
```

2. Запустите приложение с PostgreSQL:
```bash
STORAGE_TYPE=postgres \
POSTGRES_CONN_STRING="postgresql://postgres:postgres@localhost:5432/goshorten?sslmode=disable" \
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

### Переход по короткой ссылке

```bash
curl -L http://localhost:8080/{shortcode}
```

## Swagger документация

Swagger UI доступен по адресу: http://localhost:8080/swagger/

## Тестирование

Запуск всех тестов:
```bash
go test ./...
```

## Метрики

Метрики доступны по адресу: http://localhost:8080/metrics

## Конфигурация

Основные настройки в `config.yaml`:
- Тип хранилища (memory/postgres)
- Настройки сервера (порт, таймауты)
- Настройки базы данных
- Rate limiting
- CORS
- Логирование
- Метрики

## Лицензия

MIT License - см. файл [LICENSE](LICENSE)

## Инструкции по тестированию

Это руководство поможет вам протестировать весь функционал сервиса GoShorten.

### Требования
- Go 1.21 или выше
- Docker и Docker Compose (для тестирования PostgreSQL)
- curl или Postman для тестирования API

### 1. Установка и настройка

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
cd api
swag init -g router.go
cd ..
```

### 2. Тестирование In-Memory хранилища

1. Запустите сервер с in-memory хранилищем:
```bash
go run cmd/server/main.go
```

2. Протестируйте сокращение URL:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

3. Протестируйте редирект:
```bash
curl -L http://localhost:8080/{shortcode}
```

4. Протестируйте ограничение запросов:
```bash
# Выполните несколько раз подряд
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

5. Проверьте метрики:
```bash
curl http://localhost:8080/metrics
```

### 3. Тестирование PostgreSQL хранилища

1. Запустите PostgreSQL через Docker:
```bash
docker-compose up -d postgres
```

2. Запустите сервер с PostgreSQL хранилищем:
```bash
STORAGE_TYPE=postgres \
POSTGRES_CONN_STRING="postgresql://postgres:postgres@localhost:5432/goshorten?sslmode=disable" \
go run cmd/server/main.go
```

3. Повторите все тесты из раздела 2.

4. Протестируйте истечение срока действия URL:
```bash
# Создайте URL со сроком действия 1 секунда
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "ttl": 1}'

# Подождите 2 секунды и попробуйте получить доступ
curl -L http://localhost:8080/{shortcode}
```

### 4. Тестирование Swagger документации

1. Откройте Swagger UI:
```bash
# Откройте в браузере
http://localhost:8080/docs
```

2. Протестируйте все эндпоинты через Swagger UI.

### 5. Тестирование обработки ошибок

1. Протестируйте неверный URL:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "invalid-url"}'
```

2. Протестируйте пустой URL:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'
```

3. Протестируйте несуществующий shortcode:
```bash
curl -L http://localhost:8080/nonexistent
```

### 6. Тестирование graceful shutdown

1. Запустите сервер
2. Отправьте запрос
3. Нажмите Ctrl+C для остановки сервера
4. Убедитесь, что сервер завершает обработку запроса перед выключением

### 7. Тестирование CORS

1. Создайте простой HTML файл с AJAX запросом
2. Попробуйте получить доступ к API с другого домена
3. Проверьте CORS заголовки в ответе

### 8. Тестирование метрик

1. Проверьте эндпоинт метрик:
```bash
curl http://localhost:8080/metrics
```

2. Убедитесь, что метрики обновляются после:
   - Успешного создания URL
   - Ошибки создания URL
   - Редиректов
   - Срабатывания ограничения запросов

### 9. Тестирование конфигурации

1. Протестируйте различные настройки:
   - Изменение порта
   - Настройки ограничения запросов
   - Настройки TTL по умолчанию
   - Настройки CORS

### 10. Тестирование Docker развертывания

1. Соберите Docker образ:
```bash
docker build -t goshorten .
```

2. Запустите контейнер:
```bash
docker run -p 8080:8080 goshorten
```

3. Повторите все тесты из раздела 2.

### 11. Тестирование CI/CD пайплайна

1. Отправьте изменения в репозиторий
2. Убедитесь, что тесты запускаются автоматически
3. Проверьте статус деплоя

### 12. Тестирование функций безопасности

1. Протестируйте ограничение запросов
2. Проверьте CORS заголовки
3. Проверьте наличие стандартных заголовков безопасности
4. Протестируйте валидацию входных данных
5. Убедитесь, что сообщения об ошибках не раскрывают конфиденциальную информацию 