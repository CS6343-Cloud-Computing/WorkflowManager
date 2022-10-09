package controllers

import (
	"taiyaki-server/models"

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