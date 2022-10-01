package manager

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
	"github.com/CS6343-Cloud-Computing/WorkflowManager/Task"
)

type Manager struct{
	Pending	queue.Queue
	TaskDb	map[string][]task.Task
	EventDb map[string][]task.TaskEvent
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