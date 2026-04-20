package dungeon

import (
	controller "dungeons/app/controllers/dungeon"
	repo "dungeons/app/repositories/mongodb"
	"dungeons/app/server"
	bsService "dungeons/app/services/bossstep"
	dService "dungeons/app/services/dungeon"

	"dungeons/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	srv := server.GetServer()

	// Repositories
	dRepo := repo.NewDungeonRepository(srv.Database)
	bsRepo := repo.NewBossStepRepository(srv.Database)

	// Services
	bossStepService := bsService.New(bsRepo)
	dungeonService := dService.New(dRepo, bossStepService)

	dungeonController := controller.New(dungeonService, bossStepService)

	v1 := g.Group("/v1")
	{
		// Game Master endpoints
		mj := v1.Group("/mj/dungeons")
		mj.Use(middleware.AuthMiddleware()) // Require login for GM actions
		{
			mj.GET("", dungeonController.GetByMJ)
			mj.POST("", dungeonController.Create)
			mj.PUT("/:id", dungeonController.Update)
			mj.PUT("/:id/full", dungeonController.UpdateFull)
			mj.POST("/:id/publish", dungeonController.Publish)
			mj.POST("/:id/steps", dungeonController.CreateStep)
			mj.PUT("/:id/steps/:stepId", dungeonController.UpdateStep)
			mj.DELETE("/:id/steps/:stepId", dungeonController.DeleteStep)
			mj.PUT("/:id/steps/reorder", dungeonController.ReorderSteps)
		}

		// Player-facing dungeon endpoints
		dungeons := v1.Group("/dungeons")
		{
			dungeons.GET("", dungeonController.GetPublished)
			dungeons.GET("/:id", dungeonController.GetByID)
		}
	}
}
