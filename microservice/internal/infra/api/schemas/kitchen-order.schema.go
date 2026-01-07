package schemas

import (
	"time"
)

type UpdateKitchenOrderSchema struct {
	StatusID string `json:"status_id" binding:"required"`
}

type OrderItemResponseSchema struct {
	ID        string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrderID   string  `json:"order_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ProductID string  `json:"product_id" example:"product-123"`
	Quantity  int     `json:"quantity" example:"2"`
	UnitPrice float64 `json:"unit_price" example:"25.90"`
}

type KitchenOrderResponseSchema struct {
	ID         string                    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrderID    string                    `json:"order_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CustomerID *string                   `json:"customer_id" example:"customer-123"`
	Amount     float64                   `json:"amount" example:"51.80"`
	Slug       string                    `json:"slug" example:"001"`
	Status     string                    `json:"status" example:"Pronto"`
	Items      []OrderItemResponseSchema `json:"items"`
	CreatedAt  time.Time                 `json:"created_at" example:"2023-10-01T12:00:00Z"`
	UpdatedAt  *time.Time                `json:"updated_at" example:"2023-10-01T12:00:00Z"`
}

type KitchenOrderNotFoundErrorSchema struct {
	Error string `json:"error" example:"Kitchen order not found"`
}

type InvalidKitchenOrderDataErrorSchema struct {
	Error string `json:"error" example:"Invalid kitchen order data"`
}
