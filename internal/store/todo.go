package store

import (
	domain "github.com/jneurock/todo-go/internal/domain"
)

type TodoStore interface {
	Create(*domain.Todo) error
	Delete(id string) error
	FindAll() ([]*domain.Todo, error)
	Update(id string, todo *domain.Todo) error
}
