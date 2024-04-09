package main

import (
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

func setUpStore(useLocalDB bool) store.TodoStore {
	if useLocalDB {
		return store.NewTodoInMemoryStore()
	}

	user := os.Getenv("DB_USER")
	pw := os.Getenv("DB_PW")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL_MODE")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", user, name, pw, sslmode)

	return store.NewTodoPostsgresStore(connStr)
}

func main() {
	flag.Bool("localdb", false, "Use a temporary local database for development")

	flag.Parse()

	localDB := flagExists("localdb")
	todoStore := setUpStore(localDB)

	templatePath := "internal/web/views"

	server, err := web.NewServer(todoStore, templatePath)

	if err != nil {
		panic(err)
	}

	server.Start()
}
