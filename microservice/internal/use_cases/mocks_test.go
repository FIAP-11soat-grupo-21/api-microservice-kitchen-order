package use_cases

import (
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/application/gateways"
	"tech_challenge/internal/daos"
	"tech_challenge/internal/domain/entities"
	"tech_challenge/internal/domain/exceptions"
	"tech_challenge/internal/shared/config/constants"
)

// MockDataStore simula um banco de dados em memória para testes
type MockDataStore struct {
	kitchenOrders              []entities.KitchenOrder
	orderStatuses              []entities.OrderStatus
	shouldReturnError          bool
	errorToReturn              error
	shouldReturnErrorOnUpdate  bool
	updateErrorToReturn        error
}

func NewMockDataStore() *MockDataStore {
	statuses := []struct{ id, name string }{
		{constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido"},
		{constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação"},
		{constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto"},
		{constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado"},
	}

	orderStatuses := make([]entities.OrderStatus, len(statuses))
	for i, s := range statuses {
		status, _ := entities.NewOrderStatus(s.id, s.name)
		orderStatuses[i] = *status
	}

	return &MockDataStore{
		kitchenOrders:              []entities.KitchenOrder{},
		orderStatuses:              orderStatuses,
		shouldReturnError:          false,
		errorToReturn:              nil,
		shouldReturnErrorOnUpdate:  false,
		updateErrorToReturn:        nil,
	}
}

// Mock DataSource para KitchenOrder
type MockKitchenOrderDataSource struct {
	dataStore *MockDataStore
}

func (ds *MockKitchenOrderDataSource) Insert(kitchenOrder daos.KitchenOrderDAO) error {
	if ds.dataStore.shouldReturnError {
		return ds.dataStore.errorToReturn
	}

	status, _ := entities.NewOrderStatus(kitchenOrder.Status.ID, kitchenOrder.Status.Name)
	order, _ := entities.NewKitchenOrder(
		kitchenOrder.ID, kitchenOrder.OrderID, kitchenOrder.Slug,
		*status, kitchenOrder.CreatedAt, kitchenOrder.UpdatedAt,
	)
	ds.dataStore.kitchenOrders = append(ds.dataStore.kitchenOrders, *order)
	return nil
}

func (ds *MockKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	if ds.dataStore.shouldReturnError {
		return daos.KitchenOrderDAO{}, ds.dataStore.errorToReturn
	}

	for _, order := range ds.dataStore.kitchenOrders {
		if order.ID == id {
			return ds.entityToDAO(order), nil
		}
	}
	return daos.KitchenOrderDAO{}, &exceptions.KitchenOrderNotFoundException{}
}

func (ds *MockKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	if ds.dataStore.shouldReturnError {
		return nil, ds.dataStore.errorToReturn
	}

	var result []daos.KitchenOrderDAO
	for _, order := range ds.dataStore.kitchenOrders {
		if ds.matchesFilter(order, filter) {
			result = append(result, ds.entityToDAO(order))
		}
	}
	return result, nil
}

func (ds *MockKitchenOrderDataSource) Update(kitchenOrder daos.KitchenOrderDAO) error {
	if ds.dataStore.shouldReturnErrorOnUpdate {
		return ds.dataStore.updateErrorToReturn
	}

	for i, order := range ds.dataStore.kitchenOrders {
		if order.ID == kitchenOrder.ID {
			status, _ := entities.NewOrderStatus(kitchenOrder.Status.ID, kitchenOrder.Status.Name)
			updatedOrder, _ := entities.NewKitchenOrder(
				kitchenOrder.ID, kitchenOrder.OrderID, kitchenOrder.Slug,
				*status, kitchenOrder.CreatedAt, kitchenOrder.UpdatedAt,
			)
			updatedOrder.Amount = kitchenOrder.Amount
			updatedOrder.CustomerID = kitchenOrder.CustomerID
			
			// Adiciona os itens
			for _, itemDAO := range kitchenOrder.Items {
				item, _ := entities.NewOrderItem(
					itemDAO.ID, itemDAO.OrderID, itemDAO.ProductID,
					itemDAO.Quantity, itemDAO.UnitPrice,
				)
				updatedOrder.AddItem(*item)
			}
			
			ds.dataStore.kitchenOrders[i] = *updatedOrder
			return nil
		}
	}
	return &exceptions.KitchenOrderNotFoundException{}
}

func (ds *MockKitchenOrderDataSource) entityToDAO(order entities.KitchenOrder) daos.KitchenOrderDAO {
	items := make([]daos.OrderItemDAO, len(order.Items))
	for i, item := range order.Items {
		items[i] = daos.OrderItemDAO{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return daos.KitchenOrderDAO{
		ID:         order.ID,
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
		Slug:       order.Slug.Value(),
		Status: daos.OrderStatusDAO{
			ID:   order.Status.ID,
			Name: order.Status.Name.Value(),
		},
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func (ds *MockKitchenOrderDataSource) matchesFilter(order entities.KitchenOrder, filter dtos.KitchenOrderFilter) bool {
	if filter.CreatedAtFrom != nil && order.CreatedAt.Before(*filter.CreatedAtFrom) {
		return false
	}
	if filter.CreatedAtTo != nil && order.CreatedAt.After(*filter.CreatedAtTo) {
		return false
	}
	if filter.StatusID != nil && order.Status.ID != constants.KITCHEN_ORDER_STATUS_RECEIVED_ID {
		// Simula filtro por status - aqui simplificamos para o teste
		return false
	}
	return true
}

// Mock DataSource para OrderStatus
type MockOrderStatusDataSource struct {
	dataStore *MockDataStore
}

func (ds *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
	if ds.dataStore.shouldReturnError {
		return daos.OrderStatusDAO{}, ds.dataStore.errorToReturn
	}

	for _, status := range ds.dataStore.orderStatuses {
		if status.ID == id {
			return daos.OrderStatusDAO{
				ID:   status.ID,
				Name: status.Name.Value(),
			}, nil
		}
	}
	return daos.OrderStatusDAO{}, &exceptions.OrderStatusNotFoundException{}
}

func (ds *MockOrderStatusDataSource) FindAll() ([]daos.OrderStatusDAO, error) {
	if ds.dataStore.shouldReturnError {
		return nil, ds.dataStore.errorToReturn
	}

	result := make([]daos.OrderStatusDAO, len(ds.dataStore.orderStatuses))
	for i, status := range ds.dataStore.orderStatuses {
		result[i] = daos.OrderStatusDAO{
			ID:   status.ID,
			Name: status.Name.Value(),
		}
	}
	return result, nil
}

// Funções helper para criar gateways com mocks
func NewMockKitchenOrderGateway(dataStore *MockDataStore) gateways.KitchenOrderGateway {
	dataSource := &MockKitchenOrderDataSource{dataStore: dataStore}
	return *gateways.NewKitchenOrderGateway(dataSource)
}

func NewMockOrderStatusGateway(dataStore *MockDataStore) gateways.OrderStatusGateway {
	dataSource := &MockOrderStatusDataSource{dataStore: dataStore}
	return *gateways.NewOrderStatusGateway(dataSource)
}