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

func TestToResponseOrderStatus_PreparingStatus(t *testing.T) {
	// Arrange
	orderStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")

	// Act
	response := ToResponseOrderStatus(*orderStatus)

	// Assert
	if response.ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, response.ID)
	}

	if response.Name != "Em preparação" {
		t.Errorf("Expected Name 'Em preparação', got %s", response.Name)
	}
}

func TestToResponseOrderStatus_ReadyStatus(t *testing.T) {
	// Arrange
	orderStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto")

	// Act
	response := ToResponseOrderStatus(*orderStatus)

	// Assert
	if response.ID != constants.KITCHEN_ORDER_STATUS_READY_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_READY_ID, response.ID)
	}

	if response.Name != "Pronto" {
		t.Errorf("Expected Name 'Pronto', got %s", response.Name)
	}
}

func TestToResponseOrderStatus_FinishedStatus(t *testing.T) {
	// Arrange
	orderStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado")

	// Act
	response := ToResponseOrderStatus(*orderStatus)

	// Assert
	if response.ID != constants.KITCHEN_ORDER_STATUS_FINISHED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_FINISHED_ID, response.ID)
	}

	if response.Name != "Finalizado" {
		t.Errorf("Expected Name 'Finalizado', got %s", response.Name)
	}
}

func TestToResponseListOrderStatus_SingleStatus(t *testing.T) {
	// Arrange
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	statuses := []entities.OrderStatus{*status}

	// Act
	responses := ToResponseListOrderStatus(statuses)

	// Assert
	if len(responses) != 1 {
		t.Errorf("Expected 1 response, got %d", len(responses))
	}

	if responses[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, responses[0].ID)
	}

	if responses[0].Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", responses[0].Name)
	}
}

func TestToResponseListOrderStatus_TwoStatuses(t *testing.T) {
	// Arrange
	status1, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	status2, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	statuses := []entities.OrderStatus{*status1, *status2}

	// Act
	responses := ToResponseListOrderStatus(statuses)

	// Assert
	if len(responses) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(responses))
	}

	if responses[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected first response ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, responses[0].ID)
	}

	if responses[0].Name != "Recebido" {
		t.Errorf("Expected first response Name 'Recebido', got %s", responses[0].Name)
	}

	if responses[1].ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected second response ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, responses[1].ID)
	}

	if responses[1].Name != "Em preparação" {
		t.Errorf("Expected second response Name 'Em preparação', got %s", responses[1].Name)
	}
}

func TestToResponseListOrderStatus_ThreeStatuses(t *testing.T) {
	// Arrange
	status1, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	status2, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	status3, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto")
	statuses := []entities.OrderStatus{*status1, *status2, *status3}

	// Act
	responses := ToResponseListOrderStatus(statuses)

	// Assert
	if len(responses) != 3 {
		t.Errorf("Expected 3 responses, got %d", len(responses))
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
		if responses[i].ID != expected.id {
			t.Errorf("Expected response[%d].ID %s, got %s", i, expected.id, responses[i].ID)
		}

		if responses[i].Name != expected.name {
			t.Errorf("Expected response[%d].Name '%s', got '%s'", i, expected.name, responses[i].Name)
		}
	}
}

func TestToResponseListOrderStatus_AllStatuses(t *testing.T) {
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

	// Verifica cada status individualmente
	if responses[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID || responses[0].Name != "Recebido" {
		t.Error("First status mapping failed")
	}

	if responses[1].ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID || responses[1].Name != "Em preparação" {
		t.Error("Second status mapping failed")
	}

	if responses[2].ID != constants.KITCHEN_ORDER_STATUS_READY_ID || responses[2].Name != "Pronto" {
		t.Error("Third status mapping failed")
	}

	if responses[3].ID != constants.KITCHEN_ORDER_STATUS_FINISHED_ID || responses[3].Name != "Finalizado" {
		t.Error("Fourth status mapping failed")
	}
}
