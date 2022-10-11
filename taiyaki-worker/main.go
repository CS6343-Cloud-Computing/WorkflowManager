package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	task "taiyaki-worker/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"time"
)

//Resp - Generic response
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
		w.Db[taskQueued.ID] = taskQueued
	}

	var result task.DockerResult
	fmt.Println(task.ValidStateTransition(taskPersisted.State, taskQueued.State))
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
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
	d := task.NewDocker(config)
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
	d := task.NewDocker(config)
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

func (w *Worker) GetTasks() []byte {
	jsonStr, err := json.Marshal(w.Db)
	if err != nil {
		log.Printf("Error marshalling the queue\n")
	}

	return jsonStr
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
				log.Fatal(err)
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
			log.Fatal("No server to connect to... :/")
		} else {
			workerPort = env
		}
	}

	if serverIP == "" {
		if env := os.Getenv("SERVER_IP"); env == "" {
			log.Fatal("No server to connect to... :/")
		} else {
			serverIP = env
		}
	}

	if serverPort == "" {
		if env := os.Getenv("SERVER_PORT"); env == "" {
			log.Fatal("No server to connect to... :/")
		} else {
			serverPort = env
		}
	}

	if joinKey == "" {
		if env := os.Getenv("JOIN_KEY"); env == "" {
			log.Fatal("Please provide join key for this cluster")
		} else {
			joinKey = env
		}
	}

	NodeJoin(workerIP, workerPort, serverIP, serverPort, joinKey)

	worker := Worker{Queue: *queue.New(), Db: make(map[uuid.UUID]task.Task)}
	go runTasks(&worker)
	go worker.CollectStats()
	api := Api{NodeIP: workerIP, NodePort: workerPort, Worker: &worker}

	api.Start()
}
