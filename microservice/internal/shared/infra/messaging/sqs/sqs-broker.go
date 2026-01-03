package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"tech_challenge/internal/shared/interfaces"
)

type SQSBroker struct {
	client *sqs.Client
	config SQSConfig
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

type SQSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	QueueURL        string
}

func NewSQSBroker(config SQSConfig) *SQSBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &SQSBroker{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *SQSBroker) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.config.Region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s.client = sqs.NewFromConfig(cfg)
	log.Println("Connected to SQS successfully")
	return nil
}

func (s *SQSBroker) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

func (s *SQSBroker) Publish(ctx context.Context, queue string, message interfaces.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return fmt.Errorf("not connected to SQS")
	}

	messageAttributes := make(map[string]types.MessageAttributeValue)
	for k, v := range message.Headers {
		messageAttributes[k] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}

	if message.ID == "" {
		message.ID = fmt.Sprintf("msg-%d", len(message.Body))
	}

	_, err := s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:          aws.String(s.config.QueueURL),
		MessageBody:       aws.String(string(message.Body)),
		MessageAttributes: messageAttributes,
		MessageGroupId:    aws.String(message.ID),
	})

	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return nil
}

func (s *SQSBroker) Subscribe(ctx context.Context, queue string, handler interfaces.MessageHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return fmt.Errorf("not connected to SQS")
	}

	go s.pollMessages(ctx, queue, handler)

	log.Printf("Subscribed to queue: %s", queue)
	return nil
}

func (s *SQSBroker) pollMessages(ctx context.Context, queue string, handler interfaces.MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ctx.Done():
			return
		default:
			result, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(s.config.QueueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20, // Long polling
				MessageAttributeNames: []string{
					"All",
				},
			})

			if err != nil {
				log.Printf("Error receiving messages from SQS: %v", err)
				continue
			}

			for _, msg := range result.Messages {
				headers := make(map[string]string)
				for k, v := range msg.MessageAttributes {
					if v.StringValue != nil {
						headers[k] = *v.StringValue
					}
				}

				message := interfaces.Message{
					ID:      *msg.MessageId,
					Body:    []byte(*msg.Body),
					Headers: headers,
				}

				if err := handler(ctx, message); err != nil {
					log.Printf("Error processing message: %v", err)
					continue
				}

				_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(s.config.QueueURL),
					ReceiptHandle: msg.ReceiptHandle,
				})

				if err != nil {
					log.Printf("Error deleting message from SQS: %v", err)
				}
			}
		}
	}
}

func (s *SQSBroker) Start(ctx context.Context) error {
	return nil
}

func (s *SQSBroker) Stop() error {
	return s.Close()
}

func SerializeMessage(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func DeserializeMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
