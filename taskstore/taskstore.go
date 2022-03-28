package taskstore

type Taskstore interface {
	CreateTask(text string, tags []string, due string) (string, error)
	GetTaskById(id string) (Task, error)
	DeleteTask(id string) error
	DeleteAll() error
	GetAllTasks() ([]Task, error)
}

type Task struct {
	Id   string   `json:"id" bson:"_id,omitempty"`
	Text string   `json:"text" bson:"text"`
	Tags []string `json:"tags" bson:"tags"`
	Due  string   `json:"due" bson:"due"`
}
