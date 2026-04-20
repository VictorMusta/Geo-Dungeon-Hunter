package player

import (
	controller "dungeons/app/controllers/player"
	repo "dungeons/app/repositories/mongodb"
	service "dungeons/app/services/player"
	"dungeons/app/server"

	"dungeons/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()
	playerRepo := repo.NewPlayerRepository(srv.Database)
	servicesPlayer := service.New(playerRepo)
	playerController := controller.New(servicesPlayer)

	v1 := g.Group("/v1")
	{
		players := v1.Group("/players")
		{
			players.POST("", playerController.Create)
			players.POST("/login", playerController.Login)
			
			// Protected routes
			protected := players.Group("")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.GET("", playerController.Get)
				protected.GET("/:id", playerController.GetByID)
				protected.POST("/:id", playerController.Update)
			}
			
			players.GET("/IDS/:ids", playerController.GetByIDs)
		}
	}
}
