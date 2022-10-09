package main

import (
	mysql "taiyaki-server/mysql"
	apiserver "taiyaki-server/apiserver"
	models "taiyaki-server/models"
	"sync"
)

func main() {
	db := mysql.InitDb()
	db.AutoMigrate(&models.Worker{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	api := apiserver.APIConfig{ServerIP: "192.168.1.92", ServerPort: "8080", WorkerJoinKey: "1234"}
	go api.Start(&wg, db)
	wg.Wait()
}
