package presenters

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/domain/entities"
)

func ToResponse(kitchenOrder entities.KitchenOrder) dtos.KitchenOrderResponseDTO {

	status := dtos.OrderStatusDTO{
		ID:   kitchenOrder.Status.ID,
		Name: kitchenOrder.Status.Name.Value(),
	}

	return dtos.KitchenOrderResponseDTO{
		ID:        kitchenOrder.ID,
		OrderID:   kitchenOrder.OrderID,
		Slug:      kitchenOrder.Slug.Value(),
		Status:    status,
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	}
}

func ToResponseList(kitchenOrders []entities.KitchenOrder) []dtos.KitchenOrderResponseDTO {
	kitchenOrderResponse := make([]dtos.KitchenOrderResponseDTO, len(kitchenOrders))

	for i, kitchenOrder := range kitchenOrders {
		kitchenOrderResponse[i] = ToResponse(kitchenOrder)
	}

	return kitchenOrderResponse
}
