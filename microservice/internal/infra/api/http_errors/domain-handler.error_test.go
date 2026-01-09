package http_errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"tech_challenge/internal/domain/exceptions"
)

func TestHandleDomainErrors_InvalidKitchenOrderDataException(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.InvalidKitchenOrderDataException{
		Message: "Invalid data",
	}

	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("Expected error to be handled, got false")
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleDomainErrors_KitchenOrderNotFoundException(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := &exceptions.KitchenOrderNotFoundException{}

	handled := HandleDomainErrors(err, ctx)

	if !handled {
		t.Error("Expected error to be handled, got false")
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleDomainErrors_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := errors.New("unknown error")

	handled := HandleDomainErrors(err, ctx)

	if handled {
		t.Error("Expected error not to be handled, got true")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for unhandled error, got %d", http.StatusOK, w.Code)
	}
}