package main

import (
	"fmt"
	"sync"
	apiserver "taiyaki-server/apiserver"
	Manager "taiyaki-server/manager"
	models "taiyaki-server/models"
	mysql "taiyaki-server/mysql"
	"time"
)

func main() {
	db := mysql.InitDb()
	db.AutoMigrate(&models.Worker{})
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.Workflow{})
	m := Manager.NewManager()
	m.DB = db
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := apiserver.APIConfig{ServerIP: "192.168.1.92", ServerPort: "8080", WorkerJoinKey: "1234"}

	go api.Start(&wg, m)
	go func() {
		for {
			fmt.Println("Before Sendwork")
			apiserver.SendWork(m)
			fmt.Println("After Sendwork")
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
