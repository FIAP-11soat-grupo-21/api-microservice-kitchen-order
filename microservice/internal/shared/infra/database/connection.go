package database

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tech_challenge/internal/infra/database/models"
	"tech_challenge/internal/shared/config/env"
	"tech_challenge/internal/shared/infra/database/seed"
)

var (
	dbConnection *gorm.DB
	instance     *gorm.DB
	once         sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		instance = dbConnection
	})
	return instance
}

func Connect() {
	if dbConnection != nil {
		log.Println("Database connection already established")
		return
	}

	config := env.GetConfig()

	dsn := "host=" + config.Database.Host +
		" user=" + config.Database.Username +
		" dbname=" + config.Database.Name +
		" password=" + config.Database.Password +
		" port=" + config.Database.Port

	queryLogLevel := logger.Info

	if config.IsProduction() {
		queryLogLevel = logger.Error
	}

	queryLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,   // Limite para considerar uma query como lenta
			LogLevel:                  queryLogLevel, // Nível de log (Info mostra todas as queries, Error mostra apenas erros)
			IgnoreRecordNotFoundError: false,         // Mostrar erro para registros não encontrados
			Colorful:                  true,          // Saída colorida no terminal
		},
	)

	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 2 * time.Second

	for i := range maxRetries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: queryLogger,
		})

		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	dbConnection = db
}

func Close() {
	if dbConnection == nil {
		log.Println("Database connection already closed")
		return
	}

	sqlDriver, err := dbConnection.DB()

	if err != nil {
		log.Fatal("Failed to close database")
	}

	sqlDriver.Close()
}

func RunMigrations() {
	dbConnection.AutoMigrate(
		&models.KitchenOrderModel{},
		&models.OrderStatusModel{},
	)
}

func SeedDefaults() {
	seed.SeedOrderStatus(dbConnection)
}
