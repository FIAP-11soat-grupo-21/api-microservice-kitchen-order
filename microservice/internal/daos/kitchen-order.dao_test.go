package daos

import (
	"testing"
	"time"
)

func TestKitchenOrderDAO_Struct(t *testing.T) {
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	customerID := "customer-123"

	status := OrderStatusDAO{
		ID:   "status-1",
		Name: "Recebido",
	}

	items := []OrderItemDAO{
		{
			ID:        "item-1",
			OrderID:   "order-123",
			ProductID: "product-1",
			Quantity:  2,
			UnitPrice: 12.75,
		},
	}

	dao := KitchenOrderDAO{
		ID:         "kitchen-order-1",
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     25.50,
		Status:     status,
		Slug:       "001",
		Items:      items,
		CreatedAt:  now,
		UpdatedAt:  &updatedAt,
	}

	if dao.ID != "kitchen-order-1" {
		t.Errorf("Expected ID 'kitchen-order-1', got %s", dao.ID)
	}

	if dao.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", dao.OrderID)
	}

	if dao.CustomerID == nil || *dao.CustomerID != "customer-123" {
		t.Errorf("Expected CustomerID 'customer-123', got %v", dao.CustomerID)
	}

	if dao.Amount != 25.50 {
		t.Errorf("Expected Amount 25.50, got %f", dao.Amount)
	}

	if dao.Status.ID != "status-1" {
		t.Errorf("Expected Status ID 'status-1', got %s", dao.Status.ID)
	}

	if dao.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", dao.Slug)
	}

	if len(dao.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(dao.Items))
	}

	if dao.Items[0].ID != "item-1" {
		t.Errorf("Expected Item ID 'item-1', got %s", dao.Items[0].ID)
	}

	if dao.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, dao.CreatedAt)
	}

	if dao.UpdatedAt == nil || *dao.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, dao.UpdatedAt)
	}
}

func TestOrderItemDAO_Struct(t *testing.T) {
	dao := OrderItemDAO{
		ID:        "item-1",
		OrderID:   "order-123",
		ProductID: "product-1",
		Quantity:  2,
		UnitPrice: 12.75,
	}

	if dao.ID != "item-1" {
		t.Errorf("Expected ID 'item-1', got %s", dao.ID)
	}

	if dao.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", dao.OrderID)
	}

	if dao.ProductID != "product-1" {
		t.Errorf("Expected ProductID 'product-1', got %s", dao.ProductID)
	}

	if dao.Quantity != 2 {
		t.Errorf("Expected Quantity 2, got %d", dao.Quantity)
	}

	if dao.UnitPrice != 12.75 {
		t.Errorf("Expected UnitPrice 12.75, got %f", dao.UnitPrice)
	}
}