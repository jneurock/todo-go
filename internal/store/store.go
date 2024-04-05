package store

import (
	domain "github.com/jneurock/todo-go/internal/domain"
)

type TodoStore interface {
	Create(description string) error
	Delete(id string) error
	Find(id string) (*domain.Todo, error)
	FindAll() ([]*domain.Todo, error)
	Update(id, description string, complete bool) error
}
