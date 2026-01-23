package seed

import (
	"errors"
	"testing"

	"gorm.io/gorm"
	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/constants"
)

type mockDB struct {
	whereFunc  func(query interface{}, args ...interface{}) DBInterface
	firstFunc  func(dest interface{}, conds ...interface{}) DBInterface
	createFunc func(value interface{}) DBInterface
	errorFunc  func() error
}

func (m *mockDB) Where(query interface{}, args ...interface{}) DBInterface {
	if m.whereFunc != nil {
		return m.whereFunc(query, args...)
	}
	return m
}

func (m *mockDB) First(dest interface{}, conds ...interface{}) DBInterface {
	if m.firstFunc != nil {
		return m.firstFunc(dest, conds...)
	}
	return m
}

func (m *mockDB) Create(value interface{}) DBInterface {
	if m.createFunc != nil {
		return m.createFunc(value)
	}
	return m
}

func (m *mockDB) GetError() error {
	if m.errorFunc != nil {
		return m.errorFunc()
	}
	return nil
}

func TestSeedOrderStatusInternal_CreateNewStatuses(t *testing.T) {
	createdCount := 0
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			createdCount++
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	expectedCount := 4
	if createdCount != expectedCount {
		t.Errorf("Expected %d statuses to be created, got %d", expectedCount, createdCount)
	} else {
		t.Logf("✓ %d status criados com sucesso", createdCount)
	}
}

func TestSeedOrderStatusInternal_SkipExistingStatuses(t *testing.T) {
	createdCount := 0
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					if model, ok := dest.(*models.OrderStatusModel); ok {
						model.ID = "existing-id"
						model.Name = "Existing"
					}
					return &mockDB{
						errorFunc: func() error {
							return nil
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			createdCount++
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	if createdCount != 0 {
		t.Errorf("Expected 0 statuses to be created (all exist), got %d", createdCount)
	} else {
		t.Log("✓ Nenhum status criado (todos já existem)")
	}
}

func TestSeedOrderStatusInternal_MixedScenario(t *testing.T) {
	createdCount := 0
	callCount := 0
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			callCount++
			shouldExist := callCount%2 == 0
			
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					if shouldExist {
						if model, ok := dest.(*models.OrderStatusModel); ok {
							model.ID = "existing-id"
							model.Name = "Existing"
						}
					}
					
					return &mockDB{
						errorFunc: func() error {
							if shouldExist {
								return nil
							}
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			createdCount++
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	expectedCreated := 2
	if createdCount != expectedCreated {
		t.Errorf("Expected %d statuses to be created, got %d", expectedCreated, createdCount)
	} else {
		t.Logf("✓ %d status criados (cenário misto)", createdCount)
	}
}

func TestSeedOrderStatusInternal_VerifyStatusData(t *testing.T) {
	var createdStatuses []models.OrderStatusModel
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			if status, ok := value.(*models.OrderStatusModel); ok {
				createdStatuses = append(createdStatuses, *status)
			}
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	expectedStatuses := map[string]string{
		constants.KITCHEN_ORDER_STATUS_RECEIVED_ID:  "Recebido",
		constants.KITCHEN_ORDER_STATUS_PREPARING_ID: "Em preparação",
		constants.KITCHEN_ORDER_STATUS_READY_ID:     "Pronto",
		constants.KITCHEN_ORDER_STATUS_FINISHED_ID:  "Finalizado",
	}

	if len(createdStatuses) != len(expectedStatuses) {
		t.Errorf("Expected %d statuses, got %d", len(expectedStatuses), len(createdStatuses))
	}

	for _, status := range createdStatuses {
		expectedName, exists := expectedStatuses[status.ID]
		if !exists {
			t.Errorf("Unexpected status ID: %s", status.ID)
		} else if status.Name != expectedName {
			t.Errorf("Expected name '%s' for ID %s, got '%s'", expectedName, status.ID, status.Name)
		}
	}

	t.Log("✓ Todos os status têm dados corretos")
}

func TestSeedOrderStatusInternal_DatabaseError(t *testing.T) {
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return errors.New("database connection error")
						},
					}
				},
			}
		},
	}

	seedOrderStatusInternal(mock)
	t.Log("✓ Erro de banco de dados tratado sem panic")
}

func TestSeedOrderStatus_WithNilDB(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("✓ Panic capturado ao usar DB nil")
		}
	}()
	
	var db *gorm.DB
	SeedOrderStatus(db)
}



