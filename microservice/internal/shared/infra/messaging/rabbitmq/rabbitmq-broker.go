package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"

	"tech_challenge/internal/shared/interfaces"
)

// AMQPChannel interface para permitir mocks
type AMQPChannel interface {
	Close() error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, deleteWhenUnused, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

// AMQPConnection interface para permitir mocks
type AMQPConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
}

type RabbitMQBroker struct {
	conn    AMQPConnection
	channel AMQPChannel
	config  RabbitMQConfig
	mu      sync.Mutex
}

type RabbitMQConfig struct {
	URL      string
	Exchange string
}

func NewRabbitMQBroker(config RabbitMQConfig) *RabbitMQBroker {
	return &RabbitMQBroker{
		config: config,
	}
}

func (r *RabbitMQBroker) Connect(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if r.config.Exchange != "" {
		err = r.channel.ExchangeDeclare(
			r.config.Exchange, // name
			"topic",           // type
			true,              // durable
			false,             // auto-deleted
			false,             // internal
			false,             // no-wait
			nil,               // arguments
		)
		if err != nil {
			r.channel.Close()
			r.conn.Close()
			return fmt.Errorf("failed to declare exchange: %w", err)
		}
	}

	log.Println("Connected to RabbitMQ successfully")
	return nil
}

func (r *RabbitMQBroker) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQBroker) Publish(ctx context.Context, queue string, message interfaces.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	// Declara a fila
	_, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	headers := amqp.Table{}
	for k, v := range message.Headers {
		headers[k] = v
	}

	err = r.channel.Publish(
		r.config.Exchange, // exchange
		queue,             // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message.Body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
			MessageId:    message.ID,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *RabbitMQBroker) Subscribe(ctx context.Context, queue string, handler interfaces.MessageHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	// Declara a fila
	_, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange if exchange is configured
	if r.config.Exchange != "" {
		err = r.channel.QueueBind(
			queue,             // queue name
			queue,             // routing key
			r.config.Exchange, // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	msgs, err := r.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}

				headers := make(map[string]string)
				if d.Headers != nil {
					for k, v := range d.Headers {
						if str, ok := v.(string); ok {
							headers[k] = str
						}
					}
				}

				msg := interfaces.Message{
					ID:      d.MessageId,
					Body:    d.Body,
					Headers: headers,
				}

				if err := handler(ctx, msg); err != nil {
					log.Printf("Error processing message: %v", err)
					if nackErr := d.Nack(false, true); nackErr != nil {
						log.Printf("Error nacking message: %v", nackErr)
					}
				} else {
					if ackErr := d.Ack(false); ackErr != nil {
						log.Printf("Error acking message: %v", ackErr)
					}
				}
			}
		}
	}()

	log.Printf("Subscribed to queue: %s", queue)
	return nil
}

func (r *RabbitMQBroker) Start(ctx context.Context) error {
	// RabbitMQ já inicia automaticamente quando Subscribe é chamado
	return nil
}

func (r *RabbitMQBroker) Stop() error {
	return r.Close()
}

// Helper function para serializar mensagens
func SerializeMessage(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// Helper function para deserializar mensagens
func DeserializeMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
