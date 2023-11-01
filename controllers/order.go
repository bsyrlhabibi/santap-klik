package controllers

import (
	"net/http"
	"santapKlik/helper"
	"santapKlik/models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type OrderController struct {
	orderModel   *models.OrderModel
	jajananModel *models.JajananModel
	makananModel *models.MakananModel
	userModel    *models.UserModel
}

func NewOrderController(om *models.OrderModel, jm *models.JajananModel, mm *models.MakananModel, um *models.UserModel) *OrderController {
	return &OrderController{
		orderModel:   om,
		jajananModel: jm,
		makananModel: mm,
		userModel:    um,
	}
}

func (oc *OrderController) CreateOrder(c echo.Context) error {
	var request struct {
		Username        string `json:"username"`
		JajananName     string `json:"jajanan_name"`
		MakananName     string `json:"makanan_name"`
		JajananQuantity int    `json:"jajanan_quantity"`
		MakananQuantity int    `json:"makanan_quantity"`
		EstimatedTime   string `json:"estimated_time"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	if request.Username == "" || (request.JajananName == "" && request.MakananName == "") {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Username and either JajananName or MakananName must not be empty"})
	}

	var estimatedTime time.Time

	if request.EstimatedTime != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", request.EstimatedTime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}
		estimatedTime = parsedTime
	} else {
		estimatedTime = time.Now()
	}

	var createdOrder *models.Order
	var createErr error

	if request.JajananName != "" && request.MakananName != "" {
		createdOrder, createErr = oc.orderModel.CreateOrderByUsernameAndJajananAndMakananName(request.Username, request.JajananName, request.MakananName, request.JajananQuantity, request.MakananQuantity, estimatedTime)
	} else if request.JajananName != "" {
		createdOrder, createErr = oc.orderModel.CreateOrderByUsernameAndJajananName(request.Username, request.JajananName, request.JajananQuantity, estimatedTime)
	} else {
		createdOrder, createErr = oc.orderModel.CreateOrderByUsernameAndMakananName(request.Username, request.MakananName, request.MakananQuantity, estimatedTime)
	}

	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": createErr.Error()})
	}

	response := models.CreateOrderJajananResponse(createdOrder)

	return c.JSON(http.StatusCreated, response)
}

func (oc *OrderController) GetOrderByPage(c echo.Context) error {
	page, errPage := strconv.Atoi(c.QueryParam("page"))
	if errPage != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page value")
	}

	limit, errLimit := strconv.Atoi(c.QueryParam("pageSize"))
	if errLimit != nil {
		return c.JSON(http.StatusBadRequest, "Invalid pageSize value")
	}

	// Menggunakan variabel orders dan err untuk mengambil hasil GetPageOrder
	orders, err := oc.orderModel.GetPageOrder(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mengkonversi data models.Order menjadi helper.OrderResponse
	orderResponses := []helper.OrderResponse{}
	for _, order := range orders {
		orderResponse := models.CreateOrderResponse(&order)
		orderResponses = append(orderResponses, orderResponse)
	}

	// Menghitung total pesanan
	total := len(orderResponses)

	response := helper.OrderPageResponse{
		Orders: orderResponses,
		Total:  total,
	}

	return c.JSON(http.StatusOK, response)
}
