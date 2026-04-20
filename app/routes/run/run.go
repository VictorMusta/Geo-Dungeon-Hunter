package run

import (
	controller "dungeons/app/controllers/run"
	repo "dungeons/app/repositories/mongodb"
	"dungeons/app/server"
	service "dungeons/app/services/run"

	"dungeons/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()

	// Repositories
	runRepo := repo.NewRunRepository(srv.Database)
	dungeonRepo := repo.NewDungeonRepository(srv.Database)
	bsRepo := repo.NewBossStepRepository(srv.Database)
	playerRepo := repo.NewPlayerRepository(srv.Database)

	// Service
	runService := service.New(runRepo, dungeonRepo, bsRepo, playerRepo)
	runController := controller.New(runService)

	v1 := g.Group("/v1")
	{
		runs := v1.Group("/runs")
		runs.Use(middleware.AuthMiddleware()) // Authentication required for all run state changes
		{
			runs.POST("", runController.Create)
			runs.GET("", runController.Get)
			runs.GET("/:id", runController.GetByID)
			runs.POST("/:id/abandon", runController.Abandon)
			runs.POST("/:id/steps/:stepId/attempt", runController.AttemptBoss)
		}
	}
}
