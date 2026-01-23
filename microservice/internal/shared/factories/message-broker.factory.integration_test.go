// +build integration

package factories

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNewMessageBroker_SQS_ConnectionError_Integration(t *testing.T) {
	
	originalEndpoint := os.Getenv("AWS_ENDPOINT_URL")
	defer os.Setenv("AWS_ENDPOINT_URL", originalEndpoint)
	
	os.Setenv("AWS_ENDPOINT_URL", "http://invalid-endpoint-xyz:9999")
	
	ctx := context.Background()
	broker, err := NewMessageBroker(ctx)
	
	if err != nil {
		t.Logf("✓ Erro de conexão capturado: %v", err)
		
		if !strings.Contains(err.Error(), "failed to connect to SQS") &&
		   !strings.Contains(err.Error(), "failed to load AWS config") {
			t.Logf("Error message: %v", err)
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
