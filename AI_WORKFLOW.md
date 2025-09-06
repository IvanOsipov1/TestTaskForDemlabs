# Тестовое задание для Demlabs

## Приветствие

Приветствую технический отдел команды Demlabs! В этом репозитории вы увидите, как с помощью языковых моделей можно создавать приложения, порой даже не вникая в суть происходящего.

## Обо мне

Я data scientist с большим опытом в машинном обучении. Последние 2 года работаю в компании, занимающейся аутсорсом и IT-проектами.

За эти два года мне приходилось быстро включаться в командный процесс разработки и осваивать технологии с невероятной скоростью, что поначалу (в 22-23 годах) было довольно сложно.

## Опыт работы с LLM

С развитием LLM я начал все больше делегировать нейросетям рутинные задачи. Иногда они помогают мне писать 100% рабочий код (ни один заказчик не возвращался с правками).

В настоящее время у меня очень много работы, поэтому на тестовое задание я уделил около 2 часов (как и описано в файле).

## Мысли о проекте

У меня возникло предположение, что вашей компании может не существовать, а тестовые задания собираются для обучения LLM для более крупных компаний. Хотя доказать это я не могу, но даже если это так - я не против косвенно поучаствовать в обучении чего-то интересного.

## Технические детали

- Репозиторий создан самостоятельно
- Всю работу с Git выполнял лично (занимаюсь этим регулярно, не будем по мелочам грузить гптшку)
- Использовался только ChatGPT-4
- Не использовались другие инструменты

## Выполнение заданий (на все тестовое задание ушло 4 запроса)

### Задание 1
- **1 запрос** - создание приложения и тестов
- **2 запрос** - анализ логов и микро-фиксы

### Задание 2
- **1 запрос** - повторный анализ логов + создание README.md в markdown
- Результат получился грамотным, переделывать не потребовалось

### Задание 3
- **Особенность**: никогда ранее не писал на Go
- **1 запрос** - создание работающего приложения
- **Самостоятельно**: установил компилятор и настроил переменные среды (gpt 5 не написал инструкцию по установке компилятора, загуглил и сделал все сам за 5 минут)

## Итоги

Мне понравилось работать над заданием. Конечно, я уделил мало времени анализу всего, что сгенерировал ChatGPT, но главное - решение работает, что делает его успешным MVP.

---

# ПРОМТЫ

## Промт 1 - Приложение на Fast API
Сделай минимальное, но реалистичное FastAPI-приложение "Inventory" и полный набор тестов. Требования:

Функциональность:

Модель Item: name (уникальный, 1..64), quantity (>=0), price (>0), tags: list[str], status: ["active","archived"].

