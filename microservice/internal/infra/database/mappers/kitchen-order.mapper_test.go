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

func TestFromDAOToModelKitchenOrder_EmptyArray(t *testing.T) {
	models := []*models.KitchenOrderModel{}

	// Act
	daos := FromModelArrayToDAOArrayKitchenOrder(models)

	// Assert
	if len(daos) != 0 {
		t.Errorf("Expected 0 DAOs, got %d", len(daos))
	}
}

func TestFromDAOToModelKitchenOrder_WithItems(t *testing.T) {
	now := time.Now()
	customerID := "customer-123"

	dao := daos.KitchenOrderDAO{
		ID:         "test-id",
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     25.50,
		Slug:       "001",
		Status: daos.OrderStatusDAO{
			ID:   "status-id",
			Name: "Recebido",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-123",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 12.75,
			},
		},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	model := FromDAOToModelKitchenOrder(dao)

	if model.CustomerID == nil || *model.CustomerID != customerID {
		t.Errorf("Expected CustomerID %s, got %v", customerID, model.CustomerID)
	}

	if model.Amount != 25.50 {
		t.Errorf("Expected Amount 25.50, got %f", model.Amount)
	}

	if len(model.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(model.Items))
	}

	if model.Items[0].ID != "item-1" {
		t.Errorf("Expected Item ID 'item-1', got %s", model.Items[0].ID)
	}

	if model.Items[0].KitchenOrderID != "test-id" {
		t.Errorf("Expected Item KitchenOrderID 'test-id', got %s", model.Items[0].KitchenOrderID)
	}

	if model.UpdatedAt != nil {
		t.Errorf("Expected UpdatedAt to be nil, got %v", model.UpdatedAt)
	}
}

func TestFromModelToDAOKitchenOrder_WithItems(t *testing.T) {
	now := time.Now()
	customerID := "customer-123"

	model := &models.KitchenOrderModel{
		ID:         "test-id",
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     25.50,
		Slug:       "001",
		StatusID:   "status-id",
		Status: models.OrderStatusModel{
			ID:   "status-id",
			Name: "Recebido",
		},
		Items: []models.OrderItemModel{
			{
				ID:             "item-1",
				KitchenOrderID: "test-id",
				OrderID:        "order-123",
				ProductID:      "product-1",
				Quantity:       2,
				UnitPrice:      12.75,
			},
		},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	dao := FromModelToDAOKitchenOrder(model)

	if dao.CustomerID == nil || *dao.CustomerID != customerID {
		t.Errorf("Expected CustomerID %s, got %v", customerID, dao.CustomerID)
	}

	if dao.Amount != 25.50 {
		t.Errorf("Expected Amount 25.50, got %f", dao.Amount)
	}

	if len(dao.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(dao.Items))
	}

	if dao.Items[0].ID != "item-1" {
		t.Errorf("Expected Item ID 'item-1', got %s", dao.Items[0].ID)
	}

	if dao.UpdatedAt != nil {
		t.Errorf("Expected UpdatedAt to be nil, got %v", dao.UpdatedAt)
	}
}

