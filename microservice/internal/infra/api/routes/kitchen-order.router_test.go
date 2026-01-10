package routes

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

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

func TestRegisterKitchenOrderRoutes(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
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