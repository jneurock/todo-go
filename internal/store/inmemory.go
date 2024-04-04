package store

import (
	"errors"
	"strconv"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/util"
)

type InMemoryTodoStore struct {
	todos util.Queue[domain.Todo]
}

var lastID int64 = -1

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{}
}

func (s *InMemoryTodoStore) Create(todo *domain.Todo) error {
	todo.ID = lastID + 1
	lastID = todo.ID

	s.todos.Add(todo)

	return nil
}

func (s *InMemoryTodoStore) Delete(id string) error {
	intId, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	s.todos.Remove(func(todo *domain.Todo) bool {
		return todo.ID == int64(intId)
	})

	return nil
}

func (s *InMemoryTodoStore) FindAll() ([]*domain.Todo, error) {
	var todos []*domain.Todo
	todo := s.todos.Head

	for todo != nil {
		todos = append([]*domain.Todo{todo.Value}, todos...)
		todo = todo.Next
	}

	return todos, nil
}

func (s *InMemoryTodoStore) Update(updatedTodo *domain.Todo) error {
	todo := s.todos.Head

	for todo != nil {
		if todo.Value.ID == updatedTodo.ID {
			todo.Value.Complete = updatedTodo.Complete
			todo.Value.Description = updatedTodo.Description

			return nil
		}

		todo = todo.Next
	}

	return errors.New("todo not found")
}
