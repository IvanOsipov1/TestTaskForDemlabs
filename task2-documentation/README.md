## Inventory API (FastAPI, in-memory) — Task 2: Documentation

Минимальное, но реалистичное API “Inventory” на FastAPI с in-memory хранилищем и индексом уникальных имён (`name_index`). Поддерживает переименование на `PATCH` с корректным обновлением индекса. Валидации — Pydantic v2. Идентификатор товара — `int`.

---

## Установка и запуск

### Windows PowerShell
    cd task1-testing
    py -m venv .venv
    .\.venv\Scripts\Activate
    pip install -r requirements.txt
    uvicorn app.main:app --reload

### macOS / Linux (bash)
    cd task1-testing
    python3 -m venv .venv
    source .venv/bin/activate
    pip install -r requirements.txt
    uvicorn app.main:app --reload

Swagger UI:
    http://127.0.0.1:8000/docs

---

## Тесты, покрытие, бенчмарк

Запуск всех тестов:
    pytest -q

Покрытие:
    pytest --cov=app --cov-report=term-missing

Бенчмарк (GET /items при 2000 элементах, выборка первых 500):
    pytest -q tests/test_bench_items.py

Ожидаемо:
- Все тесты — passed
- Покрытие `app` — ≈100%
- Бенчмарк — стабилен (в пределах вашей машины)

---

## API справка

### Обзор эндпоинтов

| Метод | Путь                 | Назначение                                   | Успех              | Ошибки                  |
|------:|----------------------|----------------------------------------------|--------------------|-------------------------|
| POST  | /items               | Создать товар                                | 201 Created        | 409 (дубликат), 422     |
| GET   | /items               | Список товаров                               | 200 OK             | —                       |
| GET   | /items/{id:int}      | Получить товар по id                         | 200 OK             | 404                     |
| PATCH | /items/{id:int}      | Частичное обновление (в т.ч. смена имени)    | 200 OK             | 404, 409, 422           |
| DELETE| /items/{id:int}      | Удалить товар                                | 204 No Content     | 404                     |
| GET   | /stats/total_value   | Сумма price*quantity по всем товарам         | 200 OK             | —                       |

Параметры GET /items:
- limit: int | None (диапазон 1..10000)
- offset: int (минимум 0, по умолчанию 0)

---

## Модели (схемы)

### Item
    {
      "id": 1,
      "name": "pen",
      "quantity": 10,
      "price": 1.5,
      "tags": ["stationery"],
      "status": "active"
    }

### ItemCreate (request для POST /items)
    {
      "name": "pen",
      "quantity": 10,
      "price": 1.5,
      "tags": ["stationery"],
      "status": "active"
    }

### ItemUpdate (request для PATCH /items/{id} — все поля опциональны)
    {
      "name": "pen-pro",
      "quantity": 12,
      "price": 1.7,
      "tags": ["stationery", "premium"],
      "status": "archived"
    }

### /stats/total_value (response)
    { "total_value": 13.0 }

---

## Примеры использования

### PowerShell curl (с переносами строк через ^)

Создание:
    curl -X POST http://127.0.0.1:8000/items ^
      -H "Content-Type: application/json" ^
      -d "{\"name\":\"pen\",\"quantity\":10,\"price\":1.5,\"tags\":[\"stationery\"],\"status\":\"active\"}"

Список:
    curl "http://127.0.0.1:8000/items?limit=5&offset=0"

Переименование (PATCH, обновит name_index):
    curl -X PATCH http://127.0.0.1:8000/items/1 ^
      -H "Content-Type: application/json" ^
      -d "{\"name\":\"pen-pro\"}"

Удаление:
    curl -X DELETE http://127.0.0.1:8000/items/1

Суммарная стоимость:
    curl http://127.0.0.1:8000/stats/total_value

### HTTPie
    http GET :8000/items limit==10 offset==0

### Python requests (пример PATCH/DELETE)
    import requests

    base = "http://127.0.0.1:8000"

    r = requests.patch(f"{base}/items/1", json={"status": "archived"})
    print(r.status_code, r.json())

    d = requests.delete(f"{base}/items/1")
    print(d.status_code)  # ожидаемо 204

---

## Ошибки

- 409 Conflict — имя уже существует (POST /items и PATCH при переименовании)
- 404 Not Found — элемент с таким id не найден
- 422 Unprocessable Entity — нарушены правила валидации (обязательные поля, диапазоны)

---

## Troubleshooting

- Активация venv:
    Windows — .\.venv\Scripts\Activate
    macOS/Linux — source .venv/bin/activate

- Проблемы с pip/правами:
    python -m pip install -U pip
    pip install -r requirements.txt

- Порт занят:
    uvicorn app.main:app --reload --port 8001

- Deprecation on_event:
    проект использует lifespan; если встретите on_event — перевести на lifespan

- pytest/coverage:
    запускать из корня task1-testing;
    pytest --cov=app --cov-report=term-missing

- benchmark:
    результаты зависят от машины; гонять отдельно tests/test_bench_items.py

---
