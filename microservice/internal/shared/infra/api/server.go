package api

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"tech_challenge/internal/infra/api/routes"
	"tech_challenge/internal/shared/config/env"
	"tech_challenge/internal/shared/infra/api/middlewares"
	file_router "tech_challenge/internal/shared/infra/api/routes"
	_ "tech_challenge/internal/shared/infra/api/swagger"
	"tech_challenge/internal/shared/infra/database"
)

func Init() {
	config := env.GetConfig()

	if config.IsProduction() {
		log.Printf("Running in production mode on [%s]", config.APIUrl)
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

	v1Routes := ginRouter.Group("/v1")

	file_router.RegisterFileRoutes(v1Routes.Group("/uploads"))
	routes.RegisterKitchenOrderRoutes(v1Routes.Group("/kitchen-orders"))

	ginRouter.Run(config.APIUrl)
}
