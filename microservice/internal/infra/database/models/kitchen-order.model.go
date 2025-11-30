package models

import "time"

type KitchenOrderModel struct {
	ID       string           `gorm:"primaryKey; size:36"`
	OrderID  string           `gorm:"not null;size:100;"`
	Slug     string           `gorm:"not null;size:100;"`
	StatusID string           `json:"statusId" gorm:"not null; size:36; index"`
	Status   OrderStatusModel `json:"status" gorm:"foreignKey:StatusID;references:ID"`

	CreatedAt time.Time  `gorm:"not null; index"`
	UpdatedAt *time.Time `gorm:""`
}

func (KitchenOrderModel) TableName() string {
	return "kitchen_order"
}
