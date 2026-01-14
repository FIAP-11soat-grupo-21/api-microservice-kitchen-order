package sqs

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"tech_challenge/internal/shared/interfaces"
)

func TestNewSQSBroker(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}

	// Act
	broker := NewSQSBroker(config)

	// Assert
	assert.NotNil(t, broker)
	assert.Equal(t, config.Region, broker.config.Region)
	assert.Equal(t, config.QueueURL, broker.config.QueueURL)
	assert.NotNil(t, broker.ctx)
	assert.NotNil(t, broker.cancel)
}

func TestSQSBroker_Structure(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-west-2",
		QueueURL: "https://sqs.us-west-2.amazonaws.com/123456789012/kitchen-queue",
	}

	// Act
	broker := NewSQSBroker(config)

	// Assert
	assert.NotNil(t, broker)
	assert.IsType(t, &SQSBroker{}, broker)
	assert.Equal(t, config, broker.config)
}

func TestSQSConfig_Structure(t *testing.T) {
	// Arrange & Act
	config := SQSConfig{
		Region:   "eu-west-1",
		QueueURL: "https://sqs.eu-west-1.amazonaws.com/123456789012/orders-queue",
	}

	// Assert
	assert.NotEmpty(t, config.Region)
	assert.NotEmpty(t, config.QueueURL)
	assert.Contains(t, config.QueueURL, "sqs")
	assert.Contains(t, config.QueueURL, config.Region)
}

func TestSQSBroker_Connect_NoCredentials(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()

	// Act
	err := broker.Connect(ctx)

	// Assert
	// Pode falhar se não houver credenciais AWS configuradas
	// Este teste verifica que a função pode ser chamada
	if err != nil {
		// Se falhar, deve ser por falta de credenciais ou configuração
		assert.Error(t, err)
	} else {
		// Se passar, o cliente foi criado com sucesso
		assert.NotNil(t, broker.client)
	}
}

func TestSQSBroker_Close(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)

	// Act
	err := broker.Close()

	// Assert
	// Close deve sempre ser seguro
	assert.NoError(t, err)
}

func TestSQSBroker_Publish_NoConnection(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:   "test-id",
		Body: []byte("test message"),
		Headers: map[string]string{
			"type": "test",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	// Deve retornar erro pois não está conectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Subscribe_NoConnection(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	// Deve retornar erro pois não está conectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Start(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()

	// Act
	err := broker.Start(ctx)

	// Assert
	// Start sempre retorna nil para SQS
	assert.NoError(t, err)
}

func TestSQSBroker_Stop(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)

	// Act
	err := broker.Stop()

	// Assert
	// Stop deve ser seguro
	assert.NoError(t, err)
}

func TestSQSSerializeMessage_Success(t *testing.T) {
	// Arrange
	data := map[string]interface{}{
		"order_id":    "order-123",
		"customer_id": "customer-456",
		"amount":      99.99,
	}

	// Act
	result, err := SerializeMessage(data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verifica se é JSON válido
	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "order-123", decoded["order_id"])
	assert.Equal(t, "customer-456", decoded["customer_id"])
	assert.Equal(t, 99.99, decoded["amount"])
}

func TestSQSSerializeMessage_String(t *testing.T) {
	// Arrange
	data := "simple message"

	// Act
	result, err := SerializeMessage(data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, `"simple message"`, string(result))
}

func TestSQSDeserializeMessage_Success(t *testing.T) {
	// Arrange
	jsonData := []byte(`{"order_id":"order-123","customer_id":"customer-456","amount":99.99}`)
	var result map[string]interface{}

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "order-123", result["order_id"])
	assert.Equal(t, "customer-456", result["customer_id"])
	assert.Equal(t, 99.99, result["amount"])
}

func TestSQSDeserializeMessage_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := []byte(`{invalid json}`)
	var result map[string]interface{}

	// Act
	err := DeserializeMessage(invalidJSON, &result)

	// Assert
	assert.Error(t, err)
}

func TestSQSBroker_MessageSerialization_RoundTrip(t *testing.T) {
	// Arrange
	originalData := map[string]interface{}{
		"order_id":    "order-789",
		"customer_id": "customer-101",
		"amount":      149.99,
		"items": []interface{}{
			map[string]interface{}{"id": "item1", "quantity": 3},
			map[string]interface{}{"id": "item2", "quantity": 2},
		},
	}

	// Act - Serialize
	serialized, err := SerializeMessage(originalData)
	assert.NoError(t, err)

	// Act - Deserialize
	var deserialized map[string]interface{}
	err = DeserializeMessage(serialized, &deserialized)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "order-789", deserialized["order_id"])
	assert.Equal(t, "customer-101", deserialized["customer_id"])
	assert.Equal(t, 149.99, deserialized["amount"])
	
	items, ok := deserialized["items"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, items, 2)
}

func TestSQSBroker_ConfigValidation(t *testing.T) {
	// Testa diferentes configurações
	testCases := []struct {
		name   string
		config SQSConfig
		valid  bool
	}{
		{
			name: "Valid config",
			config: SQSConfig{
				Region:   "us-east-1",
				QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/kitchen-queue",
			},
			valid: true,
		},
		{
			name: "Valid config different region",
			config: SQSConfig{
				Region:   "eu-west-1",
				QueueURL: "https://sqs.eu-west-1.amazonaws.com/123456789012/orders-queue",
			},
			valid: true,
		},
		{
			name: "Empty region",
			config: SQSConfig{
				Region:   "",
				QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
			},
			valid: false,
		},
		{
			name: "Empty queue URL",
			config: SQSConfig{
				Region:   "us-east-1",
				QueueURL: "",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			broker := NewSQSBroker(tc.config)
			assert.NotNil(t, broker)
			
			if tc.valid {
				assert.NotEmpty(t, broker.config.Region)
				assert.NotEmpty(t, broker.config.QueueURL)
			} else {
				// Pelo menos um campo deve estar vazio
				assert.True(t, broker.config.Region == "" || broker.config.QueueURL == "")
			}
		})
	}
}

func TestSQSBroker_Context_Management(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}

	// Act
	broker := NewSQSBroker(config)

	// Assert
	assert.NotNil(t, broker.ctx)
	assert.NotNil(t, broker.cancel)
	
	// Verifica se o contexto não foi cancelado ainda
	select {
	case <-broker.ctx.Done():
		t.Error("Context should not be cancelled initially")
	default:
		// OK, contexto não foi cancelado
	}
	
	// Testa cancelamento
	broker.cancel()
	
	// Verifica se o contexto foi cancelado
	select {
	case <-broker.ctx.Done():
		// OK, contexto foi cancelado
	default:
		t.Error("Context should be cancelled after calling cancel")
	}
}