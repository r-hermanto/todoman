package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	tmpl := template.Must(template.ParseFiles(
		"web/templates/index.html",
	))

    fs := http.FileServer(http.Dir("web/static"))

    r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	log.Println("Start server at port 8000")
	http.ListenAndServe(":8000", r)
}
