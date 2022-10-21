package main

import (
	//"fmt"
	"log"
	"os"
	"sync"
	apiserver "taiyaki-server/apiserver"
	Manager "taiyaki-server/manager"
	models "taiyaki-server/models"
	mysql "taiyaki-server/mysql"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverIP := os.Getenv("SERVER_IP")
	serverPort := os.Getenv("SERVER_PORT")
	workerJoinKey := os.Getenv("WORKER_JOIN_KEY")
	db := mysql.InitDb()
	db.AutoMigrate(&models.Worker{})
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.Workflow{})
	m := Manager.NewManager()
	m.DB = db
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := apiserver.APIConfig{ServerIP: serverIP, ServerPort: serverPort, WorkerJoinKey: workerJoinKey}

	go api.Start(&wg, m)
	go func() {
		for {
			//fmt.Println("Before Sendwork")
			apiserver.SendWork(m)
			//fmt.Println("After Sendwork")
			time.Sleep(15 * time.Second)
		}
	}()

	wg.Wait()

	//create a new manager

	// for {
	// 	for _, t := range m.TaskDb {
	// 		fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
	// 		time.Sleep(15 * time.Second)
	// 	}
	// }
}
