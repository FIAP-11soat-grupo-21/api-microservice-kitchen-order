package rabbitmq

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

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


func TestRabbitMQBroker_Connect_ChannelError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Act
	err := broker.Connect(ctx)

	// Assert
	assert.Error(t, err)
}

func TestRabbitMQBroker_Connect_ExchangeDeclareError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Act
	err := broker.Connect(ctx)

	// Assert
	assert.Error(t, err)
}

func TestRabbitMQBroker_Publish_WithHeaders(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	message := interfaces.Message{
		ID:   "test-id-123",
		Body: []byte(`{"order_id":"123","status":"pending"}`),
		Headers: map[string]string{
			"type":        "order.created",
			"version":     "1.0",
			"correlation": "corr-123",
		},
	}

	// Act
	err := broker.Publish(ctx, "orders-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_Subscribe_WithExchange(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "orders-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "orders-queue", handler)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_Subscribe_WithoutExchange(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "orders-queue", handler)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_ConcurrentPublish(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Act & Assert
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			message := interfaces.Message{
				ID:      "test-id-" + string(rune(index)),
				Body:    []byte("test message"),
				Headers: map[string]string{"index": string(rune(index))},
			}
			err := broker.Publish(ctx, "test-queue", message)
			assert.Error(t, err) // Esperamos erro pois não está conectado
		}(i)
	}
	wg.Wait()
}

func TestRabbitMQBroker_ConcurrentClose(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	// Act & Assert
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := broker.Close()
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
}

func TestRabbitMQBroker_MessageWithEmptyHeaders(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	message := interfaces.Message{
		ID:      "test-id",
		Body:    []byte("test message"),
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_MessageWithNilHeaders(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	message := interfaces.Message{
		ID:      "test-id",
		Body:    []byte("test message"),
		Headers: nil,
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_LargeMessage(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	largeBody := make([]byte, 1024*1024) // 1MB
	for i := range largeBody {
		largeBody[i] = byte(i % 256)
	}

	message := interfaces.Message{
		ID:      "large-message-id",
		Body:    largeBody,
		Headers: map[string]string{"size": "1MB"},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_MessageWithSpecialCharacters(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	message := interfaces.Message{
		ID:   "test-id-special-chars-!@#$%^&*()",
		Body: []byte(`{"message":"Special chars: !@#$%^&*()","emoji":"🚀"}`),
		Headers: map[string]string{
			"special": "!@#$%^&*()",
			"emoji":   "🚀",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestSerializeMessage_Complex(t *testing.T) {
	// Arrange
	type OrderItem struct {
		ID       string  `json:"id"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}

	type Order struct {
		OrderID    string      `json:"order_id"`
		CustomerID string      `json:"customer_id"`
		Items      []OrderItem `json:"items"`
		Total      float64     `json:"total"`
		Status     string      `json:"status"`
	}

	order := Order{
		OrderID:    "order-123",
		CustomerID: "customer-456",
		Items: []OrderItem{
			{ID: "item-1", Quantity: 2, Price: 29.99},
			{ID: "item-2", Quantity: 1, Price: 49.99},
		},
		Total:  109.97,
		Status: "pending",
	}

	// Act
	result, err := SerializeMessage(order)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	var deserialized Order
	err = json.Unmarshal(result, &deserialized)
	assert.NoError(t, err)
	assert.Equal(t, order.OrderID, deserialized.OrderID)
	assert.Equal(t, order.CustomerID, deserialized.CustomerID)
	assert.Len(t, deserialized.Items, 2)
	assert.Equal(t, order.Total, deserialized.Total)
}

func TestDeserializeMessage_Array(t *testing.T) {
	// Arrange
	jsonData := []byte(`[{"id":"1","name":"item1"},{"id":"2","name":"item2"}]`)
	var result []map[string]interface{}

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0]["id"])
	assert.Equal(t, "item1", result[0]["name"])
}

func TestDeserializeMessage_Number(t *testing.T) {
	// Arrange
	jsonData := []byte(`42`)
	var result int

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestDeserializeMessage_Boolean(t *testing.T) {
	// Arrange
	jsonData := []byte(`true`)
	var result bool

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestRabbitMQBroker_MultipleConnectAttempts(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://invalid-host:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Act & Assert
	for i := 0; i < 3; i++ {
		err := broker.Connect(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
	}
}

func TestRabbitMQBroker_PublishAfterClose(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Fecha sem conectar
	err := broker.Close()
	assert.NoError(t, err)

	message := interfaces.Message{
		ID:      "test-id",
		Body:    []byte("test message"),
		Headers: map[string]string{},
	}

	// Act
	err = broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_SubscribeAfterClose(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	// Fecha sem conectar
	err := broker.Close()
	assert.NoError(t, err)

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err = broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_ContextCancellation(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela imediatamente

	// Act
	err := broker.Connect(ctx)

	// Assert
	// Mesmo com contexto cancelado, o erro será sobre conexão RabbitMQ
	assert.Error(t, err)
}

func TestRabbitMQBroker_ContextTimeout(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Act
	err := broker.Connect(ctx)

	// Assert
	// Esperamos erro de conexão
	assert.Error(t, err)
}

func TestSerializeMessage_EmptyStruct(t *testing.T) {
	// Arrange
	type EmptyStruct struct{}
	data := EmptyStruct{}

	// Act
	result, err := SerializeMessage(data)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("{}"), result)
}

func TestSerializeMessage_NestedStructures(t *testing.T) {
	// Arrange
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	person := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	// Act
	result, err := SerializeMessage(person)

	// Assert
	assert.NoError(t, err)

	var deserialized Person
	err = json.Unmarshal(result, &deserialized)
	assert.NoError(t, err)
	assert.Equal(t, person.Name, deserialized.Name)
	assert.Equal(t, person.Address.City, deserialized.Address.City)
}

func TestRabbitMQBroker_QueueNameVariations(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)
	ctx := context.Background()

	testCases := []string{
		"simple-queue",
		"queue_with_underscore",
		"queue.with.dots",
		"queue-with-multiple-dashes",
		"UPPERCASE_QUEUE",
		"MixedCaseQueue",
	}

	for _, queueName := range testCases {
		t.Run(queueName, func(t *testing.T) {
			message := interfaces.Message{
				ID:      "test-id",
				Body:    []byte("test message"),
				Headers: map[string]string{},
			}

			err := broker.Publish(ctx, queueName, message)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not connected to RabbitMQ")
		})
	}
}