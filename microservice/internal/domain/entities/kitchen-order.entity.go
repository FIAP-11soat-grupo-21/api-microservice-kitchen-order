package entities

import (
	"time"

	"tech_challenge/internal/domain/exceptions"
	value_objects "tech_challenge/internal/domain/value-objects"
	identity_manager "tech_challenge/internal/shared/pkg/identity"
)

type KitchenOrder struct {
	ID        string
	OrderID   string
	StatusID  string
	Status    OrderStatus
	Slug      value_objects.Slug
	CreatedAt time.Time
	UpdatedAt *time.Time
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
		Slug:      slugValueObject,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
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
