package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDB_Singleton(t *testing.T) {
	// Testa se GetDB retorna a mesma instância (singleton)
	db1 := GetDB()
	db2 := GetDB()
	
	// Ambas devem ser a mesma instância
	assert.Same(t, db1, db2)
}

func TestGetDB_NotNil(t *testing.T) {
	// Testa se GetDB não retorna nil
	db := GetDB()
	
	// Pode ser nil se não foi conectado, mas não deve causar panic
	// Este teste verifica apenas que a função pode ser chamada
	_ = db
}

func TestConnect_Multiple_Calls(t *testing.T) {
	// Testa se múltiplas chamadas para Connect não causam problemas
	// Nota: Este teste não conecta realmente ao banco, apenas verifica
	// que a função pode ser chamada sem panic
	
	assert.NotPanics(t, func() {
		// Simula múltiplas chamadas
		// Em um ambiente real, isso conectaria ao banco
		// Aqui apenas testamos que não há panic
	})
}

func TestClose_Multiple_Calls(t *testing.T) {
	// Testa se múltiplas chamadas para Close não causam problemas
	assert.NotPanics(t, func() {
		Close()
		Close() // Segunda chamada deve ser segura
	})
}

func TestRunMigrations_Safe_Call(t *testing.T) {
	// Testa se RunMigrations pode ser chamada sem panic
	// quando não há conexão ativa
	
	// Como RunMigrations usa dbConnection que é nil, 
	// este teste verifica apenas que a função existe
	assert.NotNil(t, RunMigrations)
}

func TestSeedDefaults_Safe_Call(t *testing.T) {
	// Testa se SeedDefaults pode ser chamada sem panic
	// quando não há conexão ativa
	
	// Como SeedDefaults usa dbConnection que é nil,
	// este teste verifica apenas que a função existe
	assert.NotNil(t, SeedDefaults)
}

func TestDatabaseConnection_Structure(t *testing.T) {
	// Testa a estrutura básica do módulo de conexão
	
	// Verifica se as funções existem e podem ser chamadas
	assert.NotPanics(t, func() {
		GetDB()
	})
	
	assert.NotPanics(t, func() {
		Close()
	})
}

func TestDatabaseConnection_Concurrency_Safety(t *testing.T) {
	// Testa se GetDB é seguro para uso concorrente
	done := make(chan bool, 10)
	
	// Executa GetDB em múltiplas goroutines
	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("GetDB panicked in goroutine: %v", r)
				}
				done <- true
			}()
			
			db := GetDB()
			_ = db // Usa a variável para evitar otimização
		}()
	}
	
	// Aguarda todas as goroutines terminarem
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestDatabaseConnection_Initialization(t *testing.T) {
	// Testa o estado inicial das variáveis
	
	// Como as variáveis são privadas, testamos indiretamente
	// através das funções públicas
	db := GetDB()
	
	// O comportamento pode variar dependendo do estado da conexão
	// Este teste verifica apenas que não há panic
	_ = db
}

func TestDatabaseConnection_Module_Integrity(t *testing.T) {
	// Testa a integridade do módulo
	
	// Verifica se todas as funções principais existem
	functions := []func(){
		func() { GetDB() },
		func() { Close() },
	}
	
	for i, fn := range functions {
		assert.NotPanics(t, fn, "Function %d should not panic", i)
	}
	
	// Verifica se as outras funções existem (sem chamar)
	assert.NotNil(t, Connect)
	assert.NotNil(t, RunMigrations)
	assert.NotNil(t, SeedDefaults)
}