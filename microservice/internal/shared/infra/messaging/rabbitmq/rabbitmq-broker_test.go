package rabbitmq

import (
	"context"
	"encoding/json"
	"sync"
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
	return m.Called().Error(0)
}

func (m *MockAMQPChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return m.Called(name, kind, durable, autoDelete, internal, noWait, args).Error(0)
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, deleteWhenUnused, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	callArgs := m.Called(name, durable, deleteWhenUnused, exclusive, noWait, args)
	if callArgs.Get(0) == nil {
		return amqp.Queue{}, callArgs.Error(1)
	}
	return callArgs.Get(0).(amqp.Queue), callArgs.Error(1)
}

func (m *MockAMQPChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return m.Called(name, key, exchange, noWait, args).Error(0)
}

func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return m.Called(exchange, key, mandatory, immediate, msg).Error(0)
}

func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	callArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(<-chan amqp.Delivery), callArgs.Error(1)
}

type MockAMQPConnection struct {
	mock.Mock
}

func (m *MockAMQPConnection) Channel() (*amqp.Channel, error) {
	callArgs := m.Called()
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(*amqp.Channel), callArgs.Error(1)
}

func (m *MockAMQPConnection) Close() error {
	return m.Called().Error(0)
}

func setupMockBroker(withExchange bool) (*RabbitMQBroker, *MockAMQPChannel, chan amqp.Delivery) {
	broker := NewRabbitMQBroker(RabbitMQConfig{
		URL:      "amqp://localhost:5672",
		Exchange: map[bool]string{true: "test-exchange", false: ""}[withExchange],
	})
	mockChannel := new(MockAMQPChannel)
	deliveryChan := make(chan amqp.Delivery, 10)

	mockChannel.On("QueueDeclare", "test-queue", true, false, false, false, mock.Anything).
		Return(amqp.Queue{Name: "test-queue"}, nil)
	if withExchange {
		mockChannel.On("QueueBind", "test-queue", "test-queue", "test-exchange", false, mock.Anything).Return(nil)
	}
	mockChannel.On("Consume", "test-queue", "", false, false, false, false, mock.Anything).
		Return((<-chan amqp.Delivery)(deliveryChan), nil)

	broker.channel = mockChannel
	return broker, mockChannel, deliveryChan
}

// ============ Basic Tests ============

func TestNewRabbitMQBroker(t *testing.T) {
	config := RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test-exchange"}
	broker := NewRabbitMQBroker(config)
	assert.NotNil(t, broker)
	assert.Equal(t, config.URL, broker.config.URL)
	assert.Equal(t, config.Exchange, broker.config.Exchange)
}

