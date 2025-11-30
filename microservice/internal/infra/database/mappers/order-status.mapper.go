package mappers

import (
	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
)

func FromDAOToModelOrderStatus(order daos.OrderStatusDAO) models.OrderStatusModel {
	return models.OrderStatusModel{
		ID:   order.ID,
		Name: order.Name,
	}
}

func FromModelToDAOOrderStatus(orderStatus *models.OrderStatusModel) daos.OrderStatusDAO {
	return daos.OrderStatusDAO{
		ID:   orderStatus.ID,
		Name: orderStatus.Name,
	}
}

func FromModelArrayToDAOArrayOrderStatus(models []*models.OrderStatusModel) []daos.OrderStatusDAO {
	daos := make([]daos.OrderStatusDAO, len(models))
	for i, model := range models {
		daos[i] = FromModelToDAOOrderStatus(model)
	}
	return daos
}
