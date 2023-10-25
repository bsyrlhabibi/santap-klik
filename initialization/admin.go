// File: initialization/initialize.go
package initialization

import (
	"santapKlik/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitializeAdmin(db *gorm.DB) {
	initialAdminUsername := "admin"
	adminModel := models.NewAdminModel(db)
	if adminModel.Login(initialAdminUsername, "admin") == nil {
		initialAdmin := models.Admin{
			Name:     "Monica",
			Username: initialAdminUsername,
			Password: "admin",
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(initialAdmin.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Fatal("Gagal meng-hash password admin awal")
		}
		initialAdmin.Password = string(hashedPassword)

		if err := db.Create(&initialAdmin).Error; err != nil {
			logrus.Fatal("Gagal membuat admin awal")
		}
	}
}
