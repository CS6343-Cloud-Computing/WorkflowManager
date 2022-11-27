package controllers

import (
	"errors"
	"log"
	"taiyaki-server/models"

	"gorm.io/gorm"
)

type ImageRepo struct {
	Db *gorm.DB
}

func NewEntry(db *gorm.DB) *ImageRepo {
	return &ImageRepo{Db: db}
}

// create task
func (repo *ImageRepo) CreateEntry(imageCount models.ImageCount) {
	err := models.CreateEntry(repo.Db, &imageCount)
	if err != nil {
		log.Println(err)
	}
}

// get task by uuid
func (repo *ImageRepo) GetEntry(imageName string) (models.ImageCount, bool) {
	var imageCount models.ImageCount
	err := models.GetEntry(repo.Db, &imageCount, imageName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return imageCount, false
		}
		log.Println(err)
	}
	return imageCount, true
}

func (repo *ImageRepo) GetEntries() []models.ImageCount {
	var imageCounts []models.ImageCount
	err := models.GetEntries(repo.Db, &imageCounts)
	if err != nil {
		log.Println(err)
	}
	return imageCounts
}

// update the task
func (repo *ImageRepo) UpdateEntry(imageCount models.ImageCount) {
	err := models.UpdateEntry(repo.Db, &imageCount)
	if err != nil {
		log.Println(err)
	}
}

func (repo *ImageRepo) GetImageCount(images []string) []models.ImageCount {
	var imageCount []models.ImageCount
	err := models.GetImageCount(repo.Db, &imageCount, images)
	if err != nil {
		//panic(err)
		log.Println("Error in getting count image in containers")
	}
	return imageCount
}