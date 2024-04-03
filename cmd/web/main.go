package main

import (
	"github.com/jneurock/todo-go/internal/store"
	"github.com/jneurock/todo-go/internal/web"
)

func main() {
	// TODO: Set up database connection

	todoStore := store.NewInMemoryTodoStore()

	server := web.NewServer(todoStore)
	server.Start()
}
