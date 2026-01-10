package use_cases

import (
	"context"
	"os"
	"testing"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/shared/config/constants"
	"tech_challenge/internal/shared/interfaces"
)

func setupTestEnv() {
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("MESSAGE_BROKER_TYPE", "rabbitmq")
	os.Setenv("RABBITMQ_URL", "amqp://localhost:5672")
	os.Setenv("AWS_SQS_KITCHEN_ORDERS_QUEUE", "https://sqs.us-east-1.amazonaws.com/123456789/test-queue")
	os.Setenv("AWS_SQS_ORDERS_QUEUE", "https://sqs.us-east-1.amazonaws.com/123456789/orders-queue")
}

func cleanupTestEnv() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"AWS_REGION", "MESSAGE_BROKER_TYPE", "RABBITMQ_URL",
		"AWS_SQS_KITCHEN_ORDERS_QUEUE", "AWS_SQS_ORDERS_QUEUE",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

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

func TestFindAllKitchenOrdersUseCase_Success(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	order1, _ := entities.NewKitchenOrder(
		"id1", "order1", "001", status, time.Now(), nil,
	)
	order2, _ := entities.NewKitchenOrder(
		"id2", "order2", "002", status, time.Now(), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*order1, *order2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)
	filter := dtos.KitchenOrderFilter{}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(result))
	}
}

func TestFindKitchenOrderByIDUseCase_Success(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	expectedOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", status, time.Now(), nil,
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
}

func TestUpdateKitchenOrderUseCase_Success(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	newStatusID := constants.KITCHEN_ORDER_STATUS_PREPARING_ID

	dataStore := NewMockDataStore()
	originalStatus := dataStore.orderStatuses[0] // Status "Recebido"

	existingOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", originalStatus, time.Now(), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)

	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway, &MockMessageBroker{})
	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: newStatusID,
	}

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status.ID != newStatusID {
		t.Errorf("Expected status ID %s, got %s", newStatusID, result.Status.ID)
	}
}

func TestFindAllOrderStatusUseCase_Success(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewFindAllOrdersStatusUseCase(orderStatusGateway)

	// Act
	result, err := useCase.Execute()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 statuses, got %d", len(result))
	}

	// Verifica se todos os status esperados estão presentes
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
