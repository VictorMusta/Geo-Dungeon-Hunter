package leaderboard

import (
	controller "dungeons/app/controllers/leaderboard"
	repo "dungeons/app/repositories/mongodb"
	"dungeons/app/server"
	service "dungeons/app/services/leaderboard"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()

	leaderboardRepo := repo.NewLeaderboardRepository(srv.Database)
	leaderboardService := service.New(leaderboardRepo)
	leaderboardController := controller.New(leaderboardService)

	v1 := g.Group("/v1")
	{
		v1.GET("/leaderboard", leaderboardController.Get)
	}
}
