package interfaces

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
)

type IKitchenOrderDataSource interface {
	Insert(kitchenOrder daos.KitchenOrderDAO) error
	FindByID(id string) (daos.KitchenOrderDAO, error)
	FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error)
	Update(kitchenOrder daos.KitchenOrderDAO) error
}

type IOrderStatusDataSource interface {
	FindByID(id string) (daos.OrderStatusDAO, error)
	FindAll() ([]daos.OrderStatusDAO, error)
}
