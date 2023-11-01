package models

import (
	"io"
	"santapKlik/helper"
	"time"

	"gorm.io/gorm"
)

type Jajanan struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" form:"name"`
	Description string         `json:"description" form:"description"`
	Price       int            `json:"price" form:"price"`
	Image       string         `json:"image" form:"image"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type JajananModel struct {
	db *gorm.DB
}

func NewJajananModel(db *gorm.DB) *JajananModel {
	return &JajananModel{db}
}

func (jm *JajananModel) InsertJajanan(newJajanan Jajanan, imageFile io.Reader) (uint, error) {
	// Upload the image to Cloudinary and get the public ID
	imagePublicID, err := helper.UploadImage("jajanan", imageFile)
	if err != nil {
		return 0, err
	}

	// Set the Image field to the Cloudinary Public ID
	newJajanan.Image = imagePublicID

	if err := jm.db.Create(&newJajanan).Error; err != nil {
		return 0, err
	}
	return newJajanan.ID, nil
}

func (jm *JajananModel) GetJajananByID(jajananID uint) (*Jajanan, error) {
	var jajanan Jajanan

	if err := jm.db.First(&jajanan, jajananID).Error; err != nil {
		return nil, err
	}
	return &jajanan, nil
}

func (jm *JajananModel) GetAllJajanan() ([]Jajanan, error) {
	var jajananList []Jajanan
	if err := jm.db.Find(&jajananList).Error; err != nil {
		return nil, err
	}
	return jajananList, nil
}

func (jm *JajananModel) GetPageJajanan(page, limit int) ([]Jajanan, error) {
	var jajananList []Jajanan
	offset := (page - 1) * limit
	if err := jm.db.Offset(offset).Limit(limit).Find(&jajananList).Error; err != nil {
		return nil, err
	}
	return jajananList, nil
}

func (jm *JajananModel) UpdateJajanan(jajananID uint, updatedJajanan Jajanan, imageFile io.Reader) error {
	var existingJajanan Jajanan
	if err := jm.db.First(&existingJajanan, jajananID).Error; err != nil {
		return err
	}

	imagePublicID, err := helper.UploadImage("jajanan", imageFile)
	if err != nil {
		return err
	}
	updatedJajanan.Image = imagePublicID

	if updatedJajanan.Name != "" {
		existingJajanan.Name = updatedJajanan.Name
	}
	if updatedJajanan.Description != "" {
		existingJajanan.Description = updatedJajanan.Description
	}
	if updatedJajanan.Price != 0 {
		existingJajanan.Price = updatedJajanan.Price
	}

	existingJajanan.Image = updatedJajanan.Image
	existingJajanan.UpdatedAt = time.Now()

	if err := jm.db.Save(&existingJajanan).Error; err != nil {
		return err
	}

	return nil
}

func (jm *JajananModel) UpdateWithoutImageJajanan(jajananID uint, updatedJajanan Jajanan) error {

	var existingJajanan Jajanan
	if err := jm.db.Where("id = ?", jajananID).First(&existingJajanan).Error; err != nil {
		return err
	}

	if updatedJajanan.Name != "" {
		existingJajanan.Name = updatedJajanan.Name
	}
	if updatedJajanan.Description != "" {
		existingJajanan.Description = updatedJajanan.Description
	}
	if updatedJajanan.Price != 0 {
		existingJajanan.Price = updatedJajanan.Price
	}

	if err := jm.db.Save(&existingJajanan).Error; err != nil {
		return err
	}

	return nil
}

func (jm *JajananModel) DeleteJajanan(jajananID uint) error {
	var jajanan Jajanan
	if err := jm.db.Where("id = ?", jajananID).Delete(&jajanan).Error; err != nil {
		return err
	}
	return nil
}

func (mm *JajananModel) FilterJajananByPrice(minPrice, maxPrice int) ([]Jajanan, error) {
	var filteredJajanan []Jajanan

	// Lakukan filter berdasarkan harga
	if err := mm.db.Where("price >= ? AND price <= ?", minPrice, maxPrice).Find(&filteredJajanan).Error; err != nil {
		return nil, err
	}

	return filteredJajanan, nil
}
