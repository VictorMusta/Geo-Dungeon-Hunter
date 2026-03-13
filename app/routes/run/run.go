package run

import (
	controller "dungeons/app/controllers/run"
	service "dungeons/app/services/run"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	runService := service.New()
	runController := controller.New(runService)

	v1 := g.Group("/v1")
	{
		runs := v1.Group("/runs")
		{
			runs.POST("", runController.Create)
			runs.GET("", runController.Get)
			runs.GET("/:id", runController.GetByID)
			runs.POST("/:id/abandon", runController.Abandon)
			runs.POST("/:id/steps/:stepId/attempt", runController.AttemptBoss)
		}
	}
}
