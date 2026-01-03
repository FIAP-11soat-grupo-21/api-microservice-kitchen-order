package routes

import (
	"tech_challenge/internal/shared/factories"
	"tech_challenge/internal/shared/infra/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.RouterGroup) {
	fileProvider := factories.NewFileProvider()
	fileHandler := handlers.NewFileHandler(fileProvider)

	router.GET("/:fileName", func(c *gin.Context) {
		fileName := c.Param("fileName")

		fileUrl, err := fileHandler.FindFile(fileName)

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve file"})
			return
		}

		c.JSON(200, gin.H{"fileUrl": fileUrl})
	})
}
