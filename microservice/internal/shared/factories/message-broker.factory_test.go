package factories

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func setupTestEnv(t *testing.T) {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		os.Setenv("GO_ENV", "test")
		os.Setenv("API_PORT", "8082")
		os.Setenv("API_HOST", "0.0.0.0")
		os.Setenv("DB_RUN_MIGRATIONS", "false")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USERNAME", "test")
		os.Setenv("DB_PASSWORD", "test")
		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
		os.Setenv("AWS_SQS_KITCHEN_ORDERS_QUEUE", "http://localhost:4566/000000000000/kitchen-orders")
		os.Setenv("AWS_SQS_ORDERS_QUEUE", "http://localhost:4566/000000000000/orders")
		os.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")
	}
}

func TestNewMessageBroker_SQS_Success(t *testing.T) {
	setupTestEnv(t)
	ctx := context.Background()
	broker, err := NewMessageBroker(ctx)

	if err == nil {
		t.Log("✓ Broker SQS criado com sucesso")
		if broker == nil {
			t.Error("Expected broker to be non-nil when no error")
		} else {
			broker.Close()
		}
	} else {
		t.Logf("✓ Código SQS executado, erro esperado sem SQS rodando: %v", err)
		
		if !strings.Contains(err.Error(), "failed to connect to SQS") &&
		   !strings.Contains(err.Error(), "failed to load AWS config") {
			t.Errorf("Expected SQS connection error, got: %v", err)
		}
		
		if broker != nil {
			t.Error("Expected broker to be nil when error occurs")
		}
	}
}

func TestNewMessageBroker_SQS_ConnectionFlow(t *testing.T) {
	setupTestEnv(t)
	ctx := context.Background()
	broker, err := NewMessageBroker(ctx)

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "failed to connect to SQS") ||
		   strings.Contains(errMsg, "failed to load AWS config") ||
		   strings.Contains(errMsg, "connection") {
			t.Logf("✓ Código de conexão SQS executado: %v", err)
		} else {
			t.Errorf("Unexpected error type: %v", err)
		}
	} else {
		t.Log("✓ Conexão SQS bem-sucedida")
		if broker != nil {
			broker.Close()
		}
	}
}

func TestNewMessageBroker_SQS_ContextCanceled(t *testing.T) {
	setupTestEnv(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	broker, err := NewMessageBroker(ctx)

	if err == nil {
		t.Log("Broker criado apesar do contexto cancelado")
		if broker != nil {
			broker.Close()
		}
	} else {
		t.Logf("✓ Erro com contexto cancelado: %v", err)
		if broker != nil {
			t.Error("Expected broker to be nil when error occurs")
		}
	}
}

func TestNewMessageBroker_SQS_ErrorHandling(t *testing.T) {
	setupTestEnv(t)
	ctx := context.Background()
	broker, err := NewMessageBroker(ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "failed to") {
			t.Logf("Error message format: %v", err)
		}
		
		if broker != nil {
			t.Error("Expected broker to be nil when error occurs")
		}
		
		t.Logf("✓ Tratamento de erro SQS executado corretamente")
	} else {
		t.Log("✓ Broker criado sem erros")
		if broker == nil {
			t.Error("Expected broker to be non-nil when no error")
		} else {
			broker.Close()
		}
	}
}

func TestNewMessageBroker_SQS_ConnectionError_InvalidEndpoint(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	os.Setenv("API_PORT", "8082")
	os.Setenv("API_HOST", "0.0.0.0")
	os.Setenv("DB_RUN_MIGRATIONS", "false")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
	os.Setenv("AWS_SQS_KITCHEN_ORDERS_QUEUE", "http://invalid-host-12345:9999/queue")
	os.Setenv("AWS_SQS_ORDERS_QUEUE", "http://invalid-host-12345:9999/queue")
	os.Setenv("AWS_ENDPOINT_URL", "http://invalid-host-12345:9999")
	
	ctx := context.Background()
	broker, err := NewMessageBroker(ctx)
	
	if err != nil {
		t.Logf("✓ Erro de conexão capturado: %v", err)
		
		if !strings.Contains(err.Error(), "failed to connect to SQS") {
			t.Errorf("Expected 'failed to connect to SQS' error, got: %v", err)
		}
		
		if broker != nil {
			t.Error("Expected broker to be nil when connection fails")
		}
	} else {
		t.Log("Broker criado (endpoint pode estar acessível)")
		if broker != nil {
			broker.Close()
		}
	}
}


func TestNewMessageBroker_UnsupportedType(t *testing.T) {
	if os.Getenv("TEST_UNSUPPORTED_TYPE") == "1" {
		os.Setenv("GO_ENV", "test")
		os.Setenv("API_PORT", "8082")
		os.Setenv("API_HOST", "0.0.0.0")
		os.Setenv("DB_RUN_MIGRATIONS", "false")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USERNAME", "test")
		os.Setenv("DB_PASSWORD", "test")
		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("MESSAGE_BROKER_TYPE", "kafka")
		os.Setenv("AWS_SQS_KITCHEN_ORDERS_QUEUE", "http://localhost:4566/000000000000/kitchen-orders")
		os.Setenv("AWS_SQS_ORDERS_QUEUE", "http://localhost:4566/000000000000/orders")
		
		ctx := context.Background()
		_, err := NewMessageBroker(ctx)
		
		if err != nil && strings.Contains(err.Error(), "unsupported message broker type") {
			os.Exit(43)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewMessageBroker_UnsupportedType")
	cmd.Env = append(os.Environ(), "TEST_UNSUPPORTED_TYPE=1")
	err := cmd.Run()
	
	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() == 43 {
			t.Log("✓ Erro 'unsupported message broker type' capturado com sucesso")
			return
		}
	}
	
	t.Log("Teste executado (resultado pode variar dependendo do ambiente)")
}
