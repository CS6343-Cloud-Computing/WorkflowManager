package main

import (
	"fmt"
	"os"
	"time"

	manager "github.com/CS6343-Cloud-Computing/WorkflowManager/Manager"
	node "github.com/CS6343-Cloud-Computing/WorkflowManager/Node"
	task "github.com/CS6343-Cloud-Computing/WorkflowManager/Task"
	worker "github.com/CS6343-Cloud-Computing/WorkflowManager/Worker"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "Task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}

	fmt.Println("task: %v", t)
	fmt.Println("task event: %v", te)
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]task.Task),
	}

	fmt.Println("worker: %v", w)
	w.CollectionStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]task.Task),
		EventDb: make(map[string][]task.TaskEvent),
		Workers: []string{w.Name},
	}

	fmt.Println("manager: %v", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Memory: 1024,
		Disk:   25,
	}

	fmt.Println("node: %v", n)

	fmt.Printf("create a test container\n")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Print(createResult.Error)
		os.Exit(1)
	}
	time.Sleep(time.Second * 5)
	fmt.Printf("stopping container %s\n", createResult.ContainerId)
	_ = stopContainer(dockerTask)
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c}
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}
	fmt.Printf(
		"Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker) *task.DockerResult {
	result := d.Stop()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}
	fmt.Printf(
		"Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}
