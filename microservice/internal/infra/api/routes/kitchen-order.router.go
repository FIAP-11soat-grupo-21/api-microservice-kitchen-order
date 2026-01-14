package routes

import (
	"github.com/gin-gonic/gin"

	"tech_challenge/internal/infra/api/handlers"
)

func RegisterKitchenOrderRoutes(router *gin.RouterGroup) {
	kitchenOrderHandler := handlers.NewKitchenOrderHandler()
	orderStatusHandler := handlers.NewOrderStatusHandler()

	// GET
	router.GET("/", kitchenOrderHandler.FindAll)
	router.GET("/:id", kitchenOrderHandler.FindByID)
	
	router.PUT("/:id", kitchenOrderHandler.Update)
	
	// Status endpoints
	router.GET("/status", orderStatusHandler.FindAll)
}
