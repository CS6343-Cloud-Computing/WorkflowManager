package main

import (
	"github.com/CS6343-Cloud-Computing/WorkflowManager/Node"
	"github.com/CS6343-Cloud-Computing/WorkflowManager/Task"
	"github.com/CS6343-Cloud-Computing/WorkflowManager/Manager"
	"github.com/CS6343-Cloud-Computing/WorkflowManager/Worker"
	"fmt"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main(){
	t:= task.Task{
		ID:	uuid.New(),
		Name:	"Task-1",
		State:	task.Pending,
		Image:	"Image-1",
		Memory:	1024,
		Disk:	1,
	}

	te:= task.TaskEvent{
		ID:	uuid.New(),
		State:	task.Pending,
		Timestamp:	time.Now(),
		Task:	t,
	}

	fmt.Println("task: %v",t)
	fmt.Println("task event: %v", te)
	w:= worker.Worker{
		Queue: *queue.New(),
		Db: make(map[uuid.UUID]task.Task),
	}

	fmt.Println("worker: %v",w)
	w.CollectionStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m:=	manager.Manager{
			Pending:	*queue.New(),
			TaskDb:	make(map[string][]task.Task),
			EventDb:	make(map[string][]task.TaskEvent),
			Workers:	[]string{w.Name},
	}

	fmt.Println("manager: %v", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n:= node.Node{
		Name: "Node-1",
		Ip:	"192.168.1.1",
		Memory:	1024,
		Disk:	25,
	}

	fmt.Println("node: %v",n)

}