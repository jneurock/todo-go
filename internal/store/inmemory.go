package store

import (
	"errors"
	"strconv"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/util"
)

type TodoInMemoryStore struct {
	todos util.Queue[domain.Todo]
}

var lastID int64 = -1

func NewTodoInMemoryStore() *TodoInMemoryStore {
	return &TodoInMemoryStore{}
}

func (s *TodoInMemoryStore) Create(description string) error {
	todo := domain.NewTodo(description)

	todo.ID = lastID + 1
	lastID = todo.ID

	s.todos.Add(todo)

	return nil
}

func (s *TodoInMemoryStore) Delete(id string) error {
	intId, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	s.todos.Remove(func(todo *domain.Todo) bool {
		return todo.ID == int64(intId)
	})

	return nil
}

func (s *TodoInMemoryStore) Find(id string) (*domain.Todo, error) {
	intId, err := strconv.Atoi(id)

	if err != nil {
		return nil, err
	}

	todo := s.todos.Head

	for todo != nil {
		if todo.Value.ID == int64(intId) {
			return todo.Value, nil
		}

		todo = todo.Next
	}

	return nil, errors.New("todo not found")
}

func (s *TodoInMemoryStore) FindAll() ([]*domain.Todo, error) {
	var todos []*domain.Todo
	todo := s.todos.Head

	for todo != nil {
		todos = append([]*domain.Todo{todo.Value}, todos...)
		todo = todo.Next
	}

	return todos, nil
}

func (s *TodoInMemoryStore) Update(id, description string, complete bool) error {
	todo, err := s.Find(id)

	if err != nil {
		return err
	}

	todo.Description = description
	todo.Complete = complete

	return nil
}
