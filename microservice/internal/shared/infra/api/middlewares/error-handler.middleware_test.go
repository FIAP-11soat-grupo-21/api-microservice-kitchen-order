package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"tech_challenge/internal/domain/exceptions"
)

func TestErrorHandlerMiddleware_NoErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()

	// Simula um handler que não gera erros
	handler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	// Act
	middleware(c)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestErrorHandlerMiddleware_WithKitchenOrderNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()

	// Simula um handler que gera erro de pedido não encontrado
	handler := func(c *gin.Context) {
		err := &exceptions.KitchenOrderNotFoundException{}
		c.Error(err)
	}

	// Act
	handler(c)
	middleware(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Kitchen Order not found")
}

func TestErrorHandlerMiddleware_WithInvalidKitchenOrderData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()

	// Simula um handler que gera erro de dados inválidos
	handler := func(c *gin.Context) {
		err := &exceptions.InvalidKitchenOrderDataException{}
		c.Error(err)
	}

	// Act
	handler(c)
	middleware(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Kitchen Order data")
}

func TestErrorHandlerMiddleware_WithUnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()

	// Simula um handler que gera erro desconhecido
	handler := func(c *gin.Context) {
		err := errors.New("unknown error")
		c.Error(err)
	}

	// Act
	handler(c)
	middleware(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
}

func TestErrorHandlerMiddleware_MultipleErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()

	// Simula um handler que gera múltiplos erros
	handler := func(c *gin.Context) {
		c.Error(errors.New("first error"))
		c.Error(&exceptions.KitchenOrderNotFoundException{})
	}

	// Act
	handler(c)
	middleware(c)

	// Assert
	// Deve processar apenas o último erro
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Kitchen Order not found")
}

func TestErrorHandlerMiddleware_ContextAborted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()

	// Simula um handler que gera erro
	handler := func(c *gin.Context) {
		err := &exceptions.KitchenOrderNotFoundException{}
		c.Error(err)
	}

	// Act
	handler(c)
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
}

func TestErrorHandlerMiddleware_ChainedHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	middleware := ErrorHandlerMiddleware()
	executed := false

	// Simula handlers em cadeia
	handler1 := func(c *gin.Context) {
		c.Next() // Chama o próximo handler
		executed = true
	}

	handler2 := func(c *gin.Context) {
		err := &exceptions.KitchenOrderNotFoundException{}
		c.Error(err)
	}

	// Act
	handler1(c)
	handler2(c)
	middleware(c)

	// Assert
	assert.True(t, executed)
	assert.Equal(t, http.StatusNotFound, w.Code)
}