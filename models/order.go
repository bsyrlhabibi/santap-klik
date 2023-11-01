package models

import (
	"santapKlik/helper"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id"` // Kolom kunci asing
	User            User           `json:"user" gorm:"foreignKey:UserID"`
	JajananID       *uint          `json:"jajanan_id"` // Kolom kunci asing
	Jajanan         Jajanan        `json:"jajanan" gorm:"foreignKey:JajananID"`
	MakananID       *uint          `json:"makanan_id"`
	Makanan         Makanan        `json:"makanan" gorm:"foreignKey:MakananID"`
	UserUsername    string         `json:"user_username"` // Nama pengguna (username)
	JajananName     string         `json:"jajanan_name"`
	MakananName     string         `json:"makanan_name"`
	JajananQuantity int            `json:"jajanan_quantity"`
	MakananQuantity int            `json:"makanan_quantity"`
	JajananPrice    int            `json:"jajanan_price"`
	MakananPrice    int            `json:"makanan_price"`
	Quantity        int            `json:"quantity"`
	TotalPrice      float64        `json:"total_price"`
	EstimatedTime   string         `json:"estimated_time"`
	IsPaid          bool           `json:"is_paid"`
	IsReady         bool           `json:"is_ready"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type OrderModel struct {
	db *gorm.DB
}

func NewOrderModel(db *gorm.DB) *OrderModel {
	return &OrderModel{db}
}

func (om *OrderModel) CreateOrderByUsernameAndJajananName(username string, jajananName string, jajananQuantity int, estimatedTime time.Time) (*Order, error) {
	var user User
	if err := om.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	var jajanan Jajanan
	if err := om.db.Where("name = ?", jajananName).First(&jajanan).Error; err != nil {
		return nil, err
	}

	estimatedTimeString := estimatedTime.Format("2006-01-02 15:04:05")
	order := &Order{
		UserID:          user.ID,
		UserUsername:    username,
		JajananName:     jajananName,
		JajananPrice:    jajanan.Price,
		JajananQuantity: jajananQuantity,
		Quantity:        jajananQuantity,
		TotalPrice:      float64(jajananQuantity) * float64(jajanan.Price),
		EstimatedTime:   estimatedTimeString,
		IsPaid:          false,
		IsReady:         false,
	}

	if jajanan.ID != 0 {
		order.JajananID = &jajanan.ID
	}

	if err := om.db.Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (om *OrderModel) CreateOrderByUsernameAndMakananName(username string, makananName string, makananQuantity int, estimatedTime time.Time) (*Order, error) {
	var user User
	if err := om.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	var makanan Makanan
	if err := om.db.Where("name = ?", makananName).First(&makanan).Error; err != nil {
		return nil, err
	}

	estimatedTimeString := estimatedTime.Format("2006-01-02 15:04:05")

	order := &Order{
		UserID:          user.ID,
		UserUsername:    username,
		MakananName:     makananName,
		MakananPrice:    makanan.Price,
		MakananQuantity: makananQuantity,
		Quantity:        makananQuantity,
		TotalPrice:      float64(makananQuantity) * float64(makanan.Price),
		EstimatedTime:   estimatedTimeString,
		IsPaid:          false,
		IsReady:         false,
	}

	if makanan.ID != 0 {
		order.MakananID = &makanan.ID
	}

	if err := om.db.Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (om *OrderModel) CreateOrderByUsernameAndJajananAndMakananName(username string, jajananName string, makananName string, jajananQuantity int, makananQuantity int, estimatedTime time.Time) (*Order, error) {

	var user User
	if err := om.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	var jajanan Jajanan
	if err := om.db.Where("name = ?", jajananName).First(&jajanan).Error; err != nil {
		return nil, err
	}

	var makanan Makanan
	if err := om.db.Where("name = ?", makananName).First(&makanan).Error; err != nil {
		return nil, err
	}

	estimatedTimeString := estimatedTime.Format("2006-01-02 15:04:05")

	jajananTotalPrice := float64(jajananQuantity) * float64(jajanan.Price)
	makananTotalPrice := float64(makananQuantity) * float64(makanan.Price)
	totalPrice := jajananTotalPrice + makananTotalPrice

	order := &Order{
		UserID:          user.ID,
		UserUsername:    username,
		JajananName:     jajananName,
		MakananName:     makananName,
		JajananPrice:    jajanan.Price,
		MakananPrice:    makanan.Price,
		JajananQuantity: jajananQuantity,
		MakananQuantity: makananQuantity,
		Quantity:        jajananQuantity + makananQuantity,
		TotalPrice:      totalPrice,
		EstimatedTime:   estimatedTimeString,
		IsPaid:          false,
		IsReady:         false,
	}

	if jajanan.ID != 0 {
		order.JajananID = &jajanan.ID
	}

	if makanan.ID != 0 {
		order.MakananID = &makanan.ID
	}

	if err := om.db.Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (om *OrderModel) GetOrderById(orderID uint) (*Order, error) {
	order := new(Order)
	if err := om.db.Where("id = ?", orderID).First(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (om *OrderModel) UpdateOrderEstimatedTime(orderID uint, estimatedTime time.Time) (*Order, error) {
	if err := om.db.Model(&Order{}).Where("id = ?", orderID).Update("estimated_time", estimatedTime).Error; err != nil {
		return nil, err
	}

	estimatedTimeString := estimatedTime.Format("2006-01-02 15:04:05")

	order := &Order{
		ID:            orderID,
		EstimatedTime: estimatedTimeString,
		IsPaid:        false,
		IsReady:       false,
	}

	return order, nil
}

func (mm *OrderModel) GetPageOrder(page, limit int) ([]Order, error) {
	var orderList []Order
	offset := (page - 1) * limit
	if err := mm.db.Offset(offset).Limit(limit).Find(&orderList).Error; err != nil {
		return nil, err
	}
	return orderList, nil
}

func CreateOrderJajananResponse(order *Order) helper.CreateOrderResponse {
	response := helper.CreateOrderResponse{
		Message: "berhasil melakukan order",
		Data: helper.OrderResponse{
			ID:              order.ID,
			UserID:          order.UserID,
			UserUsername:    order.UserUsername,
			TotalPrice:      order.TotalPrice,
			EstimatedTime:   order.EstimatedTime,
			MakananName:     order.MakananName,
			MakananQuantity: order.MakananQuantity,
			MakananPrice:    order.MakananPrice,
			Quantity:        order.Quantity,
			IsPaid:          order.IsPaid,
			IsReady:         order.IsReady,
			CreatedAt:       order.CreatedAt,
			JajananName:     order.JajananName,
			JajananQuantity: order.JajananQuantity,
			JajananPrice:    order.JajananPrice,
			UpdatedAt:       order.UpdatedAt,
		},
	}

	if order.JajananID != nil {
		response.Data.JajananID = order.JajananID
		response.Data.JajananName = order.JajananName
		response.Data.JajananQuantity = order.JajananQuantity
		response.Data.JajananPrice = order.JajananPrice
	} else if order.MakananID != nil {
		response.Data.MakananID = order.MakananID
		response.Data.MakananName = order.MakananName
		response.Data.MakananQuantity = order.MakananQuantity
		response.Data.MakananPrice = order.MakananPrice
	}

	return response
}

func CreateOrderResponse(order *Order) helper.OrderResponse {
	// Buat logika konversi dari models.Order ke helper.OrderResponse di sini
	// Contoh:
	return helper.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		UserUsername:    order.UserUsername,
		TotalPrice:      order.TotalPrice,
		EstimatedTime:   order.EstimatedTime,
		MakananName:     order.MakananName,
		MakananQuantity: order.MakananQuantity,
		MakananPrice:    order.MakananPrice,
		Quantity:        order.Quantity,
		IsPaid:          order.IsPaid,
		IsReady:         order.IsReady,
		CreatedAt:       order.CreatedAt,
		JajananName:     order.JajananName,
		JajananQuantity: order.JajananQuantity,
		JajananPrice:    order.JajananPrice,
		UpdatedAt:       order.UpdatedAt,
	}
}
