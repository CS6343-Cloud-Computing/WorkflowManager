package worker

import (
	"fmt"
	"github.com/CS6343-Cloud-Computing/WorkflowManager/Task"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
)

type Worker struct{
	Queue queue.Queue
	Db map[uuid.UUID]Task
	TaskCount int
}

func(w *Worker) CollectionStats(){
	fmt.Println("Collect stats")
}

func(w *Worker) RunTask(){
	fmt.Println("Start or stop task")
}

func(w *Worker) StartTask(){
	fmt.Println("Start task")
}

func(w * Worker) StopTask(){
	fmt.Println("Stop task")
}