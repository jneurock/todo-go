package service

import (
	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/store"
)

type TodoService struct {
	store store.TodoStore
}

func NewTodoService(store store.TodoStore) *TodoService {
	return &TodoService{store: store}
}

// TODO: Validate description
func (ts *TodoService) Create(description string) (*domain.Todo, error) {
	return ts.store.Create(store.TodoAttrs{
		Description: description,
	})
}

func (ts *TodoService) Delete(id string) error {
	return ts.store.Delete(id)
}

func (ts *TodoService) FindAll() ([]*domain.Todo, error) {
	return ts.store.FindAll()
}
