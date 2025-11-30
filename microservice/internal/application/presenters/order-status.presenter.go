package presenters

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/domain/entities"
)

func ToResponseOrderStatus(orderStatusEntity entities.OrderStatus) dtos.OrderStatusResponseDTO {
	return dtos.OrderStatusResponseDTO{
		ID:   orderStatusEntity.ID,
		Name: orderStatusEntity.Name.Value(),
	}
}

func ToResponseListOrderStatus(orderStatusEntities []entities.OrderStatus) []dtos.OrderStatusResponseDTO {
	orderStatusResponse := make([]dtos.OrderStatusResponseDTO, len(orderStatusEntities))

	for i, orderStatus := range orderStatusEntities {
		orderStatusResponse[i] = ToResponseOrderStatus(orderStatus)
	}

	return orderStatusResponse
}
