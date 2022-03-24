package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"taskstore"
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
			tasks := t.store.GetAllTasks()
			jsonRender(w, tasks)
			return
		} else if r.Method == http.MethodPost {
			var task taskstore.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			id := t.store.CreateTask(task.Text, task.Tags, task.Due)
			fmt.Fprintf(w, "Task created successfully with id: %d\n", id)
			return
		} else if r.Method == http.MethodDelete {
			t.store.DeleteAll()
			fmt.Fprintf(w, "Tasks deleted successfully\n")
			return
		}
	} else if len(urlPart) == 2 {
		id, _ := strconv.Atoi(urlPart[1])
		if r.Method == http.MethodGet {
			task, err := t.store.GetTaskById(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jsonRender(w, task)
			return
		} else if r.Method == http.MethodDelete {
			t.store.DeleteTask(id)
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

func main() {
	mux := http.NewServeMux()
	inmem := &taskstore.InMemory{Tasks: sync.Map{}, NextId: 0}
	router := NewServer(inmem)
	mux.HandleFunc("/task/", router.TaskHandler)

	log.Printf("Listening on 3333\n")
	log.Fatal(http.ListenAndServe(":3333", mux))
}
