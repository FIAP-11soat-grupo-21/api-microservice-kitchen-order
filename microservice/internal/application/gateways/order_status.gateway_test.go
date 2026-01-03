package gateways

import (
	"testing"

	"tech_challenge/internal/daos"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

// Mock DataSource para OrderStatus
type MockOrderStatusDataSource struct {
	findByIDFunc func(string) (daos.OrderStatusDAO, error)
	findAllFunc  func() ([]daos.OrderStatusDAO, error)
}

func (m *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.OrderStatusDAO{}, nil
}

func (m *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc()
	}
	return []daos.OrderStatusDAO{}, nil
}

func TestNewOrderStatusGateway(t *testing.T) {
	// Arrange
	mockDataSource := &MockOrderStatusDataSource{}

	// Act
	gateway := NewOrderStatusGateway(mockDataSource)

	// Assert
	if gateway == nil {
		t.Error("Expected gateway to be created, got nil")
	}
}

func TestOrderStatusGateway_FindAll(t *testing.T) {
	// Arrange
	expectedDAOs := []daos.OrderStatusDAO{
		{ID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, Name: "Recebido"},
		{ID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID, Name: "Em preparação"},
		{ID: constants.KITCHEN_ORDER_STATUS_READY_ID, Name: "Pronto"},
		{ID: constants.KITCHEN_ORDER_STATUS_FINISHED_ID, Name: "Finalizado"},
	}

	mockDataSource := &MockOrderStatusDataSource{
		findAllFunc: func() ([]daos.OrderStatusDAO, error) {
			return expectedDAOs, nil
		},
	}

	gateway := NewOrderStatusGateway(mockDataSource)

	// Act
	result, err := gateway.FindAll()

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

func TestOrderStatusGateway_FindByID(t *testing.T) {
	// Arrange
	expectedDAO := daos.OrderStatusDAO{
		ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		Name: "Recebido",
	}

	mockDataSource := &MockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			if id == constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
				return expectedDAO, nil
			}
			return daos.OrderStatusDAO{}, &exceptions.OrderStatusNotFoundException{}
		},
	}

	gateway := NewOrderStatusGateway(mockDataSource)

	// Act
	result, err := gateway.FindByID(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, result.ID)
	}

	if result.Name.Value() != "Recebido" {
		t.Errorf("Expected name 'Recebido', got %s", result.Name.Value())
	}
}

func TestOrderStatusGateway_FindByID_NotFound(t *testing.T) {
	// Arrange
	mockDataSource := &MockOrderStatusDataSource{
		findByIDFunc: func(id string) (daos.OrderStatusDAO, error) {
			return daos.OrderStatusDAO{}, &exceptions.OrderStatusNotFoundException{}
		},
	}

	gateway := NewOrderStatusGateway(mockDataSource)

	// Act
	_, err := gateway.FindByID("invalid-id")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid ID, got nil")
	}

	if _, ok := err.(*exceptions.OrderStatusNotFoundException); !ok {
		t.Errorf("Expected OrderStatusNotFoundException, got %T", err)
	}
}