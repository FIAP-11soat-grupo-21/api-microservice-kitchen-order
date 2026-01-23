package consumers

import (
	"context"
	"encoding/json"
	"log"

	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/factories"
	"tech_challenge/internal/shared/config/env"
	"tech_challenge/internal/shared/interfaces"
)

type KitchenOrderConsumer struct {
	broker                 interfaces.MessageBroker
	kitchenOrderController controllers.KitchenOrderController
}

func NewKitchenOrderConsumer(broker interfaces.MessageBroker) *KitchenOrderConsumer {
	kitchenOrderDataSource := factories.NewKitchenOrderDataSource()
	orderStatusDataSource := factories.NewOrderStatusDataSource()
	kitchenOrderController := controllers.NewKitchenOrderController(kitchenOrderDataSource, orderStatusDataSource, broker)

	return &KitchenOrderConsumer{
		broker:                 broker,
		kitchenOrderController: *kitchenOrderController,
	}
}

type CreateKitchenOrderMessage struct {
	OrderID string `json:"order_id"`
}

type KitchenOrderResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (c *KitchenOrderConsumer) Start(ctx context.Context) error {
	config := env.GetConfig()
	queueName := config.MessageBroker.SQS.QueueURL
	
	if err := c.broker.Subscribe(ctx, queueName, c.handleCreate); err != nil {
		return err
	}

	log.Printf("Kitchen order consumer started listening on SQS queue: %s", queueName)
	return nil
}

func (c *KitchenOrderConsumer) handleCreate(ctx context.Context, msg interfaces.Message) error {

	log.Printf("Message Body (raw): %s", string(msg.Body))
	var createMsg CreateKitchenOrderMessage
	if err := json.Unmarshal(msg.Body, &createMsg); err != nil {
		log.Printf("Error unmarshaling create message: %v", err)
		return err
	}

	log.Printf("Received kitchen order creation request for order: %s", createMsg.OrderID)

	kitchenOrder, err := c.kitchenOrderController.Create(dtos.CreateKitchenOrderDTO{
		OrderID: createMsg.OrderID,
	})

	response := KitchenOrderResponse{
		Success: err == nil,
	}

	if err != nil {
		response.Error = err.Error()
		log.Printf("Error creating kitchen order: %v", err)
	} else {
		response.Data = kitchenOrder
		log.Printf("Kitchen order created successfully: %s (Slug: %s)", kitchenOrder.ID, kitchenOrder.Slug)
	}

	responseBody, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		log.Printf("Error marshaling response: %v", marshalErr)
		return err
	}
	responseMsg := interfaces.Message{
		ID:      msg.ID,
		Body:    responseBody,
		Headers: map[string]string{"correlation-id": msg.ID},
	}

	if responseQueue, ok := msg.Headers["reply-to"]; ok {
		if publishErr := c.broker.Publish(ctx, responseQueue, responseMsg); publishErr != nil {
			log.Printf("Error publishing response message: %v", publishErr)
		}
	}

	return err
}