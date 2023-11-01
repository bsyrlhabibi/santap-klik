package controllers

import (
	"net/http"
	"santapKlik/helper"
	"santapKlik/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type JajananController struct {
	JajananModel *models.JajananModel
}

func NewJajananController(jm *models.JajananModel) *JajananController {
	return &JajananController{JajananModel: jm}
}

func (jc *JajananController) CreateJajanan(c echo.Context) error {
	var jajanan models.Jajanan

	if err := c.Bind(&jajanan); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Image file is required")
	}

	file, err := imageFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to open the image file")
	}
	defer file.Close()

	jajananID, err := jc.JajananModel.InsertJajanan(jajanan, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	createdJajanan, err := jc.JajananModel.GetJajananByID(uint(jajananID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, createdJajanan)
}

func (jc *JajananController) UpdateJajanan(c echo.Context) error {
	jajananIDStr := c.Param("id")

	jajananID, err := strconv.ParseUint(jajananIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	var updatedJajanan models.Jajanan
	if err := c.Bind(&updatedJajanan); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	updatedJajanan.ID = uint(jajananID)

	imageFile, err := c.FormFile("image")
	if err == nil {
		updatedJajanan.Image = "jajanan/" + imageFile.Filename
	}

	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, "Failed to get the image file")
	}

	existingJajanan, err := jc.JajananModel.GetJajananByID(uint(jajananID))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Jajanan not found")
	}

	existingJajananData := helper.JajananData{
		ID:          existingJajanan.ID,
		Name:        existingJajanan.Name,
		Description: existingJajanan.Description,
		Price:       existingJajanan.Price,
		Image:       existingJajanan.Image,
		CreatedAt:   existingJajanan.CreatedAt,
		UpdatedAt:   existingJajanan.UpdatedAt,
		DeletedAt:   existingJajanan.DeletedAt,
	}

	if imageFile != nil {

		file, err := imageFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to open the image file")
		}
		defer file.Close()

		err = jc.JajananModel.UpdateJajanan(uint(jajananID), updatedJajanan, file)
	} else {
		err = jc.JajananModel.UpdateWithoutImageJajanan(uint(jajananID), updatedJajanan)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingJajananData)
}

func (jc *JajananController) GetJajananByID(c echo.Context) error {
	jajananIDStr := c.Param("id")

	jajananID, err := strconv.ParseUint(jajananIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	jajanan, err := jc.JajananModel.GetJajananByID(uint(jajananID))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Jajanan not found")
	}

	return c.JSON(http.StatusOK, jajanan)
}

func (jc *JajananController) DeleteJajanan(c echo.Context) error {
	jajananIDStr := c.Param("id")
	jajananID, err := strconv.Atoi(jajananIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Jajanan ID")
	}

	if err := jc.JajananModel.DeleteJajanan(uint(jajananID)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Jajanan deleted successfully")
}

func (jc *JajananController) GetAllJajanan(c echo.Context) error {
	jajananList, err := jc.JajananModel.GetAllJajanan()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, jajananList)
}

func (jc *JajananController) GetJajananByPage(c echo.Context) error {
	page, errPage := strconv.Atoi(c.QueryParam("page"))
	if errPage != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page value")
	}

	limit, errLimit := strconv.Atoi(c.QueryParam("pageSize"))
	if errLimit != nil {
		return c.JSON(http.StatusBadRequest, "Invalid pageSize value")
	}

	jajananList, err := jc.JajananModel.GetPageJajanan(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, jajananList)
}
