package use_cases

import (
	"fmt"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
	identity_manager "tech_challenge/internal/shared/pkg/identity"
)

type CreateKitchenOrderUseCase struct {
	kitchenOrderGateway gateways.KitchenOrderGateway
	orderStatusGateway  gateways.OrderStatusGateway
}

func NewCreateKitchenOrderUseCase(kitchenOrderGateway gateways.KitchenOrderGateway, orderStatusGateway gateways.OrderStatusGateway) *CreateKitchenOrderUseCase {
	return &CreateKitchenOrderUseCase{
		kitchenOrderGateway: kitchenOrderGateway,
		orderStatusGateway:  orderStatusGateway,
	}
}

func (ko *CreateKitchenOrderUseCase) Execute(orderID string, customerID *string, amount float64, items []dtos.OrderItemDTO) (entities.KitchenOrder, error) {

	status, err := ko.orderStatusGateway.FindByID(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.OrderStatusNotFoundException{}
	}

	// Filters the day's orders for slug generation
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := now

	filterDailyKitchenOrder := dtos.KitchenOrderFilter{
		CreatedAtFrom: &from,
		CreatedAtTo:   &to,
	}

	orders, err := ko.kitchenOrderGateway.FindAll(filterDailyKitchenOrder)
	if err != nil {
		return entities.KitchenOrder{}, err
	}

	slug := fmt.Sprintf("%03d", len(orders)+1)

	orderItems := make([]entities.OrderItem, len(items))
	for i, item := range items {
		orderItem, err := entities.NewOrderItem(
			identity_manager.NewUUIDV4(),
			orderID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
		)
		if err != nil {
			return entities.KitchenOrder{}, err
		}
		orderItems[i] = *orderItem
	}

	kitchenOrder, err := entities.NewKitchenOrderWithOrderData(
		identity_manager.NewUUIDV4(),
		orderID,
		slug,
		customerID,
		amount,
		orderItems,
		status,
		time.Now(),
		nil,
	)

	if err != nil {
		return entities.KitchenOrder{}, err
	}

	err = ko.kitchenOrderGateway.Insert(*kitchenOrder)

	if err != nil {
		return entities.KitchenOrder{}, err
	}

	return *kitchenOrder, nil
}
