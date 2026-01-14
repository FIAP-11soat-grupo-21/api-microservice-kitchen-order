package use_cases

import (
	"errors"
	"testing"
	"time"

	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
)

func TestFindKitchenOrderByIDUseCase_InvalidID(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	invalidIDs := []string{
		"",
		"invalid-uuid",
		"123",
		"not-a-uuid-at-all",
	}

	for _, invalidID := range invalidIDs {
		// Act
		result, err := useCase.Execute(invalidID)

		// Assert
		if err == nil {
			t.Errorf("Expected error for invalid ID '%s', got nil", invalidID)
		}

		if !result.IsEmpty() {
			t.Errorf("Expected empty result for invalid ID '%s', got %v", invalidID, result)
		}
	}
}

func TestFindKitchenOrderByIDUseCase_OrderNotFound(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	dataStore.kitchenOrders = []entities.KitchenOrder{} // Lista vazia

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	validID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	result, err := useCase.Execute(validID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*exceptions.KitchenOrderNotFoundException); !ok {
		t.Errorf("Expected KitchenOrderNotFoundException, got %T", err)
	}

	if !result.IsEmpty() {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestFindKitchenOrderByIDUseCase_GatewayError(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	dataStore.shouldReturnError = true
	dataStore.errorToReturn = errors.New("database connection error")

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	validID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	result, err := useCase.Execute(validID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*exceptions.KitchenOrderNotFoundException); !ok {
		t.Errorf("Expected KitchenOrderNotFoundException, got %T", err)
	}

	if !result.IsEmpty() {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestFindKitchenOrderByIDUseCase_EmptyOrder(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	
	// Adiciona um pedido vazio
	emptyOrder := entities.KitchenOrder{}
	dataStore.kitchenOrders = []entities.KitchenOrder{emptyOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	validID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	result, err := useCase.Execute(validID)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*exceptions.KitchenOrderNotFoundException); !ok {
		t.Errorf("Expected KitchenOrderNotFoundException, got %T", err)
	}

	if !result.IsEmpty() {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestFindKitchenOrderByIDUseCase_SuccessWithItems(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	expectedOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", status, time.Now(), nil,
	)

	// Adiciona itens ao pedido
	item1, _ := entities.NewOrderItem("item1", orderID, "product1", 2, 25.50)
	item2, _ := entities.NewOrderItem("item2", orderID, "product2", 1, 15.00)
	
	expectedOrder.AddItem(*item1)
	expectedOrder.AddItem(*item2)
	expectedOrder.CalcTotalAmount()

	dataStore.kitchenOrders = []entities.KitchenOrder{*expectedOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	// Act
	result, err := useCase.Execute(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != orderID {
		t.Errorf("Expected ID %s, got %s", orderID, result.ID)
	}

	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}

	if result.Amount != 66.00 { // 2*25.50 + 1*15.00
		t.Errorf("Expected amount 66.00, got %f", result.Amount)
	}
}

func TestFindKitchenOrderByIDUseCase_ValidUUIDs(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	validUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6ba7b811-9dad-11d1-80b4-00c04fd430c8",
		"01234567-89ab-cdef-0123-456789abcdef",
	}

	for i, uuid := range validUUIDs {
		order, _ := entities.NewKitchenOrder(
			uuid, "order"+string(rune(i)), "00"+string(rune(i)), status, time.Now(), nil,
		)
		dataStore.kitchenOrders = append(dataStore.kitchenOrders, *order)
	}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	for _, uuid := range validUUIDs {
		// Act
		result, err := useCase.Execute(uuid)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for UUID %s, got %v", uuid, err)
		}

		if result.ID != uuid {
			t.Errorf("Expected ID %s, got %s", uuid, result.ID)
		}
	}
}

func TestFindKitchenOrderByIDUseCase_DifferentStatuses(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()

	testCases := []struct {
		orderID string
		status  entities.OrderStatus
	}{
		{"123e4567-e89b-12d3-a456-426614174001", dataStore.orderStatuses[0]}, // Recebido
		{"123e4567-e89b-12d3-a456-426614174002", dataStore.orderStatuses[1]}, // Em preparação
		{"123e4567-e89b-12d3-a456-426614174003", dataStore.orderStatuses[2]}, // Pronto
		{"123e4567-e89b-12d3-a456-426614174004", dataStore.orderStatuses[3]}, // Finalizado
	}

	for _, tc := range testCases {
		order, _ := entities.NewKitchenOrder(
			tc.orderID, "order-"+tc.orderID, "slug-"+tc.orderID, tc.status, time.Now(), nil,
		)
		dataStore.kitchenOrders = append(dataStore.kitchenOrders, *order)
	}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	for _, tc := range testCases {
		// Act
		result, err := useCase.Execute(tc.orderID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for order %s, got %v", tc.orderID, err)
		}

		if result.ID != tc.orderID {
			t.Errorf("Expected ID %s, got %s", tc.orderID, result.ID)
		}

		if result.Status.ID != tc.status.ID {
			t.Errorf("Expected status ID %s, got %s", tc.status.ID, result.Status.ID)
		}
	}
}

func TestFindKitchenOrderByIDUseCase_WithUpdatedAt(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	now := time.Now()
	updatedAt := now.Add(time.Hour)
	
	expectedOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", status, now, &updatedAt,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*expectedOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	// Act
	result, err := useCase.Execute(orderID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != orderID {
		t.Errorf("Expected ID %s, got %s", orderID, result.ID)
	}

	if result.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set")
	}

	if !result.UpdatedAt.Equal(updatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, *result.UpdatedAt)
	}
}

func TestFindKitchenOrderByIDUseCase_EdgeCaseUUIDs(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	// UUIDs com casos extremos
	edgeCaseUUIDs := []string{
		"00000000-0000-0000-0000-000000000000", // UUID zero
		"ffffffff-ffff-ffff-ffff-ffffffffffff", // UUID máximo
		"12345678-1234-1234-1234-123456789abc", // UUID com números e letras
	}

	for i, uuid := range edgeCaseUUIDs {
		order, _ := entities.NewKitchenOrder(
			uuid, "order"+string(rune(i)), "00"+string(rune(i)), status, time.Now(), nil,
		)
		dataStore.kitchenOrders = append(dataStore.kitchenOrders, *order)
	}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	for _, uuid := range edgeCaseUUIDs {
		// Act
		result, err := useCase.Execute(uuid)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for UUID %s, got %v", uuid, err)
		}

		if result.ID != uuid {
			t.Errorf("Expected ID %s, got %s", uuid, result.ID)
		}
	}
}