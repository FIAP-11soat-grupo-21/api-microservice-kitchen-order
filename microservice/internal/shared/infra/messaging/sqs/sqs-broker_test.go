package sqs

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tech_challenge/internal/shared/interfaces"
)

func TestNewSQSBroker(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	assert.NotNil(t, broker)
	assert.Equal(t, config.Region, broker.config.Region)
	assert.Equal(t, config.QueueURL, broker.config.QueueURL)
	assert.NotNil(t, broker.ctx)
	assert.NotNil(t, broker.cancel)
}

func TestSQSBroker_Close(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	err := broker.Close()

	assert.NoError(t, err)
}

func TestSQSBroker_PublishWithoutConnection(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	message := interfaces.Message{
		ID:      "test-id",
		Body:    []byte("test body"),
		Headers: map[string]string{"key": "value"},
	}

	err := broker.Publish(ctx, "test-queue", message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_SubscribeWithoutConnection(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	err := broker.Subscribe(ctx, "test-queue", handler)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to SQS")
}

func TestSQSBroker_Start(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	err := broker.Start(ctx)

	assert.NoError(t, err)
}

func TestSQSBroker_Stop(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	err := broker.Stop()

	assert.NoError(t, err)
}

func TestSQSBroker_Config(t *testing.T) {
	config := SQSConfig{
		Region:   "us-west-2",
		QueueURL: "https://sqs.us-west-2.amazonaws.com/987654321/another-queue",
	}

	broker := NewSQSBroker(config)

	assert.Equal(t, "us-west-2", broker.config.Region)
	assert.Equal(t, "https://sqs.us-west-2.amazonaws.com/987654321/another-queue", broker.config.QueueURL)
}

func TestSQSBroker_Message_Structure(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	data := TestData{Name: "John", Age: 30}
	serialized, err := SerializeMessage(data)

	assert.NoError(t, err)
	assert.NotNil(t, serialized)
	assert.Contains(t, string(serialized), "John")
	assert.Contains(t, string(serialized), "30")
}

func TestDeserializeMessage(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonData := []byte(`{"name":"Jane","age":25}`)
	var result TestData

	err := DeserializeMessage(jsonData, &result)

	assert.NoError(t, err)
	assert.Equal(t, "Jane", result.Name)
	assert.Equal(t, 25, result.Age)
}

func TestSerializeAndDeserializeMessage(t *testing.T) {
	type TestData struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	}

	original := TestData{ID: "123", Value: "test"}
	serialized, err := SerializeMessage(original)
	require.NoError(t, err)

	var deserialized TestData
	err = DeserializeMessage(serialized, &deserialized)
	require.NoError(t, err)

	assert.Equal(t, original.ID, deserialized.ID)
	assert.Equal(t, original.Value, deserialized.Value)
}

func TestDeserializeMessage_InvalidJSON(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
	}

	invalidJSON := []byte(`{invalid json}`)
	var result TestData

	err := DeserializeMessage(invalidJSON, &result)

	assert.Error(t, err)
}

func TestSQSBroker_ContextCancellation(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	select {
	case <-broker.ctx.Done():
		t.Fatal("Context should not be cancelled initially")
	default:
	}

	broker.cancel()

	select {
	case <-broker.ctx.Done():
	default:
		t.Fatal("Context should be cancelled after calling cancel()")
	}
}

func TestSQSBroker_PollMessages_ContextCancelled(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	ctx, cancel := context.WithCancel(context.Background())

	handlerCalled := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalled = true
		return nil
	}

	cancel()

	broker.PollMessages(ctx, "test-queue", handler)

	assert.False(t, handlerCalled)
}

func TestSQSBroker_PollMessages_BrokerContextCancelled(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	handlerCalled := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalled = true
		return nil
	}

	broker.cancel()

	broker.PollMessages(ctx, "test-queue", handler)

	assert.False(t, handlerCalled)
}

func TestSQSBroker_Message_EmptyHeaders(t *testing.T) {
	msg := interfaces.Message{
		ID:      "test-456",
		Body:    []byte("body"),
		Headers: map[string]string{},
	}

	assert.Empty(t, msg.Headers)
	assert.Equal(t, 0, len(msg.Headers))
}

