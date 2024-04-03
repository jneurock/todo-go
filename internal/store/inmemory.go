package store

import (
	"strconv"

	domain "github.com/jneurock/todo-go/internal/domain"
)

type TodoQueueItem struct {
	value *domain.Todo
	next  *TodoQueueItem
}

type TodoQueue struct {
	head *TodoQueueItem
}

func (tq *TodoQueue) Add(todo *domain.Todo) {
	queueItem := &TodoQueueItem{value: todo}

	if tq.head == nil {
		tq.head = queueItem
		return
	}

	queueItem.next = tq.head
	tq.head = queueItem
}

func (tq *TodoQueue) Remove(id int64) {
	var prevQueueItem *TodoQueueItem
	queueItem := tq.head

	for queueItem != nil {
		if queueItem.value.ID == id {
			if queueItem == tq.head {
				tq.head = queueItem.next
			} else {
				prevQueueItem.next = queueItem.next
			}

			break
		}

		prevQueueItem = queueItem
		queueItem = queueItem.next
	}
}

type InMemoryTodoStore struct {
	todos TodoQueue
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

	s.todos.Remove(int64(intId))

	return nil
}

func (s *InMemoryTodoStore) FindAll() ([]*domain.Todo, error) {
	var todos []*domain.Todo
	todo := s.todos.head

	for todo != nil {
		todos = append([]*domain.Todo{todo.value}, todos...)
		todo = todo.next
	}

	return todos, nil
}

func (s *InMemoryTodoStore) Update(id string, todo *domain.Todo) error {
	return nil
}
