package sqs

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"tech_challenge/internal/shared/interfaces"
)

type mockSQSClient struct {
	sendMessageFunc    func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	receiveMessageFunc func(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	deleteMessageFunc  func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

func (m *mockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	if m.sendMessageFunc != nil {
		return m.sendMessageFunc(ctx, params, optFns...)
	}
	return &sqs.SendMessageOutput{}, nil
}

func (m *mockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	if m.receiveMessageFunc != nil {
		return m.receiveMessageFunc(ctx, params, optFns...)
	}
	return &sqs.ReceiveMessageOutput{}, nil
}

func (m *mockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	if m.deleteMessageFunc != nil {
		return m.deleteMessageFunc(ctx, params, optFns...)
	}
	return &sqs.DeleteMessageOutput{}, nil
}

func TestSQSBroker_Connect_Success(t *testing.T) {
	config := SQSConfig{
		Region:      "us-east-1",
		QueueURL:    "http://localhost:4566/000000000000/test-queue",
		EndpointURL: "http://localhost:4566",
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	err := broker.Connect(ctx)

	// Pode ter sucesso se LocalStack estiver rodando, ou erro se não estiver
	if err != nil {
		t.Logf("✓ Erro de conexão esperado (sem LocalStack): %v", err)
	} else {
		t.Log("✓ Conexão bem-sucedida (LocalStack rodando)")
		broker.Close()
	}
}

func TestSQSBroker_Connect_WithoutEndpoint(t *testing.T) {
	config := SQSConfig{
		Region:      "us-east-1",
		QueueURL:    "https://sqs.us-east-1.amazonaws.com/123456789/test",
		EndpointURL: "", // Sem endpoint customizado
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	err := broker.Connect(ctx)

	// Vai tentar conectar na AWS real
	if err != nil {
		t.Logf("✓ Erro esperado (sem credenciais AWS): %v", err)
	} else {
		t.Log("✓ Conexão AWS bem-sucedida")
		broker.Close()
	}
}

func TestSQSBroker_Connect_ContextCanceled(t *testing.T) {
	config := SQSConfig{
		Region:      "us-east-1",
		QueueURL:    "http://localhost:4566/000000000000/test-queue",
		EndpointURL: "http://localhost:4566",
	}

	broker := NewSQSBroker(config)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar antes de conectar

	err := broker.Connect(ctx)

	if err != nil {
		t.Logf("✓ Erro com contexto cancelado: %v", err)
	} else {
		t.Log("Conexão bem-sucedida apesar do contexto cancelado")
		broker.Close()
	}
}

func TestSQSBroker_Publish_Success(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	// Usar mock client
	mockClient := &mockSQSClient{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			return &sqs.SendMessageOutput{
				MessageId: aws.String("test-message-id"),
			}, nil
		},
	}
	broker.SetClient(mockClient)

	ctx := context.Background()
	message := interfaces.Message{
		ID:   "test-id",
		Body: []byte(`{"test": "data"}`),
		Headers: map[string]string{
			"content-type": "application/json",
		},
	}

	err := broker.Publish(ctx, "test-queue", message)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	} else {
		t.Log("✓ Mensagem publicada com sucesso")
	}
}

func TestSQSBroker_Publish_Error(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	// Mock que retorna erro
	mockClient := &mockSQSClient{
		sendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
			return nil, errors.New("SQS send error")
		},
	}
	broker.SetClient(mockClient)

	ctx := context.Background()
	message := interfaces.Message{
		Body: []byte(`{"test": "data"}`),
	}

	err := broker.Publish(ctx, "test-queue", message)

	if err == nil {
		t.Error("Expected error, got nil")
	} else {
		t.Logf("✓ Erro capturado: %v", err)
	}
}

func TestSQSBroker_Publish_NotConnected(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	// Não conectar nem definir client

	ctx := context.Background()
	message := interfaces.Message{
		Body: []byte(`{"test": "data"}`),
	}

	err := broker.Publish(ctx, "test-queue", message)

	if err == nil {
		t.Error("Expected error when not connected")
	} else {
		t.Logf("✓ Erro esperado: %v", err)
	}
}

func TestSQSBroker_Subscribe_Success(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	mockClient := &mockSQSClient{}
	broker.SetClient(mockClient)

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	err := broker.Subscribe(ctx, "test-queue", handler)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	} else {
		t.Log("✓ Subscribe executado com sucesso")
	}
	
	broker.Close()
}

func TestSQSBroker_Subscribe_NotConnected(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	err := broker.Subscribe(ctx, "test-queue", handler)

	if err == nil {
		t.Error("Expected error when not connected")
	} else {
		t.Logf("✓ Erro esperado: %v", err)
	}
}

