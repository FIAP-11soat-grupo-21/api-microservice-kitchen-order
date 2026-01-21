package factories

import (
	"context"
	"fmt"

	"tech_challenge/internal/shared/config/env"
	"tech_challenge/internal/shared/infra/messaging/sqs"
	"tech_challenge/internal/shared/interfaces"
)

type MessageBrokerType string

const (
	MessageBrokerSQS MessageBrokerType = "sqs"
)

func NewMessageBroker(ctx context.Context) (interfaces.MessageBroker, error) {
	config := env.GetConfig()
	brokerType := MessageBrokerType(config.MessageBroker.Type)

	switch brokerType {
	case MessageBrokerSQS:
		broker := sqs.NewSQSBroker(sqs.SQSConfig{
			Region:   config.AWS.Region,
			QueueURL: config.MessageBroker.SQS.QueueURL,
		})
		if err := broker.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to SQS: %w", err)
		}
		return broker, nil

	default:
		return nil, fmt.Errorf("unsupported message broker type: %s", brokerType)
	}
}
