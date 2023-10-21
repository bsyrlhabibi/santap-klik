package main

import (
	"fmt"
	configs "santapKlik/config"
	"santapKlik/controllers"
	"santapKlik/models"
	"santapKlik/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	e := echo.New()
	config := configs.InitConfig()

	db := models.InitModel(*config)
	models.Migrate(db)

	adminModel := models.NewAdminModel(db)
	adminController := controllers.NewAdminController(db, adminModel)

	// Check if the initial admin exists, if not, create it
	initialAdminUsername := "admin"
	if adminModel.Login(initialAdminUsername, "admin") == nil {
		// Admin with username "admin" doesn't exist, create it
		initialAdmin := models.Admin{
			Name:     "Monica",
			Username: initialAdminUsername,
			Password: "admin", // You should hash the password here
		}

		// Hash the password (use bcrypt or a secure password hashing library)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(initialAdmin.Password), bcrypt.DefaultCost)
		if err != nil {
			// Handle error
			logrus.Fatal("Failed to hash initial admin password")
		}
		initialAdmin.Password = string(hashedPassword)

		// Create the initial admin in the database
		if err := db.Create(&initialAdmin).Error; err != nil {
			// Handle error
			logrus.Fatal("Failed to create initial admin")
		}
	}

	// Middleware Logger untuk logging permintaan HTTP
	e.Use(middleware.Logger())

	// Middleware Recovery untuk pemulihan jika terjadi panic
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())

	routes.RouteAdmin(e, adminController, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)))
}
