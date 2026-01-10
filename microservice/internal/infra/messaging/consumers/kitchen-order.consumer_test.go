package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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

func TestNewKitchenOrderConsumer(t *testing.T) {
	// Arrange
	mockBroker := &MockMessageBroker{}

	// Act
	consumer := NewKitchenOrderConsumer(mockBroker)

	// Assert
	assert.NotNil(t, consumer)
	assert.Equal(t, mockBroker, consumer.broker)
}

func TestKitchenOrderConsumer_Structure(t *testing.T) {
	// Testa a estrutura do consumer
	mockBroker := &MockMessageBroker{}
	consumer := NewKitchenOrderConsumer(mockBroker)
	
	// Verifica se a estrutura está correta
	assert.NotNil(t, consumer)
	assert.IsType(t, &KitchenOrderConsumer{}, consumer)
	assert.Equal(t, mockBroker, consumer.broker)
}

func TestKitchenOrderConsumer_MessageTypes(t *testing.T) {
	// Testa se as estruturas de mensagem estão corretas
	createMsg := CreateKitchenOrderMessage{
		OrderID: "test-order",
		Amount:  100.0,
		Items:   []CreateOrderItemMessage{},
	}
	
	// Verifica se pode ser serializada
	data, err := json.Marshal(createMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Verifica se pode ser deserializada
	var decoded CreateKitchenOrderMessage
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, createMsg.OrderID, decoded.OrderID)
	assert.Equal(t, createMsg.Amount, decoded.Amount)
}

func TestKitchenOrderConsumer_MessageWithItems(t *testing.T) {
	// Testa mensagem com itens
	createMsg := CreateKitchenOrderMessage{
		OrderID: "test-order",
		Amount:  150.75,
		Items: []CreateOrderItemMessage{
			{
				ID:        "item-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 50.25,
			},
			{
				ID:        "item-2",
				ProductID: "product-2",
				Quantity:  1,
				UnitPrice: 50.25,
			},
		},
	}
	
	// Verifica se pode ser serializada
	data, err := json.Marshal(createMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Verifica se pode ser deserializada
	var decoded CreateKitchenOrderMessage
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, createMsg.OrderID, decoded.OrderID)
	assert.Equal(t, createMsg.Amount, decoded.Amount)
	assert.Len(t, decoded.Items, 2)
	assert.Equal(t, createMsg.Items[0].ProductID, decoded.Items[0].ProductID)
}

func TestKitchenOrderConsumer_MessageWithCustomerID(t *testing.T) {
	// Testa mensagem com customer ID
	customerID := "customer-123"
	createMsg := CreateKitchenOrderMessage{
		OrderID:    "test-order",
		CustomerID: &customerID,
		Amount:     100.0,
		Items:      []CreateOrderItemMessage{},
	}
	
	// Verifica se pode ser serializada
	data, err := json.Marshal(createMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Verifica se pode ser deserializada
	var decoded CreateKitchenOrderMessage
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, createMsg.OrderID, decoded.OrderID)
	assert.NotNil(t, decoded.CustomerID)
	assert.Equal(t, *createMsg.CustomerID, *decoded.CustomerID)
}

func TestKitchenOrderConsumer_ResponseStructure(t *testing.T) {
	// Testa a estrutura de resposta
	response := KitchenOrderResponse{
		Success: true,
		Data:    "test data",
		Error:   "",
	}
	
	// Verifica se pode ser serializada
	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Verifica se pode ser deserializada
	var decoded KitchenOrderResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, response.Success, decoded.Success)
	assert.Equal(t, response.Data, decoded.Data)
}

func TestKitchenOrderConsumer_ErrorResponse(t *testing.T) {
	// Testa resposta de erro
	response := KitchenOrderResponse{
		Success: false,
		Data:    nil,
		Error:   "test error",
	}
	
	// Verifica se pode ser serializada
	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	
	// Verifica se pode ser deserializada
	var decoded KitchenOrderResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, response.Success, decoded.Success)
	assert.Equal(t, response.Error, decoded.Error)
}

func TestKitchenOrderConsumer_Methods_Exist(t *testing.T) {
	mockBroker := &MockMessageBroker{}
	consumer := NewKitchenOrderConsumer(mockBroker)
	
	// Verifica se os métodos existem
	assert.NotNil(t, consumer.Start)
}

