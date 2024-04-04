package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/store"
)

var t *template.Template

type Server struct {
	store store.TodoStore
}

func NewServer(store store.TodoStore) *Server {
	return &Server{store}
}

func (s *Server) Start() {
	static := http.Dir("internal/web/static")
	staticFs := http.FileServer(static)

	var err error
	t, err = template.ParseGlob("internal/web/views/*.html")

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.HandleFunc("DELETE /todo/{id}", s.handleDeleteTodo)
	mux.HandleFunc("PUT /todo/{id}", s.handleUpdateTodo)
	mux.HandleFunc("POST /todo", s.handleNewTodo)
	mux.HandleFunc("/", s.handleIndex)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	todos, err := s.store.FindAll()

	err = t.ExecuteTemplate(w, "index.html", &struct {
		Error error
		Todos []*domain.Todo
	}{
		Error: err,
		Todos: todos,
	})

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.store.Delete(id)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

	todos, err := s.store.FindAll()

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = t.ExecuteTemplate(w, "todos", &struct {
		Error error
		Todos []*domain.Todo
	}{
		Error: err,
		Todos: todos,
	})

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

// TODO: Improve the error handling here
func (s *Server) handleNewTodo(w http.ResponseWriter, r *http.Request) {
	todo, err := domain.NewTodo(r.FormValue("description"), false)

	if err != nil {
		todos, errFindTodos := s.store.FindAll()

		if errFindTodos != nil {
			renderTodos(w, todos, errFindTodos)
			return
		}

		renderTodos(w, todos, err)
		return
	}

	err = s.store.Create(todo)
	todos, errFindTodos := s.store.FindAll()

	if err != nil {
		if errFindTodos != nil {
			renderTodos(w, todos, errFindTodos)
			return
		}

		renderTodos(w, todos, err)
		return
	}

	renderTodos(w, todos, nil)
}

func (s *Server) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	complete := false
	description := r.FormValue("description")

	if len(r.FormValue("complete")) != 0 {
		complete = true
	}

	todo, err := domain.NewTodo(description, complete)

	// TODO: Handle todo validation better
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	todo.ID = int64(id)

	err = s.store.Update(todo)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = t.ExecuteTemplate(w, "todo", todo)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func renderTodos(w http.ResponseWriter, todos []*domain.Todo, err error) {
	err = t.ExecuteTemplate(w, "todos", &struct {
		Error error
		Todos []*domain.Todo
	}{
		Error: err,
		Todos: todos,
	})

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
