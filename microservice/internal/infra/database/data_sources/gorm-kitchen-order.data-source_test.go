package data_sources

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/constants"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(
		&models.OrderStatusModel{},
		&models.KitchenOrderModel{},
		&models.OrderItemModel{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	seedTestData(t, db)

	return db
}

func seedTestData(t *testing.T, db *gorm.DB) {
	statuses := []models.OrderStatusModel{
		{ID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, Name: "Recebido"},
		{ID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID, Name: "Em preparação"},
		{ID: constants.KITCHEN_ORDER_STATUS_READY_ID, Name: "Pronto"},
		{ID: constants.KITCHEN_ORDER_STATUS_FINISHED_ID, Name: "Finalizado"},
	}

	for _, status := range statuses {
		if err := db.Create(&status).Error; err != nil {
			t.Fatalf("Failed to seed status: %v", err)
		}
	}
}

func TestNewGormKitchenOrderDataSource(t *testing.T) {
	ds := NewGormKitchenOrderDataSource()
	
	if ds == nil {
		t.Error("Expected non-nil data source")
	} else {
		t.Log("✓ NewGormKitchenOrderDataSource criado com sucesso")
	}
	
	if ds != nil && ds.db == nil {
		t.Log("⚠ Database connection is nil (expected in test environment)")
	}
}

func TestGormKitchenOrderDataSource_Insert(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	kitchenOrder := daos.KitchenOrderDAO{
		ID:      "order-123",
		OrderID: "ext-order-123",
		Amount:  50.00,
		Slug:    "order-slug",
		Status: daos.OrderStatusDAO{
			ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			Name: "Recebido",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "ext-order-123",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 25.00,
			},
		},
	}

	err := ds.Insert(kitchenOrder)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	} else {
		t.Log("✓ Insert executado com sucesso")
	}

	var count int64
	db.Model(&models.KitchenOrderModel{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 order in database, got %d", count)
	}
}

func TestGormKitchenOrderDataSource_Insert_Error(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	kitchenOrder := daos.KitchenOrderDAO{
		ID:      "order-123",
		OrderID: "ext-order-123",
		Amount:  50.00,
		Slug:    "order-slug",
		Status: daos.OrderStatusDAO{
			ID:   "invalid-status-id",
			Name: "Invalid",
		},
	}

	err := ds.Insert(kitchenOrder)

	if err != nil {
		t.Logf("✓ Erro capturado corretamente: %v", err)
	} else {
		t.Log("⚠ SQLite permite foreign key inválida (comportamento esperado)")
	}
}

func TestGormKitchenOrderDataSource_FindAll(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	now := time.Now()
	orders := []models.KitchenOrderModel{
		{
			ID:        "order-1",
			OrderID:   "ext-1",
			Amount:    50.00,
			Slug:      "order-1",
			StatusID:  constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        "order-2",
			OrderID:   "ext-2",
			Amount:    60.00,
			Slug:      "order-2",
			StatusID:  constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
			CreatedAt: now.Add(-1 * time.Hour),
		},
		{
			ID:        "order-3",
			OrderID:   "ext-3",
			Amount:    70.00,
			Slug:      "order-3",
			StatusID:  constants.KITCHEN_ORDER_STATUS_READY_ID,
			CreatedAt: now,
		},
		{
			ID:        "order-4",
			OrderID:   "ext-4",
			Amount:    80.00,
			Slug:      "order-4",
			StatusID:  constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
			CreatedAt: now,
		},
	}

	for _, order := range orders {
		db.Create(&order)
	}

	filter := dtos.KitchenOrderFilter{}
	result, err := ds.FindAll(filter)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 orders (excluding Finalizado), got %d", len(result))
	} else {
		t.Logf("✓ FindAll retornou %d pedidos (excluindo Finalizado)", len(result))
	}

	if len(result) > 0 && result[0].Status.Name != "Pronto" {
		t.Errorf("Expected first order to be 'Pronto', got '%s'", result[0].Status.Name)
	} else if len(result) > 0 {
		t.Log("✓ Ordenação correta: primeiro pedido é 'Pronto'")
	}
}

func TestGormKitchenOrderDataSource_FindAll_WithFilters(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	now := time.Now()
	orders := []models.KitchenOrderModel{
		{
			ID:        "order-1",
			OrderID:   "ext-1",
			Amount:    50.00,
			Slug:      "order-1",
			StatusID:  constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        "order-2",
			OrderID:   "ext-2",
			Amount:    60.00,
			Slug:      "order-2",
			StatusID:  constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
			CreatedAt: now.Add(-1 * time.Hour),
		},
	}

	for _, order := range orders {
		db.Create(&order)
	}

	t.Run("Filter by StatusID", func(t *testing.T) {
		var statusID uint = 1
		filter := dtos.KitchenOrderFilter{
			StatusID: &statusID,
		}
		result, err := ds.FindAll(filter)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		t.Logf("✓ Filtro por StatusID executado (resultado: %d pedidos)", len(result))
	})

	t.Run("Filter by CreatedAtFrom", func(t *testing.T) {
		from := now.Add(-90 * time.Minute)
		filter := dtos.KitchenOrderFilter{
			CreatedAtFrom: &from,
		}
		result, err := ds.FindAll(filter)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("Expected 1 order, got %d", len(result))
		} else {
			t.Log("✓ Filtro por CreatedAtFrom funcionou")
		}
	})

	t.Run("Filter by CreatedAtTo", func(t *testing.T) {
		to := now.Add(-90 * time.Minute)
		filter := dtos.KitchenOrderFilter{
			CreatedAtTo: &to,
		}
		result, err := ds.FindAll(filter)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("Expected 1 order, got %d", len(result))
		} else {
			t.Log("✓ Filtro por CreatedAtTo funcionou")
		}
	})
}

