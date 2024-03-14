package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseGlob("internal/web/views/index.html"))
		err := t.Execute(w, nil)

		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	})

	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
