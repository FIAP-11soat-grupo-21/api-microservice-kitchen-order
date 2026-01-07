package env

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv        string
	APIPort      string
	APIHost      string
	APIUrl       string
	Database     struct {
		RunMigrations bool
		Host          string
		Name          string
		Port          string
		Username      string
		Password      string
	}
	PaymentGateway struct {
		AccessToken   string
		CollectorID   string
		ExternalPosID string
		ApiBaseURL    string
	}
	AWS struct {
		Region string
	}
	MessageBroker struct {
		Type     string
		RabbitMQ struct {
			URL      string
			Exchange string
		}
		SQS struct {
			QueueURL string
		}
	}
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.Load()
	})
	return instance
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func (c *Config) Load() {
	dotEnvPath := ".env"
	_, err := os.Stat(dotEnvPath)

	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	c.GoEnv = getEnv("GO_ENV")

	c.APIPort = getEnv("API_PORT")
	c.APIHost = getEnv("API_HOST")
	c.APIUrl = c.APIHost + ":" + c.APIPort

	c.Database.RunMigrations = getEnv("DB_RUN_MIGRATIONS") == "true"
	c.Database.Host = getEnv("DB_HOST")
	c.Database.Name = getEnv("DB_NAME")
	c.Database.Port = getEnv("DB_PORT")
	c.Database.Username = getEnv("DB_USERNAME")
	c.Database.Password = getEnv("DB_PASSWORD")

	c.AWS.Region = getEnv("AWS_REGION")
	// Message Broker configuration
	c.MessageBroker.Type = getEnv("MESSAGE_BROKER_TYPE")

	if c.MessageBroker.Type == "rabbitmq" {
		c.MessageBroker.RabbitMQ.URL = getEnv("RABBITMQ_URL")
	} else if c.MessageBroker.Type == "sqs" {
		c.MessageBroker.SQS.QueueURL = getEnv("SQS_QUEUE_URL")
	}
}

func (c *Config) IsProduction() bool {
	return c.GoEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.GoEnv == "development"
}
