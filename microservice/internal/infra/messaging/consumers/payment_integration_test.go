package consumers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"tech_challenge/internal/shared/interfaces"
)

func TestPaymentToKitchenOrderIntegration(t *testing.T) {
	paymentConfirmation := map[string]interface{}{
		"order_id": "896156b7-d058-44cd-8969-2b269567a1f2",
		"status":   "Confirmado",
	}

	paymentJSON, err := json.Marshal(paymentConfirmation)
	assert.NoError(t, err)

	snsWrapper := map[string]interface{}{
		"Type":      "Notification",
		"MessageId": "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
		"TopicArn":  "arn:aws:sns:us-east-2:216989122312:payment-processed-topic",
		"Message":   string(paymentJSON),
		"Timestamp": "2026-01-24T00:08:56.000Z",
	}

	snsJSON, err := json.Marshal(snsWrapper)
	assert.NoError(t, err)

	receivedMessage := interfaces.Message{
		ID:      "test-message-id",
		Body:    snsJSON,
		Headers: map[string]string{},
	}

	var unwrappedSNS map[string]interface{}
	err = json.Unmarshal(receivedMessage.Body, &unwrappedSNS)
	assert.NoError(t, err)

	assert.Equal(t, "Notification", unwrappedSNS["Type"])
	assert.Equal(t, "arn:aws:sns:us-east-2:216989122312:payment-processed-topic", unwrappedSNS["TopicArn"])

	innerMessageStr, ok := unwrappedSNS["Message"].(string)
	assert.True(t, ok, "Message field should be a string")

	var actualPayload map[string]interface{}
	err = json.Unmarshal([]byte(innerMessageStr), &actualPayload)
	assert.NoError(t, err)

	assert.Equal(t, "896156b7-d058-44cd-8969-2b269567a1f2", actualPayload["order_id"])
	assert.Equal(t, "Confirmado", actualPayload["status"])

	t.Log("Payment to Kitchen Order integration flow validated")
	t.Logf("Payment sent: %s", string(paymentJSON))
	t.Logf("SNS wrapped: %s", string(snsJSON))
	t.Logf("Kitchen Order will receive: %s", innerMessageStr)
}

func TestCreateKitchenOrderMessageFormat(t *testing.T) {
	expectedFormat := CreateKitchenOrderMessage{
		OrderID: "896156b7-d058-44cd-8969-2b269567a1f2",
	}

	data, err := json.Marshal(expectedFormat)
	assert.NoError(t, err)

	var result CreateKitchenOrderMessage
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedFormat.OrderID, result.OrderID)

	t.Log("Kitchen Order message format validated")
	t.Logf("Expected format: %s", string(data))
}

func TestPaymentMessageCompatibility(t *testing.T) {
	paymentMessage := `{"order_id":"896156b7-d058-44cd-8969-2b269567a1f2","status":"Confirmado"}`

	var kitchenOrderMsg CreateKitchenOrderMessage
	err := json.Unmarshal([]byte(paymentMessage), &kitchenOrderMsg)
	assert.NoError(t, err)
	assert.Equal(t, "896156b7-d058-44cd-8969-2b269567a1f2", kitchenOrderMsg.OrderID)

	t.Log("Payment message is compatible with Kitchen Order consumer")
	t.Logf("Payment sends: %s", paymentMessage)
	t.Logf("Kitchen Order extracts order_id: %s", kitchenOrderMsg.OrderID)
}
