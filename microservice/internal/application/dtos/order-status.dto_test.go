package dtos

import "testing"

func TestOrderStatusDTO_Structure(t *testing.T) {
	// Arrange
	dto := OrderStatusDTO{
		ID:   "status-id-123",
		Name: "Recebido",
	}

	// Assert
	if dto.ID != "status-id-123" {
		t.Errorf("Expected ID 'status-id-123', got %s", dto.ID)
	}

	if dto.Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", dto.Name)
	}
}

func TestCreateOrderStatusDTO_Structure(t *testing.T) {
	// Arrange
	dto := CreateOrderStatusDTO{
		Name: "Em preparação",
	}

	// Assert
	if dto.Name != "Em preparação" {
		t.Errorf("Expected Name 'Em preparação', got %s", dto.Name)
	}
}

func TestOrderStatusResponseDTO_Structure(t *testing.T) {
	// Arrange
	dto := OrderStatusResponseDTO{
		ID:   "response-status-id",
		Name: "Pronto",
	}

	// Assert
	if dto.ID != "response-status-id" {
		t.Errorf("Expected ID 'response-status-id', got %s", dto.ID)
	}

	if dto.Name != "Pronto" {
		t.Errorf("Expected Name 'Pronto', got %s", dto.Name)
	}
}

func TestOrderStatusDTO_EmptyValues(t *testing.T) {
	// Arrange
	dto := OrderStatusDTO{
		ID:   "",
		Name: "",
	}

	// Assert
	if dto.ID != "" {
		t.Errorf("Expected empty ID, got %s", dto.ID)
	}

	if dto.Name != "" {
		t.Errorf("Expected empty Name, got %s", dto.Name)
	}
}

func TestCreateOrderStatusDTO_EmptyName(t *testing.T) {
	// Arrange
	dto := CreateOrderStatusDTO{
		Name: "",
	}

	// Assert
	if dto.Name != "" {
		t.Errorf("Expected empty Name, got %s", dto.Name)
	}
}

func TestOrderStatusResponseDTO_EmptyValues(t *testing.T) {
	// Arrange
	dto := OrderStatusResponseDTO{
		ID:   "",
		Name: "",
	}

	// Assert
	if dto.ID != "" {
		t.Errorf("Expected empty ID, got %s", dto.ID)
	}

	if dto.Name != "" {
		t.Errorf("Expected empty Name, got %s", dto.Name)
	}
}

func TestOrderStatusDTOs_AllStatuses(t *testing.T) {
	// Arrange & Act & Assert
	testCases := []struct {
		id   string
		name string
	}{
		{"56d3b3c3-1801-49cd-bae7-972c78082012", "Recebido"},
		{"3f9a1c98-7b2f-4f3b-8a96-c0b7c761a123", "Em preparação"},
		{"5a8b2b16-9b47-4e35-ae27-28f7994ef456", "Pronto"},
		{"bd91a1ee-1234-4cde-9c2a-efb1d2a3a789", "Finalizado"},
	}

	for _, tc := range testCases {
		// Test OrderStatusDTO
		dto := OrderStatusDTO{
			ID:   tc.id,
			Name: tc.name,
		}

		if dto.ID != tc.id {
			t.Errorf("OrderStatusDTO: Expected ID %s, got %s", tc.id, dto.ID)
		}

		if dto.Name != tc.name {
			t.Errorf("OrderStatusDTO: Expected Name %s, got %s", tc.name, dto.Name)
		}

		// Test CreateOrderStatusDTO
		createDTO := CreateOrderStatusDTO{
			Name: tc.name,
		}

		if createDTO.Name != tc.name {
			t.Errorf("CreateOrderStatusDTO: Expected Name %s, got %s", tc.name, createDTO.Name)
		}

		// Test OrderStatusResponseDTO
		responseDTO := OrderStatusResponseDTO{
			ID:   tc.id,
			Name: tc.name,
		}

		if responseDTO.ID != tc.id {
			t.Errorf("OrderStatusResponseDTO: Expected ID %s, got %s", tc.id, responseDTO.ID)
		}

		if responseDTO.Name != tc.name {
			t.Errorf("OrderStatusResponseDTO: Expected Name %s, got %s", tc.name, responseDTO.Name)
		}
	}
}
