package use_cases

import (
	"testing"

	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/shared/config/constants"
)

func TestFindAllOrdersStatusUseCase_Success(t *testing.T) {
	dataStore := NewMockDataStore()
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)

	useCase := NewFindAllOrdersStatusUseCase(orderStatusGateway)

	result, err := useCase.Execute()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 results, got %d", len(result))
	}

	if result[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected first status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, result[0].ID)
	}
}

func TestFindAllOrdersStatusUseCase_Empty(t *testing.T) {
	dataStore := &MockDataStore{
		kitchenOrders: []entities.KitchenOrder{},
		orderStatuses: []entities.OrderStatus{},
	}
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)

	useCase := NewFindAllOrdersStatusUseCase(orderStatusGateway)

	result, err := useCase.Execute()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 results, got %d", len(result))
	}
}