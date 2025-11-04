# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из стадии сборки
COPY --from=builder /app/main .
COPY --from=builder /app/application.yaml .

# Экспортируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]