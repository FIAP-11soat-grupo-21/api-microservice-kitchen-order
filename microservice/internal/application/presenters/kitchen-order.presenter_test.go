package presenters

import (
	"testing"
	"time"

	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/shared/config/constants"
)

func TestToResponse(t *testing.T) {
	// Arrange
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	
	kitchenOrder, _ := entities.NewKitchenOrder(
		"test-id",
		"order-123",
		"001",
		*status,
		now,
		&updatedAt,
	)

	// Act
	response := ToResponse(*kitchenOrder)

	// Assert
	if response.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", response.ID)
	}

	if response.OrderID != "test-id" {
		t.Errorf("Expected OrderID 'test-id', got %s", response.OrderID)
	}

	if response.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", response.Slug)
	}

	if response.Status.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		t.Errorf("Expected Status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, response.Status.ID)
	}

	if response.Status.Name != "Recebido" {
		t.Errorf("Expected Status Name 'Recebido', got %s", response.Status.Name)
	}

	if response.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, response.CreatedAt)
	}

	if response.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set, got nil")
	}

	if *response.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, *response.UpdatedAt)
	}
}

func TestToResponse_NilUpdatedAt(t *testing.T) {
	// Arrange
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	now := time.Now()
	
	kitchenOrder, _ := entities.NewKitchenOrder(
		"test-id",
		"order-123",
		"001",
		*status,
		now,
		nil,
	)

	// Act
	response := ToResponse(*kitchenOrder)

	// Assert
	if response.UpdatedAt != nil {
		t.Errorf("Expected UpdatedAt to be nil, got %v", response.UpdatedAt)
	}
}

func TestToResponseList(t *testing.T) {
	// Arrange
	status1, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	status2, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	now := time.Now()
	
	order1, _ := entities.NewKitchenOrder("id1", "order1", "001", *status1, now, nil)
	order2, _ := entities.NewKitchenOrder("id2", "order2", "002", *status2, now, nil)
	
	orders := []entities.KitchenOrder{*order1, *order2}

	// Act
	responses := ToResponseList(orders)

	// Assert
	if len(responses) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(responses))
	}

	if responses[0].ID != "id1" {
		t.Errorf("Expected first response ID 'id1', got %s", responses[0].ID)
	}

	if responses[0].Slug != "001" {
		t.Errorf("Expected first response Slug '001', got %s", responses[0].Slug)
	}

	if responses[1].ID != "id2" {
		t.Errorf("Expected second response ID 'id2', got %s", responses[1].ID)
	}

	if responses[1].Slug != "002" {
		t.Errorf("Expected second response Slug '002', got %s", responses[1].Slug)
	}
}

func TestToResponseList_EmptyList(t *testing.T) {
	// Arrange
	orders := []entities.KitchenOrder{}

	// Act
	responses := ToResponseList(orders)

	// Assert
	if len(responses) != 0 {
		t.Errorf("Expected 0 responses, got %d", len(responses))
	}
}