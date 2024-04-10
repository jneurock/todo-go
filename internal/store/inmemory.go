package store

import (
	"errors"
	"slices"
	"strconv"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/util"
)

type TodoInMemoryStore struct {
	isAvailable bool
	lastID      int64
	todos       util.Queue[domain.Todo]
}

func NewTodoInMemoryStore(available ...bool) *TodoInMemoryStore {
	var isAvailable bool

	if len(available) == 0 {
		isAvailable = true
	} else {
		isAvailable = available[0]
	}

	return &TodoInMemoryStore{isAvailable: isAvailable, lastID: -1}
}

func (s *TodoInMemoryStore) IsAvailable() bool {
	return s.isAvailable
}

func (s *TodoInMemoryStore) Create(description string) error {
	todo := domain.NewTodo(description)

	todo.ID = s.lastID + 1
	s.lastID = todo.ID

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
		todos = append(todos, todo.Value)
		todo = todo.Next
	}

	slices.SortFunc(todos, func(a, b *domain.Todo) int {
		switch {
		case !a.Complete && b.Complete:
			return -1
		case a.Complete && b.Complete && a.ID > b.ID:
			return -1
		default:
			return 1
		}
	})

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
