package models

import (
	"time"
)

// Order представляет заказ
type Order struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Comment   string    `json:"comment" db:"comment"`
	Status    string    `json:"status" db:"status"`
	TotalCost float64   `json:"total_cost" db:"total_cost"`
	Language  string    `json:"language" db:"language"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Связанные данные
	Items []OrderItem `json:"items" db:"-"`
}

// OrderItem представляет товар в заказе
type OrderItem struct {
	ID        int64   `json:"id" db:"id"`
	OrderID   int64   `json:"-" db:"order_id"`
	ProductID int64   `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`

	// Дополнительная информация о товаре (заполняется при запросе)
	ProductName  string `json:"product_name,omitempty" db:"-"`
	ProductImage string `json:"product_image,omitempty" db:"-"`
}

// OrderRequest представляет запрос на создание заказа
type OrderRequest struct {
	Name     string      `json:"name" binding:"required"`
	Email    string      `json:"email" binding:"required,email"`
	Phone    string      `json:"phone" binding:"required"`
	Comment  string      `json:"comment"`
	Language string      `json:"language"`
	Items    []OrderItem `json:"items,omitempty"`
}

// ContactFormRequest представляет форму обратной связи
type ContactFormRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Language string `json:"language"`
}

// OrderResponse представляет ответ после создания заказа
type OrderResponse struct {
	Success bool   `json:"success"`
	OrderID int64  `json:"order_id,omitempty"`
	Message string `json:"message,omitempty"`
}
