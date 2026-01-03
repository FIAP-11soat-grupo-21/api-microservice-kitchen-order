package mappers

import (
	"testing"
	"time"

	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
)

func TestFromDAOToModelKitchenOrder(t *testing.T) {
	// Arrange
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	
	dao := daos.KitchenOrderDAO{
		ID:      "test-id",
		OrderID: "order-123",
		Slug:    "001",
		Status: daos.OrderStatusDAO{
			ID:   "status-id",
			Name: "Recebido",
		},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	// Act
	model := FromDAOToModelKitchenOrder(dao)

	// Assert
	if model.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", model.ID)
	}

	if model.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", model.OrderID)
	}

	if model.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", model.Slug)
	}

	if model.StatusID != "status-id" {
		t.Errorf("Expected StatusID 'status-id', got %s", model.StatusID)
	}

	if model.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, model.CreatedAt)
	}

	if model.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set, got nil")
	}

	if *model.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, *model.UpdatedAt)
	}
}

func TestFromModelToDAOKitchenOrder(t *testing.T) {
	// Arrange
	now := time.Now()
	updatedAt := now.Add(time.Hour)
	
	model := &models.KitchenOrderModel{
		ID:       "test-id",
		OrderID:  "order-123",
		Slug:     "001",
		StatusID: "status-id",
		Status: models.OrderStatusModel{
			ID:   "status-id",
			Name: "Recebido",
		},
		CreatedAt: now,
		UpdatedAt: &updatedAt,
	}

	// Act
	dao := FromModelToDAOKitchenOrder(model)

	// Assert
	if dao.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", dao.ID)
	}

	if dao.OrderID != "order-123" {
		t.Errorf("Expected OrderID 'order-123', got %s", dao.OrderID)
	}

	if dao.Slug != "001" {
		t.Errorf("Expected Slug '001', got %s", dao.Slug)
	}

	if dao.Status.ID != "status-id" {
		t.Errorf("Expected Status.ID 'status-id', got %s", dao.Status.ID)
	}

	if dao.Status.Name != "Recebido" {
		t.Errorf("Expected Status.Name 'Recebido', got %s", dao.Status.Name)
	}

	if dao.CreatedAt != now {
		t.Errorf("Expected CreatedAt %v, got %v", now, dao.CreatedAt)
	}

	if dao.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set, got nil")
	}

	if *dao.UpdatedAt != updatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", updatedAt, *dao.UpdatedAt)
	}
}

func TestFromModelArrayToDAOArrayKitchenOrder(t *testing.T) {
	// Arrange
	now := time.Now()
	
	models := []*models.KitchenOrderModel{
		{
			ID:       "test-id-1",
			OrderID:  "order-123",
			Slug:     "001",
			StatusID: "status-id",
			Status: models.OrderStatusModel{
				ID:   "status-id",
				Name: "Recebido",
			},
			CreatedAt: now,
		},
		{
			ID:       "test-id-2",
			OrderID:  "order-456",
			Slug:     "002",
			StatusID: "status-id-2",
			Status: models.OrderStatusModel{
				ID:   "status-id-2",
				Name: "Em preparação",
			},
			CreatedAt: now,
		},
	}

	// Act
	daos := FromModelArrayToDAOArrayKitchenOrder(models)

	// Assert
	if len(daos) != 2 {
		t.Errorf("Expected 2 DAOs, got %d", len(daos))
	}

	if daos[0].ID != "test-id-1" {
		t.Errorf("Expected first DAO ID 'test-id-1', got %s", daos[0].ID)
	}

	if daos[0].Slug != "001" {
		t.Errorf("Expected first DAO Slug '001', got %s", daos[0].Slug)
	}

	if daos[1].ID != "test-id-2" {
		t.Errorf("Expected second DAO ID 'test-id-2', got %s", daos[1].ID)
	}

	if daos[1].Slug != "002" {
		t.Errorf("Expected second DAO Slug '002', got %s", daos[1].Slug)
	}
}

func TestFromModelArrayToDAOArrayKitchenOrder_EmptyArray(t *testing.T) {
	// Arrange
	models := []*models.KitchenOrderModel{}

	// Act
	daos := FromModelArrayToDAOArrayKitchenOrder(models)

	// Assert
	if len(daos) != 0 {
		t.Errorf("Expected 0 DAOs, got %d", len(daos))
	}
}