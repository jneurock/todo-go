package store

import (
	domain "github.com/jneurock/todo-go/internal/domain"
)

type TodoStore interface {
	Create(todo *domain.Todo) error
	Delete(id string) error
	FindAll() ([]*domain.Todo, error)
	Update(todo *domain.Todo) error
}