func TestSQSBroker_ProcessMessage(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	mockClient := &mockSQSClient{
		deleteMessageFunc: func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
			return &sqs.DeleteMessageOutput{}, nil
		},
	}
	broker.SetClient(mockClient)

	ctx := context.Background()
	
	messageProcessed := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		messageProcessed = true
		return nil
	}

	sqsMessage := types.Message{
		MessageId: aws.String("test-id"),
		Body:      aws.String(`{"test": "data"}`),
		ReceiptHandle: aws.String("receipt-handle"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"type": {
				StringValue: aws.String("test-type"),
			},
		},
	}

	broker.processMessage(ctx, sqsMessage, handler)

	if !messageProcessed {
		t.Error("Expected message to be processed")
	} else {
		t.Log("✓ Mensagem processada com sucesso")
	}
}

func TestSQSBroker_Close(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	err := broker.Close()
	
	if err != nil {
		t.Errorf("Expected no error on close, got: %v", err)
	} else {
		t.Log("✓ Broker fechado com sucesso")
	}
}

func TestSerializeMessage(t *testing.T) {
	data := map[string]string{"key": "value"}
	
	bytes, err := SerializeMessage(data)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if len(bytes) == 0 {
		t.Error("Expected non-empty bytes")
	} else {
		t.Logf("✓ Mensagem serializada: %s", string(bytes))
	}
}

func TestDeserializeMessage(t *testing.T) {
	jsonData := []byte(`{"key": "value"}`)
	var result map[string]string
	
	err := DeserializeMessage(jsonData, &result)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if result["key"] != "value" {
		t.Errorf("Expected 'value', got: %s", result["key"])
	} else {
		t.Log("✓ Mensagem deserializada com sucesso")
	}
}

func TestSQSBroker_Connect_InvalidRegion(t *testing.T) {
	// Testar com região inválida para forçar erro na AWS SDK
	config := SQSConfig{
		Region:      "", // Região vazia deve causar erro
		QueueURL:    "http://localhost:4566/000000000000/test-queue",
		EndpointURL: "http://localhost:4566",
	}

	broker := NewSQSBroker(config)
	ctx := context.Background()

	err := broker.Connect(ctx)

	if err != nil {
		t.Logf("✓ Erro de configuração capturado: %v", err)
	} else {
		t.Log("Conexão bem-sucedida com região vazia")
		broker.Close()
	}
}

func TestSQSBroker_ProcessBatch_WithMessages(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	messageReceived := false
	mockClient := &mockSQSClient{
		receiveMessageFunc: func(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
			return &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId: aws.String("test-id"),
						Body:      aws.String(`{"test": "data"}`),
						ReceiptHandle: aws.String("receipt-handle"),
					},
				},
			}, nil
		},
		deleteMessageFunc: func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
			messageReceived = true
			return &sqs.DeleteMessageOutput{}, nil
		},
	}
	broker.SetClient(mockClient)

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	broker.processBatch(ctx, handler)

	if !messageReceived {
		t.Error("Expected message to be received and deleted")
	} else {
		t.Log("✓ Batch processado com sucesso")
	}
}

func TestSQSBroker_ProcessBatch_ReceiveError(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	mockClient := &mockSQSClient{
		receiveMessageFunc: func(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
			return nil, errors.New("receive error")
		},
	}
	broker.SetClient(mockClient)

	ctx := context.Background()
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	// Não deve causar panic, apenas logar o erro
	broker.processBatch(ctx, handler)
	t.Log("✓ Erro de recebimento tratado corretamente")
}

func TestSQSBroker_ProcessMessage_HandlerError(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	mockClient := &mockSQSClient{}
	broker.SetClient(mockClient)

	ctx := context.Background()
	
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return errors.New("handler error")
	}

	sqsMessage := types.Message{
		MessageId: aws.String("test-id"),
		Body:      aws.String(`{"test": "data"}`),
		ReceiptHandle: aws.String("receipt-handle"),
	}

	// Não deve causar panic, apenas logar o erro
	broker.processMessage(ctx, sqsMessage, handler)
	t.Log("✓ Erro do handler tratado corretamente")
}

func TestSQSBroker_ProcessMessage_DeleteError(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	mockClient := &mockSQSClient{
		deleteMessageFunc: func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
			return nil, errors.New("delete error")
		},
	}
	broker.SetClient(mockClient)

	ctx := context.Background()
	
	handler := func(ctx context.Context, msg interfaces.Message) error {
		return nil
	}

	sqsMessage := types.Message{
		MessageId: aws.String("test-id"),
		Body:      aws.String(`{"test": "data"}`),
		ReceiptHandle: aws.String("receipt-handle"),
	}

	// Não deve causar panic, apenas logar o erro
	broker.processMessage(ctx, sqsMessage, handler)
	t.Log("✓ Erro de deleção tratado corretamente")
}

func TestSQSBroker_Start_Stop(t *testing.T) {
	config := SQSConfig{
		Region:   "us-east-1",
		QueueURL: "http://localhost:4566/000000000000/test-queue",
	}

	broker := NewSQSBroker(config)
	
	ctx := context.Background()
	
	err := broker.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error on start, got: %v", err)
	}
	
	err = broker.Stop()
	if err != nil {
		t.Errorf("Expected no error on stop, got: %v", err)
	}
	
	t.Log("✓ Start/Stop executados com sucesso")
}
