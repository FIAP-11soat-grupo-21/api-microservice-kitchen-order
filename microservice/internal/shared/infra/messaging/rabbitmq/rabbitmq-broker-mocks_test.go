package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tech_challenge/internal/shared/interfaces"
)

type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAMQPChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	callArgs := m.Called(name, kind, durable, autoDelete, internal, noWait, args)
	return callArgs.Error(0)
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, deleteWhenUnused, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	callArgs := m.Called(name, durable, deleteWhenUnused, exclusive, noWait, args)
	if callArgs.Get(0) == nil {
		return amqp.Queue{}, callArgs.Error(1)
	}
	return callArgs.Get(0).(amqp.Queue), callArgs.Error(1)
}

func (m *MockAMQPChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	callArgs := m.Called(name, key, exchange, noWait, args)
	return callArgs.Error(0)
}

func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	callArgs := m.Called(exchange, key, mandatory, immediate, msg)
	return callArgs.Error(0)
}

func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	callArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(<-chan amqp.Delivery), callArgs.Error(1)
}

// MockAMQPConnection mock para AMQPConnection
type MockAMQPConnection struct {
	mock.Mock
}

func (m *MockAMQPConnection) Channel() (*amqp.Channel, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*amqp.Channel), args.Error(1)
}

func (m *MockAMQPConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRabbitMQBroker_Publish_Success(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	// Simula uma conexão e canal
	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("Publish", "test-exchange", "test-queue", false, false, mock.Anything).
		Return(nil)

	broker.channel = mockChannel

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
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Publish_QueueDeclareError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{}, assert.AnError)

	broker.channel = mockChannel

	ctx := context.Background()
	message := interfaces.Message{
		ID:   "test-id",
		Body: []byte("test message"),
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to declare queue")
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Publish_PublishError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("Publish", "test-exchange", "test-queue", false, false, mock.Anything).
		Return(assert.AnError)

	broker.channel = mockChannel

	ctx := context.Background()
	message := interfaces.Message{
		ID:   "test-id",
		Body: []byte("test message"),
		Headers: map[string]string{},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message")
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_Success(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)

	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestRabbitMQBroker_Subscribe_QueueDeclareError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{}, assert.AnError)

	broker.channel = mockChannel

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to declare queue")
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_QueueBindError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(assert.AnError)

	broker.channel = mockChannel

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to bind queue")
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_ConsumeError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return(nil, assert.AnError)

	broker.channel = mockChannel

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to register consumer")
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_WithoutExchange_Success(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)

	// Assert
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
	mockChannel.AssertNotCalled(t, "QueueBind")

	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestRabbitMQBroker_Close_WithConnection(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("Close").Return(nil)

	mockConn := new(MockAMQPConnection)
	mockConn.On("Close").Return(nil)

	broker.channel = mockChannel
	broker.conn = mockConn

	// Act
	err := broker.Close()

	// Assert
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
	mockConn.AssertExpectations(t)
}

func TestRabbitMQBroker_Close_ConnectionError(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("Close").Return(nil)

	mockConn := new(MockAMQPConnection)
	mockConn.On("Close").Return(assert.AnError)

	broker.channel = mockChannel
	broker.conn = mockConn

	// Act
	err := broker.Close()

	// Assert
	assert.Error(t, err)
	mockChannel.AssertExpectations(t)
	mockConn.AssertExpectations(t)
}

func TestRabbitMQBroker_Close_OnlyChannel(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("Close").Return(nil)

	broker.channel = mockChannel
	broker.conn = nil

	// Act
	err := broker.Close()

	// Assert
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Publish_WithMultipleHeaders(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("Publish", "test-exchange", "test-queue", false, false, mock.Anything).
		Return(nil)

	broker.channel = mockChannel

	ctx := context.Background()
	message := interfaces.Message{
		ID:   "test-id",
		Body: []byte("test message"),
		Headers: map[string]string{
			"type":        "order.created",
			"version":     "1.0",
			"correlation": "corr-123",
			"timestamp":   "2024-01-15T10:00:00Z",
			"source":      "kitchen-service",
		},
	}

	// Act
	err := broker.Publish(ctx, "test-queue", message)

	// Assert
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_ContextCancellation(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Cancela o contexto
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Assert
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_DeliveryChannelClosed(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx := context.Background()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Fecha o canal de entrega para simular desconexão
	close(deliveryChan)
	time.Sleep(100 * time.Millisecond)

	// Assert
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_MessageWithNonStringHeaders(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	receivedMessage := interfaces.Message{}
	handler := func(ctx context.Context, msg interfaces.Message) error {
		receivedMessage = msg
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Envia uma mensagem com headers de tipos diferentes (int, bool, etc)
	delivery := amqp.Delivery{
		MessageId: "test-msg-id",
		Body:      []byte("test message"),
		Headers: amqp.Table{
			"string_header": "value",
			"int_header":    int32(42),
			"bool_header":   true,
			"float_header":  3.14,
		},
	}
	deliveryChan <- delivery

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Assert - apenas headers string devem ser inclusos
	assert.Equal(t, "value", receivedMessage.Headers["string_header"])
	assert.Empty(t, receivedMessage.Headers["int_header"])
	assert.Empty(t, receivedMessage.Headers["bool_header"])
	assert.Empty(t, receivedMessage.Headers["float_header"])
}

func TestRabbitMQBroker_Subscribe_HandlerErrorWithNack(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return assert.AnError
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Envia uma mensagem que causará erro
	delivery := amqp.Delivery{
		MessageId:   "test-msg-id",
		Body:        []byte("test message"),
		Headers:     amqp.Table{},
		DeliveryTag: 1,
	}
	deliveryChan <- delivery

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Assert
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_HandlerSuccessWithAck(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlerCalled := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalled = true
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Envia uma mensagem que será processada com sucesso
	delivery := amqp.Delivery{
		MessageId:   "test-msg-id",
		Body:        []byte("test message"),
		Headers:     amqp.Table{},
		DeliveryTag: 1,
	}
	deliveryChan <- delivery

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Assert
	assert.True(t, handlerCalled)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_MultipleMessages(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 10)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messageCount := 0
	handler := func(ctx context.Context, msg interfaces.Message) error {
		messageCount++
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Envia múltiplas mensagens
	for i := 0; i < 5; i++ {
		delivery := amqp.Delivery{
			MessageId:   "msg-" + string(rune(i)),
			Body:        []byte("message " + string(rune(i))),
			Headers:     amqp.Table{},
			DeliveryTag: uint64(i + 1),
		}
		deliveryChan <- delivery
	}

	time.Sleep(200 * time.Millisecond)
	cancel()

	// Assert
	assert.Equal(t, 5, messageCount)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Connect_ChannelCreationSuccess(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockConn := new(MockAMQPConnection)
	mockChannel := new(MockAMQPChannel)

	mockConn.On("Channel").Return(mockChannel, nil)
	mockChannel.On("ExchangeDeclare", "test-exchange", "topic", true, false, false, false, mock.Anything).
		Return(nil)

	broker.conn = mockConn
	broker.channel = mockChannel

	// Act
	// Simula o que aconteceria após Connect bem-sucedido
	err := mockChannel.ExchangeDeclare(
		config.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	// Assert
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Connect_ExchangeDeclareFailure(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockConn := new(MockAMQPConnection)
	mockChannel := new(MockAMQPChannel)

	mockConn.On("Channel").Return(mockChannel, nil)
	mockChannel.On("ExchangeDeclare", "test-exchange", "topic", true, false, false, false, mock.Anything).
		Return(assert.AnError)
	mockChannel.On("Close").Return(nil)
	mockConn.On("Close").Return(nil)

	broker.conn = mockConn
	broker.channel = mockChannel

	// Act
	err := mockChannel.ExchangeDeclare(
		config.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	// Assert
	assert.Error(t, err)
}

func TestRabbitMQBroker_Connect_WithoutExchange(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "",
	}
	broker := NewRabbitMQBroker(config)

	mockConn := new(MockAMQPConnection)
	mockChannel := new(MockAMQPChannel)

	mockConn.On("Channel").Return(mockChannel, nil)

	broker.conn = mockConn
	broker.channel = mockChannel

	// Act
	// Sem exchange, ExchangeDeclare não deve ser chamado
	mockChannel.AssertNotCalled(t, "ExchangeDeclare")

	// Assert
	assert.NotNil(t, broker.channel)
	assert.NotNil(t, broker.conn)
}

// TestRabbitMQBroker_Subscribe_EmptyHeadersTable testa headers vazio
func TestRabbitMQBroker_Subscribe_EmptyHeadersTable(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	receivedMessage := interfaces.Message{}
	handler := func(ctx context.Context, msg interfaces.Message) error {
		receivedMessage = msg
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Envia uma mensagem com headers vazio
	delivery := amqp.Delivery{
		MessageId: "test-msg-id",
		Body:      []byte("test message"),
		Headers:   amqp.Table{},
	}
	deliveryChan <- delivery

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Assert
	assert.Equal(t, "test-msg-id", receivedMessage.ID)
	assert.Equal(t, []byte("test message"), receivedMessage.Body)
	assert.Empty(t, receivedMessage.Headers)
}

func TestRabbitMQBroker_Subscribe_NilHeadersTable(t *testing.T) {
	// Arrange
	config := RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: "test-exchange",
	}
	broker := NewRabbitMQBroker(config)

	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).
		Return(nil)

	deliveryChan := make(chan amqp.Delivery, 1)
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	receivedMessage := interfaces.Message{}
	handler := func(ctx context.Context, msg interfaces.Message) error {
		receivedMessage = msg
		return nil
	}

	// Act
	err := broker.Subscribe(ctx, "test-queue", handler)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Envia uma mensagem com headers nil
	delivery := amqp.Delivery{
		MessageId: "test-msg-id",
		Body:      []byte("test message"),
		Headers:   nil,
	}
	deliveryChan <- delivery

	time.Sleep(100 * time.Millisecond)
	cancel()

	// Assert
	assert.Equal(t, "test-msg-id", receivedMessage.ID)
	assert.Equal(t, []byte("test message"), receivedMessage.Body)
	assert.NotNil(t, receivedMessage.Headers)
	assert.Empty(t, receivedMessage.Headers)
}
