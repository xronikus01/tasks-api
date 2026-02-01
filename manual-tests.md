# Manual tests (curl, PowerShell)

База: `http://localhost:8080`  
Команды выполнялись в PowerShell. Для POST/PUT использовались JSON-файлы в UTF-8 **без BOM**.

---

## 0) Запуск сервера

```powershell
go run ./cmd/server
```
### Ожидаемо:

сервер слушает :8080

## 1) GET /health

✅ OK

```powershell
curl.exe -i "http://localhost:8080/health"
```

Ожидаемо: 200 OK, тело {"status":"ok"}

❌ 405

```powershell
curl.exe -i -X POST "http://localhost:8080/health"
```

Ожидаемо: 405 Method Not Allowed

## 2) GET /tasks

✅ OK (пустой список)

```powershell
curl.exe -i "http://localhost:8080/tasks"
```

Ожидаемо: 200 OK, тело []

❌ 405

```powershell
curl.exe -i -X PATCH "http://localhost:8080/tasks"
```

Ожидаемо: 405 Method Not Allowed

## 3) POST /tasks

✅ OK (создание)
```powershell
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText("body.json", '{"title":"Buy milk","done":false}', $utf8NoBom)

curl.exe -i -X POST "http://localhost:8080/tasks" `
  -H "Content-Type: application/json" `
  --data-binary "@body.json"

```

Ожидаемо: 201 Created, JSON с id

❌ 400 (нет title)

```powershell
[System.IO.File]::WriteAllText("bad-post.json", '{"done":false}', $utf8NoBom)

curl.exe -i -X POST "http://localhost:8080/tasks" `
  -H "Content-Type: application/json" `
  --data-binary "@bad-post.json"

```
Ожидаемо: 400 Bad Request, {"error":"title is required"}

## 4) GET /tasks/{id}

Примечание: используйте id из ответа POST (в примере — 1).

✅ OK
```powershell
curl.exe -i "http://localhost:8080/tasks/1"
```

Ожидаемо: 200 OK

❌ 404
```powershell
curl.exe -i "http://localhost:8080/tasks/999"
```

Ожидаемо: 404 Not Found, {"error":"task not found"}

## 5) PUT /tasks/{id}

✅ OK (полное обновление)
```powershell
[System.IO.File]::WriteAllText("put.json", '{"title":"Buy milk ASAP","done":true}', $utf8NoBom)

curl.exe -i -X PUT "http://localhost:8080/tasks/1" `
  -H "Content-Type: application/json" `
  --data-binary "@put.json"
```

Ожидаемо: 200 OK

❌ 400 (пустой title)
```powershell
[System.IO.File]::WriteAllText("bad-put.json", '{"title":"   ","done":true}', $utf8NoBom)

curl.exe -i -X PUT "http://localhost:8080/tasks/1" `
  -H "Content-Type: application/json" `
  --data-binary "@bad-put.json"
```

Ожидаемо: 400 Bad Request, {"error":"title is required"}

## 6) DELETE /tasks/{id}

✅ OK
```powershell
curl.exe -i -X DELETE "http://localhost:8080/tasks/1"
```

Ожидаемо: 204 No Content

❌ 404
```powershell
curl.exe -i -X DELETE "http://localhost:8080/tasks/999"
```

Ожидаемо: 404 Not Found, {"error":"task not found"}

## 7) Финальная проверка
```powershell
curl.exe -i "http://localhost:8080/tasks"
curl.exe -i "http://localhost:8080/tasks/1"
```
Ожидаемо:

/tasks -> 200 OK, []

/tasks/1 -> 404 Not Found