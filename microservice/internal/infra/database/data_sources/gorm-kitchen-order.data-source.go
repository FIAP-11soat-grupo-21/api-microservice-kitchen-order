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

	return r.db.Model(&models.KitchenOrderModel{}).Create(&kitchenOrderModel).Error
}

func (r *GormKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	var kitchenOrders []*models.KitchenOrderModel

	query := r.db.
		Joins("JOIN order_status ON kitchen_order.status_id = order_status.id").
		Preload("Status").
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

	if err := r.db.Preload("Status").First(&kitchenOrder, "id = ?", id).Error; err != nil {
		return daos.KitchenOrderDAO{}, err
	}

	return mappers.FromModelToDAOKitchenOrder(kitchenOrder), nil
}

func (r *GormKitchenOrderDataSource) Update(kitchenOrder daos.KitchenOrderDAO) error {
	return r.db.Save(mappers.FromDAOToModelKitchenOrder(kitchenOrder)).Error
}

func (r *GormKitchenOrderDataSource) Delete(id string) error {
	return r.db.Delete(&models.KitchenOrderModel{}, "id = ?", id).Error
}
