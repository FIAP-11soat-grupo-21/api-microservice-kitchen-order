package dtos

import "time"

type KitchenOrderDTO struct {
	ID      string
	OrderID string
	Slug    string
	Status  OrderStatusDTO
}

type CreateKitchenOrderDTO struct {
	OrderID string
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
	ID        string
	OrderID   string
	Slug      string
	Status    OrderStatusDTO
	CreatedAt time.Time
	UpdatedAt *time.Time
}
