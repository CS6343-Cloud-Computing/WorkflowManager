package controllers

import (
	"errors"
	"taiyaki-server/models"

	"gorm.io/gorm"
)

type WorkflowRepo struct {
	Db *gorm.DB
}

func NewWorkflow(db *gorm.DB) *WorkflowRepo {
	return &WorkflowRepo{Db: db}
}

// create workflow
func (repo *WorkflowRepo) CreateWorkflow(workflow models.Workflow) {
	err := models.CreateWorkflow(repo.Db, &workflow)
	if err != nil {
		panic(err)
	}
}

// get workflow by uuid
func (repo *WorkflowRepo) GetWorkflow(id int) (models.Workflow, bool) {
	var workflow models.Workflow
	err := models.GetWorkflow(repo.Db, &workflow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return workflow, false
		}
		panic(err)
	}
	return workflow, true
}

// update the workflow
func (repo *WorkflowRepo) UpdateWorkflow(workflow models.Workflow) {
	err := models.UpdateWorkflow(repo.Db, &workflow)
	if err != nil {
		panic(err)
	}
}
