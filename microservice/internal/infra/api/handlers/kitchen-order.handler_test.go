package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/daos"
	shared_interfaces "tech_challenge/internal/shared/interfaces"
)

type MockKitchenOrderDataSource struct {
	mock.Mock
}

func (m *MockKitchenOrderDataSource) Insert(kitchenOrder daos.KitchenOrderDAO) error {
	args := m.Called(kitchenOrder)
	return args.Error(0)
}

func (m *MockKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return daos.KitchenOrderDAO{}, args.Error(1)
	}
	return args.Get(0).(daos.KitchenOrderDAO), args.Error(1)
}

func (m *MockKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]daos.KitchenOrderDAO), args.Error(1)
}

func (m *MockKitchenOrderDataSource) Update(kitchenOrder daos.KitchenOrderDAO) error {
	args := m.Called(kitchenOrder)
	return args.Error(0)
}

type MockOrderStatusDataSource struct {
	mock.Mock
}

func (m *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return daos.OrderStatusDAO{}, args.Error(1)
	}
	return args.Get(0).(daos.OrderStatusDAO), args.Error(1)
}

func (m *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]daos.OrderStatusDAO), args.Error(1)
}

type MockMessageBroker struct {
	mock.Mock
}

func (m *MockMessageBroker) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageBroker) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMessageBroker) Publish(ctx context.Context, queue string, message shared_interfaces.Message) error {
	args := m.Called(ctx, queue, message)
	return args.Error(0)
}

func (m *MockMessageBroker) Subscribe(ctx context.Context, queue string, handler shared_interfaces.MessageHandler) error {
	args := m.Called(ctx, queue, handler)
	return args.Error(0)
}

func (m *MockMessageBroker) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageBroker) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func setupTestEnv() {
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8080")
	os.Setenv("API_HOST", "localhost")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("MESSAGE_BROKER_TYPE", "rabbitmq")
	os.Setenv("RABBITMQ_URL", "amqp://localhost:5672")
	os.Setenv("AWS_SQS_KITCHEN_ORDERS_QUEUE", "https://sqs.us-east-1.amazonaws.com/123456789/test-queue")
	os.Setenv("AWS_SQS_ORDERS_QUEUE", "https://sqs.us-east-1.amazonaws.com/123456789/orders-queue")
}

func cleanupTestEnv() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"AWS_REGION", "MESSAGE_BROKER_TYPE", "RABBITMQ_URL",
		"AWS_SQS_KITCHEN_ORDERS_QUEUE", "AWS_SQS_ORDERS_QUEUE",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestNewKitchenOrderHandler(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	t.Run("should create handler successfully", func(t *testing.T) {
		handler := NewKitchenOrderHandler()

		assert.NotNil(t, handler)
		assert.NotNil(t, handler.kitchenOrderController)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		handler1 := NewKitchenOrderHandler()
		handler2 := NewKitchenOrderHandler()

		assert.NotNil(t, handler1)
		assert.NotNil(t, handler2)
		assert.NotSame(t, handler1, handler2)
	})
}

