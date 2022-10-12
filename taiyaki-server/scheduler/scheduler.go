package scheduler

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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



func SelectWorker(m *Manager.Manager) models.Worker {
	db := m.DB
	workrCntrl := Controller.NewWorker(db)
	workers := workrCntrl.GetWorkers()
	selectedWorker := models.Worker{}
	threshold := 0.5
	for _,worker := range workers{
		respBody,err := reqWorker("stats",nil,worker.WorkerIP,worker.WorkerPort)
		if err!= nil{
			//handle error
		}
		fmt.Println(string(respBody))
		usage,err := strconv.Atoi(string(respBody))
		if err!= nil {
			//handle error
		}
		if usage<int(threshold){
			selectedWorker = worker
			break
		}
	}
	return selectedWorker
}

func reqWorker(endpoint string, reqBody io.Reader, workerIP string, workerPort string) (resBody []byte, err error) {
	url := "http://" + workerIP + ":" + workerPort + "/" + endpoint
	c := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: c}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, reqBody)
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
