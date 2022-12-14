package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	task "taiyaki-worker/task"

	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/exp/maps"
)

// Resp - Generic response
type Resp struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type NodeJoinReq struct {
	NodeIP   string
	NodePort string
	JoinKey  string
}

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]task.Task
	TaskCount int
	Stats     Stats
}

var workerIP string
var workerPort string
var serverIP string
var serverPort string
var joinKey string

func (w *Worker) CollectionStats() {
	fmt.Println("Collect stats")
}

func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return task.DockerResult{Error: nil}
	}
	taskQueued := t.(task.Task)
	log.Println("")
	taskPersisted, isPresent := w.Db[taskQueued.ID]

	if !isPresent {
		taskPersisted = taskQueued
		//w.Db[taskQueued.ID] = taskQueued
	}

	var result task.DockerResult
	fmt.Println(task.ValidStateTransition(taskPersisted.State, taskQueued.State))
	if true {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("we should not get here")
		}
	} else {
		err := fmt.Errorf("invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}
	return result
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config, t)
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = t
		return result
	}

	d.ContainerId = result.ContainerId
	t.ContainerId = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = t

	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config, t)
	d.ContainerId = t.ContainerId
	result := d.Stop()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = t
	log.Printf("Stopped and removed the container %v for task %v", d.ContainerId, t)

	return result
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) GetTasks() []task.Task {
	values := maps.Values(w.Db)
	return values
}

func (w *Worker) CollectStats() {
	for {
		log.Println("Collecting stats")
		w.Stats = *GetStats()
		time.Sleep(15 * time.Second)
	}
}

func runTasks(w *Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}

func findContainerInList(containerList []types.Container, containerId string) *types.Container {
	for _, cntr := range containerList {
		if cntr.ID == containerId {
			return &cntr
		}
	}
	return nil
}

func splitStatus(status string) int {
	re := regexp.MustCompile(`\((.*?)\)`)
	submatchall := re.FindAllString(status, -1)
	for _, statusCode := range submatchall {
		statusCode = strings.Trim(statusCode, "(")
		statusCode = strings.Trim(statusCode, ")")
		intStCode, err := strconv.Atoi(statusCode)
		if err != nil {
			log.Println(err)
		}
		return intStCode
	}
	return 0
}

func runSyncDockerStatuses(w *Worker) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println(err)
	}
	for {
		//log.Println("syncing docker status")
		if len(w.Db) > 0 {
			//log.Println("syncing since queue has entries")
			containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
			if err != nil {
				log.Println(err)
			}

			//log.Println("list of containers running: ", containerList)
			//log.Println("syncing since queue has entries before finding in the list")
			for taskId, taskRunning := range w.Db {
				cntr := findContainerInList(containerList, taskRunning.ContainerId)
				if cntr != nil && cntr.State == "exited" {
					status := splitStatus(cntr.Status)
					if status != 0 {
						taskRunning.State = task.Failed
						//fmt.Println("-----------------syncDockerStatuses task.State: ", status)
						w.Db[taskId] = taskRunning
					} else {
						taskRunning.State = task.Completed
						//fmt.Println("-----------------syncDockerStatuses task.State: ", status)
						w.Db[taskId] = taskRunning
					}
				}
			}
		}
		log.Println("syncing docker status sleeping for 3 seconds")
		time.Sleep(3 * time.Second)
	}
}

func reqServer(endpoint string, reqBody io.Reader) (resBody []byte, err error) {
	url := "http://" + serverIP + ":" + serverPort + "/" + endpoint
	c := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: c}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, reqBody)
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

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	if workerIP == "" {
		if env := os.Getenv("WORKER_IP"); env == "" {
			conn, err := net.Dial("udp", "8.8.8.8:80")
			if err != nil {
				// fmt.Println(err)
				fmt.Println("Line 260 main.go UDP dial error")
			}
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			workerIP = localAddr.IP.String()
			fmt.Println("No Worker IP provided. But I can do better:", workerIP)
			conn.Close()
		} else {
			workerIP = env
		}
	}

	if workerPort == "" {
		if env := os.Getenv("WORKER_PORT"); env == "" {
			fmt.Println("No server to connect to... :/")

		} else {
			workerPort = env
		}
	}

	if serverIP == "" {
		if env := os.Getenv("SERVER_IP"); env == "" {
			fmt.Println("No server to connect to... :/")
		} else {
			serverIP = env
		}
	}

	if serverPort == "" {
		if env := os.Getenv("SERVER_PORT"); env == "" {
			fmt.Println("No server to connect to... :/")
		} else {
			serverPort = env
		}
	}

	if joinKey == "" {
		if env := os.Getenv("JOIN_KEY"); env == "" {
			fmt.Println("Please provide join key for this cluster")
		} else {
			joinKey = env
		}
	}

	for !NodeJoin(workerIP, workerPort, serverIP, serverPort, joinKey) {
		fmt.Println("Server is not running probably")
		time.Sleep(2 * time.Second)
	}

	worker := Worker{Queue: *queue.New(), Db: make(map[uuid.UUID]task.Task)}
	go runTasks(&worker)
	go worker.CollectStats()

	go runSyncDockerStatuses(&worker)

	api := Api{NodeIP: workerIP, NodePort: workerPort, Worker: &worker}

	api.Start()
}
