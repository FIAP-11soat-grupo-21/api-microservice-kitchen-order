package entities

import (
	"tech_challenge/internal/domain/exceptions"
)

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
}

func NewOrderItem(id, orderID, productID string, quantity int, unitPrice float64) (*OrderItem, error) {
	if id == "" {
		return nil, &exceptions.InvalidKitchenOrderDataException{
			Message: "Order item ID is required",
		}
	}

	if orderID == "" {
		return nil, &exceptions.InvalidKitchenOrderDataException{
			Message: "Order ID is required",
		}
	}

	if productID == "" {
		return nil, &exceptions.InvalidKitchenOrderDataException{
			Message: "Product ID is required",
		}
	}

	if quantity <= 0 {
		return nil, &exceptions.InvalidKitchenOrderDataException{
			Message: "Quantity must be greater than zero",
		}
	}

	if unitPrice < 0 {
		return nil, &exceptions.InvalidKitchenOrderDataException{
			Message: "Unit price cannot be negative",
		}
	}

	return &OrderItem{
		ID:        id,
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}, nil
}

func (oi *OrderItem) GetTotal() float64 {
	return oi.UnitPrice * float64(oi.Quantity)
}

func (oi *OrderItem) IsEmpty() bool {
	return oi.ID == ""
}