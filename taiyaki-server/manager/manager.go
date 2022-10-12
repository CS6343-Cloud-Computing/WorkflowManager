package manager

import (
	"taiyaki-server/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manager struct {
	DB         *gorm.DB
	Pending    queue.Queue
	EventDb    map[uuid.UUID]*task.TaskEvent
	LastWorker int
}


func (m *Manager) AddTask(te task.TaskEvent) {
	m.Pending.Enqueue(te)
}

func NewManager() *Manager {
	eventDb := make(map[uuid.UUID]*task.TaskEvent)

	return &Manager{
		Pending: *queue.New(),
		EventDb: eventDb,
	}
}
