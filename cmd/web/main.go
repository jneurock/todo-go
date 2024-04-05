package main

import (
	"github.com/jneurock/todo-go/internal/store"
	"github.com/jneurock/todo-go/internal/web"
)

func main() {
	// TODO: Set up database connection
	// NOTE: Don't panic if database connection fails. Render 500.html.

	todoStore := store.NewTodoInMemoryStore()

	server := web.NewServer(todoStore)
	server.Start()
}
