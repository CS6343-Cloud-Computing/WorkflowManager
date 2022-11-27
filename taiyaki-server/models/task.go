package models

import (
	"log"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID              int
	WorkflowID      string
	UUID            string
	ContainerID     string
	Name            string
	State           string
	RestartPolicy   string
	StartTime       datatypes.Time
	FinishTime      datatypes.Time
	Config          datatypes.JSON
	DeploymentOrder int
	WorkerIpPort    string
	Expiry          time.Time
	Image           string
	Output          datatypes.JSON
	Input           datatypes.JSON
	Indegree        int
	Persistence     bool
}

type ContainerCount struct {
	ContainerId string
	Image       string
	Count       int
}

// create a task
func CreateTask(db *gorm.DB, task *Task) (err error) {
	err = db.Create(task).Error
	if err != nil {
		return err
	}
	return nil
}

// get tasks
func GetTasks(db *gorm.DB, tasks *[]Task) (err error) {
	err = db.Find(tasks).Error
	if err != nil {
		return err
	}
	return nil
}

// //get task by id
func GetTask(db *gorm.DB, task *Task, uuid string) (err error) {
	err = db.Where("uuid = ?", uuid).First(task).Error
	if err != nil {
		return err
	}
	return nil
}

// //update task
func UpdateTask(db *gorm.DB, task *Task) (err error) {
	db.Save(task)
	return nil
}

// //delete task
func DeleteTask(db *gorm.DB, task *Task, uuid string) (err error) {
	db.Where("uuid = ?", uuid).Delete(task)
	return nil
}

func GetTasksToDelete(db *gorm.DB, tasks *[]Task) (err error) {
	err = db.Raw("select * from tasks where container_id in (select t.container_id  from tasks t where t.state = \"running\" group by t.container_id having max(t.expiry) < UTC_TIMESTAMP())").Scan(&tasks).Error
	if err != nil {
		return err
	}
	return nil
}

func GetRunningTasks(db *gorm.DB, tasks *[]Task) (err error) {
	err = db.Where("state = ?", "Running").Find(tasks).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTaskWithSameImage(db *gorm.DB, task *Task, imageName string) (err error) {
	err = db.Where("state = ? and image = ? ", "Running", imageName).First(task).Error
	if err != nil {
		return err
	}
	return nil
}

func GetOldestTaskForContainer(db *gorm.DB, task *Task, containerId string) (err error) {
	err = db.Where("container_id = ?", containerId).Order("created_at").First(task).Error
	if err != nil {
		return err
	}
	return nil
}

func GetCntnrIdFromWorkflowId(db *gorm.DB, containerIds *[]string, workflowId string) (err error) {
	err = db.Raw("select container_id from tasks where workflow_id = ?", workflowId).Scan(&containerIds).Error
	if err != nil {
		return err
	}
	return nil
}

func GetCountImageInContainers(db *gorm.DB, containerIdCount *[]ContainerCount, containerIds []string) (err error) {

	err = db.Raw("select container_id, image, COUNT(container_id) as count from tasks where state = \"Running\" and container_id in ? GROUP by container_id, image", containerIds).Scan(&containerIdCount).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateTasksInWrkFlw(db *gorm.DB, workflowId string) (err error) {
	err = db.Exec("update tasks set state = \"KillBitReceived\" where workflow_id = ?", workflowId).Error
	if err != nil {
		return err
	}
	return nil
}