Эндпоинты:
POST /items (создать), GET /items, GET /items/{id}, PATCH /items/{id}, DELETE /items/{id}
GET /stats/total_value (сумма price*quantity по всем item'ам).

In-memory хранилище с индексом уникальных имен (name_index). На PATCH поддержи смену имени с корректным обновлением индекса.

Тесты:

Юнит: чистая бизнес-логика для функции calculate_total_value(items).

Интеграционные (FastAPI TestClient): полный CRUD + /stats/total_value.

Edge/error: дубликат имени (409), 404/422, конфликты при смене имени, невалидные данные.

Перф: pytest-benchmark на GET /items со списком из 2000 записей, выборка первых 500.

Технические требования:

Без внешней БД, только in-memory.

Валидации на Pydantic v2.

Структура папок: /task1-testing/app(...), /task1-testing/tests(...), requirements.txt, README.md.

Тесты должны гоняться без поднятия реального сервера (TestClient).

Отдай готовые файлы целиком (код), команды для Windows PowerShell (venv, pytest, coverage, benchmark).

Выведи:

Дерево проекта.

Полные содержимые всех файлов.

Команды запуска (Windows PowerShell).

## Промт 2 - Анализ логов
Анализируй мои логи pytest и coverage (ниже). Задача: предложи только минимально необходимые фиксы для "зелёного" состояния без варнингов и с высоким покрытием (≈98–100%). Конкретно:

Если есть DeprecationWarning для @app.on_event("startup"), переведи на FastAPI lifespan и покажи DIFF-патч.

Если coverage пропустил ветку PATCH со сменой имени (name_index), добавь точечные тесты и покажи файл(ы) целиком.

Если есть дырки в main.py по строкам lifespan, предложи один тест, который заходит в lifespan (TestClient как context manager), чтобы покрыть startup/shutdown.

Дай ожидаемый результат прогона: команды и целевые цифры (pass/coverage/benchmark).

Логи:
<<< (.venv) PS C:\Projects\TestTaskForDemlabs\task1-testing> pytest -q

----------------------------------------------------- benchmark: 1 tests -----------------------------------------------------
Name (time in ms) Min Max Mean StdDev Median IQR Outliers OPS Rounds Iterations

test_benchmark_get_items_first_500 3.6386 27.5127 6.3319 6.1653 3.8308 0.4787 10;11 157.9302 65
1

Legend:
Outliers: 1 Standard Deviation from Mean; 1.5 IQR (InterQuartile Range) from 1st Quartile and 3rd Quartile.
OPS: Operations Per Second, computed as 1 / Mean
11 passed in 1.83s
(.venv) PS C:\Projects\TestTaskForDemlabs\task1-testing> pytest --cov=app --cov-report=term-missing
================================================== test session starts ==================================================
platform win32 -- Python 3.12.7, pytest-8.4.1, pluggy-1.6.0
benchmark: 5.1.0 (defaults: timer=time.perf_counter disable_gc=False min_rounds=5 min_time=0.000005 max_time=1.0 calibration_precision=10 warmup=False warmup_iterations=100000)
rootdir: C:\Projects\TestTaskForDemlabs\task1-testing
plugins: anyio-4.10.0, benchmark-5.1.0, cov-6.2.1
collected 11 items

tests\test_bench_items.py . [ 9%]
tests\test_integration_inventory.py ........ [ 81%]
tests\test_unit_total_value.py .. [100%]

==================================================== tests coverage =====================================================
____________________________________ coverage: platform win32, python 3.12.7-final-0 ____________________________________

Name Stmts Miss Cover Missing

app_init_.py 0 0 100%
app\main.py 39 0 100%
app\models.py 18 0 100%
app\storage.py 76 3 96% 61, 97, 100

----------------------------------------------------- benchmark: 1 tests ----------------------------------------------------
Name (time in ms) Min Max Mean StdDev Median IQR Outliers OPS Rounds Iterations

test_benchmark_get_items_first_500 6.4958 8.1127 6.7484 0.3109 6.6289 0.2743 4;1 148.1830 34
1

Legend:
Outliers: 1 Standard Deviation from Mean; 1.5 IQR (InterQuartile Range) from 1st Quartile and 3rd Quartile.
OPS: Operations Per Second, computed as 1 / Mean
================================================== 11 passed in 2.79s ===================================================
(.venv) PS C:\Projects\TestTaskForDemlabs\task1-testing> pytest -q tests/test_bench_items.py
. [100%]

----------------------------------------------------- benchmark: 1 tests -----------------------------------------------------
Name (time in ms) Min Max Mean StdDev Median IQR Outliers OPS Rounds Iterations

test_benchmark_get_items_first_500 3.6021 27.5860 5.8273 5.7598 3.8270 0.4157 7;8 171.6062 62
1

Legend:
Outliers: 1 Standard Deviation from Mean; 1.5 IQR (InterQuartile Range) from 1st Quartile and 3rd Quartile.
OPS: Operations Per Second, computed as 1 / Mean
1 passed in 1.72s
(.venv) PS C:\Projects\TestTaskForDemlabs\task1-testing> >>>


## Промт 3 - Документация
Создай красивую документацию для моего FastAPI приложения. У меня все тесты проходят, coverage 100%. Нужно:

README.md с:

Инструкциями по установке и запуску (Windows/Linux)

Описанием API эндпоинтов

Примеры запросов (curl, HTTPie, Python)

Описанием ошибок

ЛОГИ с предыдущего задания (зелёные, для справки):
<<< ----------------------------------------------------- benchmark: 1 tests -----------------------------------------------------
Name (time in ms)                         Min      Max    Mean  StdDev  Median     IQR  Outliers       OPS  Rounds  Iterations
------------------------------------------------------------------------------------------------------------------------------
test_benchmark_get_items_first_500     3.6324  27.5657  6.1406  5.8436  3.7842  0.3679       8;9  162.8513      54        
   1
------------------------------------------------------------------------------------------------------------------------------

Legend:
  Outliers: 1 Standard Deviation from Mean; 1.5 IQR (InterQuartile Range) from 1st Quartile and 3rd Quartile.
  OPS: Operations Per Second, computed as 1 / Mean
13 passed in 1.89s
(.venv) PS C:\Projects\TestTaskForDemlabs\task1-testing> pytest --cov=app --cov-report=term-missing
================================================== test session starts ==================================================
platform win32 -- Python 3.12.7, pytest-8.4.1, pluggy-1.6.0
benchmark: 5.1.0 (defaults: timer=time.perf_counter disable_gc=False min_rounds=5 min_time=0.000005 max_time=1.0 calibration_precision=10 warmup=False warmup_iterations=100000)
rootdir: C:\Projects\TestTaskForDemlabs\task1-testing
plugins: anyio-4.10.0, benchmark-5.1.0, cov-6.2.1
collected 13 items

tests\test_bench_items.py .                                                                                        [  7%]
tests\test_integration_inventory.py ........                                                                       [ 69%]
tests\test_storage_edges.py ..                                                                                     [ 84%]
tests\test_unit_total_value.py ..                                                                                  [100%]

==================================================== tests coverage ===================================================== 
____________________________________ coverage: platform win32, python 3.12.7-final-0 ____________________________________ 

Name              Stmts   Miss  Cover   Missing
-----------------------------------------------
app\__init__.py       0      0   100%
app\main.py          39      0   100%
app\models.py        18      0   100%
app\storage.py       76      0   100%
-----------------------------------------------

----------------------------------------------------- benchmark: 1 tests -----------------------------------------------------
Name (time in ms)                         Min      Max     Mean  StdDev  Median     IQR  Outliers      OPS  Rounds  Iterations
------------------------------------------------------------------------------------------------------------------------------
test_benchmark_get_items_first_500     6.2573  32.0091  12.1168  9.1292  7.3578  2.8093       9;9  82.5300      38        
   1
------------------------------------------------------------------------------------------------------------------------------

Legend:
  Outliers: 1 Standard Deviation from Mean; 1.5 IQR (InterQuartile Range) from 1st Quartile and 3rd Quartile.
  OPS: Operations Per Second, computed as 1 / Mean
================================================== 13 passed in 2.90s =================================================== 
(.venv) PS C:\Projects\TestTaskForDemlabs\task1-testing> pytest -q tests/test_bench_items.py       
.                                                                                                                  [100%]

----------------------------------------------------- benchmark: 1 tests -----------------------------------------------------
Name (time in ms)                         Min      Max    Mean  StdDev  Median     IQR  Outliers       OPS  Rounds  Iterations
------------------------------------------------------------------------------------------------------------------------------
test_benchmark_get_items_first_500     3.7277  23.1360  5.4110  4.2268  4.2241  0.6655       4;4  184.8093      50        
   1
------------------------------------------------------------------------------------------------------------------------------

Legend:
  Outliers: 1 Standard Deviation from Mean; 1.5 IQR (InterQuartile Range) from 1st Quartile and 3rd Quartile.
  OPS: Operations Per Second, computed as 1 / Mean
1 passed in 1.72s >>>


## Промт 4 - Незнакомый язык

Теперь нужно быстро собрать минимальный, но реалистичный HTTP-сервис на Go с SQLite. Чистая структура, валидация входных данных, понятные ответы, пара юнит/интеграционных тестов на httptest.

Что нужно:

1) Стек и зависимости
- Go 1.22+
- Роутер: chi v5
- SQLite без CGO: modernc.org/sqlite
- Валидации: go-playground/validator/v10
- JSON-кодек стандартный. Для массива tags храним JSON в TEXT.

