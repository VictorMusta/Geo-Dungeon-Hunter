package auction

import (
	controller "dungeons/app/controllers/auction"
	service "dungeons/app/services/auction"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	auctionService := service.New()
	auctionController := controller.New(auctionService)

	v1 := g.Group("/v1")
	{
		listings := v1.Group("/auction/listings")
		{
			listings.POST("", auctionController.CreateListing)
			listings.GET("", auctionController.GetListings)
			listings.POST("/:id/buy", auctionController.Buy)
			listings.POST("/:id/cancel", auctionController.Cancel)
		}
	}
}
