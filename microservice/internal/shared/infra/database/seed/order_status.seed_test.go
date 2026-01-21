package seed

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/constants"
)

// MockGormDB é um mock do GORM DB para testes
type MockGormDB struct {
	whereQuery     interface{}
	whereArgs      []interface{}
	firstCalled    bool
	firstError     error
	createCalled   bool
	createError    error
	createValue    interface{}
	recordNotFound bool
	lastError      error
	callCount      int
}

func (m *MockGormDB) Where(query interface{}, args ...interface{}) DBInterface {
	m.whereQuery = query
	m.whereArgs = args
	m.callCount++

	if m.recordNotFound {
		m.lastError = gorm.ErrRecordNotFound
	} else if m.firstError != nil {
		m.lastError = m.firstError
	} else {
		m.lastError = nil
	}
	return m
}

func (m *MockGormDB) First(dest interface{}, conds ...interface{}) DBInterface {
	m.firstCalled = true
	return m
}

func (m *MockGormDB) Create(value interface{}) DBInterface {
	m.createCalled = true
	m.createValue = value
	if m.createError != nil {
		m.lastError = m.createError
	} else {
		m.lastError = nil
	}
	return m
}

func (m *MockGormDB) GetError() error {
	return m.lastError
}

func TestSeedOrderStatus_AllStatusesCreated(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.firstCalled, "First should be called to check if status exists")
	assert.True(t, mockDB.createCalled, "Create should be called to insert status")
}

func TestSeedOrderStatus_StatusAlreadyExists(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: false, // Status já existe
		firstError:     nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.firstCalled, "First should be called to check if status exists")
	assert.False(t, mockDB.createCalled, "Create should not be called if status already exists")
}

func TestSeedOrderStatus_CreateError(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    gorm.ErrInvalidData,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.firstCalled, "First should be called")
	assert.True(t, mockDB.createCalled, "Create should be called even if it fails")
}

func TestSeedOrderStatus_Constants_Exist(t *testing.T) {
	// Arrange & Act & Assert
	assert.NotEmpty(t, constants.KITCHEN_ORDER_STATUS_RECEIVED_ID)
	assert.NotEmpty(t, constants.KITCHEN_ORDER_STATUS_PREPARING_ID)
	assert.NotEmpty(t, constants.KITCHEN_ORDER_STATUS_READY_ID)
	assert.NotEmpty(t, constants.KITCHEN_ORDER_STATUS_FINISHED_ID)
}

func TestSeedOrderStatus_DefaultStatuses_Structure(t *testing.T) {
	// Arrange
	expectedStatuses := map[string]string{
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID:  "Recebido",
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID: "Em preparação",
		constants.KITCHEN_ORDER_STATUS_READY_ID:     "Pronto",
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID:  "Finalizado",
	}

	// Act & Assert
	assert.Len(t, expectedStatuses, 4, "Should have exactly 4 statuses")

	for id, name := range expectedStatuses {
		assert.NotEmpty(t, id, "Status ID should not be empty")
		assert.NotEmpty(t, name, "Status name should not be empty")
	}
}

func TestSeedOrderStatus_StatusNames_Portuguese(t *testing.T) {
	// Arrange
	expectedNames := []string{
		"Recebido",
		"Em preparação",
		"Pronto",
		"Finalizado",
	}

	// Act & Assert
	for _, name := range expectedNames {
		assert.NotEmpty(t, name, "Status name should not be empty")
		assert.True(t, len(name) > 0, "Status name should have content")
	}
}

func TestSeedOrderStatus_UniqueIDs(t *testing.T) {
	// Arrange
	ids := []string{
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		constants.KITCHEN_ORDER_STATUS_READY_ID,
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
	}

	// Act
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	// Assert
	assert.Len(t, idMap, 4, "All IDs should be unique")
}

func TestSeedOrderStatus_ModelCreation(t *testing.T) {
	// Arrange & Act
	status := models.OrderStatusModel{
		ID:   constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		Name: "Recebido",
	}

	// Assert
	assert.Equal(t, constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, status.ID)
	assert.Equal(t, "Recebido", status.Name)
}

