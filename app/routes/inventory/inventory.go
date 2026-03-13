package inventory

import (
	controller "dungeons/app/controllers/inventory"
	service "dungeons/app/services/inventory"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	inventoryService := service.New()
	inventoryController := controller.New(inventoryService)

	v1 := g.Group("/v1")
	{
		v1.GET("/inventory", inventoryController.Get)
	}
}
