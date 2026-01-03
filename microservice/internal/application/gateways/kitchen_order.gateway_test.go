package gateways

import (
	"testing"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

// Mock DataSource para testes
type MockKitchenOrderDataSource struct {
	insertFunc    func(daos.KitchenOrderDAO) error
	findByIDFunc  func(string) (daos.KitchenOrderDAO, error)
	findAllFunc   func(dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error)
	updateFunc    func(daos.KitchenOrderDAO) error
}

func (m *MockKitchenOrderDataSource) Insert(order daos.KitchenOrderDAO) error {
	if m.insertFunc != nil {
		return m.insertFunc(order)
	}
	return nil
}

func (m *MockKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.KitchenOrderDAO{}, nil
}

func (m *MockKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(filter)
	}
	return []daos.KitchenOrderDAO{}, nil
}

func (m *MockKitchenOrderDataSource) Update(order daos.KitchenOrderDAO) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func TestNewKitchenOrderGateway(t *testing.T) {
	// Arrange
	mockDataSource := &MockKitchenOrderDataSource{}

	// Act
	gateway := NewKitchenOrderGateway(mockDataSource)

	// Assert
	if gateway == nil {
		t.Error("Expected gateway to be created, got nil")
	}
}

func TestKitchenOrderGateway_Insert(t *testing.T) {
	// Arrange
	mockDataSource := &MockKitchenOrderDataSource{
		insertFunc: func(order daos.KitchenOrderDAO) error {
			if order.ID == "test-id" {
				return nil
			}
			return &exceptions.InvalidKitchenOrderDataException{}
		},
	}

	gateway := NewKitchenOrderGateway(mockDataSource)
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	order, _ := entities.NewKitchenOrder("test-id", "order-123", "001", *status, time.Now(), nil)

	// Act
	err := gateway.Insert(*order)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestKitchenOrderGateway_FindByID(t *testing.T) {
	// Arrange
	expectedDAO := daos.KitchenOrderDAO{
		ID:      "test-id",
		OrderID: "order-123",
		Slug:    "001",
		Status: daos.OrderStatusDAO{
			ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			Name: "Recebido",
		},
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	mockDataSource := &MockKitchenOrderDataSource{
		findByIDFunc: func(id string) (daos.KitchenOrderDAO, error) {
			if id == "test-id" {
				return expectedDAO, nil
			}
			return daos.KitchenOrderDAO{}, &exceptions.KitchenOrderNotFoundException{}
		},
	}

	gateway := NewKitchenOrderGateway(mockDataSource)

	// Act
	result, err := gateway.FindByID("test-id")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", result.ID)
	}

	if result.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", result.OrderID)
	}
}

func TestKitchenOrderGateway_FindAll(t *testing.T) {
	// Arrange
	expectedDAOs := []daos.KitchenOrderDAO{
		{
			ID:      "test-id-1",
			OrderID: "order-123",
			Slug:    "001",
			Status: daos.OrderStatusDAO{
				ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
				Name: "Recebido",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:      "test-id-2",
			OrderID: "order-456",
			Slug:    "002",
			Status: daos.OrderStatusDAO{
				ID:   constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
				Name: "Em preparação",
			},
			CreatedAt: time.Now(),
		},
	}

	mockDataSource := &MockKitchenOrderDataSource{
		findAllFunc: func(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
			return expectedDAOs, nil
		},
	}

	gateway := NewKitchenOrderGateway(mockDataSource)
	filter := dtos.KitchenOrderFilter{}

	// Act
	result, err := gateway.FindAll(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(result))
	}

	if result[0].ID != "test-id-1" {
		t.Errorf("Expected first order ID 'test-id-1', got %s", result[0].ID)
	}
}

func TestKitchenOrderGateway_Update(t *testing.T) {
	// Arrange
	mockDataSource := &MockKitchenOrderDataSource{
		updateFunc: func(order daos.KitchenOrderDAO) error {
			if order.ID == "test-id" {
				return nil
			}
			return &exceptions.KitchenOrderNotFoundException{}
		},
	}

	gateway := NewKitchenOrderGateway(mockDataSource)
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	now := time.Now()
	order, _ := entities.NewKitchenOrder("test-id", "order-123", "001", *status, time.Now(), &now)

	// Act
	err := gateway.Update(*order)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}