# TasksStream API Documentation

TasksStream — это backend-сервис таск-трекера, написанный на Go.  
Проект реализует доски, колонки и карточки с разделением ответственности.

Общая информация

Base URL: http://localhost:8080

Формат данных: JSON

Аутентификация: нет (пока)

## Архитектура

HTTP Request
↓
Handlers (API)
↓
Services (бизнес-логика)
↓
Repositories (storage interface)
↓
Memory Storage (реализация)

## Структура проекта

internal/
api/http/ # HTTP handlers
service/ # Бизнес-логика
storage/
memory/ # In-memory репозитории
domain/ # Доменные модели
cmd/
app/
main.go

---

Сущности:

Board (Доска)
{
  "id": "string",
  "name": "string",
  "created_at": "RFC3339 timestamp"
}

Column (Колонка)
{
  "id": "string",
  "title": "string",
  "board_id": "string",
  "created_at": "RFC3339 timestamp"
}

Task (Карточка)
{
  "id": "string",
  "title": "string",
  "description": "string",
  "column_id": "string",
  "created_at": "RFC3339 timestamp"
}


---

Boards API

Создать доску

POST /boards

Request:

{
  "name": "My board"
}
Response 201 Created:

{
  "id": "board-1",
  "name": "My board",
  "created_at": "2026-01-27T12:00:00Z"
}

Получить все доски

GET /boards

Response 200 OK:

[
  {
    "id": "board-1",
    "name": "My board",
    "created_at": "2026-01-27T12:00:00Z"
  }
]

Columns API

Создать колонку

POST /columns

Request:

{
  "board_id": "board-1",
  "title": "To Do"
}
Response 201 Created:

{
  "id": "column-1",
  "title": "To Do",
  "board_id": "board-1",
  "created_at": "2026-01-27T12:05:00Z"
}

Получить колонки доски

GET /columns?board_id={board_id}

Tasks API

Создать задачу

POST /tasks

Request:

{
  "title": "Implement API",
  "description": "Write handlers and services",
  "column_id": "column-1"
}
Response 201 Created:

{
  "id": "task-1",
  "title": "Implement API",
  "description": "Write handlers and services",
  "column_id": "column-1",
  "created_at": "2026-01-27T12:10:00Z"
}

Получить задачи колонки

GET /tasks?column_id={column_id}

Обработка ошибок

Ошибки возвращаются в формате:

{
  "error": "error message"
}

Основные ошибки:

HTTP Code	Error
400	invalid input
404	not found
500	internal error


Тестирование
Запуск всех тестов:

go test ./...
Тесты покрывают:

сервисы (unit)

HTTP handlers (в процессе)

бизнес-валидацию

🛠 Запуск проекта
go run cmd/app/main.go
Сервер будет доступен по адресу:

http://localhost:8080



1. Запросы Удалить, Пеместить 
2. Подключить Постгрессукеле
3. Докер подключить
4. Тесты написать 
5. Документация АПИ