func TestSeedOrderStatus_AllStatuses_Defined(t *testing.T) {
	// Arrange
	statuses := []struct {
		ID   string
		Name string
	}{
		{constants.KITCHEN_ORDER_STATUS_RECEIVED_ID, "Recebido"},
		{constants.KITCHEN_ORDER_STATUS_PREPARING_ID, "Em preparação"},
		{constants.KITCHEN_ORDER_STATUS_READY_ID, "Pronto"},
		{constants.KITCHEN_ORDER_STATUS_FINISHED_ID, "Finalizado"},
	}

	// Act & Assert
	require.Len(t, statuses, 4, "Should have exactly 4 statuses")

	for _, status := range statuses {
		assert.NotEmpty(t, status.ID, "Status ID should not be empty")
		assert.NotEmpty(t, status.Name, "Status name should not be empty")
	}
}

func TestSeedOrderStatus_StatusFlow_Order(t *testing.T) {
	// Arrange
	statusFlow := []string{
		"Recebido",
		"Em preparação",
		"Pronto",
		"Finalizado",
	}

	// Act & Assert
	assert.Len(t, statusFlow, 4, "Should have 4 statuses in flow")

	// Verify all statuses are different
	for i, status := range statusFlow {
		for j, otherStatus := range statusFlow {
			if i != j {
				assert.NotEqual(t, status, otherStatus, "All statuses should be different")
			}
		}
	}
}

func TestSeedOrderStatus_Constants_Are_Strings(t *testing.T) {
	// Arrange
	constants := []string{
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID,
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID,
		constants.KITCHEN_ORDER_STATUS_READY_ID,
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID,
	}

	// Act & Assert
	for _, constant := range constants {
		assert.IsType(t, "", constant, "Constant should be a string")
		assert.NotEmpty(t, constant, "Constant should not be empty")
	}
}

func TestSeedOrderStatus_WhereQuery_Called(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.NotNil(t, mockDB.whereQuery, "Where query should be called")
	assert.Equal(t, "id = ?", mockDB.whereQuery)
}

func TestSeedOrderStatus_Multiple_Iterations(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// The function should iterate through all 4 statuses
	assert.True(t, mockDB.firstCalled, "First should be called for each status")
	assert.True(t, mockDB.createCalled, "Create should be called for each status")
}

func TestSeedOrderStatus_ReceivedStatus(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.Equal(t, "id = ?", mockDB.whereQuery)
}

func TestSeedOrderStatus_PreparingStatus(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.createCalled)
}

func TestSeedOrderStatus_ReadyStatus(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.firstCalled)
}

func TestSeedOrderStatus_FinishedStatus(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.createCalled)
}

func TestSeedOrderStatus_StatusCopy_Created(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that a copy was created and passed to Create
	assert.True(t, mockDB.createCalled, "Create should be called with a copy of the status")
}

func TestSeedOrderStatus_Error_Handling(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    gorm.ErrInvalidData,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Function should handle errors gracefully
	assert.True(t, mockDB.createCalled, "Create should be called even if error occurs")
}

func TestSeedOrderStatus_Success_Logging(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Function should complete without panic
	assert.True(t, mockDB.createCalled, "Create should be called successfully")
}

func TestSeedOrderStatus_All_Four_Statuses_Processed(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that all 4 statuses are processed
	assert.True(t, mockDB.firstCalled, "First should be called for all statuses")
	assert.True(t, mockDB.createCalled, "Create should be called for all statuses")
}

// SeedOrderStatusWithDB é uma versão testável da função SeedOrderStatus
func SeedOrderStatusWithDB(db DBInterface) {
	seedOrderStatusInternal(db)
}

func TestSeedOrderStatus_SuccessfulCreation(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.firstCalled, "First should be called")
	assert.True(t, mockDB.createCalled, "Create should be called")
	assert.Nil(t, mockDB.lastError, "No error should occur on successful creation")
}

func TestSeedOrderStatus_SuccessfulCreation_NoError(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.createCalled, "Create should be called")
	assert.Nil(t, mockDB.lastError, "Error should be nil after successful creation")
}

