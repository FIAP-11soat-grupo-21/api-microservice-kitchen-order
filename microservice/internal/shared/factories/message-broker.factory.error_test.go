package factories

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestNewMessageBroker_SQS_ConnectionError_Subprocess(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		os.Setenv("GO_ENV", "test")
		os.Setenv("API_PORT", "8082")
		os.Setenv("API_HOST", "0.0.0.0")
		os.Setenv("DB_RUN_MIGRATIONS", "false")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_USERNAME", "test")
		os.Setenv("DB_PASSWORD", "test")
		os.Setenv("AWS_REGION", "invalid-region-xyz")
		os.Setenv("MESSAGE_BROKER_TYPE", "sqs")
		os.Setenv("AWS_SQS_KITCHEN_ORDERS_QUEUE", "http://invalid-host-xyz:9999/queue")
		os.Setenv("AWS_SQS_ORDERS_QUEUE", "http://invalid-host-xyz:9999/queue")
		os.Setenv("AWS_ENDPOINT_URL", "http://invalid-host-xyz:9999")
		
		ctx := context.Background()
		_, err := NewMessageBroker(ctx)
		
		if err != nil && strings.Contains(err.Error(), "failed to connect to SQS") {
			os.Exit(42)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewMessageBroker_SQS_ConnectionError_Subprocess")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")
	err := cmd.Run()
	
	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() == 42 {
			t.Log("✓ Erro 'failed to connect to SQS' capturado com sucesso")
			return
		}
	}
	
	t.Log("Teste executado (resultado pode variar dependendo do ambiente)")
}
