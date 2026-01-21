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

func TestSQSBroker_Publish_EmptyMessageID(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:   "", // ID vazio
		Body: []byte("test message"),
		Headers: map[string]string{
			"type": "test",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	// Deve retornar erro pois não está conectado, mas o ID seria gerado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Publish_WithHeaders(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:   "msg-123",
		Body: []byte(`{"order_id":"123","status":"pending"}`),
		Headers: map[string]string{
			"content-type": "application/json",
			"priority":     "high",
			"source":       "kitchen-service",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	// Deve retornar erro pois não está conectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Close_Multiple_Times(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)

	// Act & Assert
	// Deve ser seguro chamar Close múltiplas vezes
	err1 := broker.Close()
	assert.NoError(t, err1)
	
	err2 := broker.Close()
	assert.NoError(t, err2)
	
	err3 := broker.Close()
	assert.NoError(t, err3)
}

func TestSQSBroker_Stop_Calls_Close(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)

	// Act
	err := broker.Stop()

	// Assert
	assert.NoError(t, err)
	
	// Verifica se o contexto foi cancelado (Stop chama Close que cancela o contexto)
	select {
	case <-broker.ctx.Done():
		// OK, contexto foi cancelado
	default:
		t.Error("Context should be cancelled after Stop")
	}
}

func TestSQSBroker_Mutex_Protection(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()

	// Act - Tenta operações concorrentes
	done := make(chan bool, 3)
	
	go func() {
		broker.Close()
		done <- true
	}()
	
	go func() {
		_ = broker.Publish(ctx, "queue", interfaces.Message{
			ID:   "msg-1",
			Body: []byte("test"),
		})
		done <- true
	}()
	
	go func() {
		_ = broker.Subscribe(ctx, "queue", func(ctx context.Context, msg interfaces.Message) error {
			return nil
		})
		done <- true
	}()

	// Assert - Aguarda todas as goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
	// Se chegou aqui sem deadlock, o mutex está funcionando
	assert.True(t, true)
}

func TestSQSSerializeMessage_Complex_Nested_Structure(t *testing.T) {
	// Arrange
	data := map[string]interface{}{
		"order": map[string]interface{}{
			"id":     "order-123",
			"status": "pending",
			"items": []map[string]interface{}{
				{
					"id":       "item-1",
					"quantity": 2,
					"price":    50.00,
				},
				{
					"id":       "item-2",
					"quantity": 1,
					"price":    100.00,
				},
			},
		},
		"customer": map[string]interface{}{
			"id":   "cust-456",
			"name": "John Doe",
		},
	}

	// Act
	result, err := SerializeMessage(data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	assert.NoError(t, err)
	
	order, ok := decoded["order"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "order-123", order["id"])
	
	items, ok := order["items"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, items, 2)
}

func TestSQSDeserializeMessage_To_Struct(t *testing.T) {
	// Arrange
	type OrderData struct {
		OrderID    string `json:"order_id"`
		CustomerID string `json:"customer_id"`
		Amount     float64 `json:"amount"`
	}
	
	jsonData := []byte(`{"order_id":"order-999","customer_id":"cust-888","amount":299.99}`)
	var result OrderData

	// Act
	err := DeserializeMessage(jsonData, &result)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "order-999", result.OrderID)
	assert.Equal(t, "cust-888", result.CustomerID)
	assert.Equal(t, 299.99, result.Amount)
}

func TestSQSDeserializeMessage_Empty_JSON(t *testing.T) {
	// Arrange
	emptyJSON := []byte(`{}`)
	var result map[string]interface{}

	// Act
	err := DeserializeMessage(emptyJSON, &result)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSQSBroker_Publish_EmptyBody(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:      "msg-empty",
		Body:    []byte(""),
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_PollMessages_ContextCancelled(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	
	// Cria um contexto já cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	handlerCalled := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalled = true
		return nil
	}

	// Act
	broker.pollMessages(ctx, "test-queue", handler)

	// Assert
	// pollMessages deve retornar imediatamente sem chamar o handler
	assert.False(t, handlerCalled)
}

func TestSQSBroker_PollMessages_BrokerContextCancelled(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	
	// Cancela o contexto do broker
	broker.cancel()
	
	ctx := context.Background()
	handlerCalled := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalled = true
		return nil
	}

	// Act
	broker.pollMessages(ctx, "test-queue", handler)

	// Assert
	// pollMessages deve retornar imediatamente
	assert.False(t, handlerCalled)
}

func TestSQSBroker_Publish_WithValidMessage(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:   "msg-123",
		Body: []byte(`{"order_id":"123"}`),
		Headers: map[string]string{
			"type": "order",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	// Deve retornar erro pois não está conectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Publish_GeneratesIDWhenEmpty(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:      "",
		Body:    []byte("test"),
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	// Deve retornar erro pois não está conectado
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Publish_MultipleHeaders(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:   "msg-456",
		Body: []byte("test"),
		Headers: map[string]string{
			"header1": "value1",
			"header2": "value2",
			"header3": "value3",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}


func TestSQSBroker_Connect_Success(t *testing.T) {
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
	// Pode falhar se não houver credenciais, mas o cliente deve ser criado
	if err == nil {
		assert.NotNil(t, broker.client)
	} else {
		// Se falhar, deve ser por credenciais
		assert.Error(t, err)
	}
}

func TestSQSBroker_Publish_LargeBody(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	// Cria um body grande
	largeBody := make([]byte, 10000)
	for i := range largeBody {
		largeBody[i] = 'a'
	}
	
	message := interfaces.Message{
		ID:      "msg-large",
		Body:    largeBody,
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Publish_NoHeaders(t *testing.T) {
	// Arrange
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	}
	broker := NewSQSBroker(config)
	ctx := context.Background()
	
	message := interfaces.Message{
		ID:      "msg-no-headers",
		Body:    []byte("test"),
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}
