#!/bin/sh
set -e

echo "Running pre-start operations..."

# Проливка миграций
echo "Running database migrations..."
goose -dir migrations mysql "${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}?parseTime=true" up

# Запуск основного приложения
echo "Starting application..."
exec "$@"