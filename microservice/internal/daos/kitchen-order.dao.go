package daos

import "time"

type KitchenOrderDAO struct {
	ID         string
	OrderID    string
	CustomerID *string
	Amount     float64
	Status     OrderStatusDAO
	Slug       string
	Items      []OrderItemDAO
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type OrderItemDAO struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
}
