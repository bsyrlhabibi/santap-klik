package helper

import (
	"time"

	"gorm.io/gorm"
)

func FormatResponse(message string, data any) map[string]any {
	var response = map[string]any{}
	response["message"] = message
	if data != nil {
		response["data"] = data
	}
	return response
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RegistrationResponse struct {
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type UserData struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Alamat   string `json:"alamat"`
	HP       string `json:"hp"`
	Password string `json:"password"`
}

type JajananData struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       int            `json:"price"`
	Image       string         `json:"image"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}

type MakananData struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       int            `json:"price"`
	Image       string         `json:"image"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}

type OrderResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	UserUsername string `json:"user_username"`
	// Tambahkan bidang sesuai dengan jenis pesanan
	JajananID       *uint  `json:"jajanan_id,omitempty"`
	JajananName     string `json:"jajanan_name,omitempty"`
	MakananID       *uint  `json:"makanan_id,omitempty"`
	MakananName     string `json:"makanan_name,omitempty"`
	JajananPrice    int    `json:"jajanan_price,omitempty"`
	MakananPrice    int    `json:"makanan_price,omitempty"`
	JajananQuantity int    `json:"jajanan_quantity,omitempty"`
	MakananQuantity int    `json:"makanan_quantity,omitempty"`

	Quantity      int       `json:"quantity"`
	TotalPrice    float64   `json:"total_price"`
	EstimatedTime string    `json:"estimated_time"`
	IsPaid        bool      `json:"is_paid"`
	IsReady       bool      `json:"is_ready"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateOrderResponse struct {
	Message string        `json:"berhasil melakukan order"`
	Data    OrderResponse `json:"data"`
}

type UpdateOrderResponse struct {
	Message string        `json:"berhasil update order"`
	Data    OrderResponse `json:"data"`
}

type OrderPageResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
}
