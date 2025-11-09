FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN swag init -g ./cmd/server/main.go -o ./docs --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add netcat-openbsd
RUN apk --no-cache add mysql-client 

WORKDIR /root/

# Копируем из builder стадии
COPY --from=builder /app/main .
COPY --from=builder /app/application.yaml .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
COPY --from=builder /go/bin/goose /usr/local/bin/

# ЯВНО копируем entrypoint с хоста в текущую директорию (/root/)
COPY migrations-entrypoint.sh ./
RUN chmod +x ./migrations-entrypoint.sh

# ДЛЯ ОТЛАДКИ - проверим что файлы на месте
RUN ls -la ./
RUN pwd

EXPOSE 8080

ENTRYPOINT ["./migrations-entrypoint.sh"]
CMD ["./main"]