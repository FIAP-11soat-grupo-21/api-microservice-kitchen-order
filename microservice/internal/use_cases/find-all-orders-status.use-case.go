package use_cases

import (
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/domain/entities"
)

type FindAllOrderStatusUseCase struct {
	gateway gateways.OrderStatusGateway
}

func NewFindAllOrdersStatusUseCase(gateway gateways.OrderStatusGateway) *FindAllOrderStatusUseCase {
	return &FindAllOrderStatusUseCase{
		gateway: gateway,
	}
}

func (uc *FindAllOrderStatusUseCase) Execute() ([]entities.OrderStatus, error) {
	statusList, err := uc.gateway.FindAll()

	if err != nil {
		return nil, err
	}

	return statusList, nil
}
