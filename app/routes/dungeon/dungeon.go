package dungeon

import (
	controller "dungeons/app/controllers/dungeon"
	bsService "dungeons/app/services/bossstep"
	dService "dungeons/app/services/dungeon"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	dungeonService := dService.New()
	bossStepService := bsService.New()
	dungeonController := controller.New(dungeonService, bossStepService)

	v1 := g.Group("/v1")
	{
		// Game Master endpoints
		mj := v1.Group("/mj/dungeons")
		{
			mj.POST("", dungeonController.Create)
			mj.PUT("/:id", dungeonController.Update)
			mj.POST("/:id/publish", dungeonController.Publish)
			mj.POST("/:id/steps", dungeonController.CreateStep)
			mj.PUT("/:id/steps/:stepId", dungeonController.UpdateStep)
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
