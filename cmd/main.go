package main

import (
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var idSeq int = 1

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

		task := Task{
			ID:          idSeq,
			Description: description,
			Priority:    priority,
			Status:      status,
		}

		taskList = append(taskList, task)

        idSeq++
		taskAddTmpl.Execute(w, task)
	})

	r.Delete("/task/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		domID := chi.URLParam(r, "id")
		taskIDStr := strings.TrimPrefix(domID, "task-")
		id, err := strconv.Atoi(taskIDStr)
		if err != nil {
			panic(err)
		}

		idx := slices.IndexFunc(taskList, func(t Task) bool {
			return t.ID == id
		})

		if idx == -1 {
			return
		}

		taskList = slices.Delete(taskList, idx, idx+1)
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
