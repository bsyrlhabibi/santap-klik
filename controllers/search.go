package controllers

import (
	"net/http"
	"santapKlik/models"

	"github.com/labstack/echo/v4"
)

type SearchController struct {
	SearchModel *models.SearchModel
}

func NewSearchController(sm *models.SearchModel) *SearchController {
	return &SearchController{SearchModel: sm}
}

func (sc *SearchController) SearchByKeyword(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Parameter 'keyword' harus diisi"})
	}

	result, err := sc.SearchModel.SearchByKeyword(keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (sc *SearchController) SearchJajananByKeyword(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Parameter 'keyword' harus diisi"})
	}

	jajananResults, err := sc.SearchModel.SearchJajananByKeyword(keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, jajananResults)
}

func (sc *SearchController) SearchMakananByKeyword(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Parameter 'keyword' harus diisi"})
	}

	makananResults, err := sc.SearchModel.SearchMakananByKeyword(keyword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, makananResults)
}
