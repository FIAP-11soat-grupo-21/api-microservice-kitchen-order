package seed

import (
	"gorm.io/gorm"

	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/constants"
)

// DBInterface define a interface para operações de banco de dados
type DBInterface interface {
	Where(query interface{}, args ...interface{}) DBInterface
	First(dest interface{}, conds ...interface{}) DBInterface
	Create(value interface{}) DBInterface
	GetError() error
}

// GormDBWrapper implementa DBInterface para *gorm.DB
type GormDBWrapper struct {
	db *gorm.DB
}

func (w *GormDBWrapper) Where(query interface{}, args ...interface{}) DBInterface {
	return &GormDBWrapper{db: w.db.Where(query, args...)}
}

func (w *GormDBWrapper) First(dest interface{}, conds ...interface{}) DBInterface {
	return &GormDBWrapper{db: w.db.First(dest, conds...)}
}

func (w *GormDBWrapper) Create(value interface{}) DBInterface {
	return &GormDBWrapper{db: w.db.Create(value)}
}

func (w *GormDBWrapper) GetError() error {
	return w.db.Error
}

func SeedOrderStatus(db *gorm.DB) {
	wrapper := &GormDBWrapper{db: db}
	seedOrderStatusInternal(wrapper)
}

func seedOrderStatusInternal(db DBInterface) {
	defaults := []models.OrderStatusModel{
		{ID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, Name: "Recebido"},
		{ID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID, Name: "Em preparação"},
		{ID: constants.KITCHEN_ORDER_STATUS_READY_ID, Name: "Pronto"},
		{ID: constants.KITCHEN_ORDER_STATUS_FINISHED_ID, Name: "Finalizado"},
	}

	for _, status := range defaults {
		var existing models.OrderStatusModel
		if err := db.Where("id = ?", status.ID).First(&existing).GetError(); err == gorm.ErrRecordNotFound {
			statusCopy := status
			db.Create(&statusCopy)
		}
	}
}
