package value_objects

import "testing"

func TestNewName_Success(t *testing.T) {
	// Arrange
	validName := "Valid Name"

	// Act
	name, err := NewName(validName)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if name.Value() != validName {
		t.Errorf("Expected name %s, got %s", validName, name.Value())
	}
}

func TestNewName_TooShort(t *testing.T) {
	// Arrange
	shortName := "ab" // Menos de 3 caracteres

	// Act
	_, err := NewName(shortName)

	// Assert
	if err == nil {
		t.Error("Expected error for short name, got nil")
	}

	expectedMessage := "name must be at least 3 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewName_TooLong(t *testing.T) {
	// Arrange
	longName := ""
	for i := 0; i <= 100; i++ { // Mais de 100 caracteres
		longName += "a"
	}

	// Act
	_, err := NewName(longName)

	// Assert
	if err == nil {
		t.Error("Expected error for long name, got nil")
	}

	expectedMessage := "name must be at most 100 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewName_EmptyString(t *testing.T) {
	// Arrange
	emptyName := ""

	// Act
	_, err := NewName(emptyName)

	// Assert
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}

	expectedMessage := "name must be at least 3 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewName_MinimumLength(t *testing.T) {
	// Arrange
	minName := "abc" // Exatamente 3 caracteres

	// Act
	name, err := NewName(minName)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for minimum length name, got %v", err)
	}

	if name.Value() != minName {
		t.Errorf("Expected name %s, got %s", minName, name.Value())
	}
}

func TestNewName_MaximumLength(t *testing.T) {
	// Arrange
	maxName := ""
	for i := 0; i < 100; i++ { // Exatamente 100 caracteres
		maxName += "a"
	}

	// Act
	name, err := NewName(maxName)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for maximum length name, got %v", err)
	}

	if name.Value() != maxName {
		t.Errorf("Expected name length %d, got %d", len(maxName), len(name.Value()))
	}
}

func TestNewName_WithSpaces(t *testing.T) {
	// Arrange
	nameWithSpaces := "Name With Spaces"

	// Act
	name, err := NewName(nameWithSpaces)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for name with spaces, got %v", err)
	}

	if name.Value() != nameWithSpaces {
		t.Errorf("Expected name %s, got %s", nameWithSpaces, name.Value())
	}
}

func TestNewName_WithSpecialCharacters(t *testing.T) {
	// Arrange
	nameWithSpecialChars := "Name-123_@#"

	// Act
	name, err := NewName(nameWithSpecialChars)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for name with special characters, got %v", err)
	}

	if name.Value() != nameWithSpecialChars {
		t.Errorf("Expected name %s, got %s", nameWithSpecialChars, name.Value())
	}
}

func TestName_Value(t *testing.T) {
	// Arrange
	originalValue := "Test Name"
	name, _ := NewName(originalValue)

	// Act
	value := name.Value()

	// Assert
	if value != originalValue {
		t.Errorf("Expected Value() to return %s, got %s", originalValue, value)
	}
}