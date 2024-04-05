package domain

import (
	"errors"
	"strings"
)

var ErrInvalidDescription = errors.New("description cannot be empty")

type Todo struct {
	Complete    bool
	ID          int64
	Description string
}

func NewDescription(description string) (string, error) {
	if strings.TrimSpace(description) == "" {
		return "", ErrInvalidDescription
	}

	return description, nil
}

func NewTodo(description string) *Todo {
	return &Todo{Description: description}
}
