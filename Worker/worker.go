package worker

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CS6343-Cloud-Computing/WorkflowManager/Task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct{
	Name	string
	Queue queue.Queue
	Db map[uuid.UUID]task.Task
	TaskCount int
}

func(w *Worker) CollectionStats(){
	fmt.Println("Collect stats")
}

func (w *Worker) RunTask() task.DockerResult {
	fmt.Println("reached 1 ")
	t := w.Queue.Dequeue()
	if t == nil {
			log.Println("No tasks in the queue")
			return task.DockerResult{Error: nil}
	}

	fmt.Println("reached 2 ")

	taskQueued := t.(task.Task)

	fmt.Println("reached 3 ")

	taskPersisted,isPresent := w.Db[taskQueued.ID]

	fmt.Println("reached 4 ")

	if !isPresent  {
			taskPersisted = taskQueued
			w.Db[taskQueued.ID] = taskQueued
	}

	fmt.Println("reached 5 ")

	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
			switch taskQueued.State {
			case task.Scheduled:
					fmt.Println("reached 7 ")

					result = w.StartTask(taskQueued)
			case task.Completed:
					fmt.Println("reached 8 ")

					result = w.StopTask(taskQueued)
			default:
					result.Error = errors.New("we should not get here")
			}
	} else {
			err := fmt.Errorf("invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
			result.Error = err
	}

	fmt.Println("reached 6 ")

	return result
}

func(w *Worker) StartTask(t task.Task) task.DockerResult{
	fmt.Println("reached 9 ")

	config := task.NewConfig(&t)

	fmt.Println("reached 10 ")

	d := task.NewDocker(config)

	fmt.Println("reached 11 ")

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = t
		return result
	}

	d.ContainerId = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = t

	return result
}

func(w * Worker) StopTask(t task.Task) task.DockerResult{
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	d.ContainerId = t.ContainerId
	result := d.Stop()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = t
	log.Printf("Stopped and removed the container %v for task %v", d.ContainerId, t)

	return result
}

func(w *Worker) AddTask(t task.Task){
	w.Queue.Enqueue(t)
}