package domain

import (
	"errors"
	"strings"
)

var ErrInvalidDescription = errors.New("invalid description")

type Todo struct {
	Complete    bool
	ID          int64
	Description string
}

func NewTodo(description string, complete bool) (*Todo, error) {
	if strings.TrimSpace(description) == "" {
		return nil, ErrInvalidDescription
	}

	return &Todo{
		Complete:    complete,
		Description: description,
	}, nil
}
