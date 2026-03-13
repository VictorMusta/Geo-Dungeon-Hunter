package leaderboard

import (
	controller "dungeons/app/controllers/leaderboard"
	service "dungeons/app/services/leaderboard"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	leaderboardService := service.New()
	leaderboardController := controller.New(leaderboardService)

	v1 := g.Group("/v1")
	{
		v1.GET("/leaderboard", leaderboardController.Get)
	}
}
