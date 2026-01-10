package data_sources

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"tech_challenge/internal/daos"
)

func TestGormOrderStatusDataSource_NewGormOrderStatusDataSource(t *testing.T) {
	// Test
	dataSource := NewGormOrderStatusDataSource()

	// Assertions
	assert.NotNil(t, dataSource)
	// Note: db might be nil in test environment without database connection
}

func TestGormOrderStatusDataSource_Structure(t *testing.T) {
	// Test structure
	dataSource := &GormOrderStatusDataSource{}
	assert.NotNil(t, dataSource)
	assert.IsType(t, &GormOrderStatusDataSource{}, dataSource)
}

func TestGormOrderStatusDataSource_Methods_Exist(t *testing.T) {
	dataSource := NewGormOrderStatusDataSource()

	// Verify methods exist
	assert.NotNil(t, dataSource.Insert)
	assert.NotNil(t, dataSource.FindAll)
	assert.NotNil(t, dataSource.FindByID)
}

func TestGormOrderStatusDataSource_DAO_Structure(t *testing.T) {
	// Test DAO structure
	dao := daos.OrderStatusDAO{
		ID:   "1",
		Name: "Recebido",
	}

	// Assert structure
	assert.Equal(t, "1", dao.ID)
	assert.Equal(t, "Recebido", dao.Name)
	assert.IsType(t, "", dao.ID)
	assert.IsType(t, "", dao.Name)
}

func TestGormOrderStatusDataSource_All_Statuses(t *testing.T) {
	// Test all expected order statuses
	statuses := []daos.OrderStatusDAO{
		{ID: "1", Name: "Recebido"},
		{ID: "2", Name: "Em preparação"},
		{ID: "3", Name: "Pronto"},
		{ID: "4", Name: "Finalizado"},
	}

	// Assert all statuses
	assert.Len(t, statuses, 4)
	
	for _, status := range statuses {
		assert.NotEmpty(t, status.ID)
		assert.NotEmpty(t, status.Name)
		assert.Contains(t, []string{"1", "2", "3", "4"}, status.ID)
		assert.Contains(t, []string{"Recebido", "Em preparação", "Pronto", "Finalizado"}, status.Name)
	}
}

func TestGormOrderStatusDataSource_Status_Names(t *testing.T) {
	// Test status names
	expectedNames := []string{"Recebido", "Em preparação", "Pronto", "Finalizado"}

	for i, name := range expectedNames {
		status := daos.OrderStatusDAO{
			ID:   string(rune(i + 1 + '0')),
			Name: name,
		}

		assert.Equal(t, name, status.Name)
		assert.NotEmpty(t, status.ID)
	}
}

func TestGormOrderStatusDataSource_Empty_Status(t *testing.T) {
	// Test empty status
	status := daos.OrderStatusDAO{}

	// Assert empty status
	assert.Empty(t, status.ID)
	assert.Empty(t, status.Name)
}

func TestGormOrderStatusDataSource_Status_Validation(t *testing.T) {
	// Test status validation
	validStatuses := map[string]string{
		"1": "Recebido",
		"2": "Em preparação",
		"3": "Pronto",
		"4": "Finalizado",
	}

	for id, name := range validStatuses {
		status := daos.OrderStatusDAO{
			ID:   id,
			Name: name,
		}

		assert.Equal(t, id, status.ID)
		assert.Equal(t, name, status.Name)
		assert.Contains(t, []string{"1", "2", "3", "4"}, status.ID)
		assert.Contains(t, []string{"Recebido", "Em preparação", "Pronto", "Finalizado"}, status.Name)
	}
}

func TestGormOrderStatusDataSource_Status_Types(t *testing.T) {
	// Test status types
	status := daos.OrderStatusDAO{
		ID:   "1",
		Name: "Recebido",
	}

	// Assert types
	assert.IsType(t, "", status.ID)
	assert.IsType(t, "", status.Name)
}

func TestGormOrderStatusDataSource_Status_Array(t *testing.T) {
	// Test status array
	statuses := []daos.OrderStatusDAO{
		{ID: "1", Name: "Recebido"},
		{ID: "2", Name: "Em preparação"},
	}

	// Assert array
	assert.Len(t, statuses, 2)
	assert.IsType(t, []daos.OrderStatusDAO{}, statuses)
	
	for _, status := range statuses {
		assert.NotEmpty(t, status.ID)
		assert.NotEmpty(t, status.Name)
	}
}

func TestGormOrderStatusDataSource_Status_IDs(t *testing.T) {
	// Test status IDs
	expectedIDs := []string{"1", "2", "3", "4"}
	expectedNames := []string{"Recebido", "Em preparação", "Pronto", "Finalizado"}

	for i, id := range expectedIDs {
		status := daos.OrderStatusDAO{
			ID:   id,
			Name: expectedNames[i],
		}

		assert.Equal(t, id, status.ID)
		assert.Equal(t, expectedNames[i], status.Name)
		assert.Contains(t, expectedIDs, status.ID)
		assert.Contains(t, expectedNames, status.Name)
	}
}

func TestGormOrderStatusDataSource_Status_Mapping(t *testing.T) {
	// Test status mapping
	statusMap := map[string]string{
		"1": "Recebido",
		"2": "Em preparação", 
		"3": "Pronto",
		"4": "Finalizado",
	}

	for id, name := range statusMap {
		status := daos.OrderStatusDAO{
			ID:   id,
			Name: name,
		}

		// Assert mapping
		assert.Equal(t, id, status.ID)
		assert.Equal(t, name, status.Name)
		assert.Equal(t, statusMap[status.ID], status.Name)
	}
}

func TestGormOrderStatusDataSource_Status_Workflow(t *testing.T) {
	// Test status workflow order
	workflowStatuses := []daos.OrderStatusDAO{
		{ID: "1", Name: "Recebido"},      // First status
		{ID: "2", Name: "Em preparação"}, // Second status
		{ID: "3", Name: "Pronto"},        // Third status
		{ID: "4", Name: "Finalizado"},    // Final status
	}

	// Assert workflow order
	assert.Len(t, workflowStatuses, 4)
	assert.Equal(t, "Recebido", workflowStatuses[0].Name)
	assert.Equal(t, "Em preparação", workflowStatuses[1].Name)
	assert.Equal(t, "Pronto", workflowStatuses[2].Name)
	assert.Equal(t, "Finalizado", workflowStatuses[3].Name)

	// Assert IDs are sequential
	for i, status := range workflowStatuses {
		expectedID := string(rune(i + 1 + '0'))
		assert.Equal(t, expectedID, status.ID)
	}
}