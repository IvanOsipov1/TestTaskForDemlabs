# Inventory Go API (SQLite, chi, validator)

Минимальный HTTP-сервис на Go 1.22+:
- Роутер: **chi v5**
- SQLite **без CGO**: `modernc.org/sqlite`
- Валидации: `go-playground/validator/v10`
- Хранение `tags` как JSON в `TEXT`
- Без ORM — чистый `database/sql`
- Эндпоинты: `POST /items`, `GET /items`, `GET /items/{id}`, `PATCH /items/{id}`, `DELETE /items/{id}`

## Быстрый старт

### Windows PowerShell
```powershell
cd path\to\inventory-go
go mod init github.com/yourname/inventory-go
go mod tidy
go run .\cmd\api
# Сервер слушает :8080
