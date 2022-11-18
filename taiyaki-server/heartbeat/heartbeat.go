package heartbeat

import (
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

			fmt.Println("Worker heart:", worker)

			RequestWorker(workrCntrl, &worker)
		}

		time.Sleep(20 * time.Second)
	}

}

func RequestWorker(workCntrl *Controller.WorkerRepo, worker *models.Worker) {
	_, err := Client.ReqServer(worker.WorkerIP, worker.WorkerPort, "heartbeat/", nil)

	if err != nil {
		// log.Fatal("Error in getting heartbeat of worker", err)

		worker.Status = "Inactive"
		_, err1 := workCntrl.UpdateWorker(*worker)
		if err1 != nil {
			log.Fatal("Error while updating the Status to inactive", err1)
		}
	}

	// respBody := Client.Resp{}
	// fmt.Println(resp)
	// err = json.Unmarshal(resp, &respBody)
	// if err != nil {
	// 	log.Fatal("Unmarshal error while updating a worker to inactive", err)
	// 	worker.Status = "Inactive"
	// 	_, err1 := workCntrl.UpdateWorker(*worker)
	// 	if err1 != nil {
	// 		log.Fatal("Error while updating the Status to inactive", err1)
	// 	}
	// }
}
