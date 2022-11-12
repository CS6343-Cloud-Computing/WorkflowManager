package scheduler

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"strconv"
	Controller "taiyaki-server/controllers"
	Manager "taiyaki-server/manager"
	"taiyaki-server/models"

	"github.com/c9s/goprocinfo/linux"
)

type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
}

func CpuUsage(s Stats) float64 {

	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle


	if total == 0 {
			return 0.00
	}

	return (float64(total) - float64(idle)) / float64(total)
}


func MemAvailablePercent(s Stats) float64 {
	return float64(s.MemStats.MemAvailable) / float64(s.MemStats.MemTotal)
}

func SelectWorker(m *Manager.Manager) models.Worker {
	//db := m.DB
	//workrCntrl := Controller.NewWorker(db)
	//workers := workrCntrl.GetWorkers()
	workers := WorkerWithMinTasks(m)
	selectedWorker := models.Worker{}
	cpuThreshold := 0.90
	memThreshhold := 0.90
	workerFound := false
	for _,worker := range workers{
		resp,err := ReqWorker("stats","GET",nil,worker.WorkerIP,worker.WorkerPort)
		if err!= nil{
			//handle error
		}
		respBody := Stats{}
		err = json.Unmarshal(resp, &respBody)

		if err!= nil {
			//handle error
		}

		availMem := MemAvailablePercent(respBody)
		cpuUsage := CpuUsage(respBody)

		fmt.Println("------------------Available Mem : ", availMem)
		fmt.Println("------------------CpuUsage CPU : ", cpuUsage)

		if (cpuUsage<cpuThreshold && availMem<memThreshhold){
			fmt.Println("Got a useful worker", worker)
			selectedWorker = worker
			workerFound = true
			break
		}
	}
	if !workerFound {
	println("No userful worker")
	}
	return selectedWorker
}

//returns the list of workers with minimum tasks in asc order
func WorkerWithMinTasks(m *Manager.Manager) []models.Worker{
	db := m.DB
	workrCntrl := Controller.NewWorker(db)
	workers := workrCntrl.GetMinTaskWorkers()
	return workers
}

func ReqWorker(endpoint string,method string, reqBody io.Reader, workerIP string, workerPort string) (resBody []byte, err error) {
	url := "http://" + workerIP + ":" + workerPort + "/" + endpoint
	c := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: c}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "close")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(resp.Body)
	resBody, _ = ioutil.ReadAll(reader)
	resp.Body.Close()

	return resBody, nil
}
