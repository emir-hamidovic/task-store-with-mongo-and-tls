package taskstore

import (
	"errors"
	"log"
	"sync"
	"time"
)

type InMemory struct {
	mtx    sync.Mutex
	Tasks  sync.Map
	NextId int
}

func (t *InMemory) CreateTask(text string, tags []string, due time.Time) int {
	log.Printf("Create task in progress\n")
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.NextId++
	task := Task{Id: t.NextId, Text: text, Due: due}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	t.Tasks.Store(t.NextId, task)
	return t.NextId
}

func (t *InMemory) GetTaskById(id int) (Task, error) {
	log.Printf("Get task by ID in progress\n")
	if x, found := t.Tasks.Load(id); found {
		v, ok := x.(Task)
		if ok {
			return v, nil
		}

		return Task{}, errors.New("returned value is not of type Task")
	}

	return Task{}, errors.New("couldn't find task with specified ID")
}

func (t *InMemory) DeleteTask(id int) {
	log.Printf("Delete in progress\n")
	t.Tasks.Delete(id)
}

func (t *InMemory) DeleteAll() {
	log.Printf("Delete all in progress\n")
	t.Tasks.Range(func(key interface{}, value interface{}) bool {
		t.Tasks.Delete(key)
		return true
	})
}

func (t *InMemory) GetAllTasks() []Task {
	log.Printf("Get all tasks in progress\n")

	tasks := make([]Task, 0)
	t.Tasks.Range(func(key interface{}, value interface{}) bool {
		v, ok := value.(Task)
		if ok {
			tasks = append(tasks, v)
			return true
		}

		return false
	})

	return tasks
}
