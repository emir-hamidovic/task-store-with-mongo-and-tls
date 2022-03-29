package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rest/taskstore"
	"rest/taskstore/mongotaskstore"

	"github.com/gorilla/mux"
)

type TaskServer struct {
	store taskstore.Taskstore
}

func (t *TaskServer) HandlePostTask(w http.ResponseWriter, r *http.Request) {
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
}
func (t *TaskServer) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get all tasks in progress\n")
	tasks, err := t.store.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonRender(w, tasks)
}
func (t *TaskServer) HandleDeleteTasks(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delete all in progress\n")
	err := t.store.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Tasks deleted successfully\n")
}
func (t *TaskServer) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get task by ID in progress\n")
	vars := mux.Vars(r)
	task, err := t.store.GetTaskById(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonRender(w, task)
}
func (t *TaskServer) HandleDeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Delete in progress\n")
	vars := mux.Vars(r)
	err := t.store.DeleteTask(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Task deleted successfully\n")
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

// rate limiting
// middleware
// later, i want this to be able on https (eli bednersky website)
// port needs to be configurable
// maybe implement a Makefile and/or Dockerfile for this after all this?

// side note: research context package, tests overall, gopkg.in/check.v1
func main() {
	mng, err := mongotaskstore.NewMongoServer("", "", "")
	defer mng.CloseMongoServer()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	router := NewServer(mng)

	r := mux.NewRouter()
	r.HandleFunc("/task", router.HandlePostTask).Methods("POST")
	r.HandleFunc("/task", router.HandleGetTasks).Methods("GET")
	r.HandleFunc("/task", router.HandleDeleteTasks).Methods("DELETE")
	r.HandleFunc("/task/{id:[0-9a-zA-z-]+}", router.HandleGetTaskByID).Methods("GET")
	r.HandleFunc("/task/{id:[0-9a-zA-z-]+}", router.HandleDeleteTaskByID).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:3333",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on 3333.\n")
	log.Fatal(srv.ListenAndServe())
}
