#!/bin/bash

# Скрипт для запуска миграций базы данных

echo "=== Запуск миграций QuantumDiary ==="

# Проверяем, установлен ли goose
if ! command -v goose &> /dev/null; then
    echo "Установка goose..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Определяем тип базы данных из .env
DB_TYPE=$(grep DB_TYPE .env | cut -d '=' -f2 | head -c -1)
DB_PATH=$(grep DB_SQLITE_PATH .env | cut -d '=' -f2 | head -c -1)

if [ "$DB_TYPE" = "sqlite" ]; then
    echo "Используется SQLite: $DB_PATH"
    goose sqlite3 "$DB_PATH" up
elif [ "$DB_TYPE" = "postgres" ]; then
    DB_HOST=$(grep DB_HOST .env | cut -d '=' -f2 | head -c -1)
    DB_PORT=$(grep DB_PORT .env | cut -d '=' -f2 | head -c -1)
    DB_NAME=$(grep DB_NAME .env | cut -d '=' -f2 | head -c -1)
    DB_USER=$(grep DB_USER .env | cut -d '=' -f2 | head -c -1)
    DB_PASS=$(grep DB_PASSWORD .env | cut -d '=' -f2 | head -c -1)
    
    echo "Используется PostgreSQL: $DB_HOST:$DB_PORT/$DB_NAME"
    goose postgres "user=$DB_USER password=$DB_PASS dbname=$DB_NAME host=$DB_HOST port=$DB_PORT sslmode=disable" up
else
    echo "Неподдерживаемый тип базы данных: $DB_TYPE"
    exit 1
fi

echo "Миграции успешно применены!"