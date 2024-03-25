package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Start() {
	static := http.Dir("internal/web/static")
	staticFs := http.FileServer(static)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))

	mux.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseGlob("internal/web/views/todo/new.html"))
		err := t.Execute(w, nil)

		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		} else {
			t := template.Must(template.ParseGlob("internal/web/views/index.html"))
			err := t.Execute(w, nil)

			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
			}
		}
	})

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
