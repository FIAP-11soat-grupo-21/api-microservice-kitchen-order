package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	kitchen_order_http_errors "tech_challenge/internal/infra/api/http_errors"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			errorHasBinHandled := kitchen_order_http_errors.HandleDomainErrors(err, ctx)

			if !errorHasBinHandled {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}

			ctx.Abort()
		}
	}
}