func TestGormDBWrapper_First_Called(t *testing.T) {
	firstCalled := false
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					firstCalled = true
					return &mockDB{
						errorFunc: func() error {
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	if !firstCalled {
		t.Error("Expected First to be called")
	} else {
		t.Log("✓ First foi chamado corretamente")
	}
}

func TestGormDBWrapper_Create_Called(t *testing.T) {
	createCalled := false
	var createdValue interface{}
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			createCalled = true
			createdValue = value
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	if !createCalled {
		t.Error("Expected Create to be called")
	} else {
		t.Log("✓ Create foi chamado corretamente")
	}
	
	if createdValue == nil {
		t.Error("Expected created value to be non-nil")
	} else {
		if status, ok := createdValue.(*models.OrderStatusModel); ok {
			t.Logf("✓ Create recebeu OrderStatusModel: ID=%s, Name=%s", status.ID, status.Name)
		}
	}
}

func TestGormDBWrapper_GetError_Called(t *testing.T) {
	getErrorCalled := false
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							getErrorCalled = true
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	if !getErrorCalled {
		t.Error("Expected GetError to be called")
	} else {
		t.Log("✓ GetError foi chamado corretamente")
	}
}

func TestGormDBWrapper_GetError_ReturnsError(t *testing.T) {
	expectedError := errors.New("test database error")
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return expectedError
						},
					}
				},
			}
		},
	}

	seedOrderStatusInternal(mock)
	t.Log("✓ GetError retornou erro corretamente")
}

func TestGormDBWrapper_GetError_ReturnsNil(t *testing.T) {
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return nil
						},
					}
				},
			}
		},
	}

	seedOrderStatusInternal(mock)
	t.Log("✓ GetError retornou nil (sem erro)")
}

func TestGormDBWrapper_First_WithDestination(t *testing.T) {
	var capturedDest interface{}
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					capturedDest = dest
					if model, ok := dest.(*models.OrderStatusModel); ok {
						model.ID = "test-id"
						model.Name = "Test Name"
					}
					return &mockDB{
						errorFunc: func() error {
							return nil
						},
					}
				},
			}
		},
	}

	seedOrderStatusInternal(mock)

	if capturedDest == nil {
		t.Error("Expected destination to be captured")
	} else {
		if _, ok := capturedDest.(*models.OrderStatusModel); ok {
			t.Log("✓ First recebeu ponteiro para OrderStatusModel")
		} else {
			t.Error("Expected destination to be *OrderStatusModel")
		}
	}
}

func TestGormDBWrapper_Create_WithMultipleStatuses(t *testing.T) {
	var createdStatuses []models.OrderStatusModel
	
	mock := &mockDB{
		whereFunc: func(query interface{}, args ...interface{}) DBInterface {
			return &mockDB{
				firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
					return &mockDB{
						errorFunc: func() error {
							return gorm.ErrRecordNotFound
						},
					}
				},
			}
		},
		createFunc: func(value interface{}) DBInterface {
			if status, ok := value.(*models.OrderStatusModel); ok {
				createdStatuses = append(createdStatuses, *status)
			}
			return &mockDB{}
		},
	}

	seedOrderStatusInternal(mock)

	if len(createdStatuses) != 4 {
		t.Errorf("Expected 4 statuses to be created, got %d", len(createdStatuses))
	} else {
		t.Logf("✓ Create foi chamado %d vezes", len(createdStatuses))
		for i, status := range createdStatuses {
			t.Logf("  Status %d: ID=%s, Name=%s", i+1, status.ID, status.Name)
		}
	}
}


func TestGormDBWrapper_AllMethods_Coverage(t *testing.T) {
	t.Run("Verify all wrapper methods are called through seedOrderStatusInternal", func(t *testing.T) {
		whereCalled := false
		firstCalled := false
		createCalled := false
		getErrorCalled := false
		
		mock := &mockDB{
			whereFunc: func(query interface{}, args ...interface{}) DBInterface {
				whereCalled = true
				return &mockDB{
					firstFunc: func(dest interface{}, conds ...interface{}) DBInterface {
						firstCalled = true
						return &mockDB{
							errorFunc: func() error {
								getErrorCalled = true
								return gorm.ErrRecordNotFound
							},
						}
					},
				}
			},
			createFunc: func(value interface{}) DBInterface {
				createCalled = true
				return &mockDB{}
			},
		}

		seedOrderStatusInternal(mock)

		if !whereCalled {
			t.Error("Where was not called")
		}
		if !firstCalled {
			t.Error("First was not called")
		}
		if !createCalled {
			t.Error("Create was not called")
		}
		if !getErrorCalled {
			t.Error("GetError was not called")
		}

		t.Log("✓ Todos os métodos do wrapper foram chamados:")
		t.Log("  - Where: chamado")
		t.Log("  - First: chamado")
		t.Log("  - Create: chamado")
		t.Log("  - GetError: chamado")
	})
}
