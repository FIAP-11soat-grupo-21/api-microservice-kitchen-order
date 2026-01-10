package use_cases

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/env"
	"tech_challenge/internal/shared/interfaces"
)

type UpdateKitchenOrderUseCase struct {
	gateway       gateways.KitchenOrderGateway
	statusGateway gateways.OrderStatusGateway
	messageBroker interfaces.MessageBroker
}

func NewUpdateKitchenOrderUseCase(
	gateway gateways.KitchenOrderGateway, 
	statusGateway gateways.OrderStatusGateway,
	messageBroker interfaces.MessageBroker,
) *UpdateKitchenOrderUseCase {
	return &UpdateKitchenOrderUseCase{
		gateway:       gateway,
		statusGateway: statusGateway,
		messageBroker: messageBroker,
	}
}

type KitchenOrderStatusUpdateMessage struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

func (ko *UpdateKitchenOrderUseCase) Execute(kitchenOrderDTO dtos.UpdateKitchenOrderDTO) (entities.KitchenOrder, error) {
	err := entities.ValidateID(kitchenOrderDTO.ID)

	if err != nil {
		return entities.KitchenOrder{}, err
	}

	kitchenOrder, err := ko.gateway.FindByID(kitchenOrderDTO.ID)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.KitchenOrderNotFoundException{}
	}

	kitchenOrderStatus, err := ko.statusGateway.FindByID(kitchenOrderDTO.StatusID)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.OrderStatusNotFoundException{}
	}

	kitchenOrder.Status.ID = kitchenOrderDTO.StatusID
	kitchenOrder.Status.Name = kitchenOrderStatus.Name

	now := time.Now()
	kitchenOrder.UpdatedAt = &now

	err = ko.gateway.Update(kitchenOrder)

	if err != nil {
		return entities.KitchenOrder{}, &exceptions.InvalidKitchenOrderDataException{}
	}

	// Notificar Orders apenas para status específicos
	if ko.shouldNotifyOrders(kitchenOrderStatus.Name.Value()) {
		ko.notifyOrdersService(kitchenOrder, kitchenOrderStatus.Name.Value())
	}

	return kitchenOrder, nil
}

func (ko *UpdateKitchenOrderUseCase) shouldNotifyOrders(status string) bool {
	notifiableStatuses := []string{"Em preparação", "Pronto", "Finalizado"}
	for _, s := range notifiableStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func (ko *UpdateKitchenOrderUseCase) notifyOrdersService(kitchenOrder entities.KitchenOrder, status string) {
	message := KitchenOrderStatusUpdateMessage{
		OrderID:   kitchenOrder.OrderID,
		Status:    status,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling kitchen order status update message: %v", err)
		return
	}

	config := env.GetConfig()
	var queueName string
	
	if config.MessageBroker.Type == "rabbitmq" {
		queueName = config.MessageBroker.RabbitMQ.OrdersQueue
	} else if config.MessageBroker.Type == "sqs" {
		queueName = config.MessageBroker.SQS.OrdersQueueURL
	} else {
		log.Printf("Unsupported message broker type: %s", config.MessageBroker.Type)
		return
	}

	msg := interfaces.Message{
		Body:    messageBody,
		Headers: map[string]string{"message-type": "kitchen-order-status-update"},
	}

	if err := ko.messageBroker.Publish(context.Background(), queueName, msg); err != nil {
		log.Printf("Failed to send kitchen order status update to orders service: %v", err)
	} else {
		log.Printf("Kitchen order status update sent to orders service: OrderID=%s, Status=%s", kitchenOrder.OrderID, status)
	}
}
