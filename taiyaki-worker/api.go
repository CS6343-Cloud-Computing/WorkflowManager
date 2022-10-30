package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Api struct {
	NodeIP   string
	NodePort string
	Worker   *Worker
	Router   *chi.Mux
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTasksHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
	a.Router.Route("/heartbeat", func(r chi.Router) {
		r.Get("/", a.GetHeartbeat)
	})
}

func (a *Api) Start() {
	a.initRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%s", a.NodeIP, a.NodePort), a.Router)
}