func TestSQSBroker_Concurrent_Operations(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	done := make(chan bool, 3)

	go func() {
		_ = broker.Close()
		done <- true
	}()

	go func() {
		_ = broker.Stop()
		done <- true
	}()

	go func() {
		ctx := context.Background()
		_ = broker.Start(ctx)
		done <- true
	}()

	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestSQSBroker_Message_With_Multiple_Headers(t *testing.T) {
	msg := interfaces.Message{
		ID:   "msg-001",
		Body: []byte("test"),
		Headers: map[string]string{
			"header1": "value1",
			"header2": "value2",
			"header3": "value3",
		},
	}

	assert.Equal(t, 3, len(msg.Headers))
	assert.Equal(t, "value1", msg.Headers["header1"])
	assert.Equal(t, "value2", msg.Headers["header2"])
	assert.Equal(t, "value3", msg.Headers["header3"])
}

func TestSQSBroker_Large_Message_Body(t *testing.T) {
	largeBody := make([]byte, 10000)
	for i := range largeBody {
		largeBody[i] = byte(i % 256)
	}

	msg := interfaces.Message{
		ID:      "large-msg",
		Body:    largeBody,
		Headers: map[string]string{},
	}

	assert.Equal(t, 10000, len(msg.Body))
	assert.Equal(t, "large-msg", msg.ID)
}

func TestSQSBroker_Config_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config SQSConfig
	}{
		{
			name: "valid config",
			config: SQSConfig{
				Region:   "us-east-1",
				QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/queue",
			},
		},
		{
			name: "different region",
			config: SQSConfig{
				Region:   "eu-west-1",
				QueueURL: "https://sqs.eu-west-1.amazonaws.com/987654321/queue",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broker := NewSQSBroker(tt.config)
			assert.Equal(t, tt.config.Region, broker.config.Region)
			assert.Equal(t, tt.config.QueueURL, broker.config.QueueURL)
		})
	}
}

func TestSQSBroker_Handler_Signature(t *testing.T) {
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	msg := interfaces.Message{
		ID:      "test",
		Body:    []byte("test"),
		Headers: map[string]string{},
	}

	err := handler(context.Background(), msg)

	assert.NoError(t, err)
}

func TestSQSBroker_Handler_With_Error(t *testing.T) {
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return fmt.Errorf("test error")
	}

	msg := interfaces.Message{
		ID:      "test",
		Body:    []byte("test"),
		Headers: map[string]string{},
	}

	err := handler(context.Background(), msg)

	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}


// MockSQSClient completo para testar o fluxo completo
type MockSQSClient struct {
	receiveMessageCalls int
	deleteMessageCalls  int
	sendMessageCalls    int
	messages            []types.Message
	shouldFailReceive   bool
	shouldFailDelete    bool
	shouldFailSend      bool
	receivedDeleteCalls []string
}

func (m *MockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	m.receiveMessageCalls++

	if m.shouldFailReceive {
		return nil, fmt.Errorf("mock receive error")
	}

	return &sqs.ReceiveMessageOutput{
		Messages: m.messages,
	}, nil
}

func (m *MockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	m.deleteMessageCalls++

	if m.shouldFailDelete {
		return nil, fmt.Errorf("mock delete error")
	}

	if params.ReceiptHandle != nil {
		m.receivedDeleteCalls = append(m.receivedDeleteCalls, *params.ReceiptHandle)
	}

	return &sqs.DeleteMessageOutput{}, nil
}

func (m *MockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.sendMessageCalls++

	if m.shouldFailSend {
		return nil, fmt.Errorf("mock send error")
	}

	return &sqs.SendMessageOutput{
		MessageId: aws.String("mock-message-id"),
	}, nil
}

func TestSQSBroker_SetClient(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)
	mockClient := &MockSQSClient{}

	broker.SetClient(mockClient)

	assert.Equal(t, mockClient, broker.client)
}

func TestSQSBroker_ProcessBatch_Success(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	messages := []types.Message{
		{
			MessageId:     aws.String("msg-1"),
			Body:          aws.String("body-1"),
			ReceiptHandle: aws.String("receipt-1"),
		},
		{
			MessageId:     aws.String("msg-2"),
			Body:          aws.String("body-2"),
			ReceiptHandle: aws.String("receipt-2"),
		},
	}

	mockClient := &MockSQSClient{
		messages: messages,
	}

	broker.SetClient(mockClient)
	ctx := context.Background()

	handlerCalls := 0
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalls++
		return nil
	}

	broker.processBatch(ctx, handler)

	assert.Equal(t, 1, mockClient.receiveMessageCalls)
	assert.Equal(t, 2, handlerCalls)
	assert.Equal(t, 2, mockClient.deleteMessageCalls)
}

func TestSQSBroker_ProcessBatch_ReceiveError(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	mockClient := &MockSQSClient{
		shouldFailReceive: true,
	}

	broker.SetClient(mockClient)
	ctx := context.Background()

	handlerCalls := 0
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalls++
		return nil
	}

	broker.processBatch(ctx, handler)

	assert.Equal(t, 1, mockClient.receiveMessageCalls)
	assert.Equal(t, 0, handlerCalls)
	assert.Equal(t, 0, mockClient.deleteMessageCalls)
}

