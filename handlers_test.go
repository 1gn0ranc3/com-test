package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCRUD(t *testing.T) {
	storage := NewStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", storage.todosHandler)
	mux.HandleFunc("/todos/", storage.todoByIDHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// 1. Создание задачи (успешно)
	createPayload := `{"title":"Test Todo","description":"desc","completed":false}`
	resp, body := testRequest(t, server, "POST", "/todos", bytes.NewBufferString(createPayload))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201, got %d: %s", resp.StatusCode, body)
	}
	var created Todo
	json.Unmarshal([]byte(body), &created)
	if created.Title != "Test Todo" || created.ID == 0 {
		t.Fatal("Created todo invalid")
	}

	id := created.ID

	// 2. Создание с пустым title (валидация)
	resp, body = testRequest(t, server, "POST", "/todos", bytes.NewBufferString(`{"title":"","description":""}`))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 on empty title, got %d: %s", resp.StatusCode, body)
	}

	// 3. Получение списка
	resp, body = testRequest(t, server, "GET", "/todos", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}
	var list []Todo
	json.Unmarshal([]byte(body), &list)
	if len(list) != 1 || list[0].ID != id {
		t.Fatal("List invalid")
	}

	// 4. Получение по ID
	resp, body = testRequest(t, server, "GET", "/todos/"+strconv.Itoa(id), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// 5. Получение несуществующего
	resp, _ = testRequest(t, server, "GET", "/todos/999", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("Expected 404 for non-existing")
	}

	// 6. Обновление
	updatePayload := `{"title":"Updated","description":"new desc","completed":true}`
	resp, body = testRequest(t, server, "PUT", "/todos/"+strconv.Itoa(id), bytes.NewBufferString(updatePayload))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 on update, got %d: %s", resp.StatusCode, body)
	}
	var updated Todo
	json.Unmarshal([]byte(body), &updated)
	if updated.Title != "Updated" || !updated.Completed {
		t.Fatal("Update failed")
	}

	// 7. Обновление с пустым title
	resp, _ = testRequest(t, server, "PUT", "/todos/"+strconv.Itoa(id), bytes.NewBufferString(`{"title":""}`))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Expected 400 on empty title update")
	}

	// 8. Удаление
	resp, _ = testRequest(t, server, "DELETE", "/todos/"+strconv.Itoa(id), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatal("Expected 204 on delete")
	}

	// 9. Удаление несуществующего
	resp, _ = testRequest(t, server, "DELETE", "/todos/999", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("Expected 404 on delete non-existing")
	}
}

func testRequest(t *testing.T, server *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, server.URL+path, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return resp, string(respBody)
}
