package controllers

import (
	"errors"
	"log"
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
		log.Println(err)
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
		log.Println(err)
	}
	return task, true
}

func (repo *TaskRepo) GetTasks() []models.Task {
	var tasks []models.Task
	err := models.GetTasks(repo.Db, &tasks)
	if err != nil {
		log.Println(err)
	}
	return tasks
}

// update the task
func (repo *TaskRepo) UpdateTask(task models.Task) {
	err := models.UpdateTask(repo.Db, &task)
	if err != nil {
		log.Println(err)
	}
}

func (repo *TaskRepo) GetTasksToDelete() []models.Task {
	var tasks []models.Task
	err := models.GetTasksToDelete(repo.Db, &tasks)
	if err != nil {
		log.Println(err)
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

func (repo *TaskRepo) GetTasksWithSameImage(imageName string) []models.Task {
	var tasks []models.Task
	err := models.GetTasksWithSameImage(repo.Db, &tasks, imageName)
	if err != nil {
		//panic(err)
	}
	return tasks
}

func (repo *TaskRepo) GetOldestTaskForContainer(containerId string) models.Task {
	var task models.Task
	err := models.GetOldestTaskForContainer(repo.Db, &task, containerId)
	if err != nil {
		//panic(err)
	}
	return task
}

func (repo *TaskRepo) GetCntnrIdFromWorkflowId(workflowId string) []string {
	var containerIds []string
	err := models.GetCntnrIdFromWorkflowId(repo.Db, &containerIds, workflowId)
	if err != nil {
		//panic(err)
		log.Println("error in GetCntnrIdFromWorkflowId ", err)
	}
	return containerIds
}

func (repo *TaskRepo) GetCountImageInContainers(containerIds []string) []models.ContainerCount {
	var containerIdCount []models.ContainerCount
	err := models.GetCountImageInContainers(repo.Db, &containerIdCount, containerIds)
	if err != nil {
		//panic(err)
		log.Println("Error in getting count image in containers")
	}
	return containerIdCount
}

func (repo *TaskRepo) UpdateTasksInWrkFlw(workflowId string) {
	err := models.UpdateTasksInWrkFlw(repo.Db, workflowId)
	if err != nil {
		log.Println(err)
	}
}

func (repo *TaskRepo) GetLatestTaskWithImage(imageName string) models.Task {
	var task models.Task
	err := models.GetLatestTaskWithImage(repo.Db, &task, imageName)
	if err != nil {
		//panic(err)
	}
	return task
}