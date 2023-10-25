package routes

import (
	"santapKlik/configs"
	"santapKlik/controllers"

	"github.com/labstack/echo/v4"
)

func RouteAdmin(e *echo.Echo, ac *controllers.AdminController, cfg configs.ProgramConfig) {
	adminGroup := e.Group("/admin")

	adminGroup.POST("/login", ac.LoginAdmin)
}
