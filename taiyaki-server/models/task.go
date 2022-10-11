package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID            int
	Uuid          string
	ContainerId   string
	UserName      string
	Name          string
	State         string
	RestartPolicy string
	StartTime     datatypes.Time
	FinishTime    datatypes.Time
	Config        datatypes.JSON
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
