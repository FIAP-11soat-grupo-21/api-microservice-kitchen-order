package handlers

import (
	"testing"
)

func TestNewOrderStatusHandler(t *testing.T) {
	handler := NewOrderStatusHandler()

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}