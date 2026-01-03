package use_cases

import (
	"testing"
	"time"

	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

func TestCreateKitchenOrderUseCase_Success(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)

	useCase := NewCreateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)
	orderID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	result, err := useCase.Execute(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.OrderID != orderID {
		t.Errorf("Expected OrderID %s, got %s", orderID, result.OrderID)
	}

	if result.Slug.Value() != "001" {
		t.Errorf("Expected slug '001', got %s", result.Slug.Value())
	}

	if result.Status.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, result.Status.ID)
	}
}

func TestCreateKitchenOrderUseCase_StatusNotFound(t *testing.T) {
	// Arrange
	dataStore := &MockDataStore{
		kitchenOrders: []entities.KitchenOrder{},
		orderStatuses: []entities.OrderStatus{}, // Sem status disponíveis
	}
	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)

	useCase := NewCreateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)
	orderID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	_, err := useCase.Execute(orderID)

	// Assert
	if err == nil {
		t.Error("Expected error for status not found, got nil")
	}

	if _, ok := err.(*exceptions.OrderStatusNotFoundException); !ok {
		t.Errorf("Expected OrderStatusNotFoundException, got %T", err)
	}
}

func TestCreateKitchenOrderUseCase_SlugGeneration(t *testing.T) {
	// Arrange - Simula 2 pedidos já existentes no dia
	dataStore := NewMockDataStore()

	// Adiciona pedidos existentes
	existingOrder1, _ := entities.NewKitchenOrder(
		"id1", "order1", "001",
		dataStore.orderStatuses[0], // Status "Recebido"
		time.Now(), nil,
	)
	existingOrder2, _ := entities.NewKitchenOrder(
		"id2", "order2", "002",
		dataStore.orderStatuses[0], // Status "Recebido"
		time.Now(), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder1, *existingOrder2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)

	useCase := NewCreateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)
	orderID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	result, err := useCase.Execute(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Slug.Value() != "003" {
		t.Errorf("Expected slug '003', got %s", result.Slug.Value())
	}
}
