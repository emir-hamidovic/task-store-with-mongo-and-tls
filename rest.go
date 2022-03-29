package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rest/taskstore"
	"rest/taskstore/mongotaskstore"

	"github.com/gorilla/handlers"
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

	id, err := t.store.CreateTask(task.Text, task.Tags, task.Due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Task created successfully with id: %s\n", id)
}

func (t *TaskServer) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.store.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonRender(w, tasks)
}

func (t *TaskServer) HandleDeleteTasks(w http.ResponseWriter, r *http.Request) {
	err := t.store.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Tasks deleted successfully\n")
}

func (t *TaskServer) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	task, err := t.store.GetTaskById(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonRender(w, task)
}

func (t *TaskServer) HandleDeleteTaskByID(w http.ResponseWriter, r *http.Request) {
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

func limitNumClients(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	limiter := make(chan struct{}, maxClients)
	return func(w http.ResponseWriter, r *http.Request) {
		limiter <- struct{}{}
		defer func() { <-limiter }()
		f(w, r)
	}
}

// later, i want this to be able on https (eli bednersky website)
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
	r.StrictSlash(true)
	r.HandleFunc("/task", limitNumClients(router.HandlePostTask, 5)).Methods("POST")
	r.HandleFunc("/task", limitNumClients(router.HandleGetTasks, 5)).Methods("GET")
	r.HandleFunc("/task", limitNumClients(router.HandleDeleteTasks, 5)).Methods("DELETE")
	r.HandleFunc("/task/{id:[0-9a-zA-z-]+}", limitNumClients(router.HandleGetTaskByID, 5)).Methods("GET")
	r.HandleFunc("/task/{id:[0-9a-zA-z-]+}", limitNumClients(router.HandleDeleteTaskByID, 5)).Methods("DELETE")
	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	r.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})

	port := os.Getenv("SERVERPORT")
	if port == "" {
		port = "3333"
	}

	srv := &http.Server{
		Handler:      http.TimeoutHandler(r, 1*time.Second, "Timeout!\n"),
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on %s.\n", port)
	log.Fatal(srv.ListenAndServe())
}