func TestRabbitMQBroker_Connect_NoConnection(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://invalid:5672", Exchange: "test"})
	err := broker.Connect(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
}

func TestRabbitMQBroker_Close_NoConnection(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	err := broker.Close()
	assert.NoError(t, err)
}

func TestRabbitMQBroker_Publish_NoConnection(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	msg := interfaces.Message{ID: "test", Body: []byte("test"), Headers: map[string]string{}}
	err := broker.Publish(context.Background(), "queue", msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_Subscribe_NoConnection(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	err := broker.Subscribe(context.Background(), "queue", func(ctx context.Context, msg interfaces.Message) error { return nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to RabbitMQ")
}

func TestRabbitMQBroker_Start(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	err := broker.Start(context.Background())
	assert.NoError(t, err)
}

func TestRabbitMQBroker_Stop(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	err := broker.Stop()
	assert.NoError(t, err)
}

func TestSerializeMessage_Success(t *testing.T) {
	data := map[string]interface{}{"id": "123", "message": "test", "count": 42}
	result, err := SerializeMessage(data)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "123", decoded["id"])
}

func TestDeserializeMessage_Success(t *testing.T) {
	jsonData := []byte(`{"id":"123","message":"test","count":42}`)
	var result map[string]interface{}
	err := DeserializeMessage(jsonData, &result)
	assert.NoError(t, err)
	assert.Equal(t, "123", result["id"])
}

func TestDeserializeMessage_InvalidJSON(t *testing.T) {
	err := DeserializeMessage([]byte(`{invalid}`), &map[string]interface{}{})
	assert.Error(t, err)
}

func TestRabbitMQBroker_Publish_WithMock(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	mockChannel := new(MockAMQPChannel)
	mockChannel.On("QueueDeclare", "queue", true, false, false, false, mock.Anything).Return(amqp.Queue{Name: "queue"}, nil)
	mockChannel.On("Publish", "test", "queue", false, false, mock.Anything).Return(nil)
	broker.channel = mockChannel

	msg := interfaces.Message{ID: "id", Body: []byte("body"), Headers: map[string]string{"key": "value"}}
	err := broker.Publish(context.Background(), "queue", msg)
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_WithMock(t *testing.T) {
	broker, mockChannel, _ := setupMockBroker(true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := broker.Subscribe(ctx, "test-queue", func(ctx context.Context, msg interfaces.Message) error { return nil })
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	cancel()
	mockChannel.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_MessageProcessing(t *testing.T) {
	broker, _, deliveryChan := setupMockBroker(true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	receivedMsg := interfaces.Message{}
	err := broker.Subscribe(ctx, "test-queue", func(ctx context.Context, msg interfaces.Message) error {
		receivedMsg = msg
		return nil
	})
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	deliveryChan <- amqp.Delivery{
		MessageId: "msg-1",
		Body:      []byte("test"),
		Headers:   amqp.Table{"key": "value", "num": int32(42)},
	}
	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Equal(t, "msg-1", receivedMsg.ID)
	assert.Equal(t, "value", receivedMsg.Headers["key"])
	assert.Empty(t, receivedMsg.Headers["num"])
}

func TestRabbitMQBroker_Close_WithMock(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	mockChannel := new(MockAMQPChannel)
	mockConn := new(MockAMQPConnection)
	mockChannel.On("Close").Return(nil)
	mockConn.On("Close").Return(nil)

	broker.channel = mockChannel
	broker.conn = mockConn

	err := broker.Close()
	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
	mockConn.AssertExpectations(t)
}

func TestRabbitMQBroker_Subscribe_MultipleMessages(t *testing.T) {
	broker, _, deliveryChan := setupMockBroker(true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count := 0
	err := broker.Subscribe(ctx, "test-queue", func(ctx context.Context, msg interfaces.Message) error {
		count++
		return nil
	})
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 5; i++ {
		deliveryChan <- amqp.Delivery{MessageId: "msg", Body: []byte("test"), Headers: amqp.Table{}, DeliveryTag: uint64(i + 1)}
	}
	time.Sleep(200 * time.Millisecond)
	cancel()

	assert.Equal(t, 5, count)
}

func TestRabbitMQBroker_Subscribe_HandlerError(t *testing.T) {
	broker, _, deliveryChan := setupMockBroker(true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := broker.Subscribe(ctx, "test-queue", func(ctx context.Context, msg interfaces.Message) error { return assert.AnError })
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	deliveryChan <- amqp.Delivery{MessageId: "msg", Body: []byte("test"), Headers: amqp.Table{}, DeliveryTag: 1}
	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestRabbitMQBroker_Subscribe_WithoutExchange(t *testing.T) {
	broker, mockChannel, _ := setupMockBroker(false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := broker.Subscribe(ctx, "test-queue", func(ctx context.Context, msg interfaces.Message) error { return nil })
	assert.NoError(t, err)
	mockChannel.AssertNotCalled(t, "QueueBind")
	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestRabbitMQBroker_ConcurrentPublish(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			msg := interfaces.Message{ID: "id", Body: []byte("test"), Headers: map[string]string{}}
			err := broker.Publish(context.Background(), "queue", msg)
			assert.Error(t, err)
		}(i)
	}
	wg.Wait()
}

func TestRabbitMQBroker_ConcurrentClose(t *testing.T) {
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
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

func TestRabbitMQBroker_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	err := broker.Connect(ctx)
	assert.Error(t, err)
}

func TestRabbitMQBroker_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	broker := NewRabbitMQBroker(RabbitMQConfig{URL: "amqp://localhost:5672", Exchange: "test"})
	err := broker.Connect(ctx)
	assert.Error(t, err)
}
