package daos

import "testing"

func TestOrderStatusDAO_Struct(t *testing.T) {
	dao := OrderStatusDAO{
		ID:   "status-1",
		Name: "Recebido",
	}

	if dao.ID != "status-1" {
		t.Errorf("Expected ID 'status-1', got %s", dao.ID)
	}

	if dao.Name != "Recebido" {
		t.Errorf("Expected Name 'Recebido', got %s", dao.Name)
	}
}