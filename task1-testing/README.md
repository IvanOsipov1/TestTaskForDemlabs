# Inventory (FastAPI, in-memory)

Минимальное, но реалистичное FastAPI-приложение с in-memory хранилищем и индексом уникальных имён (`name_index`).
Поддерживает переименование на PATCH с корректным обновлением индекса и конфликт-менеджментом (409).

## Функциональность

- Модель `Item`:
  - `name` (уникальный, длина 1..64)
  - `quantity` (целое, >= 0)
  - `price` (число, > 0)
  - `tags` (список строк)
  - `status` (одно из: `"active"`, `"archived"`)
- Эндпоинты:
  - `POST /items` — создать
  - `GET /items` — список (поддерживает `limit`, `offset`)
  - `GET /items/{id}` — получить по id
  - `PATCH /items/{id}` — частичное обновление (в т.ч. смена имени с обновлением индекса и защитой от конфликтов)
  - `DELETE /items/{id}` — удалить
  - `GET /stats/total_value` — сумма `price * quantity` по всем item’ам
- Хранилище: только in-memory (`dict` + индекс `name_index`)
- Валидации — Pydantic v2 (через `Field(...)`)
- Тесты: pytest + fastapi TestClient + pytest-benchmark

## Структура папок

См. дерево в корне.

## Запуск (Windows PowerShell)

```powershell
# 1) Создать и активировать venv
py -m venv .venv
.\.venv\Scripts\Activate

# 2) Установить зависимости
pip install -r requirements.txt

# 3) Запустить тесты (юнит + интеграционные + edge)
pytest -q

# 4) Покрытие
pytest --cov=app --cov-report=term-missing

# 5) Бенчмарк (GET /items с 2000 записями, выборка первых 500)
pytest -q tests/test_bench_items.py

# (Опционально) Локально покрутить API
uvicorn app.main:app --reload

