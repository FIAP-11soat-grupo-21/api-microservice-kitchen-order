package seed

import (
	"log"

	"gorm.io/gorm"

	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/constants"
)

func SeedOrderStatus(db *gorm.DB) {
	defaults := []models.OrderStatusModel{
		{ID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, Name: "Recebido"},
		{ID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID, Name: "Em preparação"},
		{ID: constants.KITCHEN_ORDER_STATUS_READY_ID, Name: "Pronto"},
		{ID: constants.KITCHEN_ORDER_STATUS_FINISHED_ID, Name: "Finalizado"},
	}

	for _, status := range defaults {
		var existing models.OrderStatusModel
		if err := db.Where("id = ?", status.ID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&status).Error; err != nil {
				log.Printf("Erro ao criar status %s: %v", status.ID, err)
			} else {
				log.Printf("Status %s criado", status.ID)
			}
		}
	}
}
