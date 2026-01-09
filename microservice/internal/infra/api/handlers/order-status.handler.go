package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/factories"
)

type OrderStatusHandler struct {
	controller controllers.OrderStatusController
}

func NewOrderStatusHandler() *OrderStatusHandler {
	orderStatusDataSource := factories.NewOrderStatusDataSource()
	controller := controllers.NewOrderStatusController(orderStatusDataSource)

	return &OrderStatusHandler{
		controller: *controller,
	}
}

// FindAll godoc
// @Summary Get all order status
// @Description Get all available order status
// @Tags order-status
// @Accept json
// @Produce json
// @Success 200 {array} dtos.OrderStatusResponseDTO
// @Failure 500 {object} map[string]interface{}
// @Router /v1/kitchen-orders/status [get]
func (h *OrderStatusHandler) FindAll(c *gin.Context) {
	orderStatus, err := h.controller.FindAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orderStatus)
}