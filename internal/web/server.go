package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/service"
)

var t *template.Template

type Server struct {
	todoService *service.TodoService
}

func NewServer(todoService *service.TodoService) *Server {
	return &Server{
		todoService: todoService,
	}
}

func (s *Server) Start() {
	static := http.Dir("internal/web/static")
	staticFs := http.FileServer(static)

	var err error
	t, err = template.ParseGlob("internal/web/views/**/*.html")

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.HandleFunc("DELETE /todo/{id}", s.handleDeleteTodo)
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

	todos, err := s.todoService.FindAll()

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = t.ExecuteTemplate(w, "index.html", &struct{ Todos []*domain.Todo }{
		Todos: todos,
	})

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.todoService.Delete(id)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func (s *Server) handleNewTodo(w http.ResponseWriter, r *http.Request) {
	_, err := s.todoService.Create(r.FormValue("description"))

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	todos, err := s.todoService.FindAll()

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = t.ExecuteTemplate(w, "new.html", &struct{ Todos []*domain.Todo }{
		Todos: todos,
	})

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
