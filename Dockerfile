# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем CA сертификаты для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
