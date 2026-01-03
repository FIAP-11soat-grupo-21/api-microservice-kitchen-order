package exceptions

import "testing"

func TestKitchenOrderNotFoundException_DefaultMessage(t *testing.T) {
	// Arrange
	exception := &KitchenOrderNotFoundException{}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Kitchen Order not found"
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func TestKitchenOrderNotFoundException_CustomMessage(t *testing.T) {
	// Arrange
	customMessage := "Custom kitchen order not found message"
	exception := &KitchenOrderNotFoundException{Message: customMessage}

	// Act
	message := exception.Error()

	// Assert
	if message != customMessage {
		t.Errorf("Expected message '%s', got '%s'", customMessage, message)
	}
}

func TestKitchenOrderNotFoundException_EmptyMessage(t *testing.T) {
	// Arrange
	exception := &KitchenOrderNotFoundException{Message: ""}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Kitchen Order not found"
	if message != expectedMessage {
		t.Errorf("Expected default message '%s', got '%s'", expectedMessage, message)
	}
}

func TestInvalidKitchenOrderDataException_DefaultMessage(t *testing.T) {
	// Arrange
	exception := &InvalidKitchenOrderDataException{}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Invalid Kitchen Order data"
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}

func TestInvalidKitchenOrderDataException_CustomMessage(t *testing.T) {
	// Arrange
	customMessage := "Custom invalid kitchen order data message"
	exception := &InvalidKitchenOrderDataException{Message: customMessage}

	// Act
	message := exception.Error()

	// Assert
	if message != customMessage {
		t.Errorf("Expected message '%s', got '%s'", customMessage, message)
	}
}

func TestInvalidKitchenOrderDataException_EmptyMessage(t *testing.T) {
	// Arrange
	exception := &InvalidKitchenOrderDataException{Message: ""}

	// Act
	message := exception.Error()

	// Assert
	expectedMessage := "Invalid Kitchen Order data"
	if message != expectedMessage {
		t.Errorf("Expected default message '%s', got '%s'", expectedMessage, message)
	}
}

func TestKitchenOrderExceptions_AsError(t *testing.T) {
	// Arrange
	notFoundException := &KitchenOrderNotFoundException{}
	invalidDataException := &InvalidKitchenOrderDataException{}

	// Act & Assert
	var err error

	err = notFoundException
	if err.Error() != "Kitchen Order not found" {
		t.Errorf("KitchenOrderNotFoundException should implement error interface")
	}

	err = invalidDataException
	if err.Error() != "Invalid Kitchen Order data" {
		t.Errorf("InvalidKitchenOrderDataException should implement error interface")
	}
}

func TestKitchenOrderExceptions_TypeAssertion(t *testing.T) {
	// Arrange
	var err error = &KitchenOrderNotFoundException{Message: "Test message"}

	// Act
	if notFoundErr, ok := err.(*KitchenOrderNotFoundException); ok {
		// Assert
		if notFoundErr.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", notFoundErr.Message)
		}
	} else {
		t.Error("Type assertion to KitchenOrderNotFoundException failed")
	}

	// Arrange
	err = &InvalidKitchenOrderDataException{Message: "Invalid data"}

	// Act
	if invalidErr, ok := err.(*InvalidKitchenOrderDataException); ok {
		// Assert
		if invalidErr.Message != "Invalid data" {
			t.Errorf("Expected message 'Invalid data', got '%s'", invalidErr.Message)
		}
	} else {
		t.Error("Type assertion to InvalidKitchenOrderDataException failed")
	}
}