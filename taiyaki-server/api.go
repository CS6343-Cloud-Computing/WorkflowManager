package main

import (
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"gopkg.in/yaml.v3"
	"github.com/gorilla/mux"
)

type Api struct {
	ServerIP   string
	ServerPort string
}

type WorkflowTemplate struct {
	Main struct {
		Steps []StepItem
	}
}

type StepItem struct{
	Name string
	Image string
	Cmd []string
	Env []string
	Autoremove bool
}

func UnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("UnHandler: It worked but the route is not found!!!\n"))
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	resp := Resp{Result:"Server is running",Success: true, Error: ""}
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

	resp := Resp{"Successfully got the workflow",true, ""}
	json.NewEncoder(w).Encode(resp)
}

func (a Api) start(wg *sync.WaitGroup) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", UnHandler)
	router.HandleFunc("/server", UnHandler)
	router.HandleFunc("/workflow", UnHandler)
	router.HandleFunc("/server/status", serverStatusHandler)
	router.HandleFunc("/workflow/submit", workflowHandler)
	srv := &http.Server{
		Handler: router,
		Addr:    a.ServerIP + ":" + a.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	wg.Done()
}
