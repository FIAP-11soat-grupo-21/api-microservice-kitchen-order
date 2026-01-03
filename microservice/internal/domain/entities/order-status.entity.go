package entities

import value_objects "tech_challenge/internal/domain/value-objects"

type OrderStatus struct {
	ID   string
	Name value_objects.Name
}

func NewOrderStatus(id string, name string) (*OrderStatus, error) {
	nameValueObject, err := value_objects.NewName(name)
	if err != nil {
		return nil, err
	}

	return &OrderStatus{
		ID:   id,
		Name: nameValueObject,
	}, nil
}
