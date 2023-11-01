package routes

import (
	"santapKlik/configs"
	"santapKlik/controllers"
	"santapKlik/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RouteAdmin(e *echo.Echo, ac *controllers.AdminController, jc *controllers.JajananController, oc *controllers.OrderController, mc *controllers.MakananController, cfg *configs.ProgramConfig) {
	adminGroup := e.Group("/admin")

	adminGroup.Use(helper.AdminMiddleware(cfg))

	adminGroup.POST("/jajanan", jc.CreateJajanan)
	adminGroup.PUT("/jajanan/:id", jc.UpdateJajanan)
	adminGroup.DELETE("/jajanan/:id", jc.DeleteJajanan)

	adminGroup.POST("/makanan", mc.CreateMakanan)
	adminGroup.PUT("/makanan/:id", mc.UpdateMakanan)
	adminGroup.DELETE("/makanan/:id", mc.DeleteMakanan)

	authGroup := e.Group("/auth")
	authGroup.POST("/admin/login", ac.LoginAdmin)
	authGroup.POST("/admin/register", ac.RegisterAdmin)
}

func RouteUser(e *echo.Echo, uc *controllers.UserController, cfg configs.ProgramConfig) {
	userGroup := e.Group("/users")

	userGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.Secret),
	}))

	userGroup.PUT("/:id", uc.UpdateUser)
	userGroup.DELETE("/:id", uc.DeleteUser)
	userGroup.GET("/:id", uc.GetUserByID)
	userGroup.GET("/all", uc.GetAllUsers)
	userGroup.GET("/get-by-username", uc.GetUserByUsername)
	userGroup.GET("", uc.GetUserByPage)

	authGroup := e.Group("/auth")

	authGroup.POST("/users/register", uc.Register)
	authGroup.POST("/users/login", uc.LoginUser)

}

func RouteJajanan(e *echo.Echo, jc *controllers.JajananController, cfg configs.ProgramConfig) {

	jajananGroup := e.Group("jajanan")
	jajananGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.Secret),
	}))

	jajananGroup.GET("/all", jc.GetAllJajanan)
	jajananGroup.GET("", jc.GetJajananByPage)
	jajananGroup.GET("/:id", jc.GetJajananByID)

}

func RouteMakanan(e *echo.Echo, mc *controllers.MakananController, cfg configs.ProgramConfig) {
	makananGroup := e.Group("/makanan")
	makananGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.Secret),
	}))
	makananGroup.GET("/all", mc.GetAllMakanan)
	makananGroup.GET("", mc.GetMakananByPage)
	makananGroup.GET("/:id", mc.GetMakananByID)

}

func RouteFitur(e *echo.Echo, sc *controllers.SearchController, fc *controllers.FilterController, cfg configs.ProgramConfig) {
	fiturGroup := e.Group("")
	fiturGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.Secret),
	}))
	fiturGroup.GET("/search", sc.SearchByKeyword)
	fiturGroup.GET("jajanan/search", sc.SearchJajananByKeyword)
	fiturGroup.GET("makanan/search", sc.SearchMakananByKeyword)

	fiturGroup.GET("filter-by-price", fc.FilterByPrice)
}

func RouteOrder(e *echo.Echo, oc *controllers.OrderController, cfg *configs.ProgramConfig) {
	orderGroup := e.Group("/orders")
	orderGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(cfg.Secret),
	}))

	orderGroup.GET("", oc.GetOrderByPage)

	createGroup := e.Group("/orders")
	createGroup.Use(helper.JWTMiddleware(cfg))
	createGroup.POST("", oc.CreateOrder)
}
