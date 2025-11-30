package exceptions

type KitchenOrderNotFoundException struct {
	Message string
}

type InvalidKitchenOrderDataException struct {
	Message string
}

func (e *KitchenOrderNotFoundException) Error() string {
	if e.Message == "" {
		return "Kitchen Order not found"
	}

	return e.Message
}

func (e *InvalidKitchenOrderDataException) Error() string {
	if e.Message == "" {
		return "Invalid Kitchen Order data"
	}

	return e.Message
}