2) Модель и БД
- Item: id:int (PK autoincrement), name:string (уникальный, длина 1..64), quantity:int (>=0), price:float64 (>0), tags:[]string, status: "active"|"archived".
- SQLite файл: inventory.db в корне проекта.
- Схема SQL инициализируется при старте (миграция 0001_init.sql): таблица items, уникальный индекс по name, CHECK на status.

3) Эндпоинты (3–4 штуки)
- POST /items → 201; 409 при дубликате name; 422 при невалидных данных.
- GET /items → 200, с query: limit (1..1000), offset (>=0, по умолчанию 0).
- GET /items/{id} → 200 или 404.
- PATCH /items/{id} → 200; частичное обновление (можно менять name, проверяя уникальность → 409), иначе 404/422.
(если удобно — можно добавить DELETE /items/{id} → 204/404, но это опционально)

4) Ответы и ошибки
- Всегда JSON. На ошибки отдаём {"error":"…","details":{…}}.
- Чётко выставляй HTTP-коды (201/204/404/409/422).

5) Структура проекта (пример)
- cmd/api/main.go — инициализация, роуты, запуск
- internal/http/handlers.go — обработчики
- internal/store/sqlite.go — подключение к БД и методы CRUD
- internal/model/item.go — структуры DTO/модели
- internal/validate/validate.go — инициализация validator + кастомные правила
- db/migrations/0001_init.sql — схема
- go.mod
- README.md — как запустить (Windows PowerShell + Bash), примеры curl, как гонять тесты

6) Тесты
- На httptest: happy-path POST/GET, 404, 409 на дубликат имени, PATCH (в т.ч. переименование), базовая 422.
- Для тестов используй in-memory SQLite: "file::memory:?cache=shared&_pragma=foreign_keys(1)".
- Команда: go test ./... -v

7) Что отдать в ответе
- Дерево проекта.
- Полное содержимое всех файлов (код сразу компилируемый), включая go.mod и миграцию.
- В README.md: команды для Windows PowerShell и Bash (go mod tidy, запуск, curl-примеры для всех эндпоинтов, go test).

Пара заметок:
- Никаких ORM, только database/sql.
- На PATCH обновляй только те поля, что пришли (partial update).
- Конфликт уникального имени обрабатывай и на уровне БД (unique index), и в коде (маппинг ошибки в 409).

Сделай, пожалуйста, полностью рабочий пример: после копирования файлов `go mod tidy`, `go run ./cmd/api`, curl-примеры должны проходить, а `go test ./... -v` — зелёные. Если где-то есть тонкости (Windows vs Linux), кратко отметь в README.
