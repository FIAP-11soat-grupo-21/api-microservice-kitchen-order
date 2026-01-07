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

	items := make([]dtos.OrderItemDTO, len(kitchenOrder.Items))
	for i, item := range kitchenOrder.Items {
		items[i] = dtos.OrderItemDTO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return dtos.KitchenOrderResponseDTO{
		ID:         kitchenOrder.ID,
		OrderID:    kitchenOrder.OrderID,
		CustomerID: kitchenOrder.CustomerID,
		Amount:     kitchenOrder.Amount,
		Slug:       kitchenOrder.Slug.Value(),
		Status:     status,
		Items:      items,
		CreatedAt:  kitchenOrder.CreatedAt,
		UpdatedAt:  kitchenOrder.UpdatedAt,
	}
}

func ToResponseList(kitchenOrders []entities.KitchenOrder) []dtos.KitchenOrderResponseDTO {
	kitchenOrderResponse := make([]dtos.KitchenOrderResponseDTO, len(kitchenOrders))

	for i, kitchenOrder := range kitchenOrders {
		kitchenOrderResponse[i] = ToResponse(kitchenOrder)
	}

	return kitchenOrderResponse
}
