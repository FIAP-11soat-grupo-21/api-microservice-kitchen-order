package dtos

type OrderStatusDTO struct {
	ID   string
	Name string
}

type CreateOrderStatusDTO struct {
	Name string
}

type OrderStatusResponseDTO struct {
	ID   string
	Name string
}
