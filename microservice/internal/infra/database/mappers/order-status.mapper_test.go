package mappers

import (
	"testing"

	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/constants"
)

func TestFromDAOToModelOrderStatus(t *testing.T) {
	// Arrange
	dao := daos.OrderStatusDAO{
		ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		Name: "Recebido",
	}

	// Act
	model := FromDAOToModelOrderStatus(dao)

	// Assert
	if model.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, model.ID)
	}

	if model.Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", model.Name)
	}
}

func TestFromModelToDAOOrderStatus(t *testing.T) {
	// Arrange
	model := &models.OrderStatusModel{
		ID:   constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		Name: "Em preparação",
	}

	// Act
	dao := FromModelToDAOOrderStatus(model)

	// Assert
	if dao.ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, dao.ID)
	}

	if dao.Name != "Em preparação" {
		t.Errorf("Expected Name 'Em preparação', got %s", dao.Name)
	}
}

func TestFromModelArrayToDAOArrayOrderStatus(t *testing.T) {
	// Arrange
	models := []*models.OrderStatusModel{
		{
			ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			Name: "Recebido",
		},
		{
			ID:   constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
			Name: "Em preparação",
		},
		{
			ID:   constants.KITCHEN_ORDER_STATUS_READY_ID,
			Name: "Pronto",
		},
		{
			ID:   constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
			Name: "Finalizado",
		},
	}

	// Act
	daos := FromModelArrayToDAOArrayOrderStatus(models)

	// Assert
	if len(daos) != 4 {
		t.Errorf("Expected 4 DAOs, got %d", len(daos))
	}

	// Verifica se todos os status esperados estão presentes
	statusMap := make(map[string]string)
	for _, dao := range daos {
		statusMap[dao.ID] = dao.Name
	}

	expectedStatuses := map[string]string{
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID:  "Recebido",
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID: "Em preparação",
		constants.KITCHEN_ORDER_STATUS_READY_ID:     "Pronto",
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID:  "Finalizado",
	}

	for expectedID, expectedName := range expectedStatuses {
		if actualName, exists := statusMap[expectedID]; !exists {
			t.Errorf("Expected status ID %s not found", expectedID)
		} else if actualName != expectedName {
			t.Errorf("Expected status name '%s' for ID %s, got '%s'", expectedName, expectedID, actualName)
		}
	}
}

func TestFromModelArrayToDAOArrayOrderStatus_EmptyArray(t *testing.T) {
	// Arrange
	models := []*models.OrderStatusModel{}

	// Act
	daos := FromModelArrayToDAOArrayOrderStatus(models)

	// Assert
	if len(daos) != 0 {
		t.Errorf("Expected 0 DAOs, got %d", len(daos))
	}
}
