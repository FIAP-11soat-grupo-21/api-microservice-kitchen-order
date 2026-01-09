package entities

import (
	"testing"

	"tech_challenge/internal/domain/exceptions"
)

func TestNewOrderItem_Success(t *testing.T) {
	orderItem, err := NewOrderItem("item-1", "order-123", "product-1", 2, 12.75)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if orderItem.ID != "item-1" {
		t.Errorf("Expected ID 'item-1', got %s", orderItem.ID)
	}

	if orderItem.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", orderItem.OrderID)
	}

	if orderItem.ProductID != "product-1" {
		t.Errorf("Expected ProductID 'product-1', got %s", orderItem.ProductID)
	}

	if orderItem.Quantity != 2 {
		t.Errorf("Expected Quantity 2, got %d", orderItem.Quantity)
	}

	if orderItem.UnitPrice != 12.75 {
		t.Errorf("Expected UnitPrice 12.75, got %f", orderItem.UnitPrice)
	}
}

func TestNewOrderItem_EmptyID(t *testing.T) {
	_, err := NewOrderItem("", "order-123", "product-1", 2, 12.75)

	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}
}

func TestNewOrderItem_EmptyOrderID(t *testing.T) {
	_, err := NewOrderItem("item-1", "", "product-1", 2, 12.75)

	if err == nil {
		t.Error("Expected error for empty OrderID, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}
}

func TestNewOrderItem_EmptyProductID(t *testing.T) {
	_, err := NewOrderItem("item-1", "order-123", "", 2, 12.75)

	if err == nil {
		t.Error("Expected error for empty ProductID, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}
}

func TestNewOrderItem_InvalidQuantity(t *testing.T) {
	_, err := NewOrderItem("item-1", "order-123", "product-1", 0, 12.75)

	if err == nil {
		t.Error("Expected error for invalid quantity, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}
}

func TestNewOrderItem_NegativeUnitPrice(t *testing.T) {
	_, err := NewOrderItem("item-1", "order-123", "product-1", 2, -1.0)

	if err == nil {
		t.Error("Expected error for negative unit price, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}
}

func TestOrderItem_GetTotal(t *testing.T) {
	orderItem, _ := NewOrderItem("item-1", "order-123", "product-1", 3, 10.50)

	total := orderItem.GetTotal()

	expected := 31.50 
	if total != expected {
		t.Errorf("Expected total %f, got %f", expected, total)
	}
}

func TestOrderItem_IsEmpty_False(t *testing.T) {
	orderItem, _ := NewOrderItem("item-1", "order-123", "product-1", 2, 12.75)

	isEmpty := orderItem.IsEmpty()

	if isEmpty {
		t.Error("Expected IsEmpty to be false, got true")
	}
}

func TestOrderItem_IsEmpty_True(t *testing.T) {
	orderItem := &OrderItem{
		ID:        "",
		OrderID:   "order-123",
		ProductID: "product-1",
		Quantity:  2,
		UnitPrice: 12.75,
	}

	isEmpty := orderItem.IsEmpty()

	if !isEmpty {
		t.Error("Expected IsEmpty to be true, got false")
	}
}