func TestFindAll_Success(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	now := time.Now()
	kitchenOrders := []daos.KitchenOrderDAO{
		{
			ID:         "550e8400-e29b-41d4-a716-446655440000",
			OrderID:    "order-001",
			CustomerID: nil,
			Amount:     100.50,
			Slug:       "001",
			Status: daos.OrderStatusDAO{
				ID:   "1",
				Name: "Recebido",
			},
			Items: []daos.OrderItemDAO{
				{
					ID:        "item-001",
					OrderID:   "order-001",
					ProductID: "prod-001",
					Quantity:  2,
					UnitPrice: 50.25,
				},
			},
			CreatedAt: now,
			UpdatedAt: &now,
		},
	}

	mockDataSource.On("FindAll", mock.Anything).Return(kitchenOrders, nil)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.GET("/kitchen-orders", handler.FindAll)

	req, _ := http.NewRequest("GET", "/kitchen-orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response[0]["id"])
	assert.Equal(t, "order-001", response[0]["order_id"])
	assert.Equal(t, 100.50, response[0]["amount"])

	mockDataSource.AssertExpectations(t)
}

func TestFindAll_WithFilters(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	mockDataSource.On("FindAll", mock.MatchedBy(func(filter dtos.KitchenOrderFilter) bool {
		return filter.CreatedAtFrom != nil && filter.CreatedAtTo != nil
	})).Return([]daos.KitchenOrderDAO{}, nil)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.GET("/kitchen-orders", handler.FindAll)

	req, _ := http.NewRequest("GET", "/kitchen-orders?created_at_from=2024-01-01T00:00:00Z&created_at_to=2024-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}

func TestFindByID_Success(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	now := time.Now()
	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	kitchenOrder := daos.KitchenOrderDAO{
		ID:         kitchenOrderID,
		OrderID:    "order-001",
		CustomerID: nil,
		Amount:     100.50,
		Slug:       "001",
		Status: daos.OrderStatusDAO{
			ID:   "1",
			Name: "Recebido",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-001",
				OrderID:   "order-001",
				ProductID: "prod-001",
				Quantity:  2,
				UnitPrice: 50.25,
			},
		},
		CreatedAt: now,
		UpdatedAt: &now,
	}

	mockDataSource.On("FindByID", kitchenOrderID).Return(kitchenOrder, nil)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.GET("/kitchen-orders/:id", handler.FindByID)

	req, _ := http.NewRequest("GET", "/kitchen-orders/"+kitchenOrderID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, kitchenOrderID, response["id"])
	assert.Equal(t, "order-001", response["order_id"])
	assert.Equal(t, 100.50, response["amount"])
	assert.Equal(t, "Recebido", response["status"])

	mockDataSource.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	now := time.Now()
	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	updatedKitchenOrder := daos.KitchenOrderDAO{
		ID:         kitchenOrderID,
		OrderID:    "order-001",
		CustomerID: nil,
		Amount:     100.50,
		Slug:       "001",
		Status: daos.OrderStatusDAO{
			ID:   "2",
			Name: "Em preparação",
		},
		Items:     []daos.OrderItemDAO{},
		CreatedAt: now,
		UpdatedAt: &now,
	}

	mockDataSource.On("Update", mock.Anything).Return(nil)
	mockDataSource.On("FindByID", kitchenOrderID).Return(updatedKitchenOrder, nil)
	mockStatusDataSource.On("FindByID", "2").Return(daos.OrderStatusDAO{ID: "2", Name: "Em preparação"}, nil)
	mockMessageBroker.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{
		"status_id": "2",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/kitchen-orders/"+kitchenOrderID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, kitchenOrderID, response["id"])
	assert.Equal(t, "Em preparação", response["status"])

	mockDataSource.AssertExpectations(t)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.PUT("/kitchen-orders/:id", handler.Update)

	req, _ := http.NewRequest("PUT", "/kitchen-orders/ko-001", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_MissingStatusID(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/kitchen-orders/ko-001", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}



func TestFindByID_WithMultipleItems(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	now := time.Now()
	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	kitchenOrder := daos.KitchenOrderDAO{
		ID:         kitchenOrderID,
		OrderID:    "order-001",
		CustomerID: nil,
		Amount:     150.75,
		Slug:       "001",
		Status: daos.OrderStatusDAO{
			ID:   "1",
			Name: "Recebido",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-001",
				OrderID:   "order-001",
				ProductID: "prod-001",
				Quantity:  2,
				UnitPrice: 50.25,
			},
			{
				ID:        "item-002",
				OrderID:   "order-001",
				ProductID: "prod-002",
				Quantity:  1,
				UnitPrice: 50.25,
			},
			{
				ID:        "item-003",
				OrderID:   "order-001",
				ProductID: "prod-003",
				Quantity:  3,
				UnitPrice: 16.75,
			},
		},
		CreatedAt: now,
		UpdatedAt: &now,
	}

	mockDataSource.On("FindByID", kitchenOrderID).Return(kitchenOrder, nil)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.GET("/kitchen-orders/:id", handler.FindByID)

	req, _ := http.NewRequest("GET", "/kitchen-orders/"+kitchenOrderID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	items := response["items"].([]interface{})
	assert.Len(t, items, 3)

	for i, item := range items {
		itemMap := item.(map[string]interface{})
		assert.NotEmpty(t, itemMap["id"])
		assert.Equal(t, "order-001", itemMap["order_id"])
		assert.NotEmpty(t, itemMap["product_id"])
		assert.Greater(t, itemMap["quantity"], float64(0))
		assert.Greater(t, itemMap["unit_price"], 0.0)

		if i == 0 {
			assert.Equal(t, "item-001", itemMap["id"])
			assert.Equal(t, "prod-001", itemMap["product_id"])
		} else if i == 1 {
			assert.Equal(t, "item-002", itemMap["id"])
			assert.Equal(t, "prod-002", itemMap["product_id"])
		} else if i == 2 {
			assert.Equal(t, "item-003", itemMap["id"])
			assert.Equal(t, "prod-003", itemMap["product_id"])
		}
	}

	mockDataSource.AssertExpectations(t)
}





func TestUpdate_WithItems(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	now := time.Now()
	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	updatedKitchenOrder := daos.KitchenOrderDAO{
		ID:         kitchenOrderID,
		OrderID:    "order-001",
		CustomerID: nil,
		Amount:     100.50,
		Slug:       "001",
		Status: daos.OrderStatusDAO{
			ID:   "2",
			Name: "Em preparação",
		},
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-001",
				OrderID:   "order-001",
				ProductID: "prod-001",
				Quantity:  2,
				UnitPrice: 50.25,
			},
		},
		CreatedAt: now,
		UpdatedAt: &now,
	}

	mockDataSource.On("Update", mock.Anything).Return(nil)
	mockDataSource.On("FindByID", kitchenOrderID).Return(updatedKitchenOrder, nil)
	mockStatusDataSource.On("FindByID", "2").Return(daos.OrderStatusDAO{ID: "2", Name: "Em preparação"}, nil)
	mockMessageBroker.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{
		"status_id": "2",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/kitchen-orders/"+kitchenOrderID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, kitchenOrderID, response["id"])
	assert.Equal(t, "Em preparação", response["status"])

	items := response["items"].([]interface{})
	assert.Len(t, items, 1)
	item := items[0].(map[string]interface{})
	assert.Equal(t, "item-001", item["id"])
	assert.Equal(t, "prod-001", item["product_id"])
	assert.Equal(t, float64(2), item["quantity"])
	assert.Equal(t, 50.25, item["unit_price"])

	mockDataSource.AssertExpectations(t)
}


func TestFindAll_DataSourceError(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	mockDataSource.On("FindAll", mock.Anything).Return(nil, assert.AnError)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.GET("/kitchen-orders", handler.FindAll)

	req, _ := http.NewRequest("GET", "/kitchen-orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}

func TestFindByID_DataSourceError(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	mockDataSource.On("FindByID", kitchenOrderID).Return(daos.KitchenOrderDAO{}, assert.AnError)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.GET("/kitchen-orders/:id", handler.FindByID)

	req, _ := http.NewRequest("GET", "/kitchen-orders/"+kitchenOrderID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}

func TestUpdate_DataSourceError(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	router := setupTestRouter()
	mockDataSource := new(MockKitchenOrderDataSource)
	mockStatusDataSource := new(MockOrderStatusDataSource)
	mockMessageBroker := new(MockMessageBroker)

	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	mockDataSource.On("FindByID", mock.Anything).Return(daos.KitchenOrderDAO{}, assert.AnError)

	handler := &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}

	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{
		"status_id": "2",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/kitchen-orders/"+kitchenOrderID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}






