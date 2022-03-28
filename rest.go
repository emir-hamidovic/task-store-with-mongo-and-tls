package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"rest/taskstore"
	"rest/taskstore/mongotaskstore"
)

type TaskServer struct {
	store taskstore.Taskstore
}

func (t *TaskServer) TaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Task handler\n")
	trimmedURL := strings.Trim(r.URL.Path, "/")
	urlPart := strings.Split(trimmedURL, "/")
	if trimmedURL == "task" {
		if r.Method == http.MethodGet {
			log.Printf("Get all tasks in progress\n")
			tasks, err := t.store.GetAllTasks()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			jsonRender(w, tasks)
			return
		} else if r.Method == http.MethodPost {
			var task taskstore.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Printf("Create task in progress\n")
			id, err := t.store.CreateTask(task.Text, task.Tags, task.Due)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "Task created successfully with id: %s\n", id)
			return
		} else if r.Method == http.MethodDelete {
			log.Printf("Delete all in progress\n")
			err := t.store.DeleteAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "Tasks deleted successfully\n")
			return
		}
	} else if len(urlPart) == 2 {
		id := urlPart[1]
		if r.Method == http.MethodGet {
			log.Printf("Get task by ID in progress\n")
			task, err := t.store.GetTaskById(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jsonRender(w, task)
			return
		} else if r.Method == http.MethodDelete {
			log.Printf("Delete in progress\n")
			err := t.store.DeleteTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "Task deleted successfully\n")
			return
		}
	}
}

func jsonRender(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func NewServer(ts taskstore.Taskstore) *TaskServer {
	return &TaskServer{store: ts}
}

// want to implement rate limiting and use gorilla mux here
// later, i want this to be able on https (eli bednersky website)
// maybe implement a Makefile and/or Dockerfile for this after all this?

// side note: research context package, tests overall, gopkg.in/check.v1
func main() {
	mux := http.NewServeMux()
	/*inmem := &inmemory.InMemory{Tasks: sync.Map{}, NextId: 0}
	router := NewServer(inmem)*/

	mng, err := mongotaskstore.NewMongoServer("", "", "")
	defer mng.CloseMongoServer()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	router := NewServer(mng)

	mux.HandleFunc("/task/", router.TaskHandler)

	log.Printf("Listening on 3333.\n")
	log.Fatal(http.ListenAndServe(":3333", mux))
}
