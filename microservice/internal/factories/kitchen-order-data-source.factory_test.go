package factories

import (
	"testing"
)

func TestNewKitchenOrderDataSource(t *testing.T) {
	// Act
	dataSource := NewKitchenOrderDataSource()

	// Assert
	if dataSource == nil {
		t.Error("Expected data source to be created, got nil")
	}
}

func TestNewOrderStatusDataSource(t *testing.T) {
	// Act
	dataSource := NewOrderStatusDataSource()

	// Assert
	if dataSource == nil {
		t.Error("Expected data source to be created, got nil")
	}
}
