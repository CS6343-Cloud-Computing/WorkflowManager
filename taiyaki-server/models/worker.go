package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Worker struct {
	gorm.Model
	ID         int
	WorkerIP   string
	WorkerPort string
	WorkerKey  string
	Containers datatypes.JSON
	Status     string
}

// create a worker
func CreateWorker(db *gorm.DB, worker *Worker) (err error) {
	err = db.Create(worker).Error
	if err != nil {
		return err
	}
	return nil
}

// get workers
func GetWorkers(db *gorm.DB, workers *[]Worker) (err error) {
	err = db.Find(workers).Error
	if err != nil {
		return err
	}
	return nil
}

// get workers where status is active
func GetActiveWorkers(db *gorm.DB, workers *[]Worker) (err error) {
	err = db.Find(workers).Where("Status = ?", "Active").Error
	if err != nil {
		return err
	}
	return nil
}

// //get Worker by id
func GetWorker(db *gorm.DB, worker *Worker, WorkerIP string) (err error) {
	err = db.Where("worker_ip = ?", WorkerIP).First(worker).Error
	if err != nil {
		return err
	}
	return nil
}

// //update worker
func UpdateWorker(db *gorm.DB, worker *Worker) (err error) {
	db.Save(worker)
	return nil
}

// //delete Worker
func DeleteWorker(db *gorm.DB, worker *Worker, WorkerIP string) (err error) {
	db.Where("worker_ip = ?", WorkerIP).Delete(worker)
	return nil
}
