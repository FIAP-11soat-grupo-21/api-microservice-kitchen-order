package use_cases

import (
	"context"
	"testing"
	"time"

	"tech_challenge/internal"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/shared/config/constants"
	"tech_challenge/internal/shared/interfaces"
)

// Mock MessageBroker
type MockMessageBroker struct{}

func (m *MockMessageBroker) Connect(ctx context.Context) error {
	return nil
}

func (m *MockMessageBroker) Close() error {
	return nil
}

func (m *MockMessageBroker) Publish(ctx context.Context, queue string, message interfaces.Message) error {
	return nil
}

func (m *MockMessageBroker) Subscribe(ctx context.Context, queue string, handler interfaces.MessageHandler) error {
	return nil
}

func (m *MockMessageBroker) Start(ctx context.Context) error {
	return nil
}

func (m *MockMessageBroker) Stop() error {
	return nil
}

func createTestKitchenOrder(id, orderID, slug string, status entities.OrderStatus) *entities.KitchenOrder {
	order, _ := entities.NewKitchenOrder(id, orderID, slug, status, time.Now(), nil)
	return order
}

func TestFindAllKitchenOrdersUseCase_Success(t *testing.T) {
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	order1 := createTestKitchenOrder("id1", "order1", "001", status)
	order2 := createTestKitchenOrder("id2", "order2", "002", status)
	dataStore.kitchenOrders = []entities.KitchenOrder{*order1, *order2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)

	result, err := useCase.Execute(dtos.KitchenOrderFilter{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(result))
	}
}

func TestFindKitchenOrderByIDUseCase_Success(t *testing.T) {
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	expectedOrder := createTestKitchenOrder(orderID, "order123", "001", status)
	dataStore.kitchenOrders = []entities.KitchenOrder{*expectedOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindKitchenOrderByIDUseCase(kitchenOrderGateway)

	result, err := useCase.Execute(orderID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.ID != orderID {
		t.Errorf("Expected ID %s, got %s", orderID, result.ID)
	}
}

func TestUpdateKitchenOrderUseCase_Success(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	orderID := "550e8400-e29b-41d4-a716-446655440000"
	newStatusID := constants.KITCHEN_ORDER_STATUS_PREPARING_ID

	dataStore := NewMockDataStore()
	originalStatus := dataStore.orderStatuses[0]

	existingOrder := createTestKitchenOrder(orderID, "order123", "001", originalStatus)
	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway, &MockMessageBroker{})

	result, err := useCase.Execute(dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: newStatusID,
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Status.ID != newStatusID {
		t.Errorf("Expected status ID %s, got %s", newStatusID, result.Status.ID)
	}
}

func TestFindAllOrderStatusUseCase_Success(t *testing.T) {
	dataStore := NewMockDataStore()
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewFindAllOrdersStatusUseCase(orderStatusGateway)

	result, err := useCase.Execute()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 4 {
		t.Errorf("Expected 4 statuses, got %d", len(result))
	}

	statusIDs := make(map[string]bool)
	for _, status := range result {
		statusIDs[status.ID] = true
	}

	expectedIDs := []string{
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		constants.KITCHEN_ORDER_STATUS_READY_ID,
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
	}

	for _, expectedID := range expectedIDs {
		if !statusIDs[expectedID] {
			t.Errorf("Expected status ID %s not found in result", expectedID)
		}
	}
}
