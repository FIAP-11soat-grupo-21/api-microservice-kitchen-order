package http_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tech_challenge/internal/domain/exceptions"
)

func HandleDomainErrors(err error, ctx *gin.Context) bool {
	switch e := err.(type) {
	case *exceptions.InvalidKitchenOrderDataException:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return true

	case *exceptions.KitchenOrderNotFoundException:
		ctx.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		return true
	}

	return false
}
