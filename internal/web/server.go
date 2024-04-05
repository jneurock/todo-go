package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/store"
)

var t *template.Template

type UITodo struct {
	Error error
	Todo  *domain.Todo
}

func newUITodo(todo *domain.Todo, err error) *UITodo {
	return &UITodo{Todo: todo, Error: err}
}

func newUITodoSlice(todos []*domain.Todo) (uiTodos []*UITodo) {
	for _, t := range todos {
		uiTodos = append(uiTodos, newUITodo(t, nil))
	}

	return
}

type Server struct {
	store store.TodoStore
}

func NewServer(store store.TodoStore) *Server {
	return &Server{store}
}

func (s *Server) Start() {
	var err error
	t, err = template.ParseGlob("internal/web/views/*.html")

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	static := http.Dir("internal/web/static")
	staticFs := http.FileServer(static)
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))

	mux.HandleFunc("DELETE /todo/{id}", createHandler(s.handleDeleteTodo))
	mux.HandleFunc("PUT /todo/{id}", createHandler(s.handleUpdateTodo))
	mux.HandleFunc("POST /todo", createHandler(s.handleNewTodo))

	mux.HandleFunc("/", createHandler(s.handleIndex))

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type Handler func(http.ResponseWriter, *http.Request) error

func createHandler(handler Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		if err != nil {
			err = t.ExecuteTemplate(w, "500.html", nil)

			if err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		if err := t.ExecuteTemplate(w, "404.html", nil); err != nil {
			return err
		}

		return nil
	}

	todos, err := s.store.FindAll()

	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, "index.html", &struct {
		Error error
		Todos []*UITodo
	}{
		Error: nil,
		Todos: newUITodoSlice(todos),
	})
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	errDelete := s.store.Delete(id)
	todos, err := s.store.FindAll()

	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, "todos", &struct {
		Error error
		Todos []*UITodo
	}{
		Error: errDelete,
		Todos: newUITodoSlice(todos),
	})
}

func (s *Server) handleNewTodo(w http.ResponseWriter, r *http.Request) error {
	description, err := domain.NewDescription(r.FormValue("description"))

	if err == nil {
		err = s.store.Create(description)
	}

	todos, errFindAll := s.store.FindAll()

	if errFindAll != nil {
		return errFindAll
	}

	return t.ExecuteTemplate(w, "todos", &struct {
		Error error
		Todos []*UITodo
	}{
		Error: err,
		Todos: newUITodoSlice(todos),
	})
}

func (s *Server) handleUpdateTodo(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	description, err := domain.NewDescription(r.FormValue("description"))
	complete := r.FormValue("complete") != ""

	if err == nil {
		err = s.store.Update(id, description, complete)
	}

	todo, errFind := s.store.Find(id)

	if errFind != nil {
		return errFind
	}

	return t.ExecuteTemplate(w, "todo", newUITodo(todo, err))
}
