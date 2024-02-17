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

	funcMap := template.FuncMap{
		"ToLower": func(v TaskStatus) string {
			return strings.ToLower(string(v))
		},
	}

	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles(
		"web/templates/index.html",
		"web/templates/partials/task.html",
	))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tasksByStatus := map[TaskStatus][]Task{}

		for _, task := range taskList {
			tasksByStatus[task.Status] = append(tasksByStatus[task.Status], task)
		}

		templateData := struct {
			BacklogTasks []Task
			TodoTasks    []Task
			DoingTasks   []Task
			DoneTasks    []Task
		}{
			BacklogTasks: tasksByStatus[BACKLOG],
			TodoTasks:    tasksByStatus[TODO],
			DoingTasks:   tasksByStatus[DOING],
			DoneTasks:    tasksByStatus[DONE],
		}

		tmpl.Execute(w, templateData)
	})

	taskAddTmpl := template.Must(template.New("task_add.html").Funcs(funcMap).ParseFiles(
		"web/templates/partials/task_add.html",
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

	r.Post("/task/update", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			panic(err)
		}

		idx := slices.IndexFunc(taskList, func(t Task) bool {
			return t.ID == id
		})

		if idx == -1 {
			return
		}

		task := taskList[idx]

		description := r.FormValue("description")
		priority := TaskPriority(r.FormValue("priority"))

		task.Description = description
		task.Priority = priority
		taskList[idx] = task

		taskAddTmpl.Execute(w, task)
	})

	r.Post("/task/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		taskIDStr := chi.URLParam(r, "id")
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

		err = r.ParseForm()
		if err != nil {
			panic(err)
		}

		status := TaskStatus(r.FormValue("status"))

		task := taskList[idx]
		task.Status = status
		taskList[idx] = task

		taskAddTmpl.Execute(w, task)
	})

	r.Delete("/task/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		taskIDStr := chi.URLParam(r, "id")
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
