package auction

import (
	controller "dungeons/app/controllers/auction"
	repo "dungeons/app/repositories/mongodb"
	service "dungeons/app/services/auction"
	"dungeons/app/server"

	"dungeons/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()
	
	// Repositories
	auctionRepo := repo.NewAuctionRepository(srv.Database)
	itemRepo := repo.NewItemRepository(srv.Database)
	inventoryRepo := repo.NewInventoryRepository(srv.Database)
	playerRepo := repo.NewPlayerRepository(srv.Database)

	// Service
	auctionService := service.New(auctionRepo, itemRepo, inventoryRepo, playerRepo)
	auctionController := controller.New(auctionService)

	v1 := g.Group("/v1")
	{
		listings := v1.Group("/auction/listings")
		{
			listings.GET("", auctionController.GetListings)
			
			// Protected auction actions
			protected := listings.Group("")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("", auctionController.CreateListing)
				protected.POST("/:id/buy", auctionController.Buy)
				protected.POST("/:id/cancel", auctionController.Cancel)
			}
		}
	}
}
