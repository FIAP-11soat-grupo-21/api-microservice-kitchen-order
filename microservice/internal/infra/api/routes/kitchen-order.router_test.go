package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterKitchenOrderRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routerGroup := router.Group("/kitchen-orders")

	RegisterKitchenOrderRoutes(routerGroup)

	routes := router.Routes()
	if len(routes) < 3 {
		t.Errorf("Expected at least 3 routes to be registered, got %d", len(routes))
	}

	expectedRoutes := []string{
		"GET",
		"GET",
		"GET",
	}

	methodCount := make(map[string]int)
	for _, route := range routes {
		methodCount[route.Method]++
	}

	if methodCount["GET"] != len(expectedRoutes) {
		t.Errorf("Expected %d GET routes, got %d", len(expectedRoutes), methodCount["GET"])
	}
}