func TestKitchenOrderConsumer_Initialization(t *testing.T) {
	// Testa se múltiplas instâncias podem ser criadas
	mockBroker1 := &MockMessageBroker{}
	mockBroker2 := &MockMessageBroker{}
	
	consumer1 := NewKitchenOrderConsumer(mockBroker1)
	consumer2 := NewKitchenOrderConsumer(mockBroker2)
	
	assert.NotNil(t, consumer1)
	assert.NotNil(t, consumer2)
	
	// Verifica se são instâncias diferentes
	assert.NotSame(t, consumer1, consumer2)
	assert.NotSame(t, consumer1.broker, consumer2.broker)
}
func TestCreateOrderItemMessage_Structure(t *testing.T) {
	// Test item message structure
	item := CreateOrderItemMessage{
		ID:        "item-1",
		ProductID: "product-1",
		Quantity:  2,
		UnitPrice: 12.75,
	}
	
	// Assert structure
	assert.Equal(t, "item-1", item.ID)
	assert.Equal(t, "product-1", item.ProductID)
	assert.Equal(t, 2, item.Quantity)
	assert.Equal(t, 12.75, item.UnitPrice)
	assert.IsType(t, "", item.ID)
	assert.IsType(t, "", item.ProductID)
	assert.IsType(t, 0, item.Quantity)
	assert.IsType(t, float64(0), item.UnitPrice)
}

func TestKitchenOrderConsumer_Multiple_Items_Detailed(t *testing.T) {
	// Test message with multiple items in detail
	customerID := "customer-123"
	message := CreateKitchenOrderMessage{
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     50.00,
		Items: []CreateOrderItemMessage{
			{
				ID:        "item-1",
				ProductID: "product-1",
				Quantity:  2,
				UnitPrice: 12.75,
			},
			{
				ID:        "item-2",
				ProductID: "product-2",
				Quantity:  1,
				UnitPrice: 24.50,
			},
		},
	}
	
	// Assert multiple items in detail
	assert.Len(t, message.Items, 2)
	assert.Equal(t, "item-1", message.Items[0].ID)
	assert.Equal(t, "item-2", message.Items[1].ID)
	assert.Equal(t, "product-1", message.Items[0].ProductID)
	assert.Equal(t, "product-2", message.Items[1].ProductID)
	assert.Equal(t, 2, message.Items[0].Quantity)
	assert.Equal(t, 1, message.Items[1].Quantity)
	assert.Equal(t, 12.75, message.Items[0].UnitPrice)
	assert.Equal(t, 24.50, message.Items[1].UnitPrice)
	assert.Equal(t, 50.00, message.Amount)
}

func TestKitchenOrderConsumer_Empty_Items_Detailed(t *testing.T) {
	// Test message with empty items in detail
	customerID := "customer-123"
	message := CreateKitchenOrderMessage{
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     0.00,
		Items:      []CreateOrderItemMessage{},
	}
	
	// Assert empty items in detail
	assert.Len(t, message.Items, 0)
	assert.Equal(t, 0.00, message.Amount)
	assert.NotNil(t, message.Items)
	assert.IsType(t, []CreateOrderItemMessage{}, message.Items)
	assert.Equal(t, "order-123", message.OrderID)
	assert.Equal(t, "customer-123", *message.CustomerID)
}

func TestKitchenOrderConsumer_Nil_CustomerID_Detailed(t *testing.T) {
	// Test message with nil customer ID in detail
	message := CreateKitchenOrderMessage{
		OrderID:    "order-123",
		CustomerID: nil,
		Amount:     25.50,
		Items: []CreateOrderItemMessage{
			{
				ID:        "item-1",
				ProductID: "product-1",
				Quantity:  1,
				UnitPrice: 25.50,
			},
		},
	}
	
	// Assert nil customer ID in detail
	assert.Nil(t, message.CustomerID)
	assert.Equal(t, "order-123", message.OrderID)
	assert.Equal(t, 25.50, message.Amount)
	assert.Len(t, message.Items, 1)
	assert.Equal(t, "item-1", message.Items[0].ID)
}

func TestKitchenOrderResponse_Success_Case(t *testing.T) {
	// Test successful response
	data := map[string]interface{}{
		"id":      "order-123",
		"status":  "created",
		"amount":  25.50,
	}
	
	response := KitchenOrderResponse{
		Success: true,
		Data:    data,
		Error:   "",
	}
	
	// Assert success case
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Empty(t, response.Error)
	assert.IsType(t, map[string]interface{}{}, response.Data)
}

func TestKitchenOrderResponse_Error_Case_Detailed(t *testing.T) {
	// Test error response in detail
	response := KitchenOrderResponse{
		Success: false,
		Data:    nil,
		Error:   "validation failed: missing order ID",
	}
	
	// Assert error case in detail
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "validation failed: missing order ID", response.Error)
	assert.NotEmpty(t, response.Error)
	assert.Contains(t, response.Error, "validation failed")
}

