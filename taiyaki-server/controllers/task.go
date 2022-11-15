package controllers

import (
	"errors"
	"taiyaki-server/models"

	"gorm.io/gorm"
)

type TaskRepo struct {
	Db *gorm.DB
}

func NewTask(db *gorm.DB) *TaskRepo {
	return &TaskRepo{Db: db}
}

// create task
func (repo *TaskRepo) CreateTask(task models.Task) {
	err := models.CreateTask(repo.Db, &task)
	if err != nil {
		panic(err)
	}
}

// get task by uuid
func (repo *TaskRepo) GetTask(uuid string) (models.Task, bool) {
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

func (repo *TaskRepo) GetTasks() []models.Task {
	var tasks []models.Task
	err := models.GetTasks(repo.Db, &tasks)
	if err != nil {
		panic(err)
	}
	return tasks
}

// update the task
func (repo *TaskRepo) UpdateTask(task models.Task) {
	err := models.UpdateTask(repo.Db, &task)
	if err != nil {
		panic(err)
	}
}

func (repo *TaskRepo) GetTasksToDelete() []models.Task {
	var tasks []models.Task
	err := models.GetTasksToDelete(repo.Db, &tasks)
	if err != nil {
		panic(err)
	}
	return tasks
}


func (repo *TaskRepo) GetRunningTasks() []models.Task {
	var tasks []models.Task
	err := models.GetRunningTasks(repo.Db, &tasks)
	if err != nil {
		//panic(err)
	}
	return tasks
}

func (repo *TaskRepo) GetTaskWithSameImage(imageName string) models.Task {
	var task models.Task
	err := models.GetTaskWithSameImage(repo.Db, &task, imageName)
	if err != nil {
		//panic(err)
	}
	return task
}

func (repo *TaskRepo) GetOldestTaskForContainer(containerId string) models.Task {
	var task models.Task
	err := models.GetOldestTaskForContainer(repo.Db, &task, containerId)
	if err != nil {
		//panic(err)
	}
	return task
}