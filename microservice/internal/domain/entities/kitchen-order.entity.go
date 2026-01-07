package entities

import (
	"time"

	"tech_challenge/internal/domain/exceptions"
	value_objects "tech_challenge/internal/domain/value-objects"
	identity_manager "tech_challenge/internal/shared/pkg/identity"
)

type KitchenOrder struct {
	ID         string
	OrderID    string
	CustomerID *string
	Amount     float64
	StatusID   string
	Status     OrderStatus
	Slug       value_objects.Slug
	Items      []OrderItem
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewKitchenOrder(id, orderID, slug string, status OrderStatus, createdAt time.Time, updatedAt *time.Time) (*KitchenOrder, error) {
	slugValueObject, err := value_objects.NewSlug(slug)
	if err != nil {
		return nil, err
	}

	return &KitchenOrder{
		ID:        id,
		OrderID:   orderID,
		Status:    status,
		StatusID:  status.ID,
		Slug:      slugValueObject,
		Items:     []OrderItem{},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func NewKitchenOrderWithOrderData(id, orderID, slug string, customerID *string, amount float64, items []OrderItem, status OrderStatus, createdAt time.Time, updatedAt *time.Time) (*KitchenOrder, error) {
	slugValueObject, err := value_objects.NewSlug(slug)
	if err != nil {
		return nil, err
	}

	return &KitchenOrder{
		ID:         id,
		OrderID:    orderID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     status,
		StatusID:   status.ID,
		Slug:       slugValueObject,
		Items:      items,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func ValidateID(id string) error {
	if !identity_manager.IsValidUUID(id) {
		return &exceptions.InvalidKitchenOrderDataException{
			Message: "Invalid kitchen Order ID",
		}
	}

	return nil
}

func (c *KitchenOrder) IsEmpty() bool {
	return c.ID == ""
}

func (c *KitchenOrder) AddItem(item OrderItem) {
	c.Items = append(c.Items, item)
}

func (c *KitchenOrder) CalcTotalAmount() {
	total := 0.0
	for _, item := range c.Items {
		total += item.GetTotal()
	}
	c.Amount = total
}
