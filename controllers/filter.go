package controllers

import (
	"net/http"
	"santapKlik/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FilterController struct {
	JajananModel *models.JajananModel
	MakananModel *models.MakananModel
}

func NewFilterController(jm *models.JajananModel, mm *models.MakananModel) *FilterController {
	return &FilterController{JajananModel: jm, MakananModel: mm}
}

func (fc *FilterController) FilterByPrice(c echo.Context) error {
	jenis := c.FormValue("jenis")
	minPrice, errMin := strconv.Atoi(c.FormValue("minPrice"))
	maxPrice, errMax := strconv.Atoi(c.FormValue("maxPrice"))

	if errMin != nil || errMax != nil {
		return c.JSON(http.StatusBadRequest, "Invalid minPrice or maxPrice")
	}

	if jenis == "jajanan" {
		filteredJajanan, err := fc.JajananModel.FilterJajananByPrice(minPrice, maxPrice)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, filteredJajanan)
	} else if jenis == "makanan" {
		filteredMakanan, err := fc.MakananModel.FilterMakananByPrice(minPrice, maxPrice)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if len(filteredMakanan) == 0 {
			return c.JSON(http.StatusNotFound, "Produk tidak ditemukan")
		}
		return c.JSON(http.StatusOK, filteredMakanan)
	} else {
		return c.JSON(http.StatusBadRequest, "Invalid jenis")
	}
}
