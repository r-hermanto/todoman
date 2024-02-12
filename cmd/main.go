package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	taskList := []Task{}

	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles(
		"web/templates/index.html",
	))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	funcMap := template.FuncMap{
		"ToLower": func(v TaskStatus) string {
			return strings.ToLower(string(v))
		},
	}
	taskAddTmpl := template.Must(template.New("task.html").Funcs(funcMap).ParseFiles(
		"web/templates/partials/task.html",
	))

	r.Post("/task/add", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		description := r.FormValue("description")
		priority := TaskPriority(r.FormValue("priority"))
		status := TaskStatus(r.FormValue("status"))

		id := 1
		if len(taskList) > 0 {
			id = taskList[len(taskList)-1].ID + 1
		}

		task := Task{
			ID:          id,
			Description: description,
			Priority:    priority,
			Status:      status,
		}

        log.Println(task)
		taskList = append(taskList, task)
		taskAddTmpl.Execute(w, task)
	})

	log.Println("Start server at port 8000")
	http.ListenAndServe(":8000", r)
}

type Task struct {
	ID          int
	Description string
	Priority    TaskPriority
	Status      TaskStatus
}

type TaskPriority string

const (
	NO_PRIORITY TaskPriority = "NO_PRIORITY"
	LOW                      = "LOW"
	MEDIUM                   = "MEDIUM"
	HIGH                     = "HIGH"
	URGENT                   = "URGENT"
)

type TaskStatus string

const (
	BACKLOG TaskStatus = "BACKLOG"
	TODO               = "TODO"
	DOING              = "DOING"
	DONE               = "DONE"
)
