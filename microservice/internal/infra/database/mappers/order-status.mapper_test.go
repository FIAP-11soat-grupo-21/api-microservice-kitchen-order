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

func TestFromDAOToModelOrderStatus_DifferentStatus(t *testing.T) {
	// Arrange
	dao := daos.OrderStatusDAO{
		ID:   constants.KITCHEN_ORDER_STATUS_READY_ID,
		Name: "Pronto",
	}

	// Act
	model := FromDAOToModelOrderStatus(dao)

	// Assert
	if model.ID != constants.KITCHEN_ORDER_STATUS_READY_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_READY_ID, model.ID)
	}

	if model.Name != "Pronto" {
		t.Errorf("Expected Name 'Pronto', got %s", model.Name)
	}
}

func TestFromDAOToModelOrderStatus_FinishedStatus(t *testing.T) {
	// Arrange
	dao := daos.OrderStatusDAO{
		ID:   constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
		Name: "Finalizado",
	}

	// Act
	model := FromDAOToModelOrderStatus(dao)

	// Assert
	if model.ID != constants.KITCHEN_ORDER_STATUS_FINISHED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_FINISHED_ID, model.ID)
	}

	if model.Name != "Finalizado" {
		t.Errorf("Expected Name 'Finalizado', got %s", model.Name)
	}
}

func TestFromModelToDAOOrderStatus_ReceivedStatus(t *testing.T) {
	// Arrange
	model := &models.OrderStatusModel{
		ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		Name: "Recebido",
	}

	// Act
	dao := FromModelToDAOOrderStatus(model)

	// Assert
	if dao.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, dao.ID)
	}

	if dao.Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", dao.Name)
	}
}

func TestFromModelToDAOOrderStatus_ReadyStatus(t *testing.T) {
	// Arrange
	model := &models.OrderStatusModel{
		ID:   constants.KITCHEN_ORDER_STATUS_READY_ID,
		Name: "Pronto",
	}

	// Act
	dao := FromModelToDAOOrderStatus(model)

	// Assert
	if dao.ID != constants.KITCHEN_ORDER_STATUS_READY_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_READY_ID, dao.ID)
	}

	if dao.Name != "Pronto" {
		t.Errorf("Expected Name 'Pronto', got %s", dao.Name)
	}
}

func TestFromModelToDAOOrderStatus_FinishedStatus(t *testing.T) {
	// Arrange
	model := &models.OrderStatusModel{
		ID:   constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
		Name: "Finalizado",
	}

	// Act
	dao := FromModelToDAOOrderStatus(model)

	// Assert
	if dao.ID != constants.KITCHEN_ORDER_STATUS_FINISHED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_FINISHED_ID, dao.ID)
	}

	if dao.Name != "Finalizado" {
		t.Errorf("Expected Name 'Finalizado', got %s", dao.Name)
	}
}

func TestFromModelArrayToDAOArrayOrderStatus_SingleModel(t *testing.T) {
	// Arrange
	models := []*models.OrderStatusModel{
		{
			ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			Name: "Recebido",
		},
	}

	// Act
	daos := FromModelArrayToDAOArrayOrderStatus(models)

	// Assert
	if len(daos) != 1 {
		t.Errorf("Expected 1 DAO, got %d", len(daos))
	}

	if daos[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, daos[0].ID)
	}

	if daos[0].Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", daos[0].Name)
	}
}

func TestFromModelArrayToDAOArrayOrderStatus_TwoModels(t *testing.T) {
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
	}

	// Act
	daos := FromModelArrayToDAOArrayOrderStatus(models)

	// Assert
	if len(daos) != 2 {
		t.Errorf("Expected 2 DAOs, got %d", len(daos))
	}

	if daos[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected first DAO ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, daos[0].ID)
	}

	if daos[0].Name != "Recebido" {
		t.Errorf("Expected first DAO Name 'Recebido', got %s", daos[0].Name)
	}

	if daos[1].ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected second DAO ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, daos[1].ID)
	}

	if daos[1].Name != "Em preparação" {
		t.Errorf("Expected second DAO Name 'Em preparação', got %s", daos[1].Name)
	}
}

func TestFromModelArrayToDAOArrayOrderStatus_ThreeModels(t *testing.T) {
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
	}

	// Act
	daos := FromModelArrayToDAOArrayOrderStatus(models)

	// Assert
	if len(daos) != 3 {
		t.Errorf("Expected 3 DAOs, got %d", len(daos))
	}

	expectedStatuses := []struct {
		id   string
		name string
	}{
		{constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido"},
		{constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação"},
		{constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto"},
	}

	for i, expected := range expectedStatuses {
		if daos[i].ID != expected.id {
			t.Errorf("Expected DAO[%d] ID %s, got %s", i, expected.id, daos[i].ID)
		}

		if daos[i].Name != expected.name {
			t.Errorf("Expected DAO[%d] Name '%s', got '%s'", i, expected.name, daos[i].Name)
		}
	}
}
