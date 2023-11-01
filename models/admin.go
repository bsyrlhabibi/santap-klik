package models

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" form:"name"`
	Username  string    `json:"username" form:"username"`
	Password  string    `json:"password" form:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func (am *AdminModel) Register(name, username, password string) error {

	existingAdmin := am.GetAdminByUsername(am.db, username)
	if existingAdmin != nil {

		return errors.New("Admin with this username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &Admin{
		Name:     name,
		Username: username,
		Password: string(hashedPassword),
	}

	if err := am.db.Create(admin).Error; err != nil {
		return err
	}

	return nil
}

func (am *AdminModel) GetAdminByUsername(db *gorm.DB, username string) *Admin {
	var admin Admin
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}
