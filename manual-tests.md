# Manual tests (curl, PowerShell)

База: `http://localhost:8080`  
Команды запускались в PowerShell. Для POST/PUT использовались JSON-файлы в UTF-8 без BOM.

---

## 0) Запуск сервера

```powershell
go run ./cmd/server
---

## 1) GET /health

✅ OK
```powershell
curl.exe -i "http://localhost:8080/health"
❌ 405

curl.exe -i -X POST "http://localhost:8080/health"
GET /tasks

✅ OK 

curl.exe -i "http://localhost:8080/tasks"


❌ 405

curl.exe -i -X PATCH "http://localhost:8080/tasks"
POST /tasks

✅ OK 

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText("body.json", '{"title":"Buy milk","done":false}', $utf8NoBom)

curl.exe -i -X POST "http://localhost:8080/tasks" `
  -H "Content-Type: application/json" `
  --data-binary "@body.json"


❌ 400 

[System.IO.File]::WriteAllText("bad-post.json", '{"done":false}', $utf8NoBom)

curl.exe -i -X POST "http://localhost:8080/tasks" `
  -H "Content-Type: application/json" `
  --data-binary "@bad-post.json"
GET /tasks/{id}

✅ OK

curl.exe -i "http://localhost:8080/tasks/1"


❌ 404

curl.exe -i "http://localhost:8080/tasks/999"
PUT /tasks/{id}

✅ OK

[System.IO.File]::WriteAllText("put.json", '{"title":"Buy milk ASAP","done":true}', $utf8NoBom)

curl.exe -i -X PUT "http://localhost:8080/tasks/1" `
  -H "Content-Type: application/json" `
  --data-binary "@put.json"


❌ 400

[System.IO.File]::WriteAllText("bad-put.json", '{"title":"   ","done":true}', $utf8NoBom)

curl.exe -i -X PUT "http://localhost:8080/tasks/1" `
  -H "Content-Type: application/json" `
  --data-binary "@bad-put.json"
  DELETE /tasks/{id}

✅ OK

curl.exe -i -X DELETE "http://localhost:8080/tasks/1"


❌ 404

curl.exe -i -X DELETE "http://localhost:8080/tasks/999"
Финальная проверка
curl.exe -i "http://localhost:8080/tasks"
curl.exe -i "http://localhost:8080/tasks/1"