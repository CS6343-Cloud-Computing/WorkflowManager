package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	Controller "taiyaki-server/controllers"
	Manager "taiyaki-server/manager"
	"taiyaki-server/models"
	Scheduler "taiyaki-server/scheduler"
	"taiyaki-server/task"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"gorm.io/datatypes"
	//"gorm.io/gorm"
)

type NodeJoinReq struct {
	NodeIP   string
	NodePort string
	JoinKey  string
}

// Resp - Generic response
type Resp struct {
	Result  string `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type APIConfig struct {
	ServerIP      string
	ServerPort    string
	WorkerJoinKey string
}

type WorkflowTemplate struct {
	Main struct {
		Username   string
		Datasource string
		Steps      []StepItem
		Expiry     int
	}
}

type StepItem struct {
	Name       string
	Image      string
	Cmd        []string
	Env        []string
	Query      string
	Autoremove bool
}

func UnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("UnHandler: It worked but the route is not found!!!\n"))
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	resp := Resp{Result: "Server is running", Success: true, Error: ""}
	json.NewEncoder(w).Encode(resp)
}

func workflowHandler(w http.ResponseWriter, r *http.Request, taskCntrl *Controller.TaskRepo, workflowCntrl *Controller.WorkflowRepo, m *Manager.Manager) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	reqBytes, _ := ioutil.ReadAll(r.Body)

	workflow := WorkflowTemplate{}
	err := yaml.Unmarshal(reqBytes, &workflow)
	if err != nil {
		panic(err)
	}

	workflowDb := models.Workflow{}
	workflowDb.Username = workflow.Main.Username

	expiry := workflow.Main.Expiry
	workflowDb.Expiry = time.Now().Add(time.Second * time.Duration(expiry + 60))
	workflowId := uuid.New().String()
	workflowDb.WorkflowID = workflowId
	workflowDb.Datasource = workflow.Main.Datasource

	var taskIds []string
	for order, workflowTask := range workflow.Main.Steps {
		taskDb := models.Task{}
		taskOb := task.Task{}

		taskDb.WorkflowID = workflowId
		taskOb.WorkflowID = workflowId

		taskOb.ID = uuid.New()
		taskDb.UUID = taskOb.ID.String()

		taskDb.ContainerID = taskOb.ID.String()
		taskOb.ContainerId = taskOb.ID.String()

		taskDb.Name = taskOb.ID.String()
		taskOb.Name = taskOb.ID.String()

		taskDb.State = "Pending"
		taskOb.State = task.Pending

		taskDb.Order = order

		config := &task.Config{}
		config.Image = workflowTask.Image
		config.Cmd = workflowTask.Cmd
		config.Env = workflowTask.Env
		config.Query = workflowTask.Query

		taskOb.Config = *config

		taskDb.Expiry = workflowDb.Expiry
		configJson, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}
		taskDb.Config = configJson
		taskCntrl.CreateTask(taskDb)

		taskIds = append(taskIds, taskDb.UUID)
		te := task.TaskEvent{}
		te.ID = uuid.New()
		te.Timestamp = time.Now()
		te.State = 1
		te.Task = taskOb
		m.AddTask(te)
	}
	taskIdsJson, err := json.Marshal(taskIds)
	if err != nil {
		panic(err)
	}
	workflowDb.Tasks = taskIdsJson

	workflowCntrl.CreateWorkflow(workflowDb)

	if len(workflow.Main.Steps) == 0 {
		resp := Resp{"Empty workflow", true, ""}
		json.NewEncoder(w).Encode(resp)
		return
	}
	fmt.Println("Successfully got and saved the workflow")
	resp := Resp{"Successfully got the workflow with id: " + workflowId, true, ""}
	json.NewEncoder(w).Encode(resp)
}

func getTasksForUser(w http.ResponseWriter, r *http.Request, m *Manager.Manager) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()
	params := mux.Vars(r)
	userName := params["userName"]
	fmt.Println("Getting workflows for user " + userName)
	workflowCntrl := Controller.NewWorkflow(m.DB)
	workflws, valid := workflowCntrl.GetWorkflowByUserName(userName)
	if !valid {
		fmt.Println("Error getting workflows for user")
		panic(valid)
	}
	fmt.Println(workflws, " hello")
	//find the workflow id
	var tasks []datatypes.JSON

	for _, workflow := range workflws {
		fmt.Println(workflow.ID, workflow.Tasks)
		tasks = append(tasks, workflow.Tasks)
	}

	json.NewEncoder(w).Encode(tasks)
}
func deleteTask(w http.ResponseWriter, r *http.Request, taskCntrl *Controller.TaskRepo) {

	params := mux.Vars(r)
	taskID := params["taskID"]
	fmt.Println("deleting task " + taskID)

	taskObj, valid := taskCntrl.GetTask(taskID)
	if !valid {
		fmt.Println("Error while getting taskObj")
		panic(valid)
	}

	if strings.Compare(taskObj.State, "Failed") == 0 ||
		strings.Compare(taskObj.State, "Completed") == 0 {
		fmt.Println("Task already in end state. Unable to delete. ")
		return
	}

	workerIpPort := taskObj.WorkerIpPort
	workerStr := strings.Split(workerIpPort, ":")

	fmt.Println("task is in worker: " + workerIpPort)
	endPoint := "tasks/" + taskID

	respBody, err1 := Scheduler.ReqWorker(endPoint, "DELETE", nil, workerStr[0], workerStr[1])
	if err1 != nil {
		fmt.Println("error while serving the request to worker")
		panic(err1)
	}

	println(string(respBody))
}

func nodeJoinHandler(w http.ResponseWriter, r *http.Request, workerCntrl *Controller.WorkerRepo, config APIConfig) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	joinReq := NodeJoinReq{}
	json.NewDecoder(r.Body).Decode(&joinReq)

	//Verify key
	if joinReq.JoinKey != config.WorkerJoinKey {
		fmt.Println(config.WorkerJoinKey)
		fmt.Println(joinReq.JoinKey)
		resp := Resp{Success: false, Error: "Invalid join key"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//Check if worker exist or not
	worker, valid := workerCntrl.GetWorker(joinReq.NodeIP)

	if valid {
		resp := Resp{Success: false, Error: "This worker node has already joined the cluster but has been marked as active now"}
		worker.Status = "active"
		workerCntrl.UpdateWorker(worker)
		json.NewEncoder(w).Encode(resp)
		return
	}

	worker = models.Worker{WorkerIP: joinReq.NodeIP, WorkerPort: joinReq.NodePort, WorkerKey: joinReq.JoinKey, Containers: datatypes.JSON{}, Status: "active", NumContainers: 0}

	fmt.Println(worker)

	workerCntrl.CreateWorker(worker)

	resp := Resp{Success: true, Error: ""}
	json.NewEncoder(w).Encode(resp)
}

func (c APIConfig) Start(wg *sync.WaitGroup, m *Manager.Manager) {
	db := m.DB
	workerCntrl := Controller.NewWorker(db)
	taskCntrl := Controller.NewTask(db)
	workflowCntrl := Controller.NewWorkflow(db)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", UnHandler)
	router.HandleFunc("/server", UnHandler)
	router.HandleFunc("/workflow", UnHandler)
	router.HandleFunc("/node", UnHandler)
	router.HandleFunc("/server/status", serverStatusHandler)
	router.HandleFunc("/workflow/submit", func(w http.ResponseWriter, r *http.Request) { workflowHandler(w, r, taskCntrl, workflowCntrl, m) })
	router.HandleFunc("/node/join", func(w http.ResponseWriter, r *http.Request) { nodeJoinHandler(w, r, workerCntrl, c) })
	router.HandleFunc("/tasks/{taskID}", func(w http.ResponseWriter, r *http.Request) { deleteTask(w, r, taskCntrl) })
	router.HandleFunc("/tasks/User/{userName}", func(w http.ResponseWriter, r *http.Request) { getTasksForUser(w, r, m) })
	srv := &http.Server{
		Handler:      router,
		Addr:         c.ServerIP + ":" + c.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Listening on", c.ServerIP+":"+c.ServerPort)
	log.Fatal(srv.ListenAndServe())
	wg.Done()
}

func isWorkerPresent(m *Manager.Manager) bool {
	workrCntrl := Controller.NewWorker(m.DB)
	workers := workrCntrl.GetWorkers()
	return len(workers) > 0
}

func SendWork(m *Manager.Manager) {

	if !isWorkerPresent(m) {
		log.Println("No worker is present")
		return
	}
	if m.Pending.Len() > 0 {
		nilWorkr := models.Worker{}
		//get the image name

		//check if container with same image is Running

		//check if stats of that worker fits criteria, then no need to send it to worker
		w := Scheduler.SelectWorker(m)
		//if it returns no worker, return from the func
		if w.ID == nilWorkr.ID {
			return
		}
		e := m.Pending.Dequeue()
		te := e.(task.TaskEvent)

		log.Printf("Pulled %v off pending queue", te.Task)

		m.EventDb[te.ID] = &te

		te.Task.State = task.Scheduled
		taskCntrl := Controller.NewTask(m.DB)
		taskUpdate, _ := taskCntrl.GetTask(te.Task.ID.String())
		taskUpdate.State = "Scheduled"

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Unable to marshal task object: %v.", te.Task)
		}

		workerIpPort := w.WorkerIP + ":" + w.WorkerPort
		fmt.Println("Sending to a work", workerIpPort)
		url := fmt.Sprintf("http://%s/tasks", workerIpPort)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error connecting to %v: %v", w, err)
			m.Pending.Enqueue(te.Task)
			return
		}

		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			fmt.Println(d)
			// e := .ErrResponse{}
			// err := d.Decode(&e)
			// if err != nil {
			// 	fmt.Printf("Error decoding response: %s\n", err.Error())
			// 	return
			// }
			//log.Printf("Response error (%d): %s", e.HTTPStatusCode, e.Message)
			return
		}

		wrkrCntrl := Controller.NewWorker(m.DB)
		w.NumContainers = w.NumContainers + 1
		wrkrCntrl.UpdateWorker(w)

		taskUpdate.WorkerIpPort = workerIpPort
		taskCntrl.UpdateTask(taskUpdate)
		t := task.Task{}
		err = d.Decode(&t)
		if err != nil {
			fmt.Printf("Error decoding response: %s\n", err.Error())
			return
		}
		log.Printf("%#v\n", t)
	} else {
		log.Println("No work in the queue")
	}
}

// func UpdateTasks(m *Manager.Manager) {
// 	for _, worker := range m.Workers {
// 		log.Printf("Checking worker %v for task updates", worker)
// 		url := fmt.Sprintf("http://%s/tasks", worker)
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			log.Printf("Error connecting to %v: %v", worker, err)
// 		}

// 		if resp.StatusCode != http.StatusOK {
// 			log.Printf("Error sending request: %v", err)
// 		}

// 		d := json.NewDecoder(resp.Body)
// 		var tasks []*task.Task
// 		err = d.Decode(&tasks)
// 		if err != nil {
// 			log.Printf("Error unmarshalling tasks: %s", err.Error())
// 		}

// 		for _, t := range tasks {
// 			log.Printf("Attempting to update task %v", t.ID)

// 			_, ok := m.TaskDb[t.ID]
// 			if !ok {
// 				log.Printf("Task with ID %s not found\n", t.ID)
// 				return
// 			}

// 			if m.TaskDb[t.ID].State != t.State {
// 				m.TaskDb[t.ID].State = t.State
// 			}

// 			m.TaskDb[t.ID].StartTime = t.StartTime
// 			m.TaskDb[t.ID].FinishTime = t.FinishTime
// 			m.TaskDb[t.ID].ContainerId = t.ContainerId
// 		}
// 	}
// }
