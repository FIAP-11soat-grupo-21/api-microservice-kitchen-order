package factories

import (
	"tech_challenge/internal/infra/database/data_sources"
	"tech_challenge/internal/interfaces"
)

func NewKitchenOrderDataSource() interfaces.IKitchenOrderDataSource {
	return data_sources.NewGormKitchenOrderDataSource()
}

func NewOrderStatusDataSource() interfaces.IOrderStatusDataSource {
	return data_sources.NewGormOrderStatusDataSource()
}
