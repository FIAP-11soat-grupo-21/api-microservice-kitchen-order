package suites

import (
	"context"
	"tech_challenge/bdd/steps"
	mock_interfaces "tech_challenge/internal/interfaces/mocks"

	"github.com/cucumber/godog"
	"github.com/golang/mock/gomock"
)

type godogReporter struct{}

func (r *godogReporter) Errorf(_ string, _ ...interface{}) {
}
func (r *godogReporter) Fatalf(format string, _ ...interface{}) { panic("gomock fatal: " + format) }

func InitializeScenario(ctx *godog.ScenarioContext) {
	var helper *steps.KitchenOrderHelper

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ctrl := gomock.NewController(&godogReporter{})
		mockDS := mock_interfaces.NewMockIKitchenOrderDataSource(ctrl)

		helper = &steps.KitchenOrderHelper{
			Ctrl:   ctrl,
			MockDS: mockDS,
		}
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if helper != nil && helper.Ctrl != nil {
			helper.Ctrl.Finish()
		}
		return ctx, nil
	})

	// Create kitchen order steps
	ctx.Step(`^the kitchen order data is valid with order ID "([^"]*)"$`, func(orderID string) error {
		helper.SetOrderID(orderID)
		return helper.TheKitchenOrderDataIsValid()
	})
	ctx.Step(`^I send a request to create a new kitchen order$`, func() error {
		return helper.ISendARequestToCreateANewKitchenOrder()
	})
	ctx.Step(`^the kitchen order should be created successfully$`, func() error {
		return helper.KitchenOrderShouldBeCreated()
	})

	// Find kitchen order steps
	ctx.Step(`^a kitchen order exists with ID "([^"]*)"$`, func(id string) error {
		helper.AKitchenOrderExistsWithID(id)
		return nil
	})
	ctx.Step(`^I send a request to find the kitchen order by ID$`, func() error {
		return helper.ISendARequestToFindTheKitchenOrderByID()
	})
	ctx.Step(`^the kitchen order should be returned successfully$`, func() error {
		return helper.TheKitchenOrderShouldBeReturnedSuccessfully()
	})

	// Update kitchen order status steps
	ctx.Step(`^the new status is "([^"]*)"$`, func(status string) error {
		helper.TheNewStatusIs(status)
		return nil
	})
	ctx.Step(`^I send a request to update the kitchen order status$`, func() error {
		return helper.ISendARequestToUpdateTheKitchenOrderStatus()
	})
	ctx.Step(`^the kitchen order status should be updated successfully$`, func() error {
		return helper.TheKitchenOrderStatusShouldBeUpdatedSuccessfully()
	})
}
