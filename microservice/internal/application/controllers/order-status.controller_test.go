package controllers

import (
	"testing"

	"tech_challenge/internal/daos"
	"tech_challenge/internal/shared/config/constants"
)

func TestNewOrderStatusController(t *testing.T) {
	mockOrderStatusDS := &MockOrderStatusDataSource{}

	controller := NewOrderStatusController(mockOrderStatusDS)

	if controller == nil {
		t.Error("Expected controller to be created, got nil")
	}

	if controller.dataSource != mockOrderStatusDS {
		t.Error("Expected dataSource to be set correctly")
	}
}

func TestOrderStatusController_FindAll(t *testing.T) {
	mockOrderStatusDS := &MockOrderStatusDataSource{
		orderStatuses: []daos.OrderStatusDAO{
			{
				ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
				Name: "Recebido",
			},
			{
				ID:   constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
				Name: "Em preparação",
			},
		},
	}

	controller := NewOrderStatusController(mockOrderStatusDS)

	result, err := controller.FindAll()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}

	if result[0].ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected first status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, result[0].ID)
	}

	if result[1].ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected second status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, result[1].ID)
	}
}

func TestOrderStatusController_FindAll_Empty(t *testing.T) {
	mockOrderStatusDS := &MockOrderStatusDataSource{
		orderStatuses: []daos.OrderStatusDAO{},
	}

	controller := NewOrderStatusController(mockOrderStatusDS)

	result, err := controller.FindAll()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 results, got %d", len(result))
	}
}