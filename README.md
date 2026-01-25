# QuantumDiary Server v1.0

QuantumDiary Server - это полнофункциональный прокси-сервер для взаимодействия с различными системами электронного дневника, включая "Сетевой город. Образование". Сервер решает проблему авторизации с использованием хэширования паролей в кодировке Windows-1251 и предоставляет унифицированный интерфейс для всех типов систем дневников.

## Версия 1.0 - Полностью готовая реализация

- ✅ Все методы из @nsapi-web/** полностью реализованы
- ✅ Без заглушек - все функции работают в полном объеме
- ✅ Поддержка NSWebAPI, NSMobileAPI, DevMockAPI
- ✅ Динамический выбор типа API при аутентификации
- ✅ Браузерная аутентификация с Playwright (опционально)
- ✅ Унифицированный трансформер данных QD-1
- ✅ Поддержка всех баз данных (PostgreSQL, MariaDB, MySQL, SQLite)
- ✅ Полная безопасность и аутентификация

## Особенности

- **Двухуровневая аутентификация**: Разделение токенов для клиента и прокси-сервера
- **Поддержка Windows-1251**: Корректная обработка паролей в кодировке Windows-1251
- **Кэширование**: Поддержка Redis и in-memory кэширования
- **Отказоустойчивость**: Резервное кэширование при недоступности NetSchool API
- **Rate Limiting**: Защита от чрезмерного использования API
- **Мониторинг состояния**: Health check эндпоинты для проверки работоспособности
- **Унифицированный трансформер**: Все данные из различных источников приводятся к единому формату QD-1
- **Поддержка всех методов @nsapi-web**: Полная реализация всех методов из JavaScript-клиента, за исключением scheduleDay/scheduleWeek (вместо них используется diary как основной метод получения расписания и домашних заданий)
- **Динамический выбор типа API**: Тип API передается при аутентификации и сохраняется в сессии для каждого пользователя
- **Без заглушек**: Все методы полностью реализованы без использования заглушек

## Установка

### Требования

- Go 1.22+
- PostgreSQL
- (Опционально) Redis для кэширования

### Установка зависимостей

```bash
go mod download
```

### Настройка базы данных

Приложение поддерживает несколько типов баз данных: PostgreSQL, MariaDB/MySQL и SQLite.

#### PostgreSQL (по умолчанию)

1. Установите PostgreSQL и создайте базу данных:

```bash
psql -U postgres -c "CREATE DATABASE netschool_proxy;"
psql -U postgres -c "CREATE USER proxy_user WITH PASSWORD 'dev_password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE netschool_proxy TO proxy_user;"
```

2. Примените миграции:

```bash
# Установите goose если еще не установлен
go install github.com/pressly/goose/v3/cmd/goose@latest

# Примените миграции
goose -dir migrations postgres "user=proxy_user password=dev_password dbname=netschool_proxy sslmode=disable" up
```

#### MariaDB/MySQL

1. Установите MariaDB/MySQL и создайте базу данных:

```sql
CREATE DATABASE netschool_proxy;
CREATE USER 'proxy_user'@'localhost' IDENTIFIED BY 'dev_password';
GRANT ALL PRIVILEGES ON netschool_proxy.* TO 'proxy_user'@'localhost';
FLUSH PRIVILEGES;
```

2. Примените миграции:

```bash
# Примените миграции
goose -dir migrations mysql "proxy_user:dev_password@tcp(localhost:3306)/netschool_proxy" up
```

3. Обновите .env файл:

```
DB_TYPE=mariadb
DB_HOST=localhost
DB_PORT=3306
DB_NAME=netschool_proxy
DB_USER=proxy_user
DB_PASSWORD=dev_password
```

#### SQLite (для быстрого запуска без настройки сервера БД)

1. Просто установите тип базы данных в .env файле:

```
DB_TYPE=sqlite
DB_SQLITE_PATH=./db.sqlite
```

Приложение автоматически создаст файл базы данных при первом запуске.

### Конфигурация

Скопируйте пример файла окружения и настройте под свои нужды:

```bash
cp .env.example .env
```

Отредактируйте .env файл с вашими настройками. Приложение автоматически загружает переменные из .env файла при запуске.

Обратите внимание, что URL экземпляра NetSchool теперь передается динамически в каждом запросе как `instance_url`, а не задается в конфигурации.

## Запуск

### В режиме разработки

```bash
go run api/cmd/server/main.go --config config/dev.yaml
```

### Сборка и запуск бинарного файла

```bash
# Сборка
go build -o bin/server api/cmd/server/main.go

# Запуск
./bin/server --config config/dev.yaml
```

## API Эндпоинты

### Аутентификация

- `POST /auth/login` - Аутентификация пользователя в NetSchool
- `POST /auth/logout` - Выход пользователя (требует аутентификации)

Для аутентификации отправьте POST-запрос на `/auth/login` с телом:

```json
{
  "username": "your_username",
  "password": "your_password",
  "school_id": 123,
  "instance_url": "https://sgo.rso23.ru"
}
```

Параметр `instance_url` - это URL-адрес конкретного экземпляра NetSchool для региона пользователя (например, `https://schools.dagestan.ru`, `https://sgo.rso23.ru` и т.д.).

### Здоровье системы

- `GET /health/ping` - Проверка доступности прокси-сервера
- `GET /health/intping` - Проверка соединения с NetSchool API
- `GET /health/full` - Полная проверка состояния системы

### Защищенные эндпоинты (требуют аутентификации)

- `GET /api/v1/students/me` - Получить информацию о студенте
- `GET /api/v1/students/class` - Получить список студентов класса
- `GET /api/v1/grades` - Получить оценки студента
- `GET /api/v1/schedule/weekly` - Получить расписание на неделю
- `GET /api/v1/school/info` - Получить информацию о школе

Для всех защищенных эндпоинтов, кроме `/api/v1/students/me`, требуется указать параметр `instance_url` в запросе или заголовке `X-Instance-URL`. Это позволяет использовать один и тот же токен для доступа к различным экземплярам NetSchool.

## Архитектура

Проект следует принципам чистой архитектуры Go:

- `api/` - Корневой пакет всего проекта
- `api/cmd/` - Точки входа приложения
- `api/internal/` - Внутренние пакеты (не экспортируются)
- `api/internal/app/` - Инициализация приложения
- `api/internal/config/` - Работа с конфигурацией
- `api/internal/domain/` - Бизнес-логика и модели
- `api/internal/infrastructure/` - Инфраструктурные компоненты
- `api/internal/pkg/` - Переиспользуемые утилиты
- `api/pkg/` - Общедоступные пакеты

## Тестирование

Запуск unit-тестов:

```bash
go test -v ./api/...
```

Запуск интеграционных тестов:

```bash
go test -v -tags=integration ./api/...
```

## Docker (опционально)

Для запуска с использованием Docker:

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server api/cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./server", "--config", "config/dev.yaml"]
```

## Лицензия

MIT