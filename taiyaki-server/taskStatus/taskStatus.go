package taskStatus

import (
	"encoding/json"
	Client "taiyaki-server/client"
	Controller "taiyaki-server/controllers"
	Manager "taiyaki-server/manager"
	"taiyaki-server/models"
	task "taiyaki-server/task"

	"fmt"
	"time"
)

func UpdateTasks(m *Manager.Manager) {
	workrCntrl := Controller.NewWorker(m.DB)
	taskCntrl := Controller.NewTask(m.DB)

	for {
		workers := workrCntrl.GetWorkers()
		for _, worker := range workers {
			getTaskDetails(taskCntrl, workrCntrl, worker)
		}
		time.Sleep(30 * time.Second)
	}

}

func getTaskDetails(taskCntrl *Controller.TaskRepo, workCntrl *Controller.WorkerRepo, worker models.Worker) {
	resp, err := Client.ReqServer(worker.WorkerIP, worker.WorkerPort, "heartbeat/", nil)

	if err != nil {
		fmt.Println("Error while getting task status from a worker", err)
	}

	respBody := []task.Task{}

	err = json.Unmarshal(resp, &respBody)

	if err != nil {
		fmt.Println(" Unmarshal error in get task status", err)
	}

	for _, ta := range respBody {
		id := ta.ID

		t,valid := taskCntrl.GetTask(id.String())

		if !valid{
			panic(valid)
		}
		switch ta.State {
		case 0:
			t.State = "Pending"
		case 1:
			t.State = "Scheduled"
		case 2:
			t.State = "Completed"
		case 3:
			t.State = "Running"
		case 4:
			t.State = "Failed"
		default:
			t.State = "Running"
		}

		taskCntrl.UpdateTask(t)
	}
}
