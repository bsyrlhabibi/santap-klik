package controllers

import (
	"net/http"
	"santapKlik/helper"
	"santapKlik/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MakananController struct {
	MakananModel *models.MakananModel
}

func NewMakananController(mm *models.MakananModel) *MakananController {
	return &MakananController{MakananModel: mm}
}

func (mc *MakananController) CreateMakanan(c echo.Context) error {
	var makanan models.Makanan

	// Bind the request body to the Jajanan model
	if err := c.Bind(&makanan); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Get the uploaded image file
	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Image file is required")
	}

	// Open the image file
	file, err := imageFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to open the image file")
	}
	defer file.Close()

	makananID, err := mc.MakananModel.InsertMakanan(makanan, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	createdMakanan, err := mc.MakananModel.GetMakananByID(uint(makananID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, createdMakanan)
}

func (mc *MakananController) UpdateMakanan(c echo.Context) error {
	makananIDStr := c.Param("id")

	makananID, err := strconv.ParseUint(makananIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	var updatedMakanan models.Makanan
	if err := c.Bind(&updatedMakanan); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	updatedMakanan.ID = uint(makananID)

	imageFile, err := c.FormFile("image")
	if err == nil {
		updatedMakanan.Image = "makanan/" + imageFile.Filename
	}

	if err != nil && err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, "Failed to get the image file")
	}

	existingMakanan, err := mc.MakananModel.GetMakananByID(uint(makananID))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Jajanan not found")
	}

	existingMakananData := helper.MakananData{
		ID:          existingMakanan.ID,
		Name:        existingMakanan.Name,
		Description: existingMakanan.Description,
		Price:       existingMakanan.Price,
		Image:       existingMakanan.Image,
		CreatedAt:   existingMakanan.CreatedAt,
		UpdatedAt:   existingMakanan.UpdatedAt,
		DeletedAt:   existingMakanan.DeletedAt,
	}

	if imageFile != nil {

		file, err := imageFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to open the image file")
		}
		defer file.Close()

		err = mc.MakananModel.UpdateMakanan(uint(makananID), updatedMakanan, file)
	} else {
		err = mc.MakananModel.UpdateWithoutImageMakanan(uint(makananID), updatedMakanan)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingMakananData)
}

func (mc *MakananController) GetMakananByID(c echo.Context) error {
	makananIDStr := c.Param("id")

	makananID, err := strconv.ParseUint(makananIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	makanan, err := mc.MakananModel.GetMakananByID(uint(makananID))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Makanan not found")
	}

	return c.JSON(http.StatusOK, makanan)
}

func (mc *MakananController) DeleteMakanan(c echo.Context) error {
	makananIDStr := c.Param("id")
	makananID, err := strconv.Atoi(makananIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Makanan ID")
	}

	if err := mc.MakananModel.DeleteMakanan(uint(makananID)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Makanan deleted successfully")
}

func (mc *MakananController) GetAllMakanan(c echo.Context) error {
	makananList, err := mc.MakananModel.GetAllMakanan()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, makananList)
}

func (jc *MakananController) GetMakananByPage(c echo.Context) error {
	page, errPage := strconv.Atoi(c.QueryParam("page"))
	if errPage != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page value")
	}

	limit, errLimit := strconv.Atoi(c.QueryParam("pageSize"))
	if errLimit != nil {
		return c.JSON(http.StatusBadRequest, "Invalid pageSize value")
	}

	makananList, err := jc.MakananModel.GetPageMakanan(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, makananList)
}
