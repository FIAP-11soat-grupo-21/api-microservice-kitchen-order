package use_cases

import (
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
)

type UpdateKitchenOrderUseCase struct {
	gateway       gateways.KitchenOrderGateway
	statusGateway gateways.OrderStatusGateway
}

func NewUpdateKitchenOrderUseCase(gateway gateways.KitchenOrderGateway, statusGateway gateways.OrderStatusGateway) *UpdateKitchenOrderUseCase {
	return &UpdateKitchenOrderUseCase{
		gateway:       gateway,
		statusGateway: statusGateway,
	}
}

func (ko *UpdateKitchenOrderUseCase) Execute(kitchenOrderDTO dtos.UpdateKitchenOrderDTO) (entities.KitchenOrder, error) {
	err := entities.ValidateID(kitchenOrderDTO.ID)

	if err != nil {
		return entities.KitchenOrder{}, err
	}

	kitchenOrder, err := ko.gateway.FindByID(kitchenOrderDTO.ID)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.KitchenOrderNotFoundException{}
	}

	kitchenOrderStatus, err := ko.statusGateway.FindByID(kitchenOrderDTO.StatusID)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.OrderStatusNotFoundException{}
	}

	kitchenOrder.Status.ID = kitchenOrderDTO.StatusID
	kitchenOrder.Status.Name = kitchenOrderStatus.Name

	now := time.Now()
	kitchenOrder.UpdatedAt = &now

	err = ko.gateway.Update(kitchenOrder)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.InvalidKitchenOrderDataException{}
	}

	return kitchenOrder, nil
}
