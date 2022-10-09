package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Api struct {
	ServerIP   string
	ServerPort string
}

func UnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("UnHandler: It worked but the route is not found!!!\n"))
}

func (a Api) start(wg *sync.WaitGroup) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", UnHandler)
	router.HandleFunc("/server", UnHandler)
	srv := &http.Server{
		Handler: router,
		Addr:    a.ServerIP + ":" + a.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	wg.Done()
}
