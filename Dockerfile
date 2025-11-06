# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем исходный код
COPY . .

RUN swag init -g ./cmd/server/main.go -o ./docs --parseDependency --parseInternal

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add netcat-openbsd

WORKDIR /root/

# Копируем бинарник из стадии сборки
COPY --from=builder /app/main .
COPY --from=builder /app/application.yaml .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

COPY --from=builder /go/bin/goose /usr/local/bin/

COPY docker-entrypoint.sh /root/
RUN chmod +x /root/docker-entrypoint.sh && \
    ls -la /root/docker-entrypoint.sh

# Экспортируем порт
EXPOSE 8080

ENTRYPOINT ["/root/docker-entrypoint.sh"]

# Запускаем приложение
CMD ["./main"]