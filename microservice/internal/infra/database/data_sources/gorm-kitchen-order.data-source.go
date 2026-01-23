package data_sources

import (
	"gorm.io/gorm"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/mappers"
	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/infra/database"
)

type GormKitchenOrderDataSource struct {
	db *gorm.DB
}

func NewGormKitchenOrderDataSource() *GormKitchenOrderDataSource {
	return &GormKitchenOrderDataSource{
		db: database.GetDB(),
	}
}

func (r *GormKitchenOrderDataSource) Insert(kitchenOrder daos.KitchenOrderDAO) error {
	kitchenOrderModel := mappers.FromDAOToModelKitchenOrder(kitchenOrder)

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&kitchenOrderModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	var kitchenOrders []*models.KitchenOrderModel

	query := r.db.
		Joins("JOIN order_status ON kitchen_order.status_id = order_status.id").
		Preload("Status").
		Preload("Items").
		Where("order_status.name <> ?", "Finalizado").
		Order(`
            CASE order_status.name
                WHEN 'Pronto' THEN 1
                WHEN 'Em preparação' THEN 2
                WHEN 'Recebido' THEN 3
                ELSE 4
            END
        `).
		Order("kitchen_order.created_at ASC")

	if filter.CreatedAtFrom != nil {
		query = query.Where("kitchen_order.created_at >= ?", *filter.CreatedAtFrom)
	}

	if filter.CreatedAtTo != nil {
		query = query.Where("kitchen_order.created_at <= ?", *filter.CreatedAtTo)
	}

	if filter.StatusID != nil {
		query = query.Where("kitchen_order.status_id = ?", *filter.StatusID)
	}

	if err := query.Find(&kitchenOrders).Error; err != nil {
		return nil, err
	}

	return mappers.FromModelArrayToDAOArrayKitchenOrder(kitchenOrders), nil
}

func (r *GormKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	var kitchenOrder *models.KitchenOrderModel

	if err := r.db.Preload("Status").Preload("Items").First(&kitchenOrder, "id = ?", id).Error; err != nil {
		return daos.KitchenOrderDAO{}, err
	}

	return mappers.FromModelToDAOKitchenOrder(kitchenOrder), nil
}

func (r *GormKitchenOrderDataSource) Update(kitchenOrder daos.KitchenOrderDAO) error {
	var existing models.KitchenOrderModel
	if err := r.db.First(&existing, "id = ?", kitchenOrder.ID).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"status_id":  kitchenOrder.Status.ID,
		"updated_at": kitchenOrder.UpdatedAt,
	}

	return r.db.Model(&models.KitchenOrderModel{}).
		Where("id = ?", kitchenOrder.ID).
		Updates(updates).Error
}

func (r *GormKitchenOrderDataSource) Delete(id string) error {
	return r.db.Delete(&models.KitchenOrderModel{}, "id = ?", id).Error
}
