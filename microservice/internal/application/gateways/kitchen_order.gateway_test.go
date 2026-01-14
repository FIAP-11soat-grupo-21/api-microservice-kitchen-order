package gateways

import (
	"errors"
	"testing"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

type MockKitchenOrderDataSource struct {
	insertFunc   func(daos.KitchenOrderDAO) error
	findByIDFunc func(string) (daos.KitchenOrderDAO, error)
	findAllFunc  func(dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error)
	updateFunc   func(daos.KitchenOrderDAO) error
}

func (m *MockKitchenOrderDataSource) Insert(order daos.KitchenOrderDAO) error {
	if m.insertFunc != nil {
		return m.insertFunc(order)
	}
	return nil
}

func (m *MockKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return daos.KitchenOrderDAO{}, nil
}

func (m *MockKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(filter)
	}
	return []daos.KitchenOrderDAO{}, nil
}

func (m *MockKitchenOrderDataSource) Update(order daos.KitchenOrderDAO) error {
	if m.updateFunc != nil {
		return m.updateFunc(order)
	}
	return nil
}

func makeOrderItemDAO(orderID string) daos.OrderItemDAO {
	return daos.OrderItemDAO{
		ID:        "item-1",
		OrderID:   orderID,
		ProductID: "product-1",
		Quantity:  2,
		UnitPrice: 10.0,
	}
}

func strPtr(s string) *string {
	return &s
}

func TestNewKitchenOrderGateway(t *testing.T) {
	gateway := NewKitchenOrderGateway(&MockKitchenOrderDataSource{})
	if gateway == nil {
		t.Fatal("expected gateway to be created")
	}
}

func TestKitchenOrderGateway_Insert_Success(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		insertFunc: func(order daos.KitchenOrderDAO) error {
			return nil
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	status, _ := entities.NewOrderStatus(
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		"Recebido",
	)

	item, _ := entities.NewOrderItem("item-1", "order-1", "product-1", 1, 10)

	order, _ := entities.NewKitchenOrderWithOrderData(
		"id-1",
		"order-1",
		"001",
		strPtr("customer-1"),
		10,
		[]entities.OrderItem{*item},
		*status,
		time.Now(),
		nil,
	)

	err := gateway.Insert(*order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKitchenOrderGateway_Insert_Error(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		insertFunc: func(order daos.KitchenOrderDAO) error {
			return &exceptions.InvalidKitchenOrderDataException{}
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	status, _ := entities.NewOrderStatus(
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		"Recebido",
	)

	order, _ := entities.NewKitchenOrder(
		"id-1",
		"order-1",
		"001",
		*status,
		time.Now(),
		nil,
	)

	err := gateway.Insert(*order)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestKitchenOrderGateway_FindByID_Success(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		findByIDFunc: func(id string) (daos.KitchenOrderDAO, error) {
			return daos.KitchenOrderDAO{
				ID:      "id-1",
				OrderID: "order-1",
				Slug:    "001",
				CustomerID: strPtr("customer-1"),
				Status: daos.OrderStatusDAO{
					ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
					Name: "Recebido",
				},
				Items: []daos.OrderItemDAO{
					makeOrderItemDAO("order-1"),
				},
				CreatedAt: time.Now(),
			}, nil
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	order, err := gateway.FindByID("id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.ID != "id-1" {
		t.Errorf("unexpected id: %s", order.ID)
	}
}

func TestKitchenOrderGateway_FindByID_NotFound(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		findByIDFunc: func(id string) (daos.KitchenOrderDAO, error) {
			return daos.KitchenOrderDAO{}, &exceptions.KitchenOrderNotFoundException{}
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	_, err := gateway.FindByID("x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestKitchenOrderGateway_FindByID_InvalidItem(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		findByIDFunc: func(id string) (daos.KitchenOrderDAO, error) {
			return daos.KitchenOrderDAO{
				ID:      "id-1",
				OrderID: "order-1",
				Slug:    "001",
				Status: daos.OrderStatusDAO{
					ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
					Name: "Recebido",
				},
				Items: []daos.OrderItemDAO{
					{
						ID:        "",
						OrderID:   "",
						ProductID: "",
						Quantity:  -1,
						UnitPrice: -10,
					},
				},
				CreatedAt: time.Now(),
			}, nil
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	_, err := gateway.FindByID("id-1")
	if err == nil {
		t.Fatal("expected domain error, got nil")
	}
}

func TestKitchenOrderGateway_FindAll_Success(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		findAllFunc: func(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
			return []daos.KitchenOrderDAO{
				{
					ID:      "id-1",
					OrderID: "order-1",
					Slug:    "001",
					Status: daos.OrderStatusDAO{
						ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
						Name: "Recebido",
					},
					Items: []daos.OrderItemDAO{
						makeOrderItemDAO("order-1"),
					},
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	orders, err := gateway.FindAll(dtos.KitchenOrderFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("expected 1 order")
	}
}

func TestKitchenOrderGateway_FindAll_Error(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		findAllFunc: func(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
			return nil, errors.New("db error")
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	_, err := gateway.FindAll(dtos.KitchenOrderFilter{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestKitchenOrderGateway_Update_Success(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		updateFunc: func(order daos.KitchenOrderDAO) error {
			return nil
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	status, _ := entities.NewOrderStatus(
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		"Em preparação",
	)

	now := time.Now()

	order, _ := entities.NewKitchenOrderWithOrderData(
		"id-1",
		"order-1",
		"001",
		strPtr("customer-1"),
		10,
		[]entities.OrderItem{},
		*status,
		now,
		&now,
	)

	err := gateway.Update(*order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKitchenOrderGateway_Update_NotFound(t *testing.T) {
	mock := &MockKitchenOrderDataSource{
		updateFunc: func(order daos.KitchenOrderDAO) error {
			return &exceptions.KitchenOrderNotFoundException{}
		},
	}

	gateway := NewKitchenOrderGateway(mock)

	status, _ := entities.NewOrderStatus(
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		"Em preparação",
	)

	now := time.Now()

	order, _ := entities.NewKitchenOrderWithOrderData(
		"id-1",
		"order-1",
		"001",
		strPtr("customer-1"),
		10,
		[]entities.OrderItem{},
		*status,
		now,
		&now,
	)

	err := gateway.Update(*order)
	if err == nil {
		t.Fatal("expected error")
	}
}
