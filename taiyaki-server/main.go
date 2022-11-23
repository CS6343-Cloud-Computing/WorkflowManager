package main

import (
	//"fmt"
	"fmt"
	"log"
	"os"
	"sync"
	apiserver "taiyaki-server/apiserver"
	heartbeat "taiyaki-server/heartbeat"
	Manager "taiyaki-server/manager"
	models "taiyaki-server/models"
	mysql "taiyaki-server/mysql"
	taskStatus "taiyaki-server/taskStatus"

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
	db.Exec("DROP TABLE if exists workers")
	db.Exec("DROP TABLE if exists tasks")
	db.Exec("DROP TABLE if exists workflows")

	db.AutoMigrate(&models.Worker{})
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.Workflow{})
	m := Manager.NewManager()
	m.DB = db
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := apiserver.APIConfig{ServerIP: serverIP, ServerPort: serverPort, WorkerJoinKey: workerJoinKey}
	fmt.Print("hello")
	go api.Start(&wg, m)
	go func() {
		for {
			//fmt.Println("Before Sendwork")
			apiserver.SendWork(m)
			//fmt.Println("After Sendwork")
			time.Sleep(15 * time.Second)
		}
	}()
	go heartbeat.GetHeartBeat(m)

	go taskStatus.UpdateTasks(m)

	// go Manager.KillTask(m)
	wg.Wait()

	//create a new manager

	// for {
	// 	for _, t := range m.TaskDb {
	// 		fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
	// 		time.Sleep(15 * time.Second)
	// 	}
	// }
}
