#!/bin/bash

# Скрипт для проверки состояния QuantumDiary сервера

echo "=== Проверка состояния QuantumDiary сервера ==="

# Проверяем, запущен ли сервер
if pgrep -f "server_final" > /dev/null; then
    echo "✓ Сервер запущен"
    
    # Проверяем работоспособность
    if curl -s http://localhost:8080/health/ping > /dev/null; then
        echo "✓ Сервер отвечает на health check"
    else
        echo "✗ Сервер не отвечает на health check"
    fi
else
    echo "✗ Сервер не запущен"
    echo "Для запуска выполните: ./scripts/start_server.sh"
fi

echo ""
echo "=== Информация о системе ==="
echo "Версия Go:"
go version
echo ""
echo "Количество файлов Go в проекте:"
find . -name "*.go" | wc -l
echo ""
echo "Количество миграций:"
ls -la migrations/ | grep ".sql" | wc -l