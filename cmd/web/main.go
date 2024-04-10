package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/jneurock/todo-go/internal/store"
	"github.com/jneurock/todo-go/internal/web"
)

func flagExists(flagName string) (found bool) {
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			found = true
		}
	})

	return
}

func setUpStore() store.TodoStore {
	if flagExists("localdb") {
		return store.NewTodoInMemoryStore()
	}

	connStr := "user=todouser dbname=todo password=todopassword sslmode=disable"

	if host := os.Getenv("TODO_DB_HOST"); host != "" {
		connStr = fmt.Sprintf("host=%s %s", host, connStr)
	}

	db, err := sql.Open("postgres", connStr)

	return store.NewTodoPostsgresStore(db, err == nil)
}

func main() {
	flag.Bool("localdb", false, "Use a temporary local database for development")

	flag.Parse()

	todoStore := setUpStore()
	templatePath := "internal/web/views"
	server, err := web.NewServer(todoStore, templatePath)

	if err != nil {
		panic(err)
	}

	server.Start()
}
