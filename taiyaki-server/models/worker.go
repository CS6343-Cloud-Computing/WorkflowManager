package models

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
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

//create a worker
func CreateWorker(db *gorm.DB, worker *Worker) (err error) {
	err = db.Create(worker).Error
	if err != nil {
		return err
	}
	return nil
}

//get workers
// func GetWorkers(db *gorm.DB, worker *[]Worker) (err error) {
// 	err = db.Find(worker).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// //get user by id
// func GetUser(db *gorm.DB, User *User, id int) (err error) {
// 	err = db.Where("id = ?", id).First(User).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// //update user
// func UpdateUser(db *gorm.DB, User *User) (err error) {
// 	db.Save(User)
// 	return nil
// }

// //delete user
// func DeleteUser(db *gorm.DB, User *User, id int) (err error) {
// 	db.Where("id = ?", id).Delete(User)
// 	return nil
// }