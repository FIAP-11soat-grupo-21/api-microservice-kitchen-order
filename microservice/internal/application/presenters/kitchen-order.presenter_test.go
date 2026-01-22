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

	if response.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", response.OrderID)
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

func TestToResponse_WithItems(t *testing.T) {
	// Arrange
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	now := time.Now()

	item1 := entities.OrderItem{
		ID:        "item-1",
		OrderID:   "order-123",
		ProductID: "product-1",
		Quantity:  2,
		UnitPrice: 15.50,
	}

	item2 := entities.OrderItem{
		ID:        "item-2",
		OrderID:   "order-123",
		ProductID: "product-2",
		Quantity:  1,
		UnitPrice: 25.00,
	}

	kitchenOrder, _ := entities.NewKitchenOrderWithOrderData(
		"order-id-1",
		"order-123",
		"005",
		nil,
		55.00,
		[]entities.OrderItem{item1, item2},
		*status,
		now,
		nil,
	)

	// Act
	response := ToResponse(*kitchenOrder)

	// Assert
	if response.ID != "order-id-1" {
		t.Errorf("Expected ID 'order-id-1', got %s", response.ID)
	}

	if response.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", response.OrderID)
	}

	if response.Slug != "005" {
		t.Errorf("Expected Slug '005', got %s", response.Slug)
	}

	if response.Status.ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected Status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, response.Status.ID)
	}

	if response.Status.Name != "Em preparação" {
		t.Errorf("Expected Status Name 'Em preparação', got %s", response.Status.Name)
	}
}

func TestToResponse_WithCustomerID(t *testing.T) {
	// Arrange
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto")
	now := time.Now()
	customerID := "customer-456"

	kitchenOrder, _ := entities.NewKitchenOrderWithOrderData(
		"order-id-2",
		"order-456",
		"010",
		&customerID,
		100.00,
		[]entities.OrderItem{},
		*status,
		now,
		nil,
	)

	// Act
	response := ToResponse(*kitchenOrder)

	// Assert
	if response.ID != "order-id-2" {
		t.Errorf("Expected ID 'order-id-2', got %s", response.ID)
	}

	if response.OrderID != "order-456" {
		t.Errorf("Expected OrderID 'order-456', got %s", response.OrderID)
	}

	if response.Slug != "010" {
		t.Errorf("Expected Slug '010', got %s", response.Slug)
	}
}

func TestToResponse_DifferentStatuses(t *testing.T) {
	// Arrange
	now := time.Now()

	testCases := []struct {
		statusID   string
		statusName string
	}{
		{constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido"},
		{constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação"},
		{constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto"},
		{constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado"},
	}

	for _, tc := range testCases {
		status, _ := entities.NewOrderStatus(tc.statusID, tc.statusName)
		kitchenOrder, _ := entities.NewKitchenOrder("id", "order", "001", *status, now, nil)

		// Act
		response := ToResponse(*kitchenOrder)

		// Assert
		if response.Status.ID != tc.statusID {
			t.Errorf("Expected Status ID %s, got %s", tc.statusID, response.Status.ID)
		}

		if response.Status.Name != tc.statusName {
			t.Errorf("Expected Status Name %s, got %s", tc.statusName, response.Status.Name)
		}
	}
}

func TestToResponseList_MultipleOrders(t *testing.T) {
	// Arrange
	status1, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	status2, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	status3, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto")
	now := time.Now()

	order1, _ := entities.NewKitchenOrder("id1", "order1", "001", *status1, now, nil)
	order2, _ := entities.NewKitchenOrder("id2", "order2", "002", *status2, now, nil)
	order3, _ := entities.NewKitchenOrder("id3", "order3", "003", *status3, now, nil)

	orders := []entities.KitchenOrder{*order1, *order2, *order3}

	// Act
	responses := ToResponseList(orders)

	// Assert
	if len(responses) != 3 {
		t.Errorf("Expected 3 responses, got %d", len(responses))
	}

	expectedData := []struct {
		id       string
		orderID  string
		slug     string
		statusID string
	}{
		{"id1", "order1", "001", constants.KITCHEN_ORDER_STATUS_RECEIVED_ID},
		{"id2", "order2", "002", constants.KITCHEN_ORDER_STATUS_PREPARING_ID},
		{"id3", "order3", "003", constants.KITCHEN_ORDER_STATUS_READY_ID},
	}

	for i, expected := range expectedData {
		if responses[i].ID != expected.id {
			t.Errorf("Expected response[%d].ID %s, got %s", i, expected.id, responses[i].ID)
		}

		if responses[i].OrderID != expected.orderID {
			t.Errorf("Expected response[%d].OrderID %s, got %s", i, expected.orderID, responses[i].OrderID)
		}

		if responses[i].Slug != expected.slug {
			t.Errorf("Expected response[%d].Slug %s, got %s", i, expected.slug, responses[i].Slug)
		}

		if responses[i].Status.ID != expected.statusID {
			t.Errorf("Expected response[%d].Status.ID %s, got %s", i, expected.statusID, responses[i].Status.ID)
		}
	}
}

func TestToResponseList_WithItemsInOrders(t *testing.T) {
	// Arrange
	status, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	now := time.Now()

	item1 := entities.OrderItem{
		ID:        "item-1",
		OrderID:   "order-1",
		ProductID: "product-1",
		Quantity:  1,
		UnitPrice: 10.00,
	}

	item2 := entities.OrderItem{
		ID:        "item-2",
		OrderID:   "order-2",
		ProductID: "product-2",
		Quantity:  2,
		UnitPrice: 20.00,
	}

	order1, _ := entities.NewKitchenOrderWithOrderData(
		"id1",
		"order-1",
		"001",
		nil,
		10.00,
		[]entities.OrderItem{item1},
		*status,
		now,
		nil,
	)

	order2, _ := entities.NewKitchenOrderWithOrderData(
		"id2",
		"order-2",
		"002",
		nil,
		40.00,
		[]entities.OrderItem{item2},
		*status,
		now,
		nil,
	)

	orders := []entities.KitchenOrder{*order1, *order2}

	// Act
	responses := ToResponseList(orders)

	// Assert
	if len(responses) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(responses))
	}

	if responses[0].ID != "id1" {
		t.Errorf("Expected response[0].ID 'id1', got %s", responses[0].ID)
	}

	if responses[1].ID != "id2" {
		t.Errorf("Expected response[1].ID 'id2', got %s", responses[1].ID)
	}
}
