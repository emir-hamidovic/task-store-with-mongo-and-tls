package inmemory

import (
	"errors"
	"fmt"
	"rest/taskstore"
	"strconv"
	"sync"
)

type InMemory struct {
	mtx    sync.Mutex
	Tasks  sync.Map
	NextId int
}

func (t *InMemory) CreateTask(text string, tags []string, due string) (string, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.NextId++
	task := taskstore.Task{Id: strconv.Itoa(t.NextId), Text: text, Due: due}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	t.Tasks.Store(t.NextId, task)
	return strconv.Itoa(t.NextId), nil
}

func (t *InMemory) GetTaskById(id string) (taskstore.Task, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return taskstore.Task{}, err
	}

	if x, found := t.Tasks.Load(idInt); found {
		v, ok := x.(taskstore.Task)
		if ok {
			return v, nil
		}

		return taskstore.Task{}, errors.New("returned value is not of type taskstore.Task")
	}

	return taskstore.Task{}, fmt.Errorf("couldn't find task with specified ID %s", id)
}

func (t *InMemory) DeleteTask(id string) error {
	if len(id) == 0 {
		return errors.New("no ID entered")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	t.Tasks.Delete(idInt)
	return nil
}

func (t *InMemory) DeleteAll() error {
	t.Tasks.Range(func(key interface{}, value interface{}) bool {
		t.Tasks.Delete(key)
		return true
	})

	return nil
}

func (t *InMemory) GetAllTasks() ([]taskstore.Task, error) {
	tasks := make([]taskstore.Task, 0)
	t.Tasks.Range(func(key interface{}, value interface{}) bool {
		v, ok := value.(taskstore.Task)
		if ok {
			tasks = append(tasks, v)
			return true
		}

		return false
	})

	return tasks, nil
}
