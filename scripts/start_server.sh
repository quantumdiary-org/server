#!/bin/bash

# Скрипт для запуска QuantumDiary сервера

echo "=== Запуск QuantumDiary сервера ==="

# Проверяем, что файл конфигурации существует
if [ ! -f ".env" ]; then
    echo "Создание .env файла из примера..."
    cp .env.example .env
fi

# Проверяем, что сервер собран
if [ ! -f "bin/server" ]; then
    echo "Сборка сервера..."
    go build -o bin/server ./api/cmd/server/main.go
    if [ $? -ne 0 ]; then
        echo "Ошибка сборки сервера"
        exit 1
    fi
    echo "Сервер успешно собран"
fi

echo "Запуск сервера..."
bin/server
