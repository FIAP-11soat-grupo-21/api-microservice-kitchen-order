package dtos

import "time"

type KitchenOrderDTO struct {
	ID         string
	OrderID    string
	CustomerID *string
	Amount     float64
	Slug       string
	Status     OrderStatusDTO
	Items      []OrderItemDTO
}

type OrderItemDTO struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
}

type CreateKitchenOrderDTO struct {
	OrderID    string
	CustomerID *string
	Amount     float64
	Items      []OrderItemDTO
}

type UpdateKitchenOrderDTO struct {
	ID       string
	StatusID string
}

type KitchenOrderFilter struct {
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
	StatusID      *uint
}

type KitchenOrderResponseDTO struct {
	ID         string
	OrderID    string
	CustomerID *string
	Amount     float64
	Slug       string
	Status     OrderStatusDTO
	Items      []OrderItemDTO
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
