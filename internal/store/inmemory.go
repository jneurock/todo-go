package store

import (
	domain "github.com/jneurock/todo-go/internal/domain"
)

type InMemoryTodoStore struct {
	todos []*domain.Todo
}

var lastID int64 = -1

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{}
}

func (s *InMemoryTodoStore) Create(attrs TodoAttrs) (*domain.Todo, error) {
	todo := &domain.Todo{
		Description: attrs.Description,
		ID:          lastID + 1,
	}

	lastID = todo.ID

	s.todos = append([]*domain.Todo{todo}, s.todos...)

	return todo, nil
}

func (s *InMemoryTodoStore) Delete(id string) error {
	return nil
}

func (s *InMemoryTodoStore) Find(id string) (*domain.Todo, error) {
	return nil, nil
}

func (s *InMemoryTodoStore) FindAll() ([]*domain.Todo, error) {
	return s.todos, nil
}

func (s *InMemoryTodoStore) Update(id string, attrs TodoAttrs) error {
	return nil
}
