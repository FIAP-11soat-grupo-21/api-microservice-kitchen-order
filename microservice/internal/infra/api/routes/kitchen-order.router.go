package routes

import (
	"github.com/gin-gonic/gin"

	"tech_challenge/internal/infra/api/handlers"
)

func RegisterKitchenOrderRoutes(router *gin.RouterGroup) {
	kitchenOrderHandler := handlers.NewKitchenOrderHandler()

	// Apenas GET endpoints (create e update são por mensageria)
	router.GET("/", kitchenOrderHandler.FindAll)
	router.GET("/:id", kitchenOrderHandler.FindByID)
}
