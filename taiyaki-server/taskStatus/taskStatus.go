package taskStatus

import (
	"encoding/json"
	"log"
	Controller "taiyaki-server/controllers"
	Manager "taiyaki-server/manager"
	"taiyaki-server/models"
	"taiyaki-server/scheduler"
	task "taiyaki-server/task"

	"fmt"
	"time"
)

func UpdateTasks(m *Manager.Manager) {
	workrCntrl := Controller.NewWorker(m.DB)
	taskCntrl := Controller.NewTask(m.DB)

	for {
		log.Println("updating the states from worker")
		workers := workrCntrl.GetActiveWorkers()
		for _, worker := range workers {
			getTaskDetails(taskCntrl, workrCntrl, worker)
		}
		time.Sleep(6 * time.Second)
	}

}

func getTaskDetails(taskCntrl *Controller.TaskRepo, workCntrl *Controller.WorkerRepo, worker models.Worker) {
	numContainers := worker.NumContainers
	if numContainers > 0 {
		resp,err := scheduler.ReqWorker("tasks","GET",nil,worker.WorkerIP,worker.WorkerPort)
		//log.Println(string(resp))
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

			t, valid := taskCntrl.GetTask(id.String())

			if !valid {
				log.Println("valid in getTaskDetails", valid)
				//panic(valid)
			}
			if(t.State == "Failed" || t.State == "Completed" || t.State == "KillBitReceived"){
				continue;
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
			//log.Println("Updating state for ", t)
			taskCntrl.UpdateTask(t)
		}
	}
}
