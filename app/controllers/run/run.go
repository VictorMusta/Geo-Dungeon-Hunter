package run

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/run"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Run struct {
	RunService *service.Run
}

func New(runService *service.Run) *Run {
	return &Run{RunService: runService}
}

func (ctrl *Run) Create(ctx *gin.Context) {
	var in models.Run

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "run.Create.BadRequest", err))
		return
	}

	// Securely set player from JWT
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.KnownError(http.StatusUnauthorized, "run.Create.Unauthorized", errors.New("authentication required")))
		return
	}
	in.PlayerID = playerID.(string)

	r, err := ctrl.RunService.Create(&in)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "dungeon not found", "player not found":
			status = http.StatusNotFound
		case "dungeon is not published":
			status = http.StatusConflict
		case "player already has an active run for this dungeon":
			status = http.StatusConflict
		}
		common.SendResponse(ctx, status, models.KnownError(status, "run.Create.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Run", TotalCount: 1, Count: 1},
		Data: r,
	}
	common.SendResponse(ctx, http.StatusCreated, response)
}

func (ctrl *Run) Get(ctx *gin.Context) {
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.KnownError(http.StatusUnauthorized, "run.Get.Unauthorized", errors.New("authentication required")))
		return
	}
	uid := playerID.(string)

	runs, err := ctrl.RunService.GetByPlayerID(uid)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "run.Get.Error", err))
		return
	}

	if len(runs) == 0 {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "run.Get.NotFound", errors.New("no runs found")))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Run", TotalCount: len(runs), Count: len(runs)},
		Data: runs,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Run) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	r, err := ctrl.RunService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "run.Get.NotFound", errors.New("run not found")))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Run", TotalCount: 1, Count: 1},
		Data: r,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Run) Abandon(ctx *gin.Context) {
	id := ctx.Param("id")

	// Ownership check
	if ok, err := ctrl.checkRunOwnership(ctx, id); !ok {
		common.SendResponse(ctx, http.StatusForbidden, models.KnownError(http.StatusForbidden, "run.Abandon.Forbidden", err))
		return
	}

	if err := ctrl.RunService.Abandon(id); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "run not found" {
			status = http.StatusNotFound
		}
		common.SendResponse(ctx, status, models.KnownError(status, "run.Abandon.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "run.Abandon.OK", "run abandoned"))
}

func (ctrl *Run) AttemptBoss(ctx *gin.Context) {
	runID := ctx.Param("id")

	// Ownership check
	if ok, err := ctrl.checkRunOwnership(ctx, runID); !ok {
		common.SendResponse(ctx, http.StatusForbidden, models.KnownError(http.StatusForbidden, "attempt.Forbidden", err))
		return
	}

	stepID := ctx.Param("stepId")

	var body struct {
		Lat float64 `json:"lat" binding:"required"`
		Lon float64 `json:"lon" binding:"required"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "attempt.BadRequest", err))
		return
	}

	result, err := ctrl.RunService.AttemptBoss(runID, stepID, body.Lat, body.Lon)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "NOT_IN_RANGE":
			status = http.StatusConflict
		case "WRONG_STEP_ORDER":
			status = http.StatusConflict
		case "run not found", "boss step not found":
			status = http.StatusNotFound
		case "run is not active":
			status = http.StatusConflict
		}
		common.SendResponse(ctx, status, models.KnownError(status, "attempt.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "AttemptResult", TotalCount: 1, Count: 1},
		Data: result,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Run) checkRunOwnership(ctx *gin.Context, runID string) (bool, error) {
	playerID, exists := ctx.Get("playerID")
	if !exists {
		return false, errors.New("unauthorized: missing playerID in context")
	}

	r, err := ctrl.RunService.GetByID(runID)
	if err != nil {
		return false, err
	}

	if r.PlayerID != playerID.(string) {
		return false, errors.New("forbidden: you do not own this run")
	}
	return true, nil
}
