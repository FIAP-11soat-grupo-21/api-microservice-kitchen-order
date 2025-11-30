package daos

import "time"

type KitchenOrderDAO struct {
	ID        string
	OrderID   string
	Status    OrderStatusDAO
	Slug      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
