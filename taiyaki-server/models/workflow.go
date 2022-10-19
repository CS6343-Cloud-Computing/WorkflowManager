package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Workflow struct {
	gorm.Model
	ID         int
	WorkflowID string
	Username   string
	Tasks      datatypes.JSON
}

// create a workflow
func CreateWorkflow(db *gorm.DB, workflow *Workflow) (err error) {
	err = db.Create(workflow).Error
	if err != nil {
		return err
	}
	return nil
}

// get workflows
func GetWorkflows(db *gorm.DB, workflows *[]Workflow) (err error) {
	err = db.Find(workflows).Error
	if err != nil {
		return err
	}
	return nil
}

// //get workflow by id
func GetWorkflow(db *gorm.DB, workflow *Workflow, id int) (err error) {
	err = db.Where("id = ?", id).First(workflow).Error
	if err != nil {
		return err
	}
	return nil
}

// update workflow
func UpdateWorkflow(db *gorm.DB, workflow *Workflow) (err error) {
	db.Save(workflow)
	return nil
}

// delete workflow
func DeleteWorkflow(db *gorm.DB, workflow *Workflow, id int) (err error) {
	db.Where("id = ?", id).Delete(workflow)
	return nil
}

// //get workflow by userName
func GetWorkflowByUserName(db *gorm.DB, workflows *[]Workflow, userName string) (err error) {
	err = db.Where("username = ?", userName).Find(workflows).Error
	if err != nil {
		return err
	}
	return nil
}