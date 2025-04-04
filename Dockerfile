FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o goshorten ./cmd/server

# Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/goshorten .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./goshorten"] 