package gateways

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/interfaces"
)

type KitchenOrderGateway struct {
	dataSource interfaces.IKitchenOrderDataSource
}

func NewKitchenOrderGateway(dataSource interfaces.IKitchenOrderDataSource) *KitchenOrderGateway {
	return &KitchenOrderGateway{
		dataSource: dataSource,
	}
}

func (g *KitchenOrderGateway) Insert(order entities.KitchenOrder) error {

	status := daos.OrderStatusDAO{
		ID:   order.Status.ID,
		Name: order.Status.Name.Value(),
	}

	items := make([]daos.OrderItemDAO, len(order.Items))
	for i, item := range order.Items {
		items[i] = daos.OrderItemDAO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return g.dataSource.Insert(daos.KitchenOrderDAO{
		ID:         order.ID,
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
		Slug:       order.Slug.Value(),
		Status:     status,
		Items:      items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	})
}

func (g *KitchenOrderGateway) FindByID(id string) (entities.KitchenOrder, error) {
	orderDAO, err := g.dataSource.FindByID(id)
	if err != nil {
		return entities.KitchenOrder{}, err
	}

	status, err := entities.NewOrderStatus(
		orderDAO.Status.ID,
		orderDAO.Status.Name,
	)
	if err != nil {
		return entities.KitchenOrder{}, err
	}

	items := make([]entities.OrderItem, len(orderDAO.Items))
	for i, itemDAO := range orderDAO.Items {
		item, err := entities.NewOrderItem(
			itemDAO.ID,
			itemDAO.OrderID,
			itemDAO.ProductID,
			itemDAO.Quantity,
			itemDAO.UnitPrice,
		)
		if err != nil {
			return entities.KitchenOrder{}, err
		}
		items[i] = *item
	}

	order, err := entities.NewKitchenOrderWithOrderData(
		orderDAO.ID,
		orderDAO.OrderID,
		orderDAO.Slug,
		orderDAO.CustomerID,
		orderDAO.Amount,
		items,
		*status,
		orderDAO.CreatedAt,
		orderDAO.UpdatedAt,
	)
	if err != nil {
		return entities.KitchenOrder{}, err
	}

	return *order, nil
}

func (g *KitchenOrderGateway) FindAll(filter dtos.KitchenOrderFilter) ([]entities.KitchenOrder, error) {
	orderDAOs, err := g.dataSource.FindAll(filter)
	if err != nil {
		return nil, err
	}

	orders := make([]entities.KitchenOrder, 0, len(orderDAOs))
	for _, orderDAO := range orderDAOs {

		status, err := entities.NewOrderStatus(
			orderDAO.Status.ID,
			orderDAO.Status.Name,
		)
		if err != nil {
			return nil, err
		}

		items := make([]entities.OrderItem, len(orderDAO.Items))
		for i, itemDAO := range orderDAO.Items {
			item, err := entities.NewOrderItem(
				itemDAO.ID,
				itemDAO.OrderID,
				itemDAO.ProductID,
				itemDAO.Quantity,
				itemDAO.UnitPrice,
			)
			if err != nil {
				return nil, err
			}
			items[i] = *item
		}

		order, err := entities.NewKitchenOrderWithOrderData(
			orderDAO.ID,
			orderDAO.OrderID,
			orderDAO.Slug,
			orderDAO.CustomerID,
			orderDAO.Amount,
			items,
			*status,
			orderDAO.CreatedAt,
			orderDAO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil
}

func (g *KitchenOrderGateway) Update(kitchenOrder entities.KitchenOrder) error {
	return g.dataSource.Update(daos.KitchenOrderDAO{
		ID:      kitchenOrder.ID,
		OrderID: kitchenOrder.OrderID,
		Slug:    kitchenOrder.Slug.Value(),
		Status: daos.OrderStatusDAO{
			ID:   kitchenOrder.Status.ID,
			Name: kitchenOrder.Status.Name.Value(),
		},
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	})
}
