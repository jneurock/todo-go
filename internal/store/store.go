package store

import (
	"github.com/jneurock/todo-go/internal/domain"
)

type TodoStore interface {
	IsAvailable() bool
	Create(description string) error
	Delete(id string) error
	Find(id string) (*domain.Todo, error)
	FindAll() ([]*domain.Todo, error)
	Update(id, description string, complete bool) error
}