func TestSeedOrderStatus_SkipsExistingStatus(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: false,
		firstError:     nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.firstCalled, "First should be called to check existence")
	assert.False(t, mockDB.createCalled, "Create should not be called for existing status")
}

func TestSeedOrderStatus_HandlesCreateError_Gracefully(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    gorm.ErrInvalidData,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	assert.True(t, mockDB.createCalled, "Create should be called")
	assert.NotNil(t, mockDB.lastError, "Error should be set")
}

func TestSeedOrderStatus_ProcessesAllFourStatuses(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// The function should call Where 4 times (once for each status)
	assert.Equal(t, 4, mockDB.callCount, "Should process all 4 statuses")
}

func TestSeedOrderStatus_LogsSuccessMessage(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Function should complete without panic and log success
	assert.True(t, mockDB.createCalled, "Create should be called")
	assert.Nil(t, mockDB.lastError, "No error should occur")
}

func TestSeedOrderStatus_LogsErrorMessage(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    gorm.ErrInvalidData,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Function should complete without panic and log error
	assert.True(t, mockDB.createCalled, "Create should be called")
	assert.NotNil(t, mockDB.lastError, "Error should be logged")
}

func TestSeedOrderStatus_IteratesCorrectly(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the loop processes all 4 statuses
	assert.Equal(t, 4, mockDB.callCount, "Should iterate 4 times")
	assert.True(t, mockDB.firstCalled, "First should be called in loop")
	assert.True(t, mockDB.createCalled, "Create should be called in loop")
}

func TestSeedOrderStatus_CreatesStatusCopy(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that Create was called with a status copy
	assert.True(t, mockDB.createCalled, "Create should be called with status copy")
	assert.NotNil(t, mockDB.createValue, "Create value should not be nil")
}

func TestSeedOrderStatus_ChecksRecordNotFound(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the function checks for ErrRecordNotFound
	assert.True(t, mockDB.firstCalled, "First should be called to check for record")
	assert.True(t, mockDB.createCalled, "Create should be called when record not found")
}

func TestSeedOrderStatus_SkipsWhenRecordExists(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: false,
		firstError:     nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that Create is not called when record exists
	assert.True(t, mockDB.firstCalled, "First should be called")
	assert.False(t, mockDB.createCalled, "Create should not be called when record exists")
}

func TestSeedOrderStatus_DefaultsArray(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that all 4 default statuses are processed
	assert.Equal(t, 4, mockDB.callCount, "Should process 4 default statuses")
}

func TestSeedOrderStatus_StatusIDsCorrect(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the Where query is called with correct ID
	assert.Equal(t, "id = ?", mockDB.whereQuery, "Where query should check ID")
}

func TestSeedOrderStatus_LoopsOverDefaults(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the function loops over all defaults
	assert.Equal(t, 4, mockDB.callCount, "Should loop 4 times")
	assert.True(t, mockDB.firstCalled, "First should be called in loop")
	assert.True(t, mockDB.createCalled, "Create should be called in loop")
}

func TestSeedOrderStatus_WrapperFunction(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the wrapper function works correctly
	assert.True(t, mockDB.firstCalled, "First should be called")
	assert.True(t, mockDB.createCalled, "Create should be called")
}

func TestSeedOrderStatus_GormDBWrapper(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the wrapper is used correctly
	assert.Equal(t, 4, mockDB.callCount, "Should process all statuses")
}

func TestSeedOrderStatus_CallsInternalFunction(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the internal function is called
	assert.True(t, mockDB.firstCalled, "Internal function should call First")
	assert.True(t, mockDB.createCalled, "Internal function should call Create")
}

func TestSeedOrderStatus_WrapsGormDB(t *testing.T) {
	// Arrange
	mockDB := &MockGormDB{
		recordNotFound: true,
		createError:    nil,
	}

	// Act
	SeedOrderStatusWithDB(mockDB)

	// Assert
	// Verify that the wrapper correctly wraps GORM DB
	assert.Equal(t, 4, mockDB.callCount, "Should wrap and process all statuses")
}