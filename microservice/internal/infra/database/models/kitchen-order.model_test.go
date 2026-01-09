package models

import (
	"testing"
	"time"
)

func TestOrderItemModel_TableName(t *testing.T) {
	model := OrderItemModel{}

	tableName := model.TableName()

	expected := "order_item"
	if tableName != expected {
		t.Errorf("Expected table name '%s', got '%s'", expected, tableName)
	}
}

func TestKitchenOrderModel_TableName(t *testing.T) {
	model := KitchenOrderModel{}

	tableName := model.TableName()

	expected := "kitchen_order"
	if tableName != expected {
		t.Errorf("Expected table name '%s', got '%s'", expected, tableName)
	}
}

func TestOrderItemModel_Struct(t *testing.T) {
	model := OrderItemModel{
		ID:             "item-1",
		KitchenOrderID: "kitchen-order-1",
		OrderID:        "order-123",
		ProductID:      "product-1",
		Quantity:       2,
		UnitPrice:      12.75,
	}

	if model.ID != "item-1" {
		t.Errorf("Expected ID 'item-1', got %s", model.ID)
	}

	if model.KitchenOrderID != "kitchen-order-1" {
		t.Errorf("Expected KitchenOrderID 'kitchen-order-1', got %s", model.KitchenOrderID)
	}

	if model.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", model.OrderID)
	}

	if model.ProductID != "product-1" {
		t.Errorf("Expected ProductID 'product-1', got %s", model.ProductID)
	}

	if model.Quantity != 2 {
		t.Errorf("Expected Quantity 2, got %d", model.Quantity)
	}

	if model.UnitPrice != 12.75 {
		t.Errorf("Expected UnitPrice 12.75, got %f", model.UnitPrice)
	}
}

func TestKitchenOrderModel_Struct(t *testing.T) {
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	customerID := "customer-123"

	status := OrderStatusModel{
		ID:   "status-1",
		Name: "Recebido",
	}

	items := []OrderItemModel{
		{
			ID:             "item-1",
			KitchenOrderID: "kitchen-order-1",
			OrderID:        "order-123",
			ProductID:      "product-1",
			Quantity:       2,
			UnitPrice:      12.75,
		},
	}

	model := KitchenOrderModel{
		ID:         "kitchen-order-1",
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     25.50,
		Slug:       "001",
		StatusID:   "status-1",
		Status:     status,
		Items:      items,
		CreatedAt:  now,
		UpdatedAt:  &updatedAt,
	}

	if model.ID != "kitchen-order-1" {
		t.Errorf("Expected ID 'kitchen-order-1', got %s", model.ID)
	}

	if model.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", model.OrderID)
	}

	if model.CustomerID == nil || *model.CustomerID != "customer-123" {
		t.Errorf("Expected CustomerID 'customer-123', got %v", model.CustomerID)
	}

	if model.Amount != 25.50 {
		t.Errorf("Expected Amount 25.50, got %f", model.Amount)
	}

	if model.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", model.Slug)
	}

	if model.StatusID != "status-1" {
		t.Errorf("Expected StatusID 'status-1', got %s", model.StatusID)
	}

	if model.Status.ID != "status-1" {
		t.Errorf("Expected Status ID 'status-1', got %s", model.Status.ID)
	}

	if len(model.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(model.Items))
	}

	if model.Items[0].ID != "item-1" {
		t.Errorf("Expected Item ID 'item-1', got %s", model.Items[0].ID)
	}

	if model.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, model.CreatedAt)
	}

	if model.UpdatedAt == nil || *model.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, model.UpdatedAt)
	}
}