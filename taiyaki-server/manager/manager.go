package manager

import (
	"bufio"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	Controller "taiyaki-server/controllers"
	task "taiyaki-server/task"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manager struct {
	DB         *gorm.DB
	Pending    queue.Queue
	EventDb    map[uuid.UUID]*task.TaskEvent
	LastWorker int
}

func (m *Manager) AddTask(te task.TaskEvent) {
	m.Pending.Enqueue(te)
}

func NewManager() *Manager {
	eventDb := make(map[uuid.UUID]*task.TaskEvent)

	return &Manager{
		Pending: *queue.New(),
		EventDb: eventDb,
	}
}

func KillTask(m *Manager) {
	for {
		log.Println("Fetching tasks to delete if any")
		taskCntrl := Controller.NewTask(m.DB)
		tasks := taskCntrl.GetTasksToDelete()

		allTasks := taskCntrl.GetTasks()
		for _, taskU := range allTasks {
			if taskU.Expiry.Before(time.Now()) {
				if taskU.State == "Running" {
					taskU.State = "Completed"
					taskCntrl.UpdateTask(taskU)
				}
			}
		}

		log.Println("Got the tasks to delete: ", tasks)

		for _, t := range tasks {
			oldestTask := taskCntrl.GetOldestTaskForContainer(t.ContainerID)
			workerIpPort := strings.Split(oldestTask.WorkerIpPort, ":")
			log.Println(" deleting: " + oldestTask.UUID + " workerIp: " + workerIpPort[0])
			_, err := ReqWorker("tasks/"+oldestTask.UUID, "DELETE", nil, workerIpPort[0], workerIpPort[1])
			if err != nil {
				//handle error
				log.Println("Error when deleting the task")
			}
			t.State = "Completed"
			taskCntrl.UpdateTask(oldestTask)
		}
		log.Println("Delete thread sleeping")
		time.Sleep(6 * time.Second)
	}
}

func ReqWorker(endpoint string, method string, reqBody io.Reader, workerIP string, workerPort string) (resBody []byte, err error) {
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
