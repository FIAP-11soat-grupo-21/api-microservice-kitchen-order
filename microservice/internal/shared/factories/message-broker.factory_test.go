package factories

import (
	"context"
	"os"
	"testing"
)

func TestNewMessageBroker_UnsupportedType(t *testing.T) {
	setTestEnvVars()
	defer cleanupTestEnvVars()
	
	os.Setenv("MESSAGE_BROKER_TYPE", "unsupported")
	
	ctx := context.Background()

	broker, err := NewMessageBroker(ctx)

	if err == nil {
		t.Error("Expected error for unsupported message broker type, got nil")
	}

	if broker != nil {
		t.Error("Expected broker to be nil for unsupported type")
	}

	expectedError := "unsupported message broker type: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestMessageBrokerType_Constants(t *testing.T) {
	if MessageBrokerRabbitMQ != "rabbitmq" {
		t.Errorf("Expected MessageBrokerRabbitMQ to be 'rabbitmq', got %s", MessageBrokerRabbitMQ)
	}

	if MessageBrokerSQS != "sqs" {
		t.Errorf("Expected MessageBrokerSQS to be 'sqs', got %s", MessageBrokerSQS)
	}
}

func setTestEnvVars() {
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
}

func cleanupTestEnvVars() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"AWS_REGION", "MESSAGE_BROKER_TYPE", "RABBITMQ_URL",
		"RABBITMQ_KITCHEN_QUEUE", "SQS_QUEUE_URL",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}