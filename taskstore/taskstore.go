package taskstore

import (
	"time"
)

type Taskstore interface {
	CreateTask(text string, tags []string, due time.Time) int
	GetTaskById(id int) (Task, error)
	DeleteTask(id int)
	DeleteAll()
	GetAllTasks() []Task
}

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}
