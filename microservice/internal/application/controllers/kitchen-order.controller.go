package controllers

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	presenter "tech_challenge/internal/application/presenters"
	"tech_challenge/internal/interfaces"
	"tech_challenge/internal/use_cases"
)

type KitchenOrderController struct {
	kitchenOrderDataSource interfaces.IKitchenOrderDataSource
	orderStatusDataSource  interfaces.IOrderStatusDataSource
	kitchenOrderGateway    gateways.KitchenOrderGateway
	orderStatusGateway     gateways.OrderStatusGateway
}

func NewKitchenOrderController(kitchenOrderDataSource interfaces.IKitchenOrderDataSource, orderStatusDataSource interfaces.IOrderStatusDataSource) *KitchenOrderController {
	return &KitchenOrderController{
		kitchenOrderDataSource: kitchenOrderDataSource,
		orderStatusDataSource:  orderStatusDataSource,
		kitchenOrderGateway:    *gateways.NewKitchenOrderGateway(kitchenOrderDataSource),
		orderStatusGateway:     *gateways.NewOrderStatusGateway(orderStatusDataSource),
	}
}

func (c *KitchenOrderController) Create(kitchenOrderDTO dtos.CreateKitchenOrderDTO) (dtos.KitchenOrderResponseDTO, error) {
	kitchenOrderUseCase := use_cases.NewCreateKitchenOrderUseCase(c.kitchenOrderGateway, c.orderStatusGateway)

	kitchenOrder, err := kitchenOrderUseCase.Execute(
		kitchenOrderDTO.OrderID,
	)

	if err != nil {
		return dtos.KitchenOrderResponseDTO{}, err
	}

	return presenter.ToResponse(kitchenOrder), nil
}

func (c *KitchenOrderController) FindAll(filter dtos.KitchenOrderFilter) ([]dtos.KitchenOrderResponseDTO, error) {
	kitchenOrderUseCase := use_cases.NewFindAllKitchenOrderUseCase(c.kitchenOrderGateway)

	kitchenOrders, err := kitchenOrderUseCase.Execute(filter)

	if err != nil {
		return nil, err
	}

	return presenter.ToResponseList(kitchenOrders), nil
}

func (c *KitchenOrderController) FindByID(id string) (dtos.KitchenOrderResponseDTO, error) {
	kitchenOrderUseCase := use_cases.NewFindKitchenOrderByIDUseCase(c.kitchenOrderGateway)

	kitchenOrder, err := kitchenOrderUseCase.Execute(id)

	if err != nil {
		return dtos.KitchenOrderResponseDTO{}, err
	}

	return presenter.ToResponse(kitchenOrder), nil
}

func (c *KitchenOrderController) Update(kitchenOrderDTO dtos.UpdateKitchenOrderDTO) (dtos.KitchenOrderResponseDTO, error) {
	kitchenOrderUseCase := use_cases.NewUpdateKitchenOrderUseCase(c.kitchenOrderGateway, c.orderStatusGateway)

	kitchenOrder, err := kitchenOrderUseCase.Execute(kitchenOrderDTO)

	if err != nil {
		return dtos.KitchenOrderResponseDTO{}, err
	}

	return presenter.ToResponse(kitchenOrder), nil
}
