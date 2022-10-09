package controllers

import (
	"taiyaki-server/models"
	"errors"
	"gorm.io/gorm"
)

type WorkerRepo struct {
	Db *gorm.DB
}

func NewWorker(db * gorm.DB) *WorkerRepo {
	return &WorkerRepo{Db: db}
}

//create user
func (repo *WorkerRepo) CreateWorker() {
	var worker models.Worker
	// fill worker
	err := models.CreateWorker(repo.Db, &worker)
	if err != nil {
		panic(err)
	}
}

//get user by id
func (repo *WorkerRepo) GetWorker(WorkerIP string) (models.Worker, bool){
	var worker models.Worker
	err := models.GetWorker(repo.Db, &worker, WorkerIP)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return worker, false
		}
		panic(err)
	}
	return worker, true
}