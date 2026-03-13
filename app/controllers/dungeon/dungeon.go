package dungeon

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	bsService "dungeons/app/services/bossstep"
	dService "dungeons/app/services/dungeon"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Dungeon struct {
	DungeonService  *dService.Dungeon
	BossStepService *bsService.BossStep
}

func New(ds *dService.Dungeon, bs *bsService.BossStep) *Dungeon {
	return &Dungeon{DungeonService: ds, BossStepService: bs}
}

func (ctrl *Dungeon) Create(ctx *gin.Context) {
	var in models.Dungeon

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Create.BadRequest", err))
		return
	}

	d, err := ctrl.DungeonService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Create.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Dungeon", TotalCount: 1, Count: 1},
		Data: d,
	}
	common.SendResponse(ctx, http.StatusCreated, response)
}

func (ctrl *Dungeon) Update(ctx *gin.Context) {
	var in models.Dungeon

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Update.BadRequest", err))
		return
	}

	id := ctx.Param("id")
	if err := ctrl.DungeonService.Update(id, &in); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "dungeon not found" {
			status = http.StatusNotFound
		} else if err.Error() == "only draft dungeons can be modified" {
			status = http.StatusConflict
		}
		common.SendResponse(ctx, status, models.KnownError(status, "dungeon.Update.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.Update.OK", "dungeon updated"))
}

func (ctrl *Dungeon) Publish(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.DungeonService.Publish(id); err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "dungeon not found":
			status = http.StatusNotFound
		case "only draft dungeons can be published":
			status = http.StatusConflict
		case "dungeon must have at least one boss step to publish":
			status = http.StatusUnprocessableEntity
		}
		common.SendResponse(ctx, status, models.KnownError(status, "dungeon.Publish.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.Publish.OK", "dungeon published"))
}

func (ctrl *Dungeon) GetPublished(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)

	dungeons, err := ctrl.DungeonService.GetPublished(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "dungeon.Search.Error", err))
		return
	}

	if len(dungeons) == 0 {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "dungeon.Search.NotFound", errors.New("no dungeons found")))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Dungeon", TotalCount: len(dungeons), Count: len(dungeons)},
		Data: dungeons,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Dungeon) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	d, err := ctrl.DungeonService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "dungeon.Get.NotFound", errors.New("dungeon not found")))
		return
	}

	steps, _ := ctrl.BossStepService.GetByDungeonID(id)

	type DungeonDetail struct {
		models.Dungeon
		BossSteps []models.BossStep `json:"bossSteps"`
	}

	detail := DungeonDetail{Dungeon: d, BossSteps: steps}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Dungeon", TotalCount: 1, Count: 1},
		Data: detail,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Dungeon) CreateStep(ctx *gin.Context) {
	var in models.BossStep
	dungeonID := ctx.Param("id")

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "step.Create.BadRequest", err))
		return
	}

	in.DungeonID = dungeonID

	step, err := ctrl.BossStepService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "step.Create.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "BossStep", TotalCount: 1, Count: 1},
		Data: step,
	}
	common.SendResponse(ctx, http.StatusCreated, response)
}

func (ctrl *Dungeon) UpdateStep(ctx *gin.Context) {
	var in models.BossStep
	dungeonID := ctx.Param("id")
	stepID := ctx.Param("stepId")

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "step.Update.BadRequest", err))
		return
	}

	if err := ctrl.BossStepService.Update(dungeonID, stepID, &in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "step.Update.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "step.Update.OK", "step updated"))
}

func (ctrl *Dungeon) ReorderSteps(ctx *gin.Context) {
	dungeonID := ctx.Param("id")

	var body struct {
		Order []string `json:"order" binding:"required"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "step.Reorder.BadRequest", err))
		return
	}

	if err := ctrl.BossStepService.Reorder(dungeonID, body.Order); err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "step.Reorder.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "step.Reorder.OK", "steps reordered"))
}
