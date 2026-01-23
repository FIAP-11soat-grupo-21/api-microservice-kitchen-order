// +build integration

package seed

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tech_challenge/internal/infra/database/models"
)

func TestGormDBWrapper_Integration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	err = db.AutoMigrate(&models.OrderStatusModel{})
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	wrapper := &GormDBWrapper{db: db}

	t.Run("First", func(t *testing.T) {
		var model models.OrderStatusModel
		result := wrapper.First(&model)
		
		if result == nil {
			t.Error("Expected non-nil result from First")
		} else {
			t.Log("✓ First executado")
		}
		
		err := result.GetError()
		if err != gorm.ErrRecordNotFound {
			t.Logf("Expected ErrRecordNotFound, got: %v", err)
		}
	})

	t.Run("Create", func(t *testing.T) {
		model := models.OrderStatusModel{
			ID:   "test-id-123",
			Name: "Test Status",
		}
		
		result := wrapper.Create(&model)
		
		if result == nil {
			t.Error("Expected non-nil result from Create")
		} else {
			t.Log("✓ Create executado")
		}
		
		err := result.GetError()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Where", func(t *testing.T) {
		result := wrapper.Where("id = ?", "test-id-123")
		
		if result == nil {
			t.Error("Expected non-nil result from Where")
		} else {
			t.Log("✓ Where executado")
		}
	})

	t.Run("GetError", func(t *testing.T) {
		err := wrapper.GetError()
		
		if err != nil {
			t.Logf("Error: %v", err)
		} else {
			t.Log("✓ GetError retornou nil")
		}
	})

	t.Run("ChainedOperations", func(t *testing.T) {
		var model models.OrderStatusModel
		result := wrapper.Where("id = ?", "test-id-123").First(&model)
		
		if result == nil {
			t.Error("Expected non-nil result")
		}
		
		err := result.GetError()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		
		if model.ID != "test-id-123" {
			t.Errorf("Expected ID 'test-id-123', got '%s'", model.ID)
		}
		
		if model.Name != "Test Status" {
			t.Errorf("Expected Name 'Test Status', got '%s'", model.Name)
		}
		
		t.Log("✓ Operações encadeadas funcionaram corretamente")
	})
}
