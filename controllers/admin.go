package controllers

import (
	"net/http"
	"santapKlik/configs"
	"santapKlik/helper"
	"santapKlik/initialization"
	"santapKlik/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var adminTokens = make(map[string]string)
var adminInitialized bool

type AdminController struct {
	db    *gorm.DB
	admin *models.AdminModel
	cfg   *configs.ProgramConfig
}

func NewAdminController(db *gorm.DB, adminModel *models.AdminModel, cfg *configs.ProgramConfig) *AdminController {
	return &AdminController{
		db:    db,
		admin: adminModel,
		cfg:   cfg,
	}
}

func (ac *AdminController) LoginAdmin(c echo.Context) error {
	var loginRequest models.LoginModel
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	admin := ac.admin.Login(loginRequest.Username, loginRequest.Password)
	if admin == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login failed"})
	}

	tokenString, err := helper.CreateToken(int(admin.ID), admin.Username, "admin", ac.cfg)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate JWT"})
	}

	adminTokens[loginRequest.Username] = tokenString

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Login successful", "token": tokenString})
}

func (ac *AdminController) RegisterAdmin(c echo.Context) error {
	var registerRequest struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&registerRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if !adminInitialized {
		initialization.InitializeAdmin(ac.db)
		adminInitialized = true
	}

	existingAdmin := ac.admin.GetAdminByUsername(ac.db, registerRequest.Username)
	if existingAdmin != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already in use"})
	}

	err := ac.admin.Register(registerRequest.Name, registerRequest.Username, registerRequest.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register admin"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Admin registered successfully"})
}
