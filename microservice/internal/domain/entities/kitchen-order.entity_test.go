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

func TestNewKitchenOrderWithOrderData_Success(t *testing.T) {
	// Arrange
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	customerID := "customer-456"
	amount := 25.50
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	orderItem, _ := NewOrderItem("item-1", orderID, "product-1", 2, 12.75)
	items := []OrderItem{*orderItem}

	kitchenOrder, err := NewKitchenOrderWithOrderData(id, orderID, slug, &customerID, amount, items, *status, createdAt, nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if kitchenOrder.ID != id {
		t.Errorf("Expected ID %s, got %s", id, kitchenOrder.ID)
	}

	if kitchenOrder.OrderID != orderID {
		t.Errorf("Expected OrderID %s, got %s", orderID, kitchenOrder.OrderID)
	}

	if kitchenOrder.CustomerID == nil || *kitchenOrder.CustomerID != customerID {
		t.Errorf("Expected CustomerID %s, got %v", customerID, kitchenOrder.CustomerID)
	}

	if kitchenOrder.Amount != amount {
		t.Errorf("Expected Amount %f, got %f", amount, kitchenOrder.Amount)
	}

	if len(kitchenOrder.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(kitchenOrder.Items))
	}
}

func TestKitchenOrder_AddItem(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	kitchenOrder, _ := NewKitchenOrder(id, orderID, slug, *status, createdAt, nil)
	orderItem, _ := NewOrderItem("item-1", orderID, "product-1", 2, 12.75)

	kitchenOrder.AddItem(*orderItem)

	if len(kitchenOrder.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(kitchenOrder.Items))
	}

	if kitchenOrder.Items[0].ID != "item-1" {
		t.Errorf("Expected item ID 'item-1', got %s", kitchenOrder.Items[0].ID)
	}
}

func TestKitchenOrder_CalcTotalAmount(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440000"
	orderID := "order-123"
	slug := "001"
	status, _ := NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	createdAt := time.Now()

	kitchenOrder, _ := NewKitchenOrder(id, orderID, slug, *status, createdAt, nil)
	
	item1, _ := NewOrderItem("item-1", orderID, "product-1", 2, 10.00) // 20.00
	item2, _ := NewOrderItem("item-2", orderID, "product-2", 1, 15.50) // 15.50
	
	kitchenOrder.AddItem(*item1)
	kitchenOrder.AddItem(*item2)

	kitchenOrder.CalcTotalAmount()

	expected := 35.5
	if kitchenOrder.Amount != expected {
		t.Errorf("Expected amount %f, got %f", expected, kitchenOrder.Amount)
	}
}
