#!/bin/bash

# Скрипт для очистки временных файлов и бинарников

echo "=== Очистка ==="

# Удаляем бинарные файлы
find . -type f -name "server_*" -executable -delete
find . -type f -name "*.db" -not -name ".gitignore" -delete

# Удаляем временные файлы
find . -type f -name "*.tmp" -delete
find . -type f -name "*.log" -not -path "./docs/*" -delete
find . -type f -name "*.bak" -delete
find . -type f -name "*.old" -delete
find . -type f -name "*~" -delete

# Удаляем кэш Go
go clean -cache

echo "Очистка завершена"
