package main

import (
	"github.com/jneurock/todo-go/internal/service"
	"github.com/jneurock/todo-go/internal/store"
	"github.com/jneurock/todo-go/internal/web"
)

func main() {
	// TODO: Set up database connection
	// TODO: Set up app

	todoStore := store.NewInMemoryTodoStore()
	todoService := service.NewTodoService(todoStore)

	server := web.NewServer(todoService)
	server.Start()
}
