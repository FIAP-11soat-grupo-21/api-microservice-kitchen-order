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
	kitchenOrders []entities.KitchenOrder
	orderStatuses []entities.OrderStatus
}

func NewMockDataStore() *MockDataStore {
	// Inicializa com os status padrão
	receivedStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido")
	preparingStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação")
	readyStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto")
	finishedStatus, _ := entities.NewOrderStatus(constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado")

	return &MockDataStore{
		kitchenOrders: []entities.KitchenOrder{},
		orderStatuses: []entities.OrderStatus{
			*receivedStatus, *preparingStatus, *readyStatus, *finishedStatus,
		},
	}
}

// Mock DataSource para KitchenOrder
type MockKitchenOrderDataSource struct {
	dataStore *MockDataStore
}

func (ds *MockKitchenOrderDataSource) Insert(kitchenOrder daos.KitchenOrderDAO) error {
	// Converte DAO para entity e adiciona ao dataStore
	status, _ := entities.NewOrderStatus(kitchenOrder.Status.ID, kitchenOrder.Status.Name)
	order, _ := entities.NewKitchenOrder(
		kitchenOrder.ID,
		kitchenOrder.OrderID,
		kitchenOrder.Slug,
		*status,
		kitchenOrder.CreatedAt,
		kitchenOrder.UpdatedAt,
	)
	ds.dataStore.kitchenOrders = append(ds.dataStore.kitchenOrders, *order)
	return nil
}

func (ds *MockKitchenOrderDataSource) FindByID(id string) (daos.KitchenOrderDAO, error) {
	for _, order := range ds.dataStore.kitchenOrders {
		if order.ID == id {
			return daos.KitchenOrderDAO{
				ID:      order.ID,
				OrderID: order.OrderID,
				Slug:    order.Slug.Value(),
				Status: daos.OrderStatusDAO{
					ID:   order.Status.ID,
					Name: order.Status.Name.Value(),
				},
				CreatedAt: order.CreatedAt,
				UpdatedAt: order.UpdatedAt,
			}, nil
		}
	}
	return daos.KitchenOrderDAO{}, &exceptions.KitchenOrderNotFoundException{}
}

func (ds *MockKitchenOrderDataSource) FindAll(filter dtos.KitchenOrderFilter) ([]daos.KitchenOrderDAO, error) {
	var result []daos.KitchenOrderDAO

	for _, order := range ds.dataStore.kitchenOrders {
		include := true

		// Aplica filtros se fornecidos
		if filter.CreatedAtFrom != nil && order.CreatedAt.Before(*filter.CreatedAtFrom) {
			include = false
		}
		if filter.CreatedAtTo != nil && order.CreatedAt.After(*filter.CreatedAtTo) {
			include = false
		}

		if include {
			result = append(result, daos.KitchenOrderDAO{
				ID:      order.ID,
				OrderID: order.OrderID,
				Slug:    order.Slug.Value(),
				Status: daos.OrderStatusDAO{
					ID:   order.Status.ID,
					Name: order.Status.Name.Value(),
				},
				CreatedAt: order.CreatedAt,
				UpdatedAt: order.UpdatedAt,
			})
		}
	}

	return result, nil
}

func (ds *MockKitchenOrderDataSource) Update(kitchenOrder daos.KitchenOrderDAO) error {
	for i, order := range ds.dataStore.kitchenOrders {
		if order.ID == kitchenOrder.ID {
			// Converte DAO para entity e atualiza
			status, _ := entities.NewOrderStatus(kitchenOrder.Status.ID, kitchenOrder.Status.Name)
			updatedOrder, _ := entities.NewKitchenOrder(
				kitchenOrder.ID,
				kitchenOrder.OrderID,
				kitchenOrder.Slug,
				*status,
				kitchenOrder.CreatedAt,
				kitchenOrder.UpdatedAt,
			)
			ds.dataStore.kitchenOrders[i] = *updatedOrder
			return nil
		}
	}
	return &exceptions.KitchenOrderNotFoundException{}
}

// Mock DataSource para OrderStatus
type MockOrderStatusDataSource struct {
	dataStore *MockDataStore
}

func (ds *MockOrderStatusDataSource) FindByID(id string) (daos.OrderStatusDAO, error) {
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
	var result []daos.OrderStatusDAO
	for _, status := range ds.dataStore.orderStatuses {
		result = append(result, daos.OrderStatusDAO{
			ID:   status.ID,
			Name: status.Name.Value(),
		})
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
