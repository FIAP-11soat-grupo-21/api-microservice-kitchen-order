package handlers

import (
	"testing"
)

func TestNewKitchenOrderHandler(t *testing.T) {
	handler := NewKitchenOrderHandler()

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}