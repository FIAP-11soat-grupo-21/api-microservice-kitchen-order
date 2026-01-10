package use_cases

import (
	"errors"
	"testing"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

func TestUpdateKitchenOrderUseCase_InvalidID(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	invalidIDs := []string{
		"",
		"invalid-uuid",
		"123",
		"not-a-uuid-at-all",
	}

	for _, invalidID := range invalidIDs {
		updateDTO := dtos.UpdateKitchenOrderDTO{
			ID:       invalidID,
			StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		}

		// Act
		result, err := useCase.Execute(updateDTO)

		// Assert
		if err == nil {
			t.Errorf("Expected error for invalid ID '%s', got nil", invalidID)
		}

		if !result.IsEmpty() {
			t.Errorf("Expected empty result for invalid ID '%s', got %v", invalidID, result)
		}
	}
}

func TestUpdateKitchenOrderUseCase_OrderNotFound(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	dataStore.kitchenOrders = []entities.KitchenOrder{} // Lista vazia

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
	}

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*exceptions.KitchenOrderNotFoundException); !ok {
		t.Errorf("Expected KitchenOrderNotFoundException, got %T", err)
	}

	if !result.IsEmpty() {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestUpdateKitchenOrderUseCase_StatusNotFound(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	existingOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", status, time.Now(), nil,
	)
	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: "999", // Status inexistente
	}

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*exceptions.OrderStatusNotFoundException); !ok {
		t.Errorf("Expected OrderStatusNotFoundException, got %T", err)
	}

	if !result.IsEmpty() {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestUpdateKitchenOrderUseCase_GatewayUpdateError(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	existingOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", status, time.Now(), nil,
	)
	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	// Configura para retornar erro no update
	dataStore.shouldReturnErrorOnUpdate = true
	dataStore.updateErrorToReturn = errors.New("update failed")

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
	}

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if _, ok := err.(*exceptions.InvalidKitchenOrderDataException); !ok {
		t.Errorf("Expected InvalidKitchenOrderDataException, got %T", err)
	}

	if !result.IsEmpty() {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestUpdateKitchenOrderUseCase_AllStatusTransitions(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	
	// Testa todas as transições de status possíveis
	statusTransitions := []struct {
		from string
		to   string
	}{
		{constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, constants.KITCHEN_ORDER_STATUS_PREPARING_ID},
		{constants.KITCHEN_ORDER_STATUS_PREPARING_ID, constants.KITCHEN_ORDER_STATUS_READY_ID},
		{constants.KITCHEN_ORDER_STATUS_READY_ID, constants.KITCHEN_ORDER_STATUS_FINISHED_ID},
	}

	for _, transition := range statusTransitions {
		// Encontra o status inicial
		var initialStatus entities.OrderStatus
		for _, status := range dataStore.orderStatuses {
			if status.ID == transition.from {
				initialStatus = status
				break
			}
		}

		existingOrder, _ := entities.NewKitchenOrder(
			orderID, "order123", "001", initialStatus, time.Now(), nil,
		)
		dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

		kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
		orderStatusGateway := NewMockOrderStatusGateway(dataStore)
		useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

		updateDTO := dtos.UpdateKitchenOrderDTO{
			ID:       orderID,
			StatusID: transition.to,
		}

		// Act
		result, err := useCase.Execute(updateDTO)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for transition %s -> %s, got %v", transition.from, transition.to, err)
		}

		if result.Status.ID != transition.to {
			t.Errorf("Expected status ID %s, got %s", transition.to, result.Status.ID)
		}

		if result.UpdatedAt == nil {
			t.Error("Expected UpdatedAt to be set")
		}
	}
}

func TestUpdateKitchenOrderUseCase_UpdatedAtIsSet(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	originalTime := time.Now().Add(-time.Hour) // 1 hora atrás
	existingOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", status, originalTime, nil,
	)
	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
	}

	beforeUpdate := time.Now()

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set")
	}

	if result.UpdatedAt.Before(beforeUpdate) {
		t.Error("Expected UpdatedAt to be after the update operation")
	}

	if result.CreatedAt.Equal(originalTime) == false {
		t.Error("Expected CreatedAt to remain unchanged")
	}
}

func TestUpdateKitchenOrderUseCase_StatusNameIsUpdated(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	receivedStatus := dataStore.orderStatuses[0] // Status "Recebido"

	existingOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", receivedStatus, time.Now(), nil,
	)
	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
	}

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Status.ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, result.Status.ID)
	}

	if result.Status.Name.Value() != "Em preparação" {
		t.Errorf("Expected status name 'Em preparação', got %s", result.Status.Name.Value())
	}
}

func TestUpdateKitchenOrderUseCase_PreservesOrderData(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	originalOrder, _ := entities.NewKitchenOrder(
		orderID, "original-order", "original-slug", status, time.Now(), nil,
	)
	originalOrder.Amount = 123.45
	customerID := "customer-456"
	originalOrder.CustomerID = &customerID

	dataStore.kitchenOrders = []entities.KitchenOrder{*originalOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       orderID,
		StatusID: constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
	}

	// Act
	result, err := useCase.Execute(updateDTO)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verifica se os dados originais foram preservados
	if result.OrderID != "original-order" {
		t.Errorf("Expected OrderID 'original-order', got %s", result.OrderID)
	}

	if result.Amount != 123.45 {
		t.Errorf("Expected Amount 123.45, got %f", result.Amount)
	}

	if result.CustomerID == nil || *result.CustomerID != "customer-456" {
		t.Errorf("Expected CustomerID 'customer-456', got %v", result.CustomerID)
	}

	// Mas o status deve ter sido atualizado
	if result.Status.ID != constants.KITCHEN_ORDER_STATUS_PREPARING_ID {
		t.Errorf("Expected status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_PREPARING_ID, result.Status.ID)
	}
}

func TestUpdateKitchenOrderUseCase_MultipleUpdates(t *testing.T) {
	// Arrange
	orderID := "550e8400-e29b-41d4-a716-446655440000"
	dataStore := NewMockDataStore()
	receivedStatus := dataStore.orderStatuses[0]

	existingOrder, _ := entities.NewKitchenOrder(
		orderID, "order123", "001", receivedStatus, time.Now(), nil,
	)
	dataStore.kitchenOrders = []entities.KitchenOrder{*existingOrder}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	orderStatusGateway := NewMockOrderStatusGateway(dataStore)
	useCase := NewUpdateKitchenOrderUseCase(kitchenOrderGateway, orderStatusGateway)

	// Sequência de updates
	updates := []string{
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		constants.KITCHEN_ORDER_STATUS_READY_ID,
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
	}

	var lastResult entities.KitchenOrder
	for _, statusID := range updates {
		updateDTO := dtos.UpdateKitchenOrderDTO{
			ID:       orderID,
			StatusID: statusID,
		}

		// Act
		result, err := useCase.Execute(updateDTO)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for status %s, got %v", statusID, err)
		}

		if result.Status.ID != statusID {
			t.Errorf("Expected status ID %s, got %s", statusID, result.Status.ID)
		}

		// Atualiza o dataStore para o próximo update
		dataStore.kitchenOrders[0] = result
		lastResult = result
	}

	// Verifica se o último update foi aplicado corretamente
	if lastResult.Status.ID != constants.KITCHEN_ORDER_STATUS_FINISHED_ID {
		t.Errorf("Expected final status ID %s, got %s", constants.KITCHEN_ORDER_STATUS_FINISHED_ID, lastResult.Status.ID)
	}
}