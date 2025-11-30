package handlers

import (
	"net/http"
	"strconv"
	"tech_challenge/internal/application/controllers"
	"tech_challenge/internal/application/dtos"
	"tech_challenge/internal/factories"
	"tech_challenge/internal/infra/api/schemas"
	"time"

	"github.com/gin-gonic/gin"
)

type KitchenOrderHandler struct {
	kitchenOrderController controllers.KitchenOrderController
}

func NewKitchenOrderHandler() *KitchenOrderHandler {
	kitchenOrderDataSource := factories.NewKitchenOrderDataSource()
	orderStatusDataSource := factories.NewOrderStatusDataSource()

	kitchenOrderController := controllers.NewKitchenOrderController(kitchenOrderDataSource, orderStatusDataSource)

	return &KitchenOrderHandler{
		kitchenOrderController: *kitchenOrderController,
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
		ctx.Error(err)
		return
	}

	kitchenOrderResponses := make([]schemas.KitchenOrderResponseSchema, len(kitchenOrders))

	for i, kitchenOrder := range kitchenOrders {

		kitchenOrderResponses[i] = schemas.KitchenOrderResponseSchema{
			OrderID:   kitchenOrder.OrderID,
			Status:    kitchenOrder.Status.Name,
			Slug:      kitchenOrder.Slug,
			CreatedAt: kitchenOrder.CreatedAt,
			UpdatedAt: kitchenOrder.UpdatedAt,
		}
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
		ctx.Error(err)
		return
	}

	kitchenOrderResponse := schemas.KitchenOrderResponseSchema{
		OrderID:   kitchenOrder.OrderID,
		Status:    kitchenOrder.Status.Name,
		Slug:      kitchenOrder.Slug,
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, kitchenOrderResponse)
}

// @Summary Update a kitchenOrder by ID
// @Tags KitchenOrders
// @Accept json
// @Produce json
// @Param id path string true "KitchenOrder ID"
// @Param kitchenOrder body schemas.UpdateKitchenOrderSchema true "Updated kitchenOrder data"
// @Success 200 {object} schemas.KitchenOrderResponseSchema
// @Failure 400 {object} schemas.InvalidKitchenOrderDataErrorSchema
// @Failure 404 {object} schemas.KitchenOrderNotFoundErrorSchema
// @Failure 500 {object} schemas.ErrorMessageSchema
// @Router /kitchen-orders/{id} [put]
func (h *KitchenOrderHandler) Update(ctx *gin.Context) {
	kitchenOrderID := ctx.Param("id")

	var kitchenOrderRequestBody schemas.UpdateKitchenOrderSchema

	if err := ctx.ShouldBindJSON(&kitchenOrderRequestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	kitchenOrder, err := h.kitchenOrderController.Update(dtos.UpdateKitchenOrderDTO{
		ID:       kitchenOrderID,
		StatusID: kitchenOrderRequestBody.StatusID,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	kitchenOrderResponse := schemas.KitchenOrderResponseSchema{
		OrderID:   kitchenOrder.OrderID,
		Status:    kitchenOrder.Status.Name,
		Slug:      kitchenOrder.Slug,
		CreatedAt: kitchenOrder.CreatedAt,
		UpdatedAt: kitchenOrder.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, kitchenOrderResponse)
}
