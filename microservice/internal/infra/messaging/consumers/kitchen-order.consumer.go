package consumers

import (
	"context"
	"encoding/json"
	"log"

	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/factories"
	"tech_challenge/internal/shared/interfaces"
)

type KitchenOrderConsumer struct {
	broker                interfaces.MessageBroker
	kitchenOrderController controllers.KitchenOrderController
}

func NewKitchenOrderConsumer(broker interfaces.MessageBroker) *KitchenOrderConsumer {
	kitchenOrderDataSource := factories.NewKitchenOrderDataSource()
	orderStatusDataSource := factories.NewOrderStatusDataSource()
	kitchenOrderController := controllers.NewKitchenOrderController(kitchenOrderDataSource, orderStatusDataSource)

	return &KitchenOrderConsumer{
		broker:                broker,
		kitchenOrderController: *kitchenOrderController,
	}
}

type CreateKitchenOrderMessage struct {
	OrderID string `json:"order_id"`
}

type UpdateKitchenOrderMessage struct {
	ID       string `json:"id"`
	StatusID string `json:"status_id"`
}

type KitchenOrderResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (c *KitchenOrderConsumer) Start(ctx context.Context) error {
	if err := c.broker.Subscribe(ctx, "kitchen-order.create", c.handleCreate); err != nil {
		return err
	}

	if err := c.broker.Subscribe(ctx, "kitchen-order.update", c.handleUpdate); err != nil {
		return err
	}

	log.Println("Kitchen order consumers started (create and update only)")
	return nil
}

func (c *KitchenOrderConsumer) handleCreate(ctx context.Context, msg interfaces.Message) error {
	var createMsg CreateKitchenOrderMessage
	if err := json.Unmarshal(msg.Body, &createMsg); err != nil {
		log.Printf("Error unmarshaling create message: %v", err)
		return err
	}

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
		log.Printf("Kitchen order created successfully: %s", kitchenOrder.ID)
	}

	responseBody, _ := json.Marshal(response)
	responseMsg := interfaces.Message{
		ID:      msg.ID,
		Body:    responseBody,
		Headers: map[string]string{"correlation-id": msg.ID},
	}

	if responseQueue, ok := msg.Headers["reply-to"]; ok {
		c.broker.Publish(ctx, responseQueue, responseMsg)
	}

	return err
}

func (c *KitchenOrderConsumer) handleUpdate(ctx context.Context, msg interfaces.Message) error {
	var updateMsg UpdateKitchenOrderMessage
	if err := json.Unmarshal(msg.Body, &updateMsg); err != nil {
		log.Printf("Error unmarshaling update message: %v", err)
		return err
	}

	kitchenOrder, err := c.kitchenOrderController.Update(dtos.UpdateKitchenOrderDTO{
		ID:       updateMsg.ID,
		StatusID: updateMsg.StatusID,
	})

	response := KitchenOrderResponse{
		Success: err == nil,
	}

	if err != nil {
		response.Error = err.Error()
		log.Printf("Error updating kitchen order: %v", err)
	} else {
		response.Data = kitchenOrder
		log.Printf("Kitchen order updated successfully: %s", kitchenOrder.ID)
	}

	responseBody, _ := json.Marshal(response)
	responseMsg := interfaces.Message{
		ID:      msg.ID,
		Body:    responseBody,
		Headers: map[string]string{"correlation-id": msg.ID},
	}

	if responseQueue, ok := msg.Headers["reply-to"]; ok {
		c.broker.Publish(ctx, responseQueue, responseMsg)
	}

	return err
}

