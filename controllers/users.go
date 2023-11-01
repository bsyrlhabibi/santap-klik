package controllers

import (
	"net/http"
	"santapKlik/configs"
	"santapKlik/helper"
	"santapKlik/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	db   *gorm.DB
	user *models.UserModel
	cfg  *configs.ProgramConfig
}

func NewUserController(db *gorm.DB, userModel *models.UserModel, cfg *configs.ProgramConfig) *UserController {
	return &UserController{
		db:   db,
		user: userModel,
		cfg:  cfg,
	}
}

func (uc *UserController) LoginUser(c echo.Context) error {
	var userLogin models.UserLoginModel
	if err := c.Bind(&userLogin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := uc.user.LoginUser(userLogin.Username, userLogin.Password)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Nama pengguna atau kata sandi tidak valid"})
	}

	if userLogin.Username == "" || userLogin.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Silakan isi semua field yang diperlukan"})
	}
	// Set the user's role as "user"
	role := "user"

	token, err := helper.CreateToken(int(user.ID), user.Username, role, uc.cfg) // Use uc.cfg
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Gagal menghasilkan token JWT",
		})
	}

	response := struct {
		Message string `json:"message"`
		Data    struct {
			Token string          `json:"token"`
			User  helper.UserData `json:"user"`
		} `json:"data"`
	}{
		Message: "Login berhasil",
		Data: struct {
			Token string          `json:"token"`
			User  helper.UserData `json:"user"`
		}{
			Token: token,
			User: helper.UserData{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Alamat:   user.Alamat,
				HP:       user.HP,
			},
		},
	}

	return c.JSON(http.StatusOK, response)

}

func (uc *UserController) Register(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Validate required fields
	if user.Name == "" || user.Username == "" || user.Password == "" || user.Alamat == "" || user.HP == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Silakan isi semua field yang diperlukan"})
	}

	if err := user.ValidatePhoneNumber(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := uc.user.RegisterUser(user.Name, user.Username, user.Alamat, user.HP, user.Password) // Use uc.user
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	registeredUser := uc.user.GetUserByUsername(uc.db, user.Username)

	userData := helper.UserData{
		ID:       registeredUser.ID,
		Name:     user.Name,
		Username: user.Username,
		Alamat:   user.Alamat,
		HP:       user.HP,
	}

	response := struct {
		Message string          `json:"message"`
		Data    helper.UserData `json:"data"`
	}{
		Message: "Registration successful",
		Data:    userData,
	}

	return c.JSON(http.StatusOK, response)
}

func (uc *UserController) GetAllUsers(c echo.Context) error {
	users, err := uc.user.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil daftar pengguna"})
	}
	return c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUserByID(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := uc.user.GetUserByID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Pengguna tidak ditemukan"})
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUser(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var userUpdate models.UserUpdateModel
	if err := c.Bind(&userUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = uc.user.UpdateUser(uint(userID), userUpdate.Name, userUpdate.Username, userUpdate.Alamat, userUpdate.HP, userUpdate.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui pengguna"})
	}

	updatedUser, _ := uc.user.GetUserByID(uint(userID))

	userData := helper.UserData{
		ID:       updatedUser.ID,
		Name:     updatedUser.Name,
		Username: updatedUser.Username,
		Alamat:   updatedUser.Alamat,
		HP:       updatedUser.HP,
		Password: updatedUser.Password,
	}

	return c.JSON(http.StatusOK, userData)
}

func (uc *UserController) DeleteUser(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	err = uc.user.DeleteUser(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus pengguna"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Pengguna berhasil dihapus"})
}

func (uc *UserController) GetUserByUsername(c echo.Context) error {
	inputUsername := c.QueryParam("username")

	if inputUsername == "" {
		return c.JSON(http.StatusBadRequest, "Username tidak boleh kosong")
	}

	user := uc.user.GetUserByUsername(uc.db, inputUsername)
	if user == nil {
		return c.JSON(http.StatusNotFound, "Username not found")
	}

	userData := helper.UserData{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Alamat:   user.Alamat,
		HP:       user.HP,
	}

	return c.JSON(http.StatusOK, userData)
}

func (uc *UserController) GetUserByPage(c echo.Context) error {
	page, errPage := strconv.Atoi(c.QueryParam("page"))
	if errPage != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page value")
	}

	limit, errLimit := strconv.Atoi(c.QueryParam("pageSize"))
	if errLimit != nil {
		return c.JSON(http.StatusBadRequest, "Invalid pageSize value")
	}

	userList, err := uc.user.GetPageUser(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, userList)
}
