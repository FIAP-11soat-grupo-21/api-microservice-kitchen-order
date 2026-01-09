package mappers

import (
	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
)

func FromDAOToModelKitchenOrder(kitchenOrder daos.KitchenOrderDAO) *models.KitchenOrderModel {
	items := make([]models.OrderItemModel, len(kitchenOrder.Items))
	for i, item := range kitchenOrder.Items {
		items[i] = models.OrderItemModel{
			ID:             item.ID,
			KitchenOrderID: kitchenOrder.ID,
			OrderID:        item.OrderID,
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
		}
	}

	return &models.KitchenOrderModel{
		ID:         kitchenOrder.ID,
		OrderID:    kitchenOrder.OrderID,
		CustomerID: kitchenOrder.CustomerID,
		Amount:     kitchenOrder.Amount,
		StatusID:   kitchenOrder.Status.ID,
		Slug:       kitchenOrder.Slug,
		Items:      items,
		CreatedAt:  kitchenOrder.CreatedAt,
		UpdatedAt:  kitchenOrder.UpdatedAt,
	}
}

func FromModelToDAOKitchenOrder(kitchenOrder *models.KitchenOrderModel) daos.KitchenOrderDAO {

	statusDAO := daos.OrderStatusDAO{
		ID:   kitchenOrder.Status.ID,
		Name: kitchenOrder.Status.Name,
	}

	items := make([]daos.OrderItemDAO, len(kitchenOrder.Items))
	for i, item := range kitchenOrder.Items {
		items[i] = daos.OrderItemDAO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return daos.KitchenOrderDAO{
		ID:         kitchenOrder.ID,
		OrderID:    kitchenOrder.OrderID,
		CustomerID: kitchenOrder.CustomerID,
		Amount:     kitchenOrder.Amount,
		Status:     statusDAO,
		Slug:       kitchenOrder.Slug,
		Items:      items,
		CreatedAt:  kitchenOrder.CreatedAt,
		UpdatedAt:  kitchenOrder.UpdatedAt,
	}
}

func FromModelArrayToDAOArrayKitchenOrder(models []*models.KitchenOrderModel) []daos.KitchenOrderDAO {
	daos := make([]daos.KitchenOrderDAO, len(models))
	for i, model := range models {
		daos[i] = FromModelToDAOKitchenOrder(model)
	}
	return daos
}