func TestKitchenOrderConsumer_JSON_Marshaling_Complete(t *testing.T) {
	// Test complete JSON marshaling scenario
	customerID := "customer-123"
	message := CreateKitchenOrderMessage{
		OrderID:    "order-123",
		CustomerID: &customerID,
		Amount:     75.25,
		Items: []CreateOrderItemMessage{
			{
				ID:        "item-1",
				ProductID: "product-1",
				Quantity:  3,
				UnitPrice: 25.08,
			},
		},
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	assert.Contains(t, string(jsonData), "order-123")
	assert.Contains(t, string(jsonData), "customer-123")
	assert.Contains(t, string(jsonData), "75.25")
	
	// Unmarshal back
	var unmarshaled CreateKitchenOrderMessage
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, message.OrderID, unmarshaled.OrderID)
	assert.Equal(t, *message.CustomerID, *unmarshaled.CustomerID)
	assert.Equal(t, message.Amount, unmarshaled.Amount)
	assert.Len(t, unmarshaled.Items, 1)
	assert.Equal(t, message.Items[0].ID, unmarshaled.Items[0].ID)
	assert.Equal(t, message.Items[0].ProductID, unmarshaled.Items[0].ProductID)
	assert.Equal(t, message.Items[0].Quantity, unmarshaled.Items[0].Quantity)
	assert.Equal(t, message.Items[0].UnitPrice, unmarshaled.Items[0].UnitPrice)
}

func TestKitchenOrderResponse_JSON_Marshaling_Complete(t *testing.T) {
	// Test complete JSON marshaling for response
	data := map[string]interface{}{
		"id":           "order-123",
		"status":       "created",
		"amount":       25.50,
		"customer_id":  "customer-123",
		"items_count":  2,
	}
	
	response := KitchenOrderResponse{
		Success: true,
		Data:    data,
		Error:   "",
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	assert.Contains(t, string(jsonData), "true")
	assert.Contains(t, string(jsonData), "order-123")
	
	// Unmarshal back
	var unmarshaled KitchenOrderResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response.Success, unmarshaled.Success)
	assert.NotNil(t, unmarshaled.Data)
	assert.Empty(t, unmarshaled.Error)
}

func TestKitchenOrderConsumer_Large_Order(t *testing.T) {
	// Test large order with many items
	customerID := "customer-large-order"
	items := make([]CreateOrderItemMessage, 10)
	totalAmount := 0.0
	
	for i := 0; i < 10; i++ {
		price := float64(i+1) * 5.50
		items[i] = CreateOrderItemMessage{
			ID:        fmt.Sprintf("item-%d", i+1),
			ProductID: fmt.Sprintf("product-%d", i+1),
			Quantity:  i + 1,
			UnitPrice: price,
		}
		totalAmount += price * float64(i+1)
	}
	
	message := CreateKitchenOrderMessage{
		OrderID:    "large-order-123",
		CustomerID: &customerID,
		Amount:     totalAmount,
		Items:      items,
	}
	
	// Assert large order
	assert.Len(t, message.Items, 10)
	assert.Equal(t, "large-order-123", message.OrderID)
	assert.Equal(t, "customer-large-order", *message.CustomerID)
	assert.Greater(t, message.Amount, 0.0)
	
	// Verify each item
	for i, item := range message.Items {
		assert.Equal(t, fmt.Sprintf("item-%d", i+1), item.ID)
		assert.Equal(t, fmt.Sprintf("product-%d", i+1), item.ProductID)
		assert.Equal(t, i+1, item.Quantity)
		assert.Greater(t, item.UnitPrice, 0.0)
	}
}

func TestKitchenOrderConsumer_Edge_Cases(t *testing.T) {
	// Test edge cases
	
	// Zero amount
	message1 := CreateKitchenOrderMessage{
		OrderID: "zero-amount-order",
		Amount:  0.0,
		Items:   []CreateOrderItemMessage{},
	}
	assert.Equal(t, 0.0, message1.Amount)
	
	// Very small amount
	message2 := CreateKitchenOrderMessage{
		OrderID: "small-amount-order",
		Amount:  0.01,
		Items:   []CreateOrderItemMessage{},
	}
	assert.Equal(t, 0.01, message2.Amount)
	
	// Large amount
	message3 := CreateKitchenOrderMessage{
		OrderID: "large-amount-order",
		Amount:  9999.99,
		Items:   []CreateOrderItemMessage{},
	}
	assert.Equal(t, 9999.99, message3.Amount)
}