func TestFromDAOToModelKitchenOrder_MultipleItems(t *testing.T) {
	now := time.Now()
	customerID := "customer-456"

	dao := daos.KitchenOrderDAO{
		ID:         "order-id-1",
		OrderID:    "order-789",
		CustomerID: &customerID,
		Amount:     100.00,
		Slug:       "005",
		Status: daos.OrderStatusDAO{
			ID:   "status-preparing",
			Name: "Em preparação",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-789",
				ProductID: "product-1",
				Quantity:  3,
				UnitPrice: 15.50,
			},
			{
				ID:        "item-2",
				OrderID:   "order-789",
				ProductID: "product-2",
				Quantity:  2,
				UnitPrice: 20.00,
			},
			{
				ID:        "item-3",
				OrderID:   "order-789",
				ProductID: "product-3",
				Quantity:  1,
				UnitPrice: 29.00,
			},
		},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	model := FromDAOToModelKitchenOrder(dao)

	if model.ID != "order-id-1" {
		t.Errorf("Expected ID 'order-id-1', got %s", model.ID)
	}

	if model.OrderID != "order-789" {
		t.Errorf("Expected OrderID 'order-789', got %s", model.OrderID)
	}

	if model.CustomerID == nil || *model.CustomerID != customerID {
		t.Errorf("Expected CustomerID %s, got %v", customerID, model.CustomerID)
	}

	if model.Amount != 100.00 {
		t.Errorf("Expected Amount 100.00, got %f", model.Amount)
	}

	if model.Slug != "005" {
		t.Errorf("Expected Slug '005', got %s", model.Slug)
	}

	if model.StatusID != "status-preparing" {
		t.Errorf("Expected StatusID 'status-preparing', got %s", model.StatusID)
	}

	if len(model.Items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(model.Items))
	}

	for i, item := range model.Items {
		if item.KitchenOrderID != "order-id-1" {
			t.Errorf("Expected Item[%d] KitchenOrderID 'order-id-1', got %s", i, item.KitchenOrderID)
		}
		if item.OrderID != "order-789" {
			t.Errorf("Expected Item[%d] OrderID 'order-789', got %s", i, item.OrderID)
		}
	}

	if model.Items[0].ProductID != "product-1" || model.Items[0].Quantity != 3 {
		t.Error("Item 0 mapping failed")
	}

	if model.Items[1].ProductID != "product-2" || model.Items[1].Quantity != 2 {
		t.Error("Item 1 mapping failed")
	}

	if model.Items[2].ProductID != "product-3" || model.Items[2].Quantity != 1 {
		t.Error("Item 2 mapping failed")
	}
}

func TestFromModelToDAOKitchenOrder_MultipleItems(t *testing.T) {
	now := time.Now()
	customerID := "customer-789"

	model := &models.KitchenOrderModel{
		ID:         "order-id-2",
		OrderID:    "order-999",
		CustomerID: &customerID,
		Amount:     150.75,
		Slug:       "010",
		StatusID:   "status-ready",
		Status: models.OrderStatusModel{
			ID:   "status-ready",
			Name: "Pronto",
		},
		Items: []models.OrderItemModel{
			{
				ID:             "item-a",
				KitchenOrderID: "order-id-2",
				OrderID:        "order-999",
				ProductID:      "product-a",
				Quantity:       5,
				UnitPrice:      10.00,
			},
			{
				ID:             "item-b",
				KitchenOrderID: "order-id-2",
				OrderID:        "order-999",
				ProductID:      "product-b",
				Quantity:       3,
				UnitPrice:      16.92,
			},
		},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	dao := FromModelToDAOKitchenOrder(model)

	if dao.ID != "order-id-2" {
		t.Errorf("Expected ID 'order-id-2', got %s", dao.ID)
	}

	if dao.OrderID != "order-999" {
		t.Errorf("Expected OrderID 'order-999', got %s", dao.OrderID)
	}

	if dao.CustomerID == nil || *dao.CustomerID != customerID {
		t.Errorf("Expected CustomerID %s, got %v", customerID, dao.CustomerID)
	}

	if dao.Amount != 150.75 {
		t.Errorf("Expected Amount 150.75, got %f", dao.Amount)
	}

	if dao.Slug != "010" {
		t.Errorf("Expected Slug '010', got %s", dao.Slug)
	}

	if dao.Status.ID != "status-ready" {
		t.Errorf("Expected Status.ID 'status-ready', got %s", dao.Status.ID)
	}

	if dao.Status.Name != "Pronto" {
		t.Errorf("Expected Status.Name 'Pronto', got %s", dao.Status.Name)
	}

	if len(dao.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(dao.Items))
	}

	if dao.Items[0].ID != "item-a" || dao.Items[0].ProductID != "product-a" {
		t.Error("Item 0 mapping failed")
	}

	if dao.Items[1].ID != "item-b" || dao.Items[1].ProductID != "product-b" {
		t.Error("Item 1 mapping failed")
	}
}

