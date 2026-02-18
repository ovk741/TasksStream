# TasksStream API Documentation

TasksStream — это backend-сервис таск-трекера, написанный на Go.  
Проект реализует доски, колонки и карточки с разделением ответственности.

Общая информация

Base URL: http://localhost:8080

Формат данных: JSON

Аутентификация: JWT authentication (access + refresh tokens)

## Архитектура

internal/
 ├── api/http        # HTTP handlers
 ├── service         # Бизнес-логика
 ├── storage         # PostgreSQL репозитории
 ├── domain          # Доменные модели и ошибки
 ├── infra/auth      # JWT
 └── infra/security  # Password hashing (bcrypt)


---

## Аутентификация

Используется Bearer токен:

Authorization: Bearer <access_token>

Endpoints

POST /auth/register

POST /auth/login

POST /auth/refresh


## Boards API

| Метод  | Endpoint      | Описание              |
| ------ | ------------- | --------------------- |
| POST   | `/boards`     | Создать доску         |
| GET    | `/boards`     | Получить список досок |
| PUT    | `/boards?id=` | Обновить доску        |
| DELETE | `/boards?id=` | Удалить доску         |

## Board members

| Метод  | Endpoint                    | Описание                |
| ------ | --------------------------- | ----------------------- |
| POST   | `/boards/invite`            | Пригласить пользователя |
| GET    | `/boards/members?board_id=` | Получить участников     |
| DELETE | `/boards/members/remove`    | Удалить участника       |

## Columns API

| Метод  | Endpoint             | Описание               |
| ------ | -------------------- | ---------------------- |
| POST   | `/columns`           | Создать колонку        |
| GET    | `/columns?board_id=` | Получить колонки доски |
| PUT    | `/columns?id=`       | Обновить колонку       |
| DELETE | `/columns?id=`       | Удалить колонку        |
| PUT    | `/columns/move`      | Переместить колонку    |

## Tasks API

| Метод  | Endpoint            | Описание           |
| ------ | ------------------- | ------------------ |
| POST   | `/tasks`            | Создать задачу     |
| GET    | `/tasks?column_id=` | Получить задачи    |
| PUT    | `/tasks?id=`        | Обновить задачу    |
| DELETE | `/tasks?id=`        | Удалить задачу     |
| PUT    | `/tasks/move`       | Переместить задачу |










