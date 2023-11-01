package initialization

import (
	"santapKlik/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitializeAdmin(db *gorm.DB) {
	admins := []models.Admin{}
	if err := db.Find(&admins).Error; err != nil {
		logrus.Fatal("Gagal mengambil data admin")
		return
	}

	for i := range admins {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admins[i].Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Fatalf("Gagal meng-hash password admin %s", admins[i].Username)
			return
		}

		admins[i].Password = string(hashedPassword)

		if err := db.Save(&admins[i]).Error; err != nil {
			logrus.Fatalf("Gagal menyimpan data admin %s", admins[i].Username)
			return
		}
	}
}
