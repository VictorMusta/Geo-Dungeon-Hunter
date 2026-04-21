package inventory

import (
	controller "dungeons/app/controllers/inventory"
	repo "dungeons/app/repositories/mongodb"
	"dungeons/app/server"
	service "dungeons/app/services/inventory"

	"dungeons/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()

	invRepo := repo.NewInventoryRepository(srv.Database)
	itemRepo := repo.NewItemRepository(srv.Database)
	inventoryService := service.New(invRepo, itemRepo)
	inventoryController := controller.New(inventoryService)

	v1 := g.Group("/v1")
	v1.Use(middleware.AuthMiddleware()) // Protected inventory access
	{
		v1.GET("/inventory", inventoryController.Get)
	}
}
