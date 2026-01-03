package api

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"tech_challenge/internal/infra/api/routes"
	"tech_challenge/internal/infra/messaging/consumers"
	"tech_challenge/internal/shared/config/env"
	"tech_challenge/internal/shared/factories"
	"tech_challenge/internal/shared/infra/api/handlers"
	"tech_challenge/internal/shared/infra/api/middlewares"
	file_router "tech_challenge/internal/shared/infra/api/routes"
	_ "tech_challenge/internal/shared/infra/api/swagger"
	"tech_challenge/internal/shared/infra/database"
)

func Init() {
	config := env.GetConfig()

	if config.IsProduction() {
		log.Printf("Running in production mode on [%s] with message broker: %s", config.APIUrl, config.MessageBroker.Type)
		gin.SetMode(gin.ReleaseMode)
	}

	database.Connect()

	if config.Database.RunMigrations {
		database.RunMigrations()
	}

	database.SeedDefaults()

	ginRouter := gin.Default()

	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(middlewares.ErrorHandlerMiddleware())

	// Health check endpoint
	healthHandler := handlers.NewHealthHandler()
	ginRouter.GET("/health", healthHandler.Health)

	v1Routes := ginRouter.Group("/v1")

	file_router.RegisterFileRoutes(v1Routes.Group("/uploads"))
	routes.RegisterKitchenOrderRoutes(v1Routes.Group("/kitchen-orders"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker, err := factories.NewMessageBroker(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize message broker: %v", err)
	}
	defer broker.Close()

	kitchenOrderConsumer := consumers.NewKitchenOrderConsumer(broker)
	if err := kitchenOrderConsumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start kitchen order consumer: %v", err)
	}

	if err := broker.Start(ctx); err != nil {
		log.Fatalf("Failed to start message broker: %v", err)
	}

	log.Println("Message broker consumers started successfully")

	go func() {
		if err := ginRouter.Run(config.APIUrl); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	log.Printf("HTTP server started on %s", config.APIUrl)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")

	if err := broker.Stop(); err != nil {
		log.Printf("Error stopping broker: %v", err)
	}
	cancel()
}
