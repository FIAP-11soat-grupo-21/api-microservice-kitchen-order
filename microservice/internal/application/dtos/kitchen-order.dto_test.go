package dtos

import (
	"testing"
	"time"
)

func TestKitchenOrderDTO_Structure(t *testing.T) {
	// Arrange
	statusDTO := OrderStatusDTO{
		ID:   "status-id",
		Name: "Recebido",
	}

	dto := KitchenOrderDTO{
		ID:      "kitchen-order-id",
		OrderID: "order-id",
		Slug:    "001",
		Status:  statusDTO,
	}

	// Assert
	if dto.ID != "kitchen-order-id" {
		t.Errorf("Expected ID 'kitchen-order-id', got %s", dto.ID)
	}

	if dto.OrderID != "order-id" {
		t.Errorf("Expected OrderID 'order-id', got %s", dto.OrderID)
	}

	if dto.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", dto.Slug)
	}

	if dto.Status.ID != "status-id" {
		t.Errorf("Expected Status.ID 'status-id', got %s", dto.Status.ID)
	}

	if dto.Status.Name != "Recebido" {
		t.Errorf("Expected Status.Name 'Recebido', got %s", dto.Status.Name)
	}
}

func TestCreateKitchenOrderDTO_Structure(t *testing.T) {
	// Arrange
	dto := CreateKitchenOrderDTO{
		OrderID: "order-123",
	}

	// Assert
	if dto.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", dto.OrderID)
	}
}

func TestUpdateKitchenOrderDTO_Structure(t *testing.T) {
	// Arrange
	dto := UpdateKitchenOrderDTO{
		ID:       "kitchen-order-id",
		StatusID: "new-status-id",
	}

	// Assert
	if dto.ID != "kitchen-order-id" {
		t.Errorf("Expected ID 'kitchen-order-id', got %s", dto.ID)
	}

	if dto.StatusID != "new-status-id" {
		t.Errorf("Expected StatusID 'new-status-id', got %s", dto.StatusID)
	}
}

func TestKitchenOrderFilter_Structure(t *testing.T) {
	// Arrange
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := now
	statusID := uint(1)

	filter := KitchenOrderFilter{
		CreatedAtFrom: &from,
		CreatedAtTo:   &to,
		StatusID:      &statusID,
	}

	// Assert
	if filter.CreatedAtFrom == nil {
		t.Error("Expected CreatedAtFrom to be set, got nil")
	}

	if filter.CreatedAtTo == nil {
		t.Error("Expected CreatedAtTo to be set, got nil")
	}

	if filter.StatusID == nil {
		t.Error("Expected StatusID to be set, got nil")
	}

	if *filter.StatusID != statusID {
		t.Errorf("Expected StatusID %d, got %d", statusID, *filter.StatusID)
	}
}

func TestKitchenOrderFilter_NilValues(t *testing.T) {
	// Arrange
	filter := KitchenOrderFilter{
		CreatedAtFrom: nil,
		CreatedAtTo:   nil,
		StatusID:      nil,
	}

	// Assert
	if filter.CreatedAtFrom != nil {
		t.Error("Expected CreatedAtFrom to be nil, got value")
	}

	if filter.CreatedAtTo != nil {
		t.Error("Expected CreatedAtTo to be nil, got value")
	}

	if filter.StatusID != nil {
		t.Error("Expected StatusID to be nil, got value")
	}
}

func TestKitchenOrderResponseDTO_Structure(t *testing.T) {
	// Arrange
	now := time.Now()
	updatedAt := now.Add(time.Hour)

	statusDTO := OrderStatusDTO{
		ID:   "status-id",
		Name: "Recebido",
	}

	dto := KitchenOrderResponseDTO{
		ID:        "kitchen-order-id",
		OrderID:   "order-id",
		Slug:      "001",
		Status:    statusDTO,
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	// Assert
	if dto.ID != "kitchen-order-id" {
		t.Errorf("Expected ID 'kitchen-order-id', got %s", dto.ID)
	}

	if dto.OrderID != "order-id" {
		t.Errorf("Expected OrderID 'order-id', got %s", dto.OrderID)
	}

	if dto.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", dto.Slug)
	}

	if dto.Status.ID != "status-id" {
		t.Errorf("Expected Status.ID 'status-id', got %s", dto.Status.ID)
	}

	if dto.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, dto.CreatedAt)
	}

	if dto.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set, got nil")
	}

	if *dto.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, *dto.UpdatedAt)
	}
}

func TestKitchenOrderResponseDTO_NilUpdatedAt(t *testing.T) {
	// Arrange
	now := time.Now()

	statusDTO := OrderStatusDTO{
		ID:   "status-id",
		Name: "Recebido",
	}

	dto := KitchenOrderResponseDTO{
		ID:        "kitchen-order-id",
		OrderID:   "order-id",
		Slug:      "001",
		Status:    statusDTO,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Assert
	if dto.UpdatedAt != nil {
		t.Errorf("Expected UpdatedAt to be nil, got %v", dto.UpdatedAt)
	}
}