package entities

import (
	"testing"
	"time"

	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

func TestNewKitchenOrder_Success(t *testing.T) {
	// Arrange
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()
	updatedAt := &createdAt

	// Act
	kitchenOrder, err := NewKitchenOrder(id, orderID, slug, *status, createdAt, updatedAt)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if kitchenOrder.ID != id {
		t.Errorf("Expected ID %s, got %s", id, kitchenOrder.ID)
	}

	if kitchenOrder.OrderID != orderID {
		t.Errorf("Expected OrderID %s, got %s", orderID, kitchenOrder.OrderID)
	}

	if kitchenOrder.Slug.Value() != slug {
		t.Errorf("Expected slug %s, got %s", slug, kitchenOrder.Slug.Value())
	}

	if kitchenOrder.Status.ID != status.ID {
		t.Errorf("Expected status ID %s, got %s", status.ID, kitchenOrder.Status.ID)
	}

	if kitchenOrder.CreatedAt != createdAt {
		t.Errorf("Expected CreatedAt %v, got %v", createdAt, kitchenOrder.CreatedAt)
	}

	if kitchenOrder.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, kitchenOrder.UpdatedAt)
	}
}

func TestNewKitchenOrder_InvalidSlug(t *testing.T) {
	// Arrange
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	invalidSlug := "ab" // Muito curto
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	// Act
	_, err := NewKitchenOrder(id, orderID, invalidSlug, *status, createdAt, nil)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid slug, got nil")
	}
}

func TestNewKitchenOrder_NilUpdatedAt(t *testing.T) {
	// Arrange
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	// Act
	kitchenOrder, err := NewKitchenOrder(id, orderID, slug, *status, createdAt, nil)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if kitchenOrder.UpdatedAt != nil {
		t.Errorf("Expected UpdatedAt to be nil, got %v", kitchenOrder.UpdatedAt)
	}
}

func TestValidateID_ValidUUID(t *testing.T) {
	// Arrange
	validID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	err := ValidateID(validID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for valid UUID, got %v", err)
	}
}

func TestValidateID_InvalidUUID(t *testing.T) {
	// Arrange
	invalidID := "invalid-uuid"

	// Act
	err := ValidateID(invalidID)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid UUID, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}

	expectedMessage := "Invalid kitchen Order ID"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestValidateID_EmptyString(t *testing.T) {
	// Arrange
	emptyID := ""

	// Act
	err := ValidateID(emptyID)

	// Assert
	if err == nil {
		t.Error("Expected error for empty UUID, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}
}

func TestKitchenOrder_IsEmpty_True(t *testing.T) {
	// Arrange
	emptyOrder := &KitchenOrder{}

	// Act
	isEmpty := emptyOrder.IsEmpty()

	// Assert
	if !isEmpty {
		t.Error("Expected IsEmpty() to return true for empty order")
	}
}

func TestKitchenOrder_IsEmpty_False(t *testing.T) {
	// Arrange
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	kitchenOrder, _ := NewKitchenOrder(id, orderID, slug, *status, createdAt, nil)

	// Act
	isEmpty := kitchenOrder.IsEmpty()

	// Assert
	if isEmpty {
		t.Error("Expected IsEmpty() to return false for non-empty order")
	}
}

func TestKitchenOrder_StatusAssignment(t *testing.T) {
	// Arrange
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	// Act
	kitchenOrder, err := NewKitchenOrder(id, orderID, slug, *status, createdAt, nil)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if kitchenOrder.Status.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, kitchenOrder.Status.ID)
	}

	if kitchenOrder.Status.Name.Value() != "Recebido" {
		t.Errorf("Expected status name 'Recebido', got %s", kitchenOrder.Status.Name.Value())
	}
}
