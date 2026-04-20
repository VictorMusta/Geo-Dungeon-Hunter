package item

import (
	controller "dungeons/app/controllers/item"
	repo "dungeons/app/repositories/mongodb"
	"dungeons/app/server"
	service "dungeons/app/services/item"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()

	itemRepo := repo.NewItemRepository(srv.Database)
	itemService := service.New(itemRepo)
	itemController := controller.New(itemService)

	v1 := g.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", itemController.Create)
			items.GET("", itemController.Get)
			items.GET("/:id", itemController.GetByID)
			items.PUT("/:id", itemController.Update)
			items.DELETE("/:id", itemController.Delete)
		}
	}
}