func TestFromModelArrayToDAOArrayKitchenOrder_MultipleModels(t *testing.T) {
	now := time.Now()
	customerID1 := "customer-1"
	customerID2 := "customer-2"
	customerID3 := "customer-3"

	modelArray := []*models.KitchenOrderModel{
		{
			ID:         "order-1",
			OrderID:    "order-001",
			CustomerID: &customerID1,
			Amount:     50.00,
			Slug:       "001",
			StatusID:   "status-1",
			Status: models.OrderStatusModel{
				ID:   "status-1",
				Name: "Recebido",
			},
			Items: []models.OrderItemModel{
				{
					ID:             "item-1",
					KitchenOrderID: "order-1",
					OrderID:        "order-001",
					ProductID:      "product-1",
					Quantity:       1,
					UnitPrice:      50.00,
				},
			},
			CreatedAt: now,
		},
		{
			ID:         "order-2",
			OrderID:    "order-002",
			CustomerID: &customerID2,
			Amount:     75.50,
			Slug:       "002",
			StatusID:   "status-2",
			Status: models.OrderStatusModel{
				ID:   "status-2",
				Name: "Em preparação",
			},
			Items: []models.OrderItemModel{
				{
					ID:             "item-2",
					KitchenOrderID: "order-2",
					OrderID:        "order-002",
					ProductID:      "product-2",
					Quantity:       2,
					UnitPrice:      37.75,
				},
			},
			CreatedAt: now,
		},
		{
			ID:         "order-3",
			OrderID:    "order-003",
			CustomerID: &customerID3,
			Amount:     100.00,
			Slug:       "003",
			StatusID:   "status-3",
			Status: models.OrderStatusModel{
				ID:   "status-3",
				Name: "Pronto",
			},
			Items: []models.OrderItemModel{},
			CreatedAt: now,
		},
	}

	daoArray := FromModelArrayToDAOArrayKitchenOrder(modelArray)

	if len(daoArray) != 3 {
		t.Fatalf("Expected 3 DAOs, got %d", len(daoArray))
	}

	if daoArray[0].ID != "order-1" {
		t.Errorf("Expected DAO[0].ID 'order-1', got %s", daoArray[0].ID)
	}

	if daoArray[0].Amount != 50.00 {
		t.Errorf("Expected DAO[0].Amount 50.00, got %f", daoArray[0].Amount)
	}

	if len(daoArray[0].Items) != 1 {
		t.Errorf("Expected DAO[0] to have 1 item, got %d", len(daoArray[0].Items))
	}

	if daoArray[1].ID != "order-2" {
		t.Errorf("Expected DAO[1].ID 'order-2', got %s", daoArray[1].ID)
	}

	if daoArray[1].Amount != 75.50 {
		t.Errorf("Expected DAO[1].Amount 75.50, got %f", daoArray[1].Amount)
	}

	if daoArray[2].ID != "order-3" {
		t.Errorf("Expected DAO[2].ID 'order-3', got %s", daoArray[2].ID)
	}

	if len(daoArray[2].Items) != 0 {
		t.Errorf("Expected DAO[2] to have 0 items, got %d", len(daoArray[2].Items))
	}
}

func TestFromDAOToModelKitchenOrder_NoCustomerID(t *testing.T) {
	now := time.Now()

	dao := daos.KitchenOrderDAO{
		ID:         "order-no-customer",
		OrderID:    "order-555",
		CustomerID: nil,
		Amount:     45.00,
		Slug:       "015",
		Status: daos.OrderStatusDAO{
			ID:   "status-id",
			Name: "Recebido",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	model := FromDAOToModelKitchenOrder(dao)

	if model.CustomerID != nil {
		t.Errorf("Expected CustomerID to be nil, got %v", model.CustomerID)
	}

	if model.ID != "order-no-customer" {
		t.Errorf("Expected ID 'order-no-customer', got %s", model.ID)
	}

	if len(model.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(model.Items))
	}
}

func TestFromModelToDAOKitchenOrder_NoCustomerID(t *testing.T) {
	now := time.Now()

	model := &models.KitchenOrderModel{
		ID:         "order-no-customer",
		OrderID:    "order-555",
		CustomerID: nil,
		Amount:     45.00,
		Slug:       "015",
		StatusID:   "status-id",
		Status: models.OrderStatusModel{
			ID:   "status-id",
			Name: "Recebido",
		},
		Items:     []models.OrderItemModel{},
		CreatedAt: now,
		UpdatedAt: nil,
	}

	dao := FromModelToDAOKitchenOrder(model)

	if dao.CustomerID != nil {
		t.Errorf("Expected CustomerID to be nil, got %v", dao.CustomerID)
	}

	if dao.ID != "order-no-customer" {
		t.Errorf("Expected ID 'order-no-customer', got %s", dao.ID)
	}

	if len(dao.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(dao.Items))
	}
}
