package store

import (
	domain "github.com/jneurock/todo-go/internal/domain"
)

type TodoAttrs struct {
	Complete    bool
	Description string
}

type TodoStore interface {
	Create(attrs TodoAttrs) (*domain.Todo, error)
	Delete(id string) error
	Find(id string) (*domain.Todo, error)
	FindAll() ([]*domain.Todo, error)
	Update(id string, attrs TodoAttrs) error
}
