package main

import (
	"fmt"
	"log"
	"santapKlik/configs"
	"santapKlik/controllers"
	"santapKlik/models"
	"santapKlik/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := configs.InitConfig()

	db := models.InitModel(*config)
	models.Migrate(db)

	adminModel := models.NewAdminModel(db)
	adminController := controllers.NewAdminController(db, adminModel, config)

	jajananModel := models.NewJajananModel(db)
	jajananController := controllers.NewJajananController(jajananModel)

	makananModel := models.NewMakananModel(db)
	makananController := controllers.NewMakananController(makananModel)

	userModel := models.NewUserModel(db)
	userController := controllers.NewUserController(db, userModel, config)

	searchModel := models.NewSearchModel(db)
	searchController := controllers.NewSearchController(searchModel)

	filterController := controllers.NewFilterController(jajananModel, makananModel)

	orderModel := models.NewOrderModel(db)
	orderController := controllers.NewOrderController(orderModel, jajananModel, makananModel, userModel)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.CORS())

	routes.RouteAdmin(e, adminController, jajananController, orderController, makananController, config)
	routes.RouteJajanan(e, jajananController, *config)
	routes.RouteMakanan(e, makananController, *config)
	routes.RouteUser(e, userController, *config)
	routes.RouteFitur(e, searchController, filterController, *config)
	routes.RouteOrder(e, orderController, config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)))
}
