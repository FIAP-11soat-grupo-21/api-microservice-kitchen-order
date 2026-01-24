package sqs

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"

	"tech_challenge/internal/shared/interfaces"
)

func TestUnmarshalSNSMessage(t *testing.T) {
	broker := NewSQSBroker(SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	})

	paymentMessage := map[string]interface{}{
		"order_id": "896156b7-d058-44cd-8969-2b269567a1f2",
		"status":   "Confirmado",
	}

	paymentJSON, _ := json.Marshal(paymentMessage)

	snsNotification := SNSNotification{
		Type:      "Notification",
		MessageId: "test-message-id",
		TopicArn:  "arn:aws:sns:us-east-2:216989122312:payment-processed-topic",
		Message:   string(paymentJSON),
	}

	snsJSON, _ := json.Marshal(snsNotification)

	sqsMessage := types.Message{
		MessageId: aws.String("test-id"),
		Body:      aws.String(string(snsJSON)),
	}

	body, err := broker.unmarshalMessageBody(sqsMessage)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "896156b7-d058-44cd-8969-2b269567a1f2", result["order_id"])
	assert.Equal(t, "Confirmado", result["status"])
}

func TestUnmarshalDirectSQSMessage(t *testing.T) {
	broker := NewSQSBroker(SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	})

	directMessage := map[string]interface{}{
		"order_id": "direct-order-id",
		"status":   "Pending",
	}

	directJSON, _ := json.Marshal(directMessage)

	sqsMessage := types.Message{
		MessageId: aws.String("test-id"),
		Body:      aws.String(string(directJSON)),
	}

	body, err := broker.unmarshalMessageBody(sqsMessage)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "direct-order-id", result["order_id"])
	assert.Equal(t, "Pending", result["status"])
}

func TestProcessMessageWithSNS(t *testing.T) {
	mockClient := &mockSQSClient{
		deleteMessageFunc: func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
			return &sqs.DeleteMessageOutput{}, nil
		},
	}

	broker := NewSQSBroker(SQSConfig{
		Region:   "us-east-1",
		QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
	})
	broker.SetClient(mockClient)

	paymentMessage := map[string]interface{}{
		"order_id": "test-order-123",
		"status":   "Confirmado",
	}
	paymentJSON, _ := json.Marshal(paymentMessage)

	snsNotification := SNSNotification{
		Type:      "Notification",
		MessageId: "sns-msg-id",
		TopicArn:  "arn:aws:sns:us-east-2:216989122312:payment-processed-topic",
		Message:   string(paymentJSON),
	}
	snsJSON, _ := json.Marshal(snsNotification)

	sqsMessage := types.Message{
		MessageId:     aws.String("sqs-msg-id"),
		Body:          aws.String(string(snsJSON)),
		ReceiptHandle: aws.String("receipt-handle"),
	}

	handlerCalled := false
	handler := func(ctx context.Context, msg interfaces.Message) error {
		handlerCalled = true

		var result map[string]interface{}
		err := json.Unmarshal(msg.Body, &result)
		assert.NoError(t, err)
		assert.Equal(t, "test-order-123", result["order_id"])
		assert.Equal(t, "Confirmado", result["status"])

		return nil
	}

	broker.processMessage(context.Background(), sqsMessage, handler)

	assert.True(t, handlerCalled, "Handler should have been called")
}

func TestSNSNotificationStructure(t *testing.T) {
	snsJSON := `{
		"Type": "Notification",
		"MessageId": "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
		"TopicArn": "arn:aws:sns:us-east-2:216989122312:payment-processed-topic",
		"Message": "{\"order_id\":\"896156b7-d058-44cd-8969-2b269567a1f2\",\"status\":\"Confirmado\"}",
		"Timestamp": "2026-01-24T00:08:56.000Z",
		"SignatureVersion": "1",
		"Signature": "EXAMPLE",
		"SigningCertURL": "EXAMPLE",
		"UnsubscribeURL": "EXAMPLE"
	}`

	var notification SNSNotification
	err := json.Unmarshal([]byte(snsJSON), &notification)
	assert.NoError(t, err)
	assert.Equal(t, "Notification", notification.Type)
	assert.Equal(t, "arn:aws:sns:us-east-2:216989122312:payment-processed-topic", notification.TopicArn)

	var innerMessage map[string]interface{}
	err = json.Unmarshal([]byte(notification.Message), &innerMessage)
	assert.NoError(t, err)
	assert.Equal(t, "896156b7-d058-44cd-8969-2b269567a1f2", innerMessage["order_id"])
	assert.Equal(t, "Confirmado", innerMessage["status"])
}
