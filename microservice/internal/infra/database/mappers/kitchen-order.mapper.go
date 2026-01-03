package mappers

import (
	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
)

func FromDAOToModelKitchenOrder(kitchenOrder daos.KitchenOrderDAO) *models.KitchenOrderModel {
	return &models.KitchenOrderModel{
		ID:        kitchenOrder.ID,
		OrderID:   kitchenOrder.OrderID,
		StatusID:  kitchenOrder.Status.ID,
		Slug:      kitchenOrder.Slug,
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	}
}

func FromModelToDAOKitchenOrder(kitchenOrder *models.KitchenOrderModel) daos.KitchenOrderDAO {

	statusDAO := daos.OrderStatusDAO{
		ID:   kitchenOrder.Status.ID,
		Name: kitchenOrder.Status.Name,
	}

	return daos.KitchenOrderDAO{
		ID:        kitchenOrder.ID,
		OrderID:   kitchenOrder.OrderID,
		Status:    statusDAO,
		Slug:      kitchenOrder.Slug,
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	}
}

func FromModelArrayToDAOArrayKitchenOrder(models []*models.KitchenOrderModel) []daos.KitchenOrderDAO {
	daos := make([]daos.KitchenOrderDAO, len(models))
	for i, model := range models {
		daos[i] = FromModelToDAOKitchenOrder(model)
	}
	return daos
}
