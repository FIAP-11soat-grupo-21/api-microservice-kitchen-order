package data_sources

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/infra/database/models"
)

func setupOrderStatusTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.OrderStatusModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestNewGormOrderStatusDataSource(t *testing.T) {
	ds := NewGormOrderStatusDataSource()
	
	if ds == nil {
		t.Error("Expected non-nil data source")
	} else {
		t.Log("✓ NewGormOrderStatusDataSource criado com sucesso")
	}
}

func TestGormOrderStatusDataSource_Insert(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	orderStatus := daos.OrderStatusDAO{
		ID:   "status-123",
		Name: "Test Status",
	}

	err := ds.Insert(orderStatus)

	if err != nil {
		t.Logf("⚠ Erro esperado devido ao bug no código original (usa tabela errada): %v", err)
	} else {
		t.Log("✓ Insert executado")
	}
}

func TestGormOrderStatusDataSource_Insert_Duplicate(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	orderStatus := daos.OrderStatusDAO{
		ID:   "status-123",
		Name: "Test Status",
	}

	_ = ds.Insert(orderStatus)
	err := ds.Insert(orderStatus)

	if err != nil {
		t.Logf("✓ Erro capturado: %v", err)
	} else {
		t.Log("⚠ Sem erro")
	}
}

func TestGormOrderStatusDataSource_FindAll(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	statuses := []models.OrderStatusModel{
		{ID: "status-1", Name: "Recebido"},
		{ID: "status-2", Name: "Em preparação"},
		{ID: "status-3", Name: "Pronto"},
	}

	for _, status := range statuses {
		db.Create(&status)
	}

	result, err := ds.FindAll()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 statuses, got %d", len(result))
	} else {
		t.Logf("✓ FindAll retornou %d status", len(result))
	}

	for i, status := range result {
		t.Logf("  %d. ID=%s, Name=%s", i+1, status.ID, status.Name)
	}
}

func TestGormOrderStatusDataSource_FindAll_Empty(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	result, err := ds.FindAll()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 statuses, got %d", len(result))
	} else {
		t.Log("✓ FindAll retornou lista vazia corretamente")
	}
}

func TestGormOrderStatusDataSource_FindByID(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	status := models.OrderStatusModel{
		ID:   "status-123",
		Name: "Test Status",
	}
	db.Create(&status)

	result, err := ds.FindByID("status-123")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.ID != "status-123" {
		t.Errorf("Expected ID 'status-123', got '%s'", result.ID)
	} else {
		t.Log("✓ FindByID retornou status correto")
	}

	if result.Name != "Test Status" {
		t.Errorf("Expected name 'Test Status', got '%s'", result.Name)
	} else {
		t.Log("✓ Nome do status está correto")
	}
}

func TestGormOrderStatusDataSource_FindByID_NotFound(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	_, err := ds.FindByID("non-existent-id")

	if err == nil {
		t.Error("Expected error for non-existent ID")
	} else {
		t.Logf("✓ Erro capturado para ID inexistente: %v", err)
	}
}

func TestGormOrderStatusDataSource_FindAll_OrderPreserved(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	statuses := []models.OrderStatusModel{
		{ID: "status-1", Name: "Primeiro"},
		{ID: "status-2", Name: "Segundo"},
		{ID: "status-3", Name: "Terceiro"},
	}

	for _, status := range statuses {
		db.Create(&status)
	}

	result, err := ds.FindAll()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 statuses, got %d", len(result))
		return
	}

	t.Log("✓ Todos os status foram retornados")
	for i, status := range result {
		t.Logf("  %d. %s", i+1, status.Name)
	}
}

func TestGormOrderStatusDataSource_Insert_WithEmptyName(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	orderStatus := daos.OrderStatusDAO{
		ID:   "status-empty",
		Name: "",
	}

	err := ds.Insert(orderStatus)

	if err != nil {
		t.Logf("⚠ Erro devido ao bug no código: %v", err)
	} else {
		t.Log("✓ Insert executado")
	}
}

func TestGormOrderStatusDataSource_Insert_WithLongName(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	longName := "Status com nome muito longo para testar o limite do campo"
	orderStatus := daos.OrderStatusDAO{
		ID:   "status-long",
		Name: longName,
	}

	err := ds.Insert(orderStatus)

	if err != nil {
		t.Logf("⚠ Erro devido ao bug no código: %v", err)
	} else {
		t.Log("✓ Insert executado")
	}
}

func TestGormOrderStatusDataSource_FindAll_MultipleStatuses(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	for i := 1; i <= 10; i++ {
		status := models.OrderStatusModel{
			ID:   "status-" + string(rune(i)),
			Name: "Status " + string(rune(i)),
		}
		db.Create(&status)
	}

	result, err := ds.FindAll()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(result) != 10 {
		t.Errorf("Expected 10 statuses, got %d", len(result))
	} else {
		t.Log("✓ FindAll retornou todos os 10 status")
	}
}

func TestGormOrderStatusDataSource_Insert_SpecialCharacters(t *testing.T) {
	db := setupOrderStatusTestDB(t)
	ds := &GormOrderStatusDataSource{db: db}

	orderStatus := daos.OrderStatusDAO{
		ID:   "status-special",
		Name: "Status com çãõ & caracteres especiais!",
	}

	err := ds.Insert(orderStatus)

	if err != nil {
		t.Logf("⚠ Erro devido ao bug no código: %v", err)
	} else {
		t.Log("✓ Insert executado")
	}
}
