package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	Controller "taiyaki-server/controllers"
	"taiyaki-server/models"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
		Steps []StepItem
	}
}

type StepItem struct {
	Name       string
	Image      string
	Cmd        []string
	Env        []string
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

func workflowHandler(w http.ResponseWriter, r *http.Request, taskCntrl *Controller.TaskRepo, workflowCntrl *Controller.WorkflowRepo) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	reqBytes, _ := ioutil.ReadAll(r.Body)

	workflow := WorkflowTemplate{}
	err := yaml.Unmarshal(reqBytes, &workflow)
	if err != nil {
		panic(err)
	}

	if len(workflow.Main.Steps) == 0 {
		resp := Resp{"Empty workflow", true, ""}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := Resp{"Successfully got the workflow", true, ""}
	json.NewEncoder(w).Encode(resp)
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

	worker = models.Worker{WorkerIP: joinReq.NodeIP, WorkerPort: joinReq.NodePort, WorkerKey: joinReq.JoinKey, Containers: datatypes.JSON{}, Status: "active"}

	fmt.Println(worker)

	workerCntrl.CreateWorker(worker)

}

func (c APIConfig) Start(wg *sync.WaitGroup, db *gorm.DB) {
	workerCntrl := Controller.NewWorker(db)
	taskCntrl := Controller.NewTask(db)
	workflowCntrl := Controller.NewWorkflow(db)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", UnHandler)
	router.HandleFunc("/server", UnHandler)
	router.HandleFunc("/workflow", UnHandler)
	router.HandleFunc("/node", UnHandler)
	router.HandleFunc("/server/status", serverStatusHandler)
	router.HandleFunc("/workflow/submit", func(w http.ResponseWriter, r *http.Request) { workflowHandler(w, r, taskCntrl, workflowCntrl) })
	router.HandleFunc("/node/join", func(w http.ResponseWriter, r *http.Request) { nodeJoinHandler(w, r, workerCntrl, c) })
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
