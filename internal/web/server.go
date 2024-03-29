package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jneurock/todo-go/internal/domain"
	"github.com/jneurock/todo-go/internal/service"
)

type Server struct {
	templates   *Templates
	todoService *service.TodoService
}

func NewServer(templates *Templates, todoService *service.TodoService) *Server {
	return &Server{
		templates:   templates,
		todoService: todoService,
	}
}

func (s *Server) Start() {
	static := http.Dir("internal/web/static")
	staticFs := http.FileServer(static)

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
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

	t := s.templates.Load("index.html")
	err = t.Execute(w, &struct{ Todos []*domain.Todo }{Todos: todos})

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

	t := s.templates.Load("todo/new.html")
	err = t.Execute(w, &struct{ Todos []*domain.Todo }{Todos: todos})

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
