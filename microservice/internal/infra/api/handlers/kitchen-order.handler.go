package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/factories"
	"tech_challenge/internal/infra/api/schemas"
	shared_factories "tech_challenge/internal/shared/factories"
)

type KitchenOrderHandler struct {
	kitchenOrderController controllers.KitchenOrderController
}

func NewKitchenOrderHandler() *KitchenOrderHandler {
	kitchenOrderDataSource := factories.NewKitchenOrderDataSource()
	orderStatusDataSource := factories.NewOrderStatusDataSource()
	
	messageBroker, err := shared_factories.NewMessageBroker(context.Background())
	if err != nil {
		log.Printf("Warning: Failed to create message broker: %v", err)
		// Continue without message broker for now
		messageBroker = nil
	}

	kitchenOrderController := controllers.NewKitchenOrderController(kitchenOrderDataSource, orderStatusDataSource, messageBroker)

	return &KitchenOrderHandler{
		kitchenOrderController: *kitchenOrderController,
	}
}

func (h *KitchenOrderHandler) toKitchenOrderResponseSchema(kitchenOrder dtos.KitchenOrderResponseDTO) schemas.KitchenOrderResponseSchema {
	return schemas.KitchenOrderResponseSchema{
		ID:        kitchenOrder.ID,
		OrderID:   kitchenOrder.OrderID,
		Status:    kitchenOrder.Status.Name,
		Slug:      kitchenOrder.Slug,
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	}
}

// @Summary List all kitchenOrders
// @Tags KitchenOrders
// @Produce json
// @Success 200 {array} schemas.KitchenOrderResponseSchema
// @Failure 500 {object} schemas.ErrorMessageSchema
// @Router /kitchen-orders/ [get]
func (h *KitchenOrderHandler) FindAll(ctx *gin.Context) {
	var filter dtos.KitchenOrderFilter

	if createdAtFromStr := ctx.Query("created_at_from"); createdAtFromStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtFromStr); err == nil {
			filter.CreatedAtFrom = &t
		}
	}

	if createdAtToStr := ctx.Query("created_at_to"); createdAtToStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtToStr); err == nil {
			filter.CreatedAtTo = &t
		}
	}

	if statusIDStr := ctx.Query("status_id"); statusIDStr != "" {
		if id, err := strconv.ParseUint(statusIDStr, 10, 64); err == nil {
			u := uint(id)
			filter.StatusID = &u
		}
	}

	kitchenOrders, err := h.kitchenOrderController.FindAll(filter)

	if err != nil {
		if ctxErr := ctx.Error(err); ctxErr != nil {
			log.Printf("Error setting context error: %v", ctxErr)
		}
		return
	}

	kitchenOrderResponses := make([]schemas.KitchenOrderResponseSchema, len(kitchenOrders))
	for i, kitchenOrder := range kitchenOrders {
		kitchenOrderResponses[i] = h.toKitchenOrderResponseSchema(kitchenOrder)
	}

	ctx.JSON(http.StatusOK, kitchenOrderResponses)
}

// @Summary Get a kitchenOrder by ID
// @Tags KitchenOrders
// @Produce json
// @Param id path string true "KitchenOrder ID"
// @Success 200 {object} schemas.KitchenOrderResponseSchema
// @Failure 400 {object} schemas.InvalidKitchenOrderDataErrorSchema
// @Failure 404 {object} schemas.KitchenOrderNotFoundErrorSchema
// @Router /kitchen-orders/{id} [get]
func (h *KitchenOrderHandler) FindByID(ctx *gin.Context) {
	kitchenOrderID := ctx.Param("id")

	kitchenOrder, err := h.kitchenOrderController.FindByID(kitchenOrderID)

	if err != nil {
		if ctxErr := ctx.Error(err); ctxErr != nil {
			log.Printf("Error setting context error: %v", ctxErr)
		}
		return
	}

	ctx.JSON(http.StatusOK, h.toKitchenOrderResponseSchema(kitchenOrder))
}

// @Summary Update a kitchenOrder status
// @Tags KitchenOrders
// @Accept json
// @Produce json
// @Param id path string true "KitchenOrder ID"
// @Param request body schemas.UpdateKitchenOrderRequestSchema true "Update request"
// @Success 200 {object} schemas.KitchenOrderResponseSchema
// @Failure 400 {object} schemas.InvalidKitchenOrderDataErrorSchema
// @Failure 404 {object} schemas.KitchenOrderNotFoundErrorSchema
// @Router /kitchen-orders/{id} [put]
func (h *KitchenOrderHandler) Update(ctx *gin.Context) {
	kitchenOrderID := ctx.Param("id")

	var request schemas.UpdateKitchenOrderRequestSchema
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateDTO := dtos.UpdateKitchenOrderDTO{
		ID:       kitchenOrderID,
		StatusID: request.StatusID,
	}

	kitchenOrder, err := h.kitchenOrderController.Update(updateDTO)

	if err != nil {
		if ctxErr := ctx.Error(err); ctxErr != nil {
			log.Printf("Error setting context error: %v", ctxErr)
		}
		return
	}

	ctx.JSON(http.StatusOK, h.toKitchenOrderResponseSchema(kitchenOrder))
}
