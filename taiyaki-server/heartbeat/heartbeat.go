package heartbeat

import (
	"encoding/json"
	"fmt"
	"log"
	Client "taiyaki-server/client"
	Controller "taiyaki-server/controllers"
	Manager "taiyaki-server/manager"
	"taiyaki-server/models"
	"time"
)

func GetHeartBeat(m *Manager.Manager) {
	workrCntrl := Controller.NewWorker(m.DB)

	for {
		workers := workrCntrl.GetActiveWorkers()

		for _, worker := range workers {

			fmt.Println("Worker:", worker)

			go RequestWorker(workrCntrl, &worker)
		}

		time.Sleep(10000)
	}
	

}

func RequestWorker(workCntrl *Controller.WorkerRepo, worker *models.Worker) {
	resp, err := Client.ReqServer(worker.WorkerIP, worker.WorkerPort, "/heartbeat/", nil)

	if err != nil {
		log.Fatal("Error in getting heartbeat of worker", err)
	}

	respBody := Client.Resp{}
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		log.Fatal("Unmarshal error while updating a worker to inactive", err)
		worker.Status = "Inactive"
		_, err := workCntrl.UpdateWorker(*worker)
		if err != nil {
			log.Fatal("Error while updating the Status to inactive", err)
		}
	}
}
