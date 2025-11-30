package controllers

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/application/presenters"
	"tech_challenge/internal/interfaces"
	"tech_challenge/internal/use_cases"
)

type OrderStatusController struct {
	dataSource interfaces.IOrderStatusDataSource
	gateway    gateways.OrderStatusGateway
}

func NewOrderStatusController(dataSource interfaces.IOrderStatusDataSource) *OrderStatusController {
	return &OrderStatusController{
		dataSource: dataSource,
		gateway:    *gateways.NewOrderStatusGateway(dataSource),
	}
}

func (c *OrderStatusController) FindAll() ([]dtos.OrderStatusResponseDTO, error) {
	orderStatusUseCase := use_cases.NewFindAllOrdersStatusUseCase(c.gateway)

	orderStatus, err := orderStatusUseCase.Execute()

	if err != nil {
		return nil, err
	}

	return presenters.ToResponseListOrderStatus(orderStatus), nil
}
