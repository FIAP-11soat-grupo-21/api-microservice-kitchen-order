package use_cases

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/domain/entities"
)

type FindAllKitchenOrdersUseCase struct {
	gateway gateways.KitchenOrderGateway
}

func NewFindAllKitchenOrderUseCase(gateway gateways.KitchenOrderGateway) *FindAllKitchenOrdersUseCase {
	return &FindAllKitchenOrdersUseCase{
		gateway: gateway,
	}
}

func (uc *FindAllKitchenOrdersUseCase) Execute(filter dtos.KitchenOrderFilter) ([]entities.KitchenOrder, error) {
	kitchenOrders, err := uc.gateway.FindAll(filter)

	if err != nil {
		return nil, err
	}

	return kitchenOrders, nil
}
