package rabbitmq

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"tech_challenge/internal/shared/interfaces"
)

func TestNewRabbitMQBroker(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}

	// Act
	broker := NewRabbitMQBroker(config)

	// Assert
	assert.NotNil(t, broker)
	assert.Equal(t, config.URL, broker.config.URL)
	assert.Equal(t, config.Exchange, broker.config.Exchange)
}

func TestRabbitMQBroker_Structure(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://test:5672",
		Exchange: "test",
	}

	// Act
	broker := NewRabbitMQBroker(config)

	// Assert
	assert.NotNil(t, broker)
	assert.IsType(t, &RabbitMQBroker{}, broker)
	assert.Equal(t, config, broker.config)
}

func TestRabbitMQConfig_Structure(t *testing.T) {
	// Arrange & Act
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "kitchen-exchange",
	}

	// Assert
	assert.NotEmpty(t, config.URL)
	assert.NotEmpty(t, config.Exchange)
	assert.Contains(t, config.URL, "amqp://")
}

func TestRabbitMQBroker_Connect_NoConnection(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://invalid:5672", // URL inválida
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Act
	err := broker.Connect(ctx)

	// Assert
	// Deve retornar erro pois não há RabbitMQ rodando
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
}

func TestRabbitMQBroker_Close_NoConnection(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	// Act
	err := broker.Close()

	// Assert
	// Deve ser seguro fechar mesmo sem conexão
	assert.NoError(t, err)
}

func TestRabbitMQBroker_Publish_NoConnection(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
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
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_Subscribe_NoConnection(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()
	
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	// Deve retornar erro pois não está conectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_Start(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Act
	err := broker.Start(ctx)

	// Assert
	// Start sempre retorna nil para RabbitMQ
	assert.NoError(t, err)
}

func TestRabbitMQBroker_Stop(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	// Act
	err := broker.Stop()

	// Assert
	// Stop deve ser seguro mesmo sem conexão
	assert.NoError(t, err)
}

func TestSerializeMessage_Success(t *testing.T) {
	// Arrange
	data := map[string]interface{}{
		"id":      "123",
		"message": "test message",
		"count":   42,
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
	assert.Equal(t, "123", decoded["id"])
	assert.Equal(t, "test message", decoded["message"])
	assert.Equal(t, float64(42), decoded["count"]) // JSON unmarshals numbers as float64
}

func TestSerializeMessage_String(t *testing.T) {
	// Arrange
	data := "simple string message"

	// Act
	result, err := SerializeMessage(data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, `"simple string message"`, string(result))
}

func TestSerializeMessage_Nil(t *testing.T) {
	// Arrange
	var data interface{} = nil

	// Act
	result, err := SerializeMessage(data)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), result)
}

func TestDeserializeMessage_Success(t *testing.T) {
	// Arrange
	jsonData := []byte(`{"id":"123","message":"test message","count":42}`)
	var result map[string]interface{}

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "123", result["id"])
	assert.Equal(t, "test message", result["message"])
	assert.Equal(t, float64(42), result["count"])
}

func TestDeserializeMessage_String(t *testing.T) {
	// Arrange
	jsonData := []byte(`"simple string message"`)
	var result string

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "simple string message", result)
}

func TestDeserializeMessage_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := []byte(`{invalid json}`)
	var result map[string]interface{}

	// Act
	err := DeserializeMessage(invalidJSON, &result)

	// Assert
	assert.Error(t, err)
}

func TestDeserializeMessage_Struct(t *testing.T) {
	// Arrange
	type TestStruct struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Count   int    `json:"count"`
	}
	
	jsonData := []byte(`{"id":"123","message":"test message","count":42}`)
	var result TestStruct

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "123", result.ID)
	assert.Equal(t, "test message", result.Message)
	assert.Equal(t, 42, result.Count)
}

func TestRabbitMQBroker_MessageSerialization_RoundTrip(t *testing.T) {
	// Arrange
	originalData := map[string]interface{}{
		"order_id":    "order-123",
		"customer_id": "customer-456",
		"amount":      99.99,
		"items": []interface{}{
			map[string]interface{}{"id": "item1", "quantity": 2},
			map[string]interface{}{"id": "item2", "quantity": 1},
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
	assert.Equal(t, "order-123", deserialized["order_id"])
	assert.Equal(t, "customer-456", deserialized["customer_id"])
	assert.Equal(t, 99.99, deserialized["amount"])
	
	items, ok := deserialized["items"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, items, 2)
}

func TestRabbitMQBroker_ConfigValidation(t *testing.T) {
	// Testa diferentes configurações
	testCases := []struct {
		name   string
		config RabbitMQConfig
		valid  bool
	}{
		{
			name: "Valid config with exchange",
			config: RabbitMQConfig{
				URL:      "amqp://localhost:5672",
				Exchange: "kitchen-exchange",
			},
			valid: true,
		},
		{
			name: "Valid config without exchange",
			config: RabbitMQConfig{
				URL:      "amqp://localhost:5672",
				Exchange: "",
			},
			valid: true,
		},
		{
			name: "Empty URL",
			config: RabbitMQConfig{
				URL:      "",
				Exchange: "test-exchange",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			broker := NewRabbitMQBroker(tc.config)
			assert.NotNil(t, broker)
			
			if tc.valid {
				assert.NotEmpty(t, broker.config.URL)
			} else {
				assert.Empty(t, broker.config.URL)
			}
		})
	}
}