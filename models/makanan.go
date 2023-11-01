package models

import (
	"io"
	"santapKlik/helper"
	"time"

	"gorm.io/gorm"
)

type Makanan struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" form:"name"`
	Description string         `json:"description" form:"description"`
	Price       int            `json:"price" form:"price"`
	Image       string         `json:"image" form:"image"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type MakananModel struct {
	db *gorm.DB
}

func NewMakananModel(db *gorm.DB) *MakananModel {
	return &MakananModel{db}
}

func (mm *MakananModel) InsertMakanan(newMakanan Makanan, imageFile io.Reader) (uint, error) {

	imagePublicID, err := helper.UploadImage("makanan", imageFile)
	if err != nil {
		return 0, err
	}

	newMakanan.Image = imagePublicID

	if err := mm.db.Create(&newMakanan).Error; err != nil {
		return 0, err
	}
	return newMakanan.ID, nil
}

func (mm *MakananModel) GetMakananByID(makanaID uint) (*Makanan, error) {
	var makanan Makanan

	if err := mm.db.First(&makanan, makanaID).Error; err != nil {
		return nil, err
	}
	return &makanan, nil
}

func (mm *MakananModel) GetAllMakanan() ([]Makanan, error) {
	var makananList []Makanan
	if err := mm.db.Find(&makananList).Error; err != nil {
		return nil, err
	}
	return makananList, nil
}

func (jm *MakananModel) GetPageMakanan(page, limit int) ([]Makanan, error) {
	var makananList []Makanan
	offset := (page - 1) * limit
	if err := jm.db.Offset(offset).Limit(limit).Find(&makananList).Error; err != nil {
		return nil, err
	}
	return makananList, nil
}

func (mm *MakananModel) UpdateMakanan(makananID uint, updatedMakanan Makanan, imageFile io.Reader) error {
	var existingMakanan Makanan
	if err := mm.db.First(&existingMakanan, makananID).Error; err != nil {
		return err
	}

	imagePublicID, err := helper.UploadImage("jajanan", imageFile)
	if err != nil {
		return err
	}
	updatedMakanan.Image = imagePublicID

	if updatedMakanan.Name != "" {
		existingMakanan.Name = updatedMakanan.Name
	}
	if updatedMakanan.Description != "" {
		existingMakanan.Description = updatedMakanan.Description
	}
	if updatedMakanan.Price != 0 {
		existingMakanan.Price = updatedMakanan.Price
	}

	existingMakanan.Image = updatedMakanan.Image
	existingMakanan.UpdatedAt = time.Now()

	if err := mm.db.Save(&existingMakanan).Error; err != nil {
		return err
	}

	return nil
}

func (mm *MakananModel) UpdateWithoutImageMakanan(makananID uint, updatedMakanan Makanan) error {

	var existingMakanan Makanan
	if err := mm.db.Where("id = ?", makananID).First(&existingMakanan).Error; err != nil {
		return err
	}

	if updatedMakanan.Name != "" {
		existingMakanan.Name = updatedMakanan.Name
	}
	if updatedMakanan.Description != "" {
		existingMakanan.Description = updatedMakanan.Description
	}
	if updatedMakanan.Price != 0 {
		existingMakanan.Price = updatedMakanan.Price
	}

	if err := mm.db.Save(&existingMakanan).Error; err != nil {
		return err
	}

	return nil
}

func (mm *MakananModel) DeleteMakanan(makananID uint) error {
	var makanan Makanan
	if err := mm.db.Where("id = ?", makananID).Delete(&makanan).Error; err != nil {
		return err
	}
	return nil
}

func (mm *MakananModel) FilterMakananByPrice(minPrice, maxPrice int) ([]Makanan, error) {
	var filteredMakanan []Makanan

	// Lakukan filter berdasarkan harga
	if err := mm.db.Where("price >= ? AND price <= ?", minPrice, maxPrice).Find(&filteredMakanan).Error; err != nil {
		return nil, err
	}

	return filteredMakanan, nil
}
