package main

import (
	"sync"
	"sync/atomic"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type Storage struct {
	mu     sync.RWMutex
	todos  map[int]Todo
	nextID int64
}

func NewStorage() *Storage {
	return &Storage{
		todos:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *Storage) Create(todo *Todo) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := int(atomic.AddInt64(&s.nextID, 1) - 1)
	todo.ID = id
	s.todos[id] = *todo
	return id
}

func (s *Storage) GetAll() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		list = append(list, t)
	}
	return list
}

func (s *Storage) GetByID(id int) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.todos[id]
	return t, ok
}

func (s *Storage) Update(id int, todo Todo) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return false
	}
	todo.ID = id
	s.todos[id] = todo
	return true
}

func (s *Storage) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return false
	}
	delete(s.todos, id)
	return true
}
