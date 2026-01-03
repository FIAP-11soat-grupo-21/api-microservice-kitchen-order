package value_objects

import "testing"

func TestNewSlug_Success(t *testing.T) {
	// Arrange
	validSlug := "001"

	// Act
	slug, err := NewSlug(validSlug)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if slug.Value() != validSlug {
		t.Errorf("Expected slug %s, got %s", validSlug, slug.Value())
	}
}

func TestNewSlug_TooShort(t *testing.T) {
	// Arrange
	shortSlug := "ab" // Menos de 3 caracteres

	// Act
	_, err := NewSlug(shortSlug)

	// Assert
	if err == nil {
		t.Error("Expected error for short slug, got nil")
	}

	expectedMessage := "name must be at least 3 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewSlug_TooLong(t *testing.T) {
	// Arrange
	longSlug := ""
	for i := 0; i <= 100; i++ { // Mais de 100 caracteres
		longSlug += "a"
	}

	// Act
	_, err := NewSlug(longSlug)

	// Assert
	if err == nil {
		t.Error("Expected error for long slug, got nil")
	}

	expectedMessage := "name must be at most 100 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewSlug_EmptyString(t *testing.T) {
	// Arrange
	emptySlug := ""

	// Act
	_, err := NewSlug(emptySlug)

	// Assert
	if err == nil {
		t.Error("Expected error for empty slug, got nil")
	}

	expectedMessage := "name must be at least 3 characters long"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestNewSlug_MinimumLength(t *testing.T) {
	// Arrange
	minSlug := "001" // Exatamente 3 caracteres

	// Act
	slug, err := NewSlug(minSlug)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for minimum length slug, got %v", err)
	}

	if slug.Value() != minSlug {
		t.Errorf("Expected slug %s, got %s", minSlug, slug.Value())
	}
}

func TestNewSlug_MaximumLength(t *testing.T) {
	// Arrange
	maxSlug := ""
	for i := 0; i < 100; i++ { // Exatamente 100 caracteres
		maxSlug += "a"
	}

	// Act
	slug, err := NewSlug(maxSlug)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for maximum length slug, got %v", err)
	}

	if slug.Value() != maxSlug {
		t.Errorf("Expected slug length %d, got %d", len(maxSlug), len(slug.Value()))
	}
}

func TestNewSlug_NumericSlug(t *testing.T) {
	// Arrange
	numericSlug := "123"

	// Act
	slug, err := NewSlug(numericSlug)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for numeric slug, got %v", err)
	}

	if slug.Value() != numericSlug {
		t.Errorf("Expected slug %s, got %s", numericSlug, slug.Value())
	}
}

func TestNewSlug_AlphanumericSlug(t *testing.T) {
	// Arrange
	alphanumericSlug := "abc123"

	// Act
	slug, err := NewSlug(alphanumericSlug)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for alphanumeric slug, got %v", err)
	}

	if slug.Value() != alphanumericSlug {
		t.Errorf("Expected slug %s, got %s", alphanumericSlug, slug.Value())
	}
}

func TestNewSlug_WithHyphens(t *testing.T) {
	// Arrange
	slugWithHyphens := "abc-123"

	// Act
	slug, err := NewSlug(slugWithHyphens)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for slug with hyphens, got %v", err)
	}

	if slug.Value() != slugWithHyphens {
		t.Errorf("Expected slug %s, got %s", slugWithHyphens, slug.Value())
	}
}

func TestSlug_Value(t *testing.T) {
	// Arrange
	originalValue := "test-slug"
	slug, _ := NewSlug(originalValue)

	// Act
	value := slug.Value()

	// Assert
	if value != originalValue {
		t.Errorf("Expected Value() to return %s, got %s", originalValue, value)
	}
}

func TestNewSlug_KitchenOrderSlugFormats(t *testing.T) {
	// Arrange & Act & Assert
	testCases := []string{
		"001", "002", "010", "100", "999",
	}

	for _, tc := range testCases {
		slug, err := NewSlug(tc)

		if err != nil {
			t.Errorf("Expected no error for kitchen order slug %s, got %v", tc, err)
		}

		if slug.Value() != tc {
			t.Errorf("Expected slug %s, got %s", tc, slug.Value())
		}
	}
}
