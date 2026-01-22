package controllers

import (
	"context"
	"testing"
	"time"

	"tech_challenge/internal"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/shared/config/constants"
	"tech_challenge/internal/shared/interfaces"
)

// Mock implementations
type MockKitchenOrderDataSource struct {
	kitchenOrders []daos.KitchenOrderDAO
}

func (m *MockKitchenOrderDataSource) Insert(kitchenOrder daos.KitchenOrderDAO) error {
	m.kitchenOrders = append(m.kitchenOrders, kitchenOrder)
	return nil
}

func (m *MockKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	return m.kitchenOrders, nil
}

func (m *MockKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	for _, order := range m.kitchenOrders {
		if order.ID == id {
			return order, nil
		}
	}
	return daos.KitchenOrderDAO{}, nil
}

func (m *MockKitchenOrderDataSource) Update(kitchenOrder daos.KitchenOrderDAO) error {
	for i, order := range m.kitchenOrders {
		if order.ID == kitchenOrder.ID {
			m.kitchenOrders[i] = kitchenOrder
			return nil
		}
	}
	return nil
}

func (m *MockKitchenOrderDataSource) Delete(id string) error {
	return nil
}

type MockOrderStatusDataSource struct {
	orderStatuses []daos.OrderStatusDAO
}

func (m *MockOrderStatusDataSource) Insert(orderStatus daos.OrderStatusDAO) error {
	return nil
}

func (m *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	return m.orderStatuses, nil
}

func (m *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	for _, status := range m.orderStatuses {
		if status.ID == id {
			return status, nil
		}
	}
	return daos.OrderStatusDAO{}, nil
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

// Test helpers
func createTestController() (*KitchenOrderController, *MockKitchenOrderDataSource, *MockOrderStatusDataSource) {
	mockKitchenOrderDS := &MockKitchenOrderDataSource{kitchenOrders: []daos.KitchenOrderDAO{}}
	mockOrderStatusDS := &MockOrderStatusDataSource{
		orderStatuses: []daos.OrderStatusDAO{
			{ID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, Name: "Recebido"},
			{ID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID, Name: "Em preparação"},
		},
	}
	mockMessageBroker := &MockMessageBroker{}
	controller := NewKitchenOrderController(mockKitchenOrderDS, mockOrderStatusDS, mockMessageBroker)
	return controller, mockKitchenOrderDS, mockOrderStatusDS
}

func createTestKitchenOrder(id string) daos.KitchenOrderDAO {
	return daos.KitchenOrderDAO{
		ID:      id,
		OrderID: "order-123",
		Slug:    "001",
		Status: daos.OrderStatusDAO{
			ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			Name: "Recebido",
		},
		CreatedAt: time.Now(),
	}
}

// Tests
func TestNewKitchenOrderController(t *testing.T) {
	controller, mockKitchenOrderDS, mockOrderStatusDS := createTestController()

	if controller == nil {
		t.Fatal("Expected controller to be created, got nil")
	}
	if controller.kitchenOrderDataSource != mockKitchenOrderDS {
		t.Error("Expected kitchenOrderDataSource to be set correctly")
	}
	if controller.orderStatusDataSource != mockOrderStatusDS {
		t.Error("Expected orderStatusDataSource to be set correctly")
	}
}

func TestKitchenOrderController_Create(t *testing.T) {
	controller, _, _ := createTestController()

	createDTO := dtos.CreateKitchenOrderDTO{
		OrderID: "order-123",
	}

	result, err := controller.Create(createDTO)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", result.OrderID)
	}
	if result.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", result.Slug)
	}
}

func TestKitchenOrderController_FindAll(t *testing.T) {
	controller, mockKitchenOrderDS, _ := createTestController()
	
	testOrder := createTestKitchenOrder("550e8400-e29b-41d4-a716-446655440000")
	mockKitchenOrderDS.kitchenOrders = []daos.KitchenOrderDAO{testOrder}

	result, err := controller.FindAll(dtos.KitchenOrderFilter{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
	if result[0].ID != testOrder.ID {
		t.Errorf("Expected ID '%s', got %s", testOrder.ID, result[0].ID)
	}
}

func TestKitchenOrderController_FindByID(t *testing.T) {
	controller, mockKitchenOrderDS, _ := createTestController()
	
	testOrder := createTestKitchenOrder("550e8400-e29b-41d4-a716-446655440000")
	mockKitchenOrderDS.kitchenOrders = []daos.KitchenOrderDAO{testOrder}

	result, err := controller.FindByID(testOrder.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.ID != testOrder.ID {
		t.Errorf("Expected ID '%s', got %s", testOrder.ID, result.ID)
	}
}

func TestKitchenOrderController_Update(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()
	
	controller, mockKitchenOrderDS, _ := createTestController()
	
	testOrder := createTestKitchenOrder("550e8400-e29b-41d4-a716-446655440000")
	mockKitchenOrderDS.kitchenOrders = []daos.KitchenOrderDAO{testOrder}

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       testOrder.ID,
		StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
	}

	result, err := controller.Update(updateDTO)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.ID != testOrder.ID {
		t.Errorf("Expected ID '%s', got %s", testOrder.ID, result.ID)
	}
	if result.Status.ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected Status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, result.Status.ID)
	}
}