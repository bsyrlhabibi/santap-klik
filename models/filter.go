package models

import (
	"gorm.io/gorm"
)

type Filter struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Jenis    string `json:"jenis" form:"jenis"`
	MinPrice int    `json:"minPrice" form:"minPrice"`
	MaxPrice int    `json:"maxPrice" form:"maxPrice"`
}

type FilterModel struct {
	db *gorm.DB
}

func NewFilterModel(db *gorm.DB) *FilterModel {
	return &FilterModel{db}
}