func TestSQSBroker_ProcessMessage_Success(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	mockClient := &MockSQSClient{}
	broker.SetClient(mockClient)

	ctx := context.Background()

	messageID := "msg-123"
	messageBody := "test body"
	receiptHandle := "receipt-123"

	msg := types.Message{
		MessageId:     &messageID,
		Body:          &messageBody,
		ReceiptHandle: &receiptHandle,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"key1": {StringValue: aws.String("value1")},
		},
	}

	handlerCalls := 0
	var receivedMessage interfaces.Message

	handler := func(ctx context.Context, m interfaces.Message) error {
		handlerCalls++
		receivedMessage = m
		return nil
	}

	broker.processMessage(ctx, msg, handler)

	assert.Equal(t, 1, handlerCalls)
	assert.Equal(t, "msg-123", receivedMessage.ID)
	assert.Equal(t, []byte("test body"), receivedMessage.Body)
	assert.Equal(t, "value1", receivedMessage.Headers["key1"])
	assert.Equal(t, 1, mockClient.deleteMessageCalls)
	assert.Equal(t, "receipt-123", mockClient.receivedDeleteCalls[0])
}

func TestSQSBroker_ProcessMessage_HandlerError(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	mockClient := &MockSQSClient{}
	broker.SetClient(mockClient)

	ctx := context.Background()

	messageID := "msg-error"
	messageBody := "error body"
	receiptHandle := "receipt-error"

	msg := types.Message{
		MessageId:     &messageID,
		Body:          &messageBody,
		ReceiptHandle: &receiptHandle,
	}

	handler := func(ctx context.Context, m interfaces.Message) error {
		return fmt.Errorf("handler error")
	}

	broker.processMessage(ctx, msg, handler)

	assert.Equal(t, 0, mockClient.deleteMessageCalls)
}

func TestSQSBroker_ProcessMessage_DeleteError(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	mockClient := &MockSQSClient{
		shouldFailDelete: true,
	}
	broker.SetClient(mockClient)

	ctx := context.Background()

	messageID := "msg-delete-error"
	messageBody := "delete error body"
	receiptHandle := "receipt-delete-error"

	msg := types.Message{
		MessageId:     &messageID,
		Body:          &messageBody,
		ReceiptHandle: &receiptHandle,
	}

	handlerCalls := 0
	handler := func(ctx context.Context, m interfaces.Message) error {
		handlerCalls++
		return nil
	}

	broker.processMessage(ctx, msg, handler)

	assert.Equal(t, 1, handlerCalls)
	assert.Equal(t, 1, mockClient.deleteMessageCalls)
}

func TestSQSBroker_ProcessMessage_With_Multiple_Attributes(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	mockClient := &MockSQSClient{}
	broker.SetClient(mockClient)

	ctx := context.Background()

	messageID := "msg-attrs"
	messageBody := "attrs body"
	receiptHandle := "receipt-attrs"

	msg := types.Message{
		MessageId:     &messageID,
		Body:          &messageBody,
		ReceiptHandle: &receiptHandle,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"attr1": {StringValue: aws.String("value1")},
			"attr2": {StringValue: aws.String("value2")},
			"attr3": {StringValue: aws.String("value3")},
			"attr4": {StringValue: nil},
		},
	}

	var receivedMessage interfaces.Message

	handler := func(ctx context.Context, m interfaces.Message) error {
		receivedMessage = m
		return nil
	}

	broker.processMessage(ctx, msg, handler)

	assert.Equal(t, 3, len(receivedMessage.Headers))
	assert.Equal(t, "value1", receivedMessage.Headers["attr1"])
	assert.Equal(t, "value2", receivedMessage.Headers["attr2"])
	assert.Equal(t, "value3", receivedMessage.Headers["attr3"])
	assert.NotContains(t, receivedMessage.Headers, "attr4")
}

func TestSQSBroker_ProcessBatch_Empty_Messages(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	mockClient := &MockSQSClient{
		messages: []types.Message{},
	}

	broker.SetClient(mockClient)
	ctx := context.Background()

	handlerCalls := 0
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalls++
		return nil
	}

	broker.processBatch(ctx, handler)

	assert.Equal(t, 1, mockClient.receiveMessageCalls)
	assert.Equal(t, 0, handlerCalls)
	assert.Equal(t, 0, mockClient.deleteMessageCalls)
}

func TestSQSBroker_ProcessBatch_Handler_Error_Continues(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
	}

	broker := NewSQSBroker(config)

	messages := []types.Message{
		{
			MessageId:     aws.String("msg-1"),
			Body:          aws.String("body-1"),
			ReceiptHandle: aws.String("receipt-1"),
		},
		{
			MessageId:     aws.String("msg-2"),
			Body:          aws.String("body-2"),
			ReceiptHandle: aws.String("receipt-2"),
		},
		{
			MessageId:     aws.String("msg-3"),
			Body:          aws.String("body-3"),
			ReceiptHandle: aws.String("receipt-3"),
		},
	}

	mockClient := &MockSQSClient{
		messages: messages,
	}

	broker.SetClient(mockClient)
	ctx := context.Background()

	handlerCalls := 0
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalls++
		if handlerCalls == 2 {
			return fmt.Errorf("error on second message")
		}
		return nil
	}

	broker.processBatch(ctx, handler)

	assert.Equal(t, 3, handlerCalls)
	assert.Equal(t, 2, mockClient.deleteMessageCalls)
}
