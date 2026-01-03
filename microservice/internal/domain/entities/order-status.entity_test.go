package entities

import (
	"testing"

	"tech_challenge/internal/shared/config/constants"
)

func TestNewOrderStatus_Success(t *testing.T) {
	// Arrange
	id := constants.KITCHEN_ORDER_STATUS_RECEIVED_ID
	name := "Recebido"

	// Act
	orderStatus, err := NewOrderStatus(id, name)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if orderStatus.ID != id {
		t.Errorf("Expected ID %s, got %s", id, orderStatus.ID)
	}

	if orderStatus.Name.Value() != name {
		t.Errorf("Expected name %s, got %s", name, orderStatus.Name.Value())
	}
}

func TestNewOrderStatus_InvalidName_TooShort(t *testing.T) {
	// Arrange
	id := constants.KITCHEN_ORDER_STATUS_RECEIVED_ID
	shortName := "ab" // Muito curto (menos de 3 caracteres)

	// Act
	_, err := NewOrderStatus(id, shortName)

	// Assert
	if err == nil {
		t.Error("Expected error for short name, got nil")
	}

	expectedMessage := "name must be at least 3 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewOrderStatus_InvalidName_TooLong(t *testing.T) {
	// Arrange
	id := constants.KITCHEN_ORDER_STATUS_RECEIVED_ID
	longName := "a" // Cria um nome com mais de 100 caracteres
	for i := 0; i < 100; i++ {
		longName += "a"
	}

	// Act
	_, err := NewOrderStatus(id, longName)

	// Assert
	if err == nil {
		t.Error("Expected error for long name, got nil")
	}

	expectedMessage := "name must be at most 100 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewOrderStatus_EmptyName(t *testing.T) {
	// Arrange
	id := constants.KITCHEN_ORDER_STATUS_RECEIVED_ID
	emptyName := ""

	// Act
	_, err := NewOrderStatus(id, emptyName)

	// Assert
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestNewOrderStatus_AllValidStatuses(t *testing.T) {
	// Arrange & Act & Assert
	testCases := []struct {
		id   string
		name string
	}{
		{constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido"},
		{constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação"},
		{constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto"},
		{constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado"},
	}

	for _, tc := range testCases {
		orderStatus, err := NewOrderStatus(tc.id, tc.name)

		if err != nil {
			t.Errorf("Expected no error for status %s, got %v", tc.name, err)
		}

		if orderStatus.ID != tc.id {
			t.Errorf("Expected ID %s, got %s", tc.id, orderStatus.ID)
		}

		if orderStatus.Name.Value() != tc.name {
			t.Errorf("Expected name %s, got %s", tc.name, orderStatus.Name.Value())
		}
	}
}

func TestNewOrderStatus_MinimumValidName(t *testing.T) {
	// Arrange
	id := constants.KITCHEN_ORDER_STATUS_RECEIVED_ID
	minName := "abc" // Exatamente 3 caracteres (mínimo válido)

	// Act
	orderStatus, err := NewOrderStatus(id, minName)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for minimum valid name, got %v", err)
	}

	if orderStatus.Name.Value() != minName {
		t.Errorf("Expected name %s, got %s", minName, orderStatus.Name.Value())
	}
}

func TestNewOrderStatus_MaximumValidName(t *testing.T) {
	// Arrange
	id := constants.KITCHEN_ORDER_STATUS_RECEIVED_ID
	maxName := ""
	for i := 0; i < 100; i++ { // Exatamente 100 caracteres (máximo válido)
		maxName += "a"
	}

	// Act
	orderStatus, err := NewOrderStatus(id, maxName)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for maximum valid name, got %v", err)
	}

	if orderStatus.Name.Value() != maxName {
		t.Errorf("Expected name length %d, got %d", len(maxName), len(orderStatus.Name.Value()))
	}
}