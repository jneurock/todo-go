package web

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	templates *Templates
}

func NewServer(templates *Templates) *Server {
	return &Server{templates: templates}
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
	} else {
		t := s.templates.Load("index.html")
		err := t.Execute(w, nil)

		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}

func (s *Server) handleNewTodo(w http.ResponseWriter, r *http.Request) {
	t := s.templates.Load("todo/new.html")
	err := t.Execute(w, nil)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
}
