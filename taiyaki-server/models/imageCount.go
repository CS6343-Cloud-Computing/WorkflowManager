package models

import (
	"gorm.io/gorm"
)

type ImageCount struct {
	gorm.Model
	Image string
	Count int
}


// create a task
func CreateEntry(db *gorm.DB, imageCount *ImageCount) (err error) {
	err = db.Create(imageCount).Error
	if err != nil {
		return err
	}
	return nil
}

// get tasks
func GetEntries(db *gorm.DB, imageCounts *[]ImageCount) (err error) {
	err = db.Find(imageCounts).Error
	if err != nil {
		return err
	}
	return nil
}

// //get task by id
func GetEntry(db *gorm.DB, imageCount *ImageCount, imageName string) (err error) {
	err = db.Where("image = ?", imageName).First(imageCount).Error
	if err != nil {
		return err
	}
	return nil
}

// //update task
func UpdateEntry(db *gorm.DB, imageCount *ImageCount) (err error) {
	db.Save(imageCount)
	return nil
}

func GetImageCount(db *gorm.DB, imageCount *[]ImageCount, images []string) (err error) {
	err = db.Raw("select * from image_counts where image in ?", images).Scan(&imageCount).Error
	if err != nil {
		return err
	}
	return nil
}
