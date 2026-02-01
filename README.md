# tasks-api

REST API для управления списком задач (to-do).  
Хранение — in-memory (без БД). Формат обмена — только JSON.

## Требования
- Go 1.20+ (рекомендуется 1.21+)

## Запуск

```powershell
go mod tidy
go run ./cmd/server
```