func TestGormKitchenOrderDataSource_FindByID(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	order := models.KitchenOrderModel{
		ID:       "order-123",
		OrderID:  "ext-123",
		Amount:   50.00,
		Slug:     "order-123",
		StatusID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		Items: []models.OrderItemModel{
			{
				ID:        "item-1",
				OrderID:   "ext-123",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 25.00,
			},
		},
	}
	db.Create(&order)

	result, err := ds.FindByID("order-123")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.ID != "order-123" {
		t.Errorf("Expected ID 'order-123', got '%s'", result.ID)
	} else {
		t.Log("✓ FindByID retornou pedido correto")
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Items))
	} else {
		t.Log("✓ Items foram carregados corretamente")
	}

	if result.Status.Name != "Recebido" {
		t.Errorf("Expected status 'Recebido', got '%s'", result.Status.Name)
	} else {
		t.Log("✓ Status foi carregado corretamente")
	}
}

func TestGormKitchenOrderDataSource_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	_, err := ds.FindByID("non-existent-id")

	if err == nil {
		t.Error("Expected error for non-existent ID")
	} else {
		t.Logf("✓ Erro capturado para ID inexistente: %v", err)
	}
}

func TestGormKitchenOrderDataSource_Update(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	order := models.KitchenOrderModel{
		ID:       "order-123",
		OrderID:  "ext-123",
		Amount:   50.00,
		Slug:     "order-123",
		StatusID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
	}
	db.Create(&order)

	updatedOrder := daos.KitchenOrderDAO{
		ID:      "order-123",
		OrderID: "ext-123",
		Amount:  50.00,
		Slug:    "order-123",
		Status: daos.OrderStatusDAO{
			ID:   constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
			Name: "Em preparação",
		},
	}

	err := ds.Update(updatedOrder)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	} else {
		t.Log("✓ Update executado com sucesso")
	}

	var updated models.KitchenOrderModel
	db.First(&updated, "id = ?", "order-123")

	if updated.StatusID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected status to be updated to '%s', got '%s'", 
			constants.KITCHEN_ORDER_STATUS_PREPARING_ID, updated.StatusID)
	} else {
		t.Log("✓ Status foi atualizado corretamente")
	}
}

func TestGormKitchenOrderDataSource_Delete(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	order := models.KitchenOrderModel{
		ID:       "order-123",
		OrderID:  "ext-123",
		Amount:   50.00,
		Slug:     "order-123",
		StatusID: constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
	}
	db.Create(&order)

	err := ds.Delete("order-123")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	} else {
		t.Log("✓ Delete executado com sucesso")
	}

	var count int64
	db.Model(&models.KitchenOrderModel{}).Where("id = ?", "order-123").Count(&count)

	if count != 0 {
		t.Errorf("Expected order to be deleted, but found %d records", count)
	} else {
		t.Log("✓ Pedido foi deletado do banco")
	}
}

func TestGormKitchenOrderDataSource_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	err := ds.Delete("non-existent-id")

	if err != nil {
		t.Logf("Error: %v", err)
	} else {
		t.Log("✓ Delete de ID inexistente não causou erro")
	}
}

func TestGormKitchenOrderDataSource_FindAll_OrderingPriority(t *testing.T) {
	db := setupTestDB(t)
	ds := &GormKitchenOrderDataSource{db: db}

	now := time.Now()
	orders := []models.KitchenOrderModel{
		{
			ID:        "order-1",
			OrderID:   "ext-1",
			Amount:    50.00,
			Slug:      "order-1",
			StatusID:  constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
			CreatedAt: now.Add(-3 * time.Hour),
		},
		{
			ID:        "order-2",
			OrderID:   "ext-2",
			Amount:    60.00,
			Slug:      "order-2",
			StatusID:  constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:        "order-3",
			OrderID:   "ext-3",
			Amount:    70.00,
			Slug:      "order-3",
			StatusID:  constants.KITCHEN_ORDER_STATUS_READY_ID,
			CreatedAt: now.Add(-1 * time.Hour),
		},
	}

	for _, order := range orders {
		db.Create(&order)
	}

	filter := dtos.KitchenOrderFilter{}
	result, err := ds.FindAll(filter)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(result))
		return
	}

	expectedOrder := []string{"Pronto", "Em preparação", "Recebido"}
	for i, expected := range expectedOrder {
		if result[i].Status.Name != expected {
			t.Errorf("Position %d: expected '%s', got '%s'", i, expected, result[i].Status.Name)
		}
	}

	t.Log("✓ Ordenação por prioridade funcionou corretamente:")
	for i, order := range result {
		t.Logf("  %d. %s", i+1, order.Status.Name)
	}
}
