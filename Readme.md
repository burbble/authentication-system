# Техническая документация системы аутентификации

## Содержание
1. [Архитектура](#архитектура)
2. [API Endpoints](#api-endpoints)
3. [Разработка](#разработка)

## Архитектура

### Слои приложения
- **Domain Layer** (`internal/domain/`)
  - Бизнес-сущности
  - Интерфейсы репозиториев
  - Обработка ошибок домена

- **Application Layer** (`internal/application/`)
  - Use Cases
  - DTOs
  - Сервисные интерфейсы

- **Infrastructure Layer** (`internal/infrastructure/`)
  - Реализация репозиториев
  - Email сервис
  - Persistence

- **Interfaces Layer** (`internal/interfaces/`)
  - HTTP handlers
  - Middleware
  - Routing

## API Endpoints

### Регистрация пользователя
**POST** `/api/v1/register`

Request:
```json
{
    "username": "string (3-30 символов)",
    "name": "string (2-50 символов)",
    "email": "string (валидный email)",
    "password": "string (мин. 8 символов)"
}
```

Response (200 OK):
```json
{
    "id": "UUID",
    "username": "string",
    "name": "string",
    "email": "string"
}
```

### Верификация Email
**POST** `/api/v1/verify-email`

Request:
```json
{
    "email": "string",
    "code": "string (6 символов)"
}
```

Response:
- 200 OK: Верификация успешна
- 400 Bad Request: Неверный код или email

## Разработка

### Требования к системе
- Go 1.22+
- PostgreSQL 14+
- Docker & Docker Compose
- Make

### Команды Make
```bash
# Запуск приложения
make start

# Пересоздание базы данных
make reset-db

# Применение миграций
make migrate

# Создание новой миграции
make create-migration name=migration_name

# Просмотр логов
make app-logs    # логи приложения
make db-logs     # логи базы данных
```

### Тестирование Email сервиса
По умолчанию используется mock-сервис для email. Для включения реального SMTP:

1. Установите `USE_MOCK_EMAIL=false`
2. Настройте SMTP параметры в `.env`

### Безопасность
- Пароли хешируются с использованием bcrypt
- Все эндпоинты защищены CORS middleware
- Реализована защита от паник через recover middleware
- Логирование всех HTTP-запросов
