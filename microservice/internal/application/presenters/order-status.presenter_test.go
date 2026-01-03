package presenters

import (
	"testing"

	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/shared/config/constants"
)

func TestToResponseOrderStatus(t *testing.T) {
	// Arrange
	orderStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")

	// Act
	response := ToResponseOrderStatus(*orderStatus)

	// Assert
	if response.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, response.ID)
	}

	if response.Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", response.Name)
	}
}

func TestToResponseListOrderStatus(t *testing.T) {
	// Arrange
	status1, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	status2, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	status3, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto")
	status4, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado")

	statuses := []entities.OrderStatus{*status1, *status2, *status3, *status4}

	// Act
	responses := ToResponseListOrderStatus(statuses)

	// Assert
	if len(responses) != 4 {
		t.Errorf("Expected 4 responses, got %d", len(responses))
	}

	// Verifica se todos os status esperados estão presentes
	statusMap := make(map[string]string)
	for _, response := range responses {
		statusMap[response.ID] = response.Name
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

func TestToResponseListOrderStatus_EmptyList(t *testing.T) {
	// Arrange
	statuses := []entities.OrderStatus{}

	// Act
	responses := ToResponseListOrderStatus(statuses)

	// Assert
	if len(responses) != 0 {
		t.Errorf("Expected 0 responses, got %d", len(responses))
	}
}
