package manager

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
)

type Manager struct{
	Pending	queue.Queue
	TaskDb	map[string][]Task
	EventDb map[string][]TaskEvent
	Workers	[]string
	WorkerTaskMap	map[string][]uuid.UUID
	TaskWorkerMap	map[uuid.UUID]string
}

func (m *Manager) SelectWorker(){
	fmt.Println("Select Worker")
}

func (m *Manager) UpdateTasks(){
	fmt.Println("Update tasks")
}

func (m *Manager) SendWork(){
	fmt.Println("Send work to workers")
}