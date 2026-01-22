package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tech_challenge/internal"
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

// Test helpers
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func createMocks() (*MockKitchenOrderDataSource, *MockOrderStatusDataSource, *MockMessageBroker) {
	return new(MockKitchenOrderDataSource), new(MockOrderStatusDataSource), new(MockMessageBroker)
}

func createHandler(mockDataSource *MockKitchenOrderDataSource, mockStatusDataSource *MockOrderStatusDataSource, mockMessageBroker *MockMessageBroker) *KitchenOrderHandler {
	return &KitchenOrderHandler{
		kitchenOrderController: *controllers.NewKitchenOrderController(
			mockDataSource,
			mockStatusDataSource,
			mockMessageBroker,
		),
	}
}

func createTestItem(id, orderID, productID string, quantity int, unitPrice float64) daos.OrderItemDAO {
	return daos.OrderItemDAO{
		ID:        id,
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
}

func createTestKitchenOrder(id, orderID string, items []daos.OrderItemDAO) daos.KitchenOrderDAO {
	now := time.Now()
	return daos.KitchenOrderDAO{
		ID:         id,
		OrderID:    orderID,
		CustomerID: nil,
		Amount:     100.50,
		Slug:       "001",
		Status: daos.OrderStatusDAO{
			ID:   "1",
			Name: "Recebido",
		},
		Items:     items,
		CreatedAt: now,
		UpdatedAt: &now,
	}
}

// Tests
func TestNewKitchenOrderHandler(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

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
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	items := []daos.OrderItemDAO{createTestItem("item-001", "order-001", "prod-001", 2, 50.25)}
	kitchenOrder := createTestKitchenOrder("550e8400-e29b-41d4-a716-446655440000", "order-001", items)
	mockDataSource.On("FindAll", mock.Anything).Return([]daos.KitchenOrderDAO{kitchenOrder}, nil)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
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
	assert.Equal(t, "001", response[0]["slug"])
	assert.Equal(t, "Recebido", response[0]["status"])
	mockDataSource.AssertExpectations(t)
}

func TestFindAll_WithFilters(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	mockDataSource.On("FindAll", mock.MatchedBy(func(filter dtos.KitchenOrderFilter) bool {
		return filter.CreatedAtFrom != nil && filter.CreatedAtTo != nil
	})).Return([]daos.KitchenOrderDAO{}, nil)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.GET("/kitchen-orders", handler.FindAll)

	req, _ := http.NewRequest("GET", "/kitchen-orders?created_at_from=2024-01-01T00:00:00Z&created_at_to=2024-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}

func TestFindByID_Success(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	items := []daos.OrderItemDAO{createTestItem("item-001", "order-001", "prod-001", 2, 50.25)}
	kitchenOrder := createTestKitchenOrder("550e8400-e29b-41d4-a716-446655440000", "order-001", items)
	mockDataSource.On("FindByID", kitchenOrder.ID).Return(kitchenOrder, nil)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.GET("/kitchen-orders/:id", handler.FindByID)

	req, _ := http.NewRequest("GET", "/kitchen-orders/"+kitchenOrder.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, kitchenOrder.ID, response["id"])
	assert.Equal(t, "order-001", response["order_id"])
	assert.Equal(t, "001", response["slug"])
	assert.Equal(t, "Recebido", response["status"])
	mockDataSource.AssertExpectations(t)
}

func TestFindByID_WithMultipleItems(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	items := []daos.OrderItemDAO{
		createTestItem("item-001", "order-001", "prod-001", 2, 50.25),
		createTestItem("item-002", "order-001", "prod-002", 1, 50.25),
		createTestItem("item-003", "order-001", "prod-003", 3, 16.75),
	}
	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	kitchenOrder := createTestKitchenOrder(kitchenOrderID, "order-001", items)
	kitchenOrder.Amount = 150.75

	mockDataSource.On("FindByID", kitchenOrderID).Return(kitchenOrder, nil)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
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
	assert.Equal(t, "001", response["slug"])
	assert.Equal(t, "Recebido", response["status"])
	mockDataSource.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

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
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	mockDataSource.On("Update", mock.Anything).Return(nil)
	mockDataSource.On("FindByID", kitchenOrderID).Return(updatedKitchenOrder, nil)
	mockStatusDataSource.On("FindByID", "2").Return(daos.OrderStatusDAO{ID: "2", Name: "Em preparação"}, nil)
	mockMessageBroker.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{"status_id": "2"}
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

func TestUpdate_WithItems(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	items := []daos.OrderItemDAO{createTestItem("item-001", "order-001", "prod-001", 2, 50.25)}
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
		Items:     items,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	mockDataSource.On("Update", mock.Anything).Return(nil)
	mockDataSource.On("FindByID", kitchenOrderID).Return(updatedKitchenOrder, nil)
	mockStatusDataSource.On("FindByID", "2").Return(daos.OrderStatusDAO{ID: "2", Name: "Em preparação"}, nil)
	mockMessageBroker.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{"status_id": "2"}
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
	assert.Equal(t, "001", response["slug"])
	mockDataSource.AssertExpectations(t)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.PUT("/kitchen-orders/:id", handler.Update)

	req, _ := http.NewRequest("PUT", "/kitchen-orders/ko-001", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_MissingStatusID(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/kitchen-orders/ko-001", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindAll_DataSourceError(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	mockDataSource.On("FindAll", mock.Anything).Return(nil, assert.AnError)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.GET("/kitchen-orders", handler.FindAll)

	req, _ := http.NewRequest("GET", "/kitchen-orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}

func TestFindByID_DataSourceError(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	mockDataSource.On("FindByID", kitchenOrderID).Return(daos.KitchenOrderDAO{}, assert.AnError)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.GET("/kitchen-orders/:id", handler.FindByID)

	req, _ := http.NewRequest("GET", "/kitchen-orders/"+kitchenOrderID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}

func TestUpdate_DataSourceError(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := setupTestRouter()
	mockDataSource, mockStatusDataSource, mockMessageBroker := createMocks()

	kitchenOrderID := "550e8400-e29b-41d4-a716-446655440000"
	mockDataSource.On("FindByID", mock.Anything).Return(daos.KitchenOrderDAO{}, assert.AnError)

	handler := createHandler(mockDataSource, mockStatusDataSource, mockMessageBroker)
	router.PUT("/kitchen-orders/:id", handler.Update)

	requestBody := map[string]interface{}{"status_id": "2"}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/kitchen-orders/"+kitchenOrderID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockDataSource.AssertExpectations(t)
}
