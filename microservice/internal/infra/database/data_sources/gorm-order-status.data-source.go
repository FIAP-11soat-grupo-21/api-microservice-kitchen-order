package data_sources

import (
	"gorm.io/gorm"

	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/mappers"
	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/infra/database"
)

type GormOrderStatusDataSource struct {
	db *gorm.DB
}

func NewGormOrderStatusDataSource() *GormOrderStatusDataSource {
	return &GormOrderStatusDataSource{
		db: database.GetDB(),
	}
}

func (r *GormOrderStatusDataSource) Insert(orderStatus daos.OrderStatusDAO) error {
	orderStatusModel := mappers.FromDAOToModelOrderStatus(orderStatus)

	return r.db.Model(&daos.OrderStatusDAO{}).Create(&orderStatusModel).Error
}

func (r *GormOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	var orderStatus []*models.OrderStatusModel

	if err := r.db.Find(&orderStatus).Error; err != nil {
		return nil, err
	}

	return mappers.FromModelArrayToDAOArrayOrderStatus(orderStatus), nil
}

func (r *GormOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	var orderStatus *models.OrderStatusModel

	if err := r.db.First(&orderStatus, "id = ?", id).Error; err != nil {
		return daos.OrderStatusDAO{}, err
	}

	return mappers.FromModelToDAOOrderStatus(orderStatus), nil
}
