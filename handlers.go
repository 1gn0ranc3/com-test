package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (s *Storage) todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /todos
		todos := s.GetAll()
		jsonResponse(w, http.StatusOK, todos)

	case http.MethodPost:
		// POST /todos
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(todo.Title) == "" {
			http.Error(w, "Title cannot be empty", http.StatusBadRequest)
			return
		}
		id := s.Create(&todo)
		newTodo, _ := s.GetByID(id)
		jsonResponse(w, http.StatusCreated, newTodo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Storage) todoByIDHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	if path == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, ok := s.GetByID(id)
		if !ok {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		jsonResponse(w, http.StatusOK, todo)

	case http.MethodPut:
		var updated Todo
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(updated.Title) == "" {
			http.Error(w, "Title cannot be empty", http.StatusBadRequest)
			return
		}
		if ok := s.Update(id, updated); !ok {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		todo, _ := s.GetByID(id)
		jsonResponse(w, http.StatusOK, todo)

	case http.MethodDelete:
		if ok := s.Delete(id); !ok {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
