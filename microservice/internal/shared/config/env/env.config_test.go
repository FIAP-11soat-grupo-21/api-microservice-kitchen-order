package env

import (
	"os"
	"testing"
)

func TestGetConfig_Singleton(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()

	config1 := GetConfig()
	config2 := GetConfig()

	if config1 != config2 {
		t.Error("Expected GetConfig to return the same instance (singleton pattern)")
	}
}

func TestConfig_Environment(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		isProd   bool
		isDev    bool
	}{
		{"production", "production", true, false},
		{"development", "development", false, true},
		{"test", "test", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnv()
			defer cleanupTestEnv()
			
			os.Setenv("GO_ENV", tt.env)
			
			config := &Config{}
			config.Load()

			if config.IsProduction() != tt.isProd {
				t.Errorf("Expected IsProduction() = %v, got %v", tt.isProd, config.IsProduction())
			}
			
			if config.IsDevelopment() != tt.isDev {
				t.Errorf("Expected IsDevelopment() = %v, got %v", tt.isDev, config.IsDevelopment())
			}
		})
	}
}

func TestConfig_MessageBroker(t *testing.T) {
	tests := []struct {
		name       string
		brokerType string
		envVars    map[string]string
		validate   func(*testing.T, *Config)
	}{
		{
			name:       "SQS",
			brokerType: "sqs",
			envVars: map[string]string{
				"AWS_SQS_KITCHEN_ORDERS_QUEUE": "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
			},
			validate: func(t *testing.T, config *Config) {
				if config.MessageBroker.Type != "sqs" {
					t.Errorf("Expected MessageBroker.Type 'sqs', got %s", config.MessageBroker.Type)
				}
				if config.MessageBroker.SQS.QueueURL != "https://sqs.us-east-1.amazonaws.com/123456789/test-queue" {
					t.Errorf("Expected SQS QueueURL 'https://sqs.us-east-1.amazonaws.com/123456789/test-queue', got %s", config.MessageBroker.SQS.QueueURL)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnv()
			defer cleanupTestEnv()
			
			os.Setenv("MESSAGE_BROKER_TYPE", tt.brokerType)
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			
			config := &Config{}
			config.Load()
			
			tt.validate(t, config)
		})
	}
}

func TestConfig_Database(t *testing.T) {
	setupTestEnv()
	defer cleanupTestEnv()
	
	dbEnvVars := map[string]string{
		"DB_RUN_MIGRATIONS": "true",
		"DB_HOST":           "localhost",
		"DB_NAME":           "test_db",
		"DB_PORT":           "5432",
		"DB_USERNAME":       "test_user",
		"DB_PASSWORD":       "test_pass",
	}
	
	for key, value := range dbEnvVars {
		os.Setenv(key, value)
	}
	
	config := &Config{}
	config.Load()

	if !config.Database.RunMigrations {
		t.Error("Expected Database.RunMigrations to be true")
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host 'localhost', got %s", config.Database.Host)
	}
	if config.Database.Name != "test_db" {
		t.Errorf("Expected Database.Name 'test_db', got %s", config.Database.Name)
	}
}

func setupTestEnv() {
	defaultEnvVars := map[string]string{
		"GO_ENV":                        "test",
		"API_PORT":                      "8080",
		"API_HOST":                      "localhost",
		"DB_RUN_MIGRATIONS":             "false",
		"DB_HOST":                       "localhost",
		"DB_NAME":                       "test",
		"DB_PORT":                       "5432",
		"DB_USERNAME":                   "test",
		"DB_PASSWORD":                   "test",
		"AWS_REGION":                    "us-east-1",
		"MESSAGE_BROKER_TYPE":           "sqs",
		"AWS_SQS_KITCHEN_ORDERS_QUEUE":  "https://sqs.us-east-1.amazonaws.com/123456789/test-queue",
		"AWS_SQS_ORDERS_QUEUE":          "https://sqs.us-east-1.amazonaws.com/123456789/orders-queue",
	}
	
	for key, value := range defaultEnvVars {
		os.Setenv(key, value)
	}
}

func cleanupTestEnv() {
	envVars := []string{
		"GO_ENV", "API_PORT", "API_HOST", "DB_RUN_MIGRATIONS",
		"DB_HOST", "DB_NAME", "DB_PORT", "DB_USERNAME", "DB_PASSWORD",
		"AWS_REGION", "MESSAGE_BROKER_TYPE",
		"AWS_SQS_KITCHEN_ORDERS_QUEUE", "AWS_SQS_ORDERS_QUEUE",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}