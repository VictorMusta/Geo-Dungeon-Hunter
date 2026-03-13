package item

import (
	controller "dungeons/app/controllers/item"
	service "dungeons/app/services/item"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	itemService := service.New()
	itemController := controller.New(itemService)

	v1 := g.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", itemController.Create)
			items.GET("", itemController.Get)
			items.GET("/:id", itemController.GetByID)
		}
	}
}
