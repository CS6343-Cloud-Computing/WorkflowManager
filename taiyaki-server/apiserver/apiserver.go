package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	workerController "taiyaki-server/controllers"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

//Resp - Generic response
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

func workflowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	reqBytes, _ := ioutil.ReadAll(r.Body)

	workflow := WorkflowTemplate{}
	err := yaml.Unmarshal(reqBytes, &workflow)
	if err != nil {
		panic(err)
	}

	resp := Resp{"Successfully got the workflow", true, ""}
	json.NewEncoder(w).Encode(resp)
}

func nodeJoinHandler(w http.ResponseWriter, r *http.Request, worker *workerController.WorkerRepo) {
	fmt.Println(worker)
}

func (c APIConfig) Start(wg *sync.WaitGroup, db *gorm.DB) {
	workerCntrl := workerController.NewWorker(db)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", UnHandler)
	router.HandleFunc("/server", UnHandler)
	router.HandleFunc("/workflow", UnHandler)
	router.HandleFunc("/node", UnHandler)
	router.HandleFunc("/server/status", serverStatusHandler)
	router.HandleFunc("/workflow/submit", workflowHandler)
	router.HandleFunc("/node/join", func(w http.ResponseWriter, r *http.Request) { nodeJoinHandler(w, r, workerCntrl) })
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
