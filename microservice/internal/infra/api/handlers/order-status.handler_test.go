package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tech_challenge/internal"
	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/daos"
)

type MockOrderStatusDataSourceForHandler struct {
	mock.Mock
}

func (m *MockOrderStatusDataSourceForHandler) FindByID(id string) (daos.OrderStatusDAO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return daos.OrderStatusDAO{}, args.Error(1)
	}
	return args.Get(0).(daos.OrderStatusDAO), args.Error(1)
}

func (m *MockOrderStatusDataSourceForHandler) FindAll() ([]daos.OrderStatusDAO, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]daos.OrderStatusDAO), args.Error(1)
}

func createOrderStatusHandler(mockDataSource *MockOrderStatusDataSourceForHandler) *OrderStatusHandler {
	return &OrderStatusHandler{
		controller: *controllers.NewOrderStatusController(mockDataSource),
	}
}

func TestNewOrderStatusHandler(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	t.Run("should create handler successfully", func(t *testing.T) {
		handler := NewOrderStatusHandler()
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.controller)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		handler1 := NewOrderStatusHandler()
		handler2 := NewOrderStatusHandler()
		assert.NotNil(t, handler1)
		assert.NotNil(t, handler2)
		assert.NotSame(t, handler1, handler2)
	})
}

func TestOrderStatusHandler_FindAll_Empty(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := gin.Default()
	gin.SetMode(gin.TestMode)

	mockDataSource := new(MockOrderStatusDataSourceForHandler)
	mockDataSource.On("FindAll").Return([]daos.OrderStatusDAO{}, nil)

	handler := createOrderStatusHandler(mockDataSource)
	router.GET("/v1/kitchen-orders/status", handler.FindAll)

	req, _ := http.NewRequest("GET", "/v1/kitchen-orders/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)
	mockDataSource.AssertExpectations(t)
}

func TestOrderStatusHandler_FindAll_Response_Format(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := gin.Default()
	gin.SetMode(gin.TestMode)

	mockDataSource := new(MockOrderStatusDataSourceForHandler)
	statusData := []daos.OrderStatusDAO{
		{ID: "1", Name: "Recebido"},
	}
	mockDataSource.On("FindAll").Return(statusData, nil)

	handler := createOrderStatusHandler(mockDataSource)
	router.GET("/v1/kitchen-orders/status", handler.FindAll)

	req, _ := http.NewRequest("GET", "/v1/kitchen-orders/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	body := w.Body.String()
	assert.True(t, body[0] == '[', "Response should be a JSON array")
	assert.True(t, body[len(body)-1] == ']', "Response should end with ]")
	mockDataSource.AssertExpectations(t)
}

func TestOrderStatusHandler_FindAll_Consistent_Responses(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := gin.Default()
	gin.SetMode(gin.TestMode)

	mockDataSource := new(MockOrderStatusDataSourceForHandler)
	statusData := []daos.OrderStatusDAO{
		{ID: "1", Name: "Status 1"},
		{ID: "2", Name: "Status 2"},
	}
	mockDataSource.On("FindAll").Return(statusData, nil)

	handler := createOrderStatusHandler(mockDataSource)
	router.GET("/v1/kitchen-orders/status", handler.FindAll)

	// Chama múltiplas vezes
	responses := make([]string, 3)
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/v1/kitchen-orders/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		responses[i] = w.Body.String()
	}

	// Todas as respostas devem ser iguais
	assert.Equal(t, responses[0], responses[1])
	assert.Equal(t, responses[1], responses[2])
	mockDataSource.AssertExpectations(t)
}

func TestOrderStatusHandler_FindAll_No_Error_Field_On_Success(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := gin.Default()
	gin.SetMode(gin.TestMode)

	mockDataSource := new(MockOrderStatusDataSourceForHandler)
	statusData := []daos.OrderStatusDAO{
		{ID: "1", Name: "Status 1"},
	}
	mockDataSource.On("FindAll").Return(statusData, nil)

	handler := createOrderStatusHandler(mockDataSource)
	router.GET("/v1/kitchen-orders/status", handler.FindAll)

	req, _ := http.NewRequest("GET", "/v1/kitchen-orders/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	// Não deve conter campo "error" em sucesso
	assert.NotContains(t, body, "\"error\"")
	mockDataSource.AssertExpectations(t)
}

func TestOrderStatusHandler_FindAll_Error_Has_Error_Field(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	router := gin.Default()
	gin.SetMode(gin.TestMode)

	mockDataSource := new(MockOrderStatusDataSourceForHandler)
	mockDataSource.On("FindAll").Return(nil, assert.AnError)

	handler := createOrderStatusHandler(mockDataSource)
	router.GET("/v1/kitchen-orders/status", handler.FindAll)

	req, _ := http.NewRequest("GET", "/v1/kitchen-orders/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.NotEmpty(t, response["error"])
	mockDataSource.AssertExpectations(t)
}
