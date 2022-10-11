package controllers

import (
	"taiyaki-server/models"
	"errors"
	"gorm.io/gorm"
)

type TaskRepo struct {
	Db *gorm.DB
}

func NewTask(db * gorm.DB) *TaskRepo {
	return &TaskRepo{Db: db}
}

//create task
func (repo *TaskRepo) CreateTask(task models.Task) {
	err := models.CreateTask(repo.Db, &task)
	if err != nil {
		panic(err)
	}
}

//get task by uuid
func (repo *TaskRepo) GetTask(uuid string) (models.Task, bool){
	var task models.Task
	err := models.GetTask(repo.Db, &task, uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return task, false
		}
		panic(err)
	}
	return task, true
}

//update the task
func (repo *TaskRepo) UpdateTask(task models.Task) {
	err := models.UpdateTask(repo.Db, &task)
	if err != nil {
		panic(err)
	}
}