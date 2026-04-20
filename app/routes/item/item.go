package item

import (
	controller "dungeons/app/controllers/item"
	repo "dungeons/app/repositories/mongodb"
	service "dungeons/app/services/item"
	"dungeons/app/server"

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
		}
	}
}
