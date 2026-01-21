package consumers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tech_challenge/internal"
	"tech_challenge/internal/shared/interfaces"
)

// MockMessageBroker é um mock do message broker
type MockMessageBroker struct {
	mock.Mock
}

func (m *MockMessageBroker) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageBroker) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMessageBroker) Publish(ctx context.Context, queue string, message interfaces.Message) error {
	args := m.Called(ctx, queue, message)
	return args.Error(0)
}

func (m *MockMessageBroker) Subscribe(ctx context.Context, queue string, handler interfaces.MessageHandler) error {
	args := m.Called(ctx, queue, handler)
	return args.Error(0)
}

func (m *MockMessageBroker) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMessageBroker) Stop() error {
	args := m.Called()
	return args.Error(0)
}

// Helper function to create a consumer with mocked controller
func createConsumerWithMocks(mockBroker *MockMessageBroker) *KitchenOrderConsumer {
	return &KitchenOrderConsumer{
		broker: mockBroker,
	}
}

// Tests for NewKitchenOrderConsumer
func TestNewKitchenOrderConsumer(t *testing.T) {
	mockBroker := &MockMessageBroker{}
	consumer := NewKitchenOrderConsumer(mockBroker)
	assert.NotNil(t, consumer)
	assert.Equal(t, mockBroker, consumer.broker)
}

// Tests for Start() method
func TestKitchenOrderConsumer_Start_Success(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	mockBroker := new(MockMessageBroker)
	mockBroker.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	consumer := NewKitchenOrderConsumer(mockBroker)
	ctx := context.Background()

	err := consumer.Start(ctx)

	assert.NoError(t, err)
	mockBroker.AssertExpectations(t)
	mockBroker.AssertCalled(t, "Subscribe", mock.Anything, mock.Anything, mock.Anything)
}

func TestKitchenOrderConsumer_Start_SubscribeError(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	mockBroker := new(MockMessageBroker)
	mockBroker.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

	consumer := NewKitchenOrderConsumer(mockBroker)
	ctx := context.Background()

	err := consumer.Start(ctx)

	assert.Error(t, err)
	mockBroker.AssertExpectations(t)
}

// Tests for handleCreate() method
func TestKitchenOrderConsumer_HandleCreate_InvalidJSON(t *testing.T) {
	internal.SetupTestEnv()
	defer internal.CleanupTestEnv()

	mockBroker := new(MockMessageBroker)
	consumer := createConsumerWithMocks(mockBroker)
	ctx := context.Background()

	msg := interfaces.Message{
		ID:      "msg-123",
		Body:    []byte("invalid json"),
		Headers: map[string]string{},
	}

	err := consumer.handleCreate(ctx, msg)

	assert.Error(t, err)
}
