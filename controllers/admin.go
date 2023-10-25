package controllers

import (
	"net/http"
	"santapKlik/helper"
	"santapKlik/initialization"
	"santapKlik/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var adminInitialized bool

type AdminController struct {
	db    *gorm.DB
	admin *models.AdminModel
}

func NewAdminController(db *gorm.DB, adminModel *models.AdminModel) *AdminController {
	return &AdminController{
		db:    db,
		admin: adminModel,
	}
}

func (ac *AdminController) LoginAdmin(c echo.Context) error {
	var loginRequest models.LoginModel
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if !adminInitialized {
		initialization.InitializeAdmin(ac.db)
		adminInitialized = true
	}

	admin := ac.admin.Login(loginRequest.Username, loginRequest.Password)
	if admin == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login failed"})
	}

	token, err := helper.GenerateAdminJWT("your_secret_key", admin.Username) // Ganti "your_secret_key" sesuai kebutuhan Anda
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate JWT"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Login successful", "token": token})
}
