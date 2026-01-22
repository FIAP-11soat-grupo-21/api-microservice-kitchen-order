package internal

import (
	"os"
	"testing"
)

func TestSetupTestEnv(t *testing.T) {
	CleanupTestEnv()

	SetupTestEnv()

	tests := []struct {
		envVar   string
		expected string
	}{
		{"GO_ENV", "test"},
		{"API_PORT", "8080"},
		{"API_HOST", "localhost"},
		{"DB_RUN_MIGRATIONS", "false"},
		{"DB_HOST", "localhost"},
		{"DB_NAME", "test"},
		{"DB_PORT", "5432"},
		{"DB_USERNAME", "test"},
		{"DB_PASSWORD", "test"},
		{"AWS_REGION", "us-east-1"},
		{"MESSAGE_BROKER_TYPE", "sqs"},
		{"AWS_SQS_KITCHEN_ORDERS_QUEUE", "https://sqs.us-east-1.amazonaws.com/123456789/test-queue"},
		{"AWS_SQS_ORDERS_QUEUE", "https://sqs.us-east-1.amazonaws.com/123456789/orders-queue"},
	}

	for _, tt := range tests {
		t.Run(tt.envVar, func(t *testing.T) {
			value, exists := os.LookupEnv(tt.envVar)
			if !exists {
				t.Errorf("variável de ambiente %s não foi definida", tt.envVar)
			}
			if value != tt.expected {
				t.Errorf("variável %s = %q, esperado %q", tt.envVar, value, tt.expected)
			}
		})
	}

	CleanupTestEnv()
}

func TestCleanupTestEnv(t *testing.T) {
	SetupTestEnv()

	_, exists := os.LookupEnv("GO_ENV")
	if !exists {
		t.Fatal("GO_ENV deveria existir após SetupTestEnv")
	}

	CleanupTestEnv()

	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"AWS_REGION", "MESSAGE_BROKER_TYPE",
		"AWS_SQS_KITCHEN_ORDERS_QUEUE", "AWS_SQS_ORDERS_QUEUE",
	}

	for _, envVar := range envVars {
		t.Run(envVar, func(t *testing.T) {
			_, exists := os.LookupEnv(envVar)
			if exists {
				t.Errorf("variável de ambiente %s deveria ter sido removida", envVar)
			}
		})
	}
}

func TestSetupAndCleanupCycle(t *testing.T) {
	CleanupTestEnv()

	SetupTestEnv()
	value1, _ := os.LookupEnv("GO_ENV")
	if value1 != "test" {
		t.Errorf("primeiro ciclo: GO_ENV = %q, esperado 'test'", value1)
	}

	CleanupTestEnv()
	_, exists1 := os.LookupEnv("GO_ENV")
	if exists1 {
		t.Error("primeiro ciclo: GO_ENV deveria ter sido removida")
	}

	SetupTestEnv()
	value2, _ := os.LookupEnv("GO_ENV")
	if value2 != "test" {
		t.Errorf("segundo ciclo: GO_ENV = %q, esperado 'test'", value2)
	}

	CleanupTestEnv()
	_, exists2 := os.LookupEnv("GO_ENV")
	if exists2 {
		t.Error("segundo ciclo: GO_ENV deveria ter sido removida")
	}
}
