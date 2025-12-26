# Todo API

Небольшой пример REST API для управления списком задач (todo) на Go.

Кратко
- Порт по умолчанию: :8080
- Запуск: `go run .` или сборка бинаря `go build -o ecom .` и запуск `./ecom`

Требования
- Go 1.20+ (советую последнюю стабильную версию)

Как запустить локально

1) Запуск из исходников (компилируется весь пакет, не только main.go):

```bash
go run .
```

2) Или собрать бинарь и запустить:

```bash
go build -o ecom .
./ecom
```

3) Запустить тесты:

```bash
go test ./...
```

Почему не использовать `go run main.go`

Команда `go run main.go` скомпилирует только `main.go` и проигнорирует остальные файлы пакета, поэтому при таком запуске вы можете увидеть ошибки вроде `undefined: NewStorage` — используйте `go run .` или `go build`.

API (эндпоинты)

Base URL: http://localhost:8080

- POST /todos
  - Создать задачу.
  - Тело (JSON): {"title":"...", "description":"...", "completed":false}
  - title обязателен и не может быть пустым.
  - Успех: 201 Created + JSON созданной задачи (включая id).
  - Ошибки: 400 Bad Request (невалидный JSON или пустой title).

- GET /todos
  - Получить список задач.
  - Успех: 200 OK + JSON массив.

- GET /todos/{id}
  - Получить задачу по id.
  - Успех: 200 OK + JSON задачи.
  - Ошибки: 400 (неверный id в URL), 404 (не найдено).

- PUT /todos/{id}
  - Обновить задачу.
  - Тело (JSON): как для POST; title не может быть пустым.
  - Успех: 200 OK + JSON обновлённой задачи.
  - Ошибки: 400, 404.

- DELETE /todos/{id}
  - Удалить задачу.
  - Успех: 204 No Content.
  - Ошибки: 404.

Примеры запросов (curl)

Создать задачу:

```bash
curl -i -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","description":"2 liters","completed":false}'
```

Получить все:

```bash
curl -i -X GET http://localhost:8080/todos
```

Получить по id:

```bash
curl -i -X GET http://localhost:8080/todos/1
```

Обновить:

```bash
curl -i -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk and eggs","description":"Add eggs","completed":true}'
```

Удалить:

```bash
curl -i -X DELETE http://localhost:8080/todos/1
```

Примеры на других языках

JavaScript (fetch):

```js
fetch('http://localhost:8080/todos', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ title: 'Task', description: 'desc', completed: false })
}).then(r => r.json()).then(console.log)
```

Go (http client):

```go
package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

func main() {
  payload := map[string]interface{}{"title":"Go Client Task","description":"desc","completed":false}
  b, _ := json.Marshal(payload)
  resp, err := http.Post("http://localhost:8080/todos", "application/json", bytes.NewReader(b))
  if err != nil { panic(err) }
  defer resp.Body.Close()
  var res map[string]interface{}
  json.NewDecoder(resp.Body).Decode(&res)
  fmt.Println("Status:", resp.StatusCode, "Body:", res)
}
```

Python (requests):

```python
import requests
r = requests.post("http://localhost:8080/todos",
                  json={"title":"Py Task","description":"desc","completed":False})
print(r.status_code, r.json())
```

CORS и браузер

Сервер добавляет заголовок `Access-Control-Allow-Origin: *`, поэтому браузерные запросы из любого источника не должны блокироваться (предусмотрен middleware `enableCORS`).

Отладка

- Проверьте, что сервер запущен и в логах есть "Server starting on :8080".
- При получении 400 проверьте JSON и заголовок `Content-Type: application/json`.
- Если при запуске видите `undefined: NewStorage` — используйте `go run .`.

Тесты

Проект содержит базовые тесты для CRUD в `handlers_test.go`.

```bash
go test ./...
```

Контрибьют

PR и issues приветствуются. Добавляйте тесты для новых фич.

Лицензия

MIT (если хотите — укажите свою лицензию).
