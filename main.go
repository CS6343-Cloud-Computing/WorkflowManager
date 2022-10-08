package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	task "github.com/CS6343-Cloud-Computing/WorkflowManager/Task"
	worker "github.com/CS6343-Cloud-Computing/WorkflowManager/Worker"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	os.Setenv("WORKFLOWMANAGER","localhost")
	host := os.Getenv("WORKFLOWMANAGER")
	os.Setenv("WORKFLOWMANAGER_PORT","5555")
	port, _ := strconv.Atoi(os.Getenv("WORKFLOWMANAGER_PORT"))

	fmt.Println("Starting WORKFLOWMANAGER worker")

	w := worker.Worker{
			Queue: *queue.New(),
			Db:    make(map[uuid.UUID]task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
			if w.Queue.Len() != 0 {
					result := w.RunTask()
					if result.Error != nil {
							log.Printf("Error running task: %v\n", result.Error)
					}
			} else {
					log.Printf("No tasks to process currently.\n")
			}
			log.Println("Sleeping for 10 seconds.")
			time.Sleep(10 * time.Second)
	}

}