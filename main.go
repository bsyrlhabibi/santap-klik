package main

import (
	"fmt"
	"santapKlik/configs"
	"santapKlik/controllers"
	"santapKlik/initialization"
	"santapKlik/models"
	"santapKlik/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	config := configs.InitConfig()

	db := models.InitModel(*config)
	models.Migrate(db)

	adminModel := models.NewAdminModel(db)
	adminController := controllers.NewAdminController(db, adminModel)

	if config.InitializeAdmin {
		initialization.InitializeAdmin(db)
	}

	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())

	routes.RouteAdmin(e, adminController, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)))
}
