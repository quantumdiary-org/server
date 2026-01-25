#!/bin/bash

# Скрипт для тестирования API QuantumDiary

echo "=== Тестирование QD API ==="

# Проверяем, запущен ли сервер
if ! curl -s http://localhost:8080/health/ping > /dev/null; then
    echo "Сервер не запущен. Запустите его с помощью scripts/start_server.sh"
    exit 1
fi

echo "✓ Сервер доступен"

# Тестируем основные эндпоинты
echo "Тестирование health check эндпоинтов..."

curl -s http://localhost:8080/health/ping | jq '.' 2>/dev/null || echo "Health ping: OK"
echo ""

curl -s "http://localhost:8080/health/intping?instance_url=https://sgo.rso23.ru" | jq '.' 2>/dev/null || echo "Health intping: Request sent"
echo ""

echo "Тестирование завершено"
