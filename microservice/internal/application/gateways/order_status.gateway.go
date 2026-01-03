package gateways

import (
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/interfaces"
)

type OrderStatusGateway struct {
	dataSource interfaces.IOrderStatusDataSource
}

func NewOrderStatusGateway(dataSource interfaces.IOrderStatusDataSource) *OrderStatusGateway {
	return &OrderStatusGateway{
		dataSource: dataSource,
	}
}

func (o *OrderStatusGateway) FindAll() ([]entities.OrderStatus, error) {
	orderStatusDAOs, err := o.dataSource.FindAll()

	if err != nil {
		return nil, err
	}

	orderStatusList := make([]entities.OrderStatus, 0, len(orderStatusDAOs))

	for _, orderStatusDAO := range orderStatusDAOs {
		orderStatus, err := entities.NewOrderStatus(
			orderStatusDAO.ID,
			orderStatusDAO.Name,
		)

		if err != nil {
			return nil, err
		}

		orderStatusList = append(orderStatusList, *orderStatus)
	}

	return orderStatusList, nil
}

func (g *OrderStatusGateway) FindByID(id string) (entities.OrderStatus, error) {
	orderStatusDAO, err := g.dataSource.FindByID(id)
	if err != nil {
		return entities.OrderStatus{}, err
	}

	orderStatus, err := entities.NewOrderStatus(
		orderStatusDAO.ID,
		orderStatusDAO.Name,
	)

	if err != nil {
		return entities.OrderStatus{}, err
	}

	return *orderStatus, nil
}
