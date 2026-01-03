package use_cases

import (
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
)

type FindKitchenOrderByIDUseCase struct {
	gateway gateways.KitchenOrderGateway
}

func NewFindKitchenOrderByIDUseCase(gateway gateways.KitchenOrderGateway) *FindKitchenOrderByIDUseCase {
	return &FindKitchenOrderByIDUseCase{
		gateway: gateway,
	}
}

func (uc *FindKitchenOrderByIDUseCase) Execute(id string) (entities.KitchenOrder, error) {
	err := entities.ValidateID(id)

	if err != nil {
		return entities.KitchenOrder{}, err
	}

	kitchenOrder, err := uc.gateway.FindByID(id)

	if err != nil || kitchenOrder.IsEmpty() {
		return entities.KitchenOrder{}, &exceptions.KitchenOrderNotFoundException{}
	}

	return kitchenOrder, nil
}
