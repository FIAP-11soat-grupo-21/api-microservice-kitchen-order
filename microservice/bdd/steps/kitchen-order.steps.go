package steps

import (
	"fmt"
	"tech_challenge/internal/daos"
	mock_interfaces "tech_challenge/internal/interfaces/mocks"
	"time"

	"github.com/golang/mock/gomock"
)

type KitchenOrderHelper struct {
	Ctrl   *gomock.Controller
	MockDS *mock_interfaces.MockIKitchenOrderDataSource

	valid struct {
		orderID    string
		customerID *string
		amount     float64
		status     daos.OrderStatusDAO
		slug       string
		items      []daos.OrderItemDAO
	}
	existingID string
	newStatus  string
	foundOrder *daos.KitchenOrderDAO
}

func (koh *KitchenOrderHelper) TheKitchenOrderDataIsValid() error {
	if koh.valid.orderID == "" {
		return fmt.Errorf("order ID cannot be empty")
	}
	if koh.valid.amount < 0 {
		return fmt.Errorf("amount cannot be negative")
	}
	if len(koh.valid.items) == 0 {
		return fmt.Errorf("kitchen order must have at least one item")
	}
	return nil
}

func (koh *KitchenOrderHelper) SetOrderID(orderID string) {
	koh.valid.orderID = orderID
	koh.valid.amount = 50.00
	koh.valid.status = daos.OrderStatusDAO{
		ID:   "1",
		Name: "RECEIVED",
	}
	koh.valid.slug = "order-slug"
	koh.valid.items = []daos.OrderItemDAO{
		{
			ID:        "item-1",
			OrderID:   orderID,
			ProductID: "prod-1",
			Quantity:  2,
			UnitPrice: 25.00,
		},
	}
}

func (koh *KitchenOrderHelper) ISendARequestToCreateANewKitchenOrder() error {
	const generatedID = "ko-123"

	kitchenOrder := daos.KitchenOrderDAO{
		ID:         generatedID,
		OrderID:    koh.valid.orderID,
		CustomerID: koh.valid.customerID,
		Amount:     koh.valid.amount,
		Status:     koh.valid.status,
		Slug:       koh.valid.slug,
		Items:      koh.valid.items,
		CreatedAt:  time.Now(),
	}

	koh.MockDS.EXPECT().Insert(gomock.Any()).Return(nil)

	err := koh.MockDS.Insert(kitchenOrder)
	if err != nil {
		return err
	}

	koh.existingID = generatedID
	return nil
}

func (koh *KitchenOrderHelper) KitchenOrderShouldBeCreated() error {
	if koh.existingID == "" {
		return fmt.Errorf("kitchen order was not created")
	}
	return nil
}

func (koh *KitchenOrderHelper) AKitchenOrderExistsWithID(id string) {
	koh.existingID = id
	koh.foundOrder = &daos.KitchenOrderDAO{
		ID:      id,
		OrderID: "order-123",
		Amount:  50.00,
		Status: daos.OrderStatusDAO{
			ID:   "1",
			Name: "RECEIVED",
		},
		Slug: "order-slug",
		Items: []daos.OrderItemDAO{
			{
				ID:        "item-1",
				OrderID:   "order-123",
				ProductID: "prod-1",
				Quantity:  2,
				UnitPrice: 25.00,
			},
		},
		CreatedAt: time.Now(),
	}
}

func (koh *KitchenOrderHelper) ISendARequestToFindTheKitchenOrderByID() error {
	koh.MockDS.EXPECT().FindByID(koh.existingID).Return(*koh.foundOrder, nil)

	order, err := koh.MockDS.FindByID(koh.existingID)
	if err != nil {
		return err
	}

	koh.foundOrder = &order
	return nil
}

func (koh *KitchenOrderHelper) TheKitchenOrderShouldBeReturnedSuccessfully() error {
	if koh.foundOrder == nil {
		return fmt.Errorf("kitchen order was not found")
	}
	if koh.foundOrder.ID != koh.existingID {
		return fmt.Errorf("returned kitchen order ID does not match")
	}
	return nil
}

func (koh *KitchenOrderHelper) TheNewStatusIs(status string) {
	koh.newStatus = status
}

func (koh *KitchenOrderHelper) ISendARequestToUpdateTheKitchenOrderStatus() error {
	updatedOrder := *koh.foundOrder
	updatedOrder.Status = daos.OrderStatusDAO{
		ID:   "2",
		Name: koh.newStatus,
	}
	now := time.Now()
	updatedOrder.UpdatedAt = &now

	koh.MockDS.EXPECT().Update(gomock.Any()).Return(nil)

	err := koh.MockDS.Update(updatedOrder)
	if err != nil {
		return err
	}

	koh.foundOrder = &updatedOrder
	return nil
}

func (koh *KitchenOrderHelper) TheKitchenOrderStatusShouldBeUpdatedSuccessfully() error {
	if koh.foundOrder == nil {
		return fmt.Errorf("kitchen order was not updated")
	}
	if koh.foundOrder.Status.Name != koh.newStatus {
		return fmt.Errorf("kitchen order status was not updated correctly")
	}
	if koh.foundOrder.UpdatedAt == nil {
		return fmt.Errorf("kitchen order UpdatedAt was not set")
	}
	return nil
}
