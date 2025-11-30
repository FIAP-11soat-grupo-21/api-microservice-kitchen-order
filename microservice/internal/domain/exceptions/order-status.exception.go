package exceptions

type OrderStatusNotFoundException struct {
	Message string
}

type InvalidOrderStatusDataException struct {
	Message string
}

func (e *OrderStatusNotFoundException) Error() string {
	if e.Message == "" {
		return "Order Status not found"
	}

	return e.Message
}

func (e *InvalidOrderStatusDataException) Error() string {
	if e.Message == "" {
		return "Invalid Order Status data"
	}

	return e.Message
}
