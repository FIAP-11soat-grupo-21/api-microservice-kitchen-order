package exceptions

import "testing"

func TestOrderStatusNotFoundException_DefaultMessage(t *testing.T) {
	// Arrange
	exception := &OrderStatusNotFoundException{}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Order Status not found"
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func TestOrderStatusNotFoundException_CustomMessage(t *testing.T) {
	// Arrange
	customMessage := "Custom order status not found message"
	exception := &OrderStatusNotFoundException{Message: customMessage}

	// Act
	message := exception.Error()

	// Assert
	if message != customMessage {
		t.Errorf("Expected message '%s', got '%s'", customMessage, message)
	}
}

func TestOrderStatusNotFoundException_EmptyMessage(t *testing.T) {
	// Arrange
	exception := &OrderStatusNotFoundException{Message: ""}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Order Status not found"
	if message != expectedMessage {
		t.Errorf("Expected default message '%s', got '%s'", expectedMessage, message)
	}
}

func TestInvalidOrderStatusDataException_DefaultMessage(t *testing.T) {
	// Arrange
	exception := &InvalidOrderStatusDataException{}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Invalid Order Status data"
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func TestInvalidOrderStatusDataException_CustomMessage(t *testing.T) {
	// Arrange
	customMessage := "Custom invalid order status data message"
	exception := &InvalidOrderStatusDataException{Message: customMessage}

	// Act
	message := exception.Error()

	// Assert
	if message != customMessage {
		t.Errorf("Expected message '%s', got '%s'", customMessage, message)
	}
}

func TestInvalidOrderStatusDataException_EmptyMessage(t *testing.T) {
	// Arrange
	exception := &InvalidOrderStatusDataException{Message: ""}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Invalid Order Status data"
	if message != expectedMessage {
		t.Errorf("Expected default message '%s', got '%s'", expectedMessage, message)
	}
}

func TestOrderStatusExceptions_AsError(t *testing.T) {
	// Arrange
	notFoundException := &OrderStatusNotFoundException{}
	invalidDataException := &InvalidOrderStatusDataException{}

	// Act & Assert
	var err error

	err = notFoundException
	if err.Error() != "Order Status not found" {
		t.Errorf("OrderStatusNotFoundException should implement error interface")
	}

	err = invalidDataException
	if err.Error() != "Invalid Order Status data" {
		t.Errorf("InvalidOrderStatusDataException should implement error interface")
	}
}

func TestOrderStatusExceptions_TypeAssertion(t *testing.T) {
	// Arrange
	var err error = &OrderStatusNotFoundException{Message: "Test message"}

	// Act
	if notFoundErr, ok := err.(*OrderStatusNotFoundException); ok {
		// Assert
		if notFoundErr.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", notFoundErr.Message)
		}
	} else {
		t.Error("Type assertion to OrderStatusNotFoundException failed")
	}

	// Arrange
	err = &InvalidOrderStatusDataException{Message: "Invalid data"}

	// Act
	if invalidErr, ok := err.(*InvalidOrderStatusDataException); ok {
		// Assert
		if invalidErr.Message != "Invalid data" {
			t.Errorf("Expected message 'Invalid data', got '%s'", invalidErr.Message)
		}
	} else {
		t.Error("Type assertion to InvalidOrderStatusDataException failed")
	}
}