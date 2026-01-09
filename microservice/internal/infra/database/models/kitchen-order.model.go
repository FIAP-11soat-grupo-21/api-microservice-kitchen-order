package models

import "time"

type OrderItemModel struct {
	ID              string  `gorm:"primaryKey; size:36"`
	KitchenOrderID  string  `gorm:"not null;size:36;index"`
	OrderID         string  `gorm:"not null;size:36;index"`
	ProductID       string  `gorm:"not null;size:100"`
	Quantity        int     `gorm:"not null"`
	UnitPrice       float64 `gorm:"not null;type:decimal(10,2)"`
}

func (OrderItemModel) TableName() string {
	return "order_item"
}

type KitchenOrderModel struct {
	ID         string           `gorm:"primaryKey; size:36"`
	OrderID    string           `gorm:"not null;size:36;index"`
	CustomerID *string          `gorm:"size:36"`
	Amount     float64          `gorm:"not null;type:decimal(10,2)"`
	Slug       string           `gorm:"not null;size:100;"`
	StatusID   string           `json:"statusId" gorm:"not null; size:36; index"`
	Status     OrderStatusModel `json:"status" gorm:"foreignKey:StatusID;references:ID"`
	Items      []OrderItemModel `gorm:"foreignKey:KitchenOrderID;references:ID"`

	CreatedAt time.Time  `gorm:"not null; index"`
	UpdatedAt *time.Time `gorm:""`
}

func (KitchenOrderModel) TableName() string {
	return "kitchen_order"
}
