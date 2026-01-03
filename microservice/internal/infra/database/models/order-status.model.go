package models

type OrderStatusModel struct {
	ID   string `gorm:"primaryKey; size:36"`
	Name string `gorm:"not null;size:100;"`
}

func (OrderStatusModel) TableName() string {
	return "order_status"
}
