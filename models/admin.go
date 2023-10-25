package models

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type LoginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminModel struct {
	db *gorm.DB
}

func NewAdminModel(db *gorm.DB) *AdminModel {
	return &AdminModel{db}
}

func (am *AdminModel) Login(username string, password string) *Admin {
	var data Admin
	if err := am.db.Where("username = ?", username).First(&data).Error; err != nil {
		logrus.Error("Model : Login data error,", err.Error())
		return nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password))
	if err != nil {
		logrus.Error("Model : Password does not match")
		return nil
	}

	return &data
}
