package use_cases

import (
	"errors"
	"testing"
	"time"

	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/domain/entities"
)

func TestFindAllKitchenOrdersUseCase_EmptyResult(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	dataStore.kitchenOrders = []entities.KitchenOrder{} // Lista vazia

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)
	filter := dtos.KitchenOrderFilter{}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 orders, got %d", len(result))
	}
}

func TestFindAllKitchenOrdersUseCase_WithFilters(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0] // Status "Recebido"

	now := time.Now()
	order1, _ := entities.NewKitchenOrder(
		"id1", "order1", "001", status, now, nil,
	)
	order2, _ := entities.NewKitchenOrder(
		"id2", "order2", "002", status, now.Add(-time.Hour), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*order1, *order2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)

	// Filtro por data
	fromTime := now.Add(-30 * time.Minute)
	toTime := now.Add(30 * time.Minute)
	statusID := uint(1)

	filter := dtos.KitchenOrderFilter{
		CreatedAtFrom: &fromTime,
		CreatedAtTo:   &toTime,
		StatusID:      &statusID,
	}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Deve retornar apenas order1 que está dentro do range de tempo
	if len(result) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result))
	}

	if len(result) > 0 && result[0].ID != "id1" {
		t.Errorf("Expected order ID 'id1', got %s", result[0].ID)
	}
}

func TestFindAllKitchenOrdersUseCase_GatewayError(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	dataStore.shouldReturnError = true
	dataStore.errorToReturn = errors.New("database connection error")

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)
	filter := dtos.KitchenOrderFilter{}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database connection error" {
		t.Errorf("Expected 'database connection error', got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestFindAllKitchenOrdersUseCase_NilFilter(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	order, _ := entities.NewKitchenOrder(
		"id1", "order1", "001", status, time.Now(), nil,
	)
	dataStore.kitchenOrders = []entities.KitchenOrder{*order}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)

	// Act - passando filtro vazio
	result, err := useCase.Execute(dtos.KitchenOrderFilter{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result))
	}
}

func TestFindAllKitchenOrdersUseCase_MultipleStatuses(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	receivedStatus := dataStore.orderStatuses[0] // "Recebido"
	preparingStatus := dataStore.orderStatuses[1] // "Em preparação"

	order1, _ := entities.NewKitchenOrder(
		"id1", "order1", "001", receivedStatus, time.Now(), nil,
	)
	order2, _ := entities.NewKitchenOrder(
		"id2", "order2", "002", preparingStatus, time.Now(), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*order1, *order2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)
	filter := dtos.KitchenOrderFilter{}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(result))
	}

	// Verifica se ambos os status estão presentes
	statusFound := make(map[string]bool)
	for _, order := range result {
		statusFound[order.Status.Name.Value()] = true
	}

	if !statusFound["Recebido"] {
		t.Error("Expected to find 'Recebido' status")
	}

	if !statusFound["Em preparação"] {
		t.Error("Expected to find 'Em preparação' status")
	}
}

func TestFindAllKitchenOrdersUseCase_LargeDataset(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	// Cria 100 pedidos
	for i := 0; i < 100; i++ {
		order, _ := entities.NewKitchenOrder(
			"id"+string(rune(i)), "order"+string(rune(i)), "00"+string(rune(i)), status, time.Now(), nil,
		)
		dataStore.kitchenOrders = append(dataStore.kitchenOrders, *order)
	}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)
	filter := dtos.KitchenOrderFilter{}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 100 {
		t.Errorf("Expected 100 orders, got %d", len(result))
	}
}

func TestFindAllKitchenOrdersUseCase_FilterByStatusOnly(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	receivedStatus := dataStore.orderStatuses[0] // "Recebido"
	preparingStatus := dataStore.orderStatuses[1] // "Em preparação"

	order1, _ := entities.NewKitchenOrder(
		"id1", "order1", "001", receivedStatus, time.Now(), nil,
	)
	order2, _ := entities.NewKitchenOrder(
		"id2", "order2", "002", preparingStatus, time.Now(), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*order1, *order2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)

	statusID := uint(1)
	filter := dtos.KitchenOrderFilter{
		StatusID: &statusID,
	}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// O mock filtra por status, então deve retornar apenas 1
	if len(result) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result))
	}
}

func TestFindAllKitchenOrdersUseCase_FilterByDateOnly(t *testing.T) {
	// Arrange
	dataStore := NewMockDataStore()
	status := dataStore.orderStatuses[0]

	now := time.Now()
	order1, _ := entities.NewKitchenOrder(
		"id1", "order1", "001", status, now, nil,
	)
	order2, _ := entities.NewKitchenOrder(
		"id2", "order2", "002", status, now.Add(-2*time.Hour), nil,
	)

	dataStore.kitchenOrders = []entities.KitchenOrder{*order1, *order2}

	kitchenOrderGateway := NewMockKitchenOrderGateway(dataStore)
	useCase := NewFindAllKitchenOrderUseCase(kitchenOrderGateway)

	fromTime := now.Add(-time.Hour)
	filter := dtos.KitchenOrderFilter{
		CreatedAtFrom: &fromTime,
	}

	// Act
	result, err := useCase.Execute(filter)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Deve retornar apenas order1
	if len(result) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result))
	}

	if len(result) > 0 && result[0].ID != "id1" {
		t.Errorf("Expected order ID 'id1', got %s", result[0].ID)
	}
}