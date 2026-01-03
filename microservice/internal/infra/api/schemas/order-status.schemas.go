package schemas

type CreateOrderStatusSchema struct {
	Name int `json:"name" example:"Pronto" binding:"required"`
}

type OrderStatusResponseSchema struct {
	ID   string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Recebido"`
}

type ErrorMessageSchema struct {
	Error string `json:"error" example:"Internal server error"`
}
