package controllers

import (
	"errors"
	"log"
	"taiyaki-server/models"

	"gorm.io/gorm"
)

type WorkerRepo struct {
	Db *gorm.DB
}

func NewWorker(db *gorm.DB) *WorkerRepo {
	return &WorkerRepo{Db: db}
}

// create worker
func (repo *WorkerRepo) CreateWorker(worker models.Worker) {
	err := models.CreateWorker(repo.Db, &worker)
	if err != nil {
		log.Println(err)
	}
}

// get worker by ip address
func (repo *WorkerRepo) GetWorker(WorkerIP string) (models.Worker, bool) {
	var worker models.Worker
	err := models.GetWorker(repo.Db, &worker, WorkerIP)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return worker, false
		}
		log.Println(err)
	}
	return worker, true
}

// get worker by ip address
func (repo *WorkerRepo) GetWorkers() []models.Worker {
	var workers []models.Worker
	err := models.GetWorkers(repo.Db, &workers)
	if err != nil {
		//panic(err)
		log.Println("Error in GetWorkers")
	}
	return workers
}

func (repo *WorkerRepo) GetActiveWorkers() []models.Worker {
	var workers []models.Worker
	err := models.GetActiveWorkers(repo.Db, &workers)
	if err != nil {
		//panic(err)
		log.Println("Error in GetActiveWorkers")
	}
	return workers
}

// update the worker
func (repo *WorkerRepo) UpdateWorker(worker models.Worker) (*models.Worker, error) {
	err := models.UpdateWorker(repo.Db, &worker)
	if err != nil {
		return nil, err
	}
	return &worker, nil
}

func (repo *WorkerRepo) GetMinTaskWorkers() []models.Worker {
	var workers []models.Worker
	err := models.GetMinTaskWorkers(repo.Db, &workers)
	if err != nil {
		//panic(err)
		log.Println("Error in GetMinTaskWorkers")
	}
	return workers
}
