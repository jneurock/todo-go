package main

import "github.com/jneurock/todo-go/internal/web"

func main() {
	// TODO: Set up database connection
	// TODO: Set up app

	templates := web.NewTemplates("internal/web/views")
	server := web.NewServer(templates)
	server.Start()
}
