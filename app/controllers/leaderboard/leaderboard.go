package leaderboard

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	"dungeons/app/repositories"
	service "dungeons/app/services/leaderboard"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Leaderboard struct {
	LeaderboardService *service.Leaderboard
}

func New(leaderboardService *service.Leaderboard) *Leaderboard {
	return &Leaderboard{LeaderboardService: leaderboardService}
}

func (ctrl *Leaderboard) Get(ctx *gin.Context) {
	leaderboardType := ctx.Query("type")
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	var (
		entries []repositories.LeaderboardEntry
		err     error
	)

	switch leaderboardType {
	case "completions":
		entries, err = ctrl.LeaderboardService.GetByCompletions(limit)
	case "gold":
		entries, err = ctrl.LeaderboardService.GetByGold(limit)
	case "speed":
		dungeonID := ctx.Query("dungeonId")
		if dungeonID == "" {
			common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "leaderboard.Get.BadRequest", errors.New("dungeonId query parameter is required for speed leaderboard")))
			return
		}
		entries, err = ctrl.LeaderboardService.GetBySpeed(dungeonID, limit)
	default:
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "leaderboard.Get.BadRequest", errors.New("type must be one of: completions, gold, speed")))
		return
	}

	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "leaderboard.Get.Error", err))
		return
	}

	if entries == nil {
		entries = []repositories.LeaderboardEntry{}
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Leaderboard", TotalCount: len(entries), Count: len(entries)},
		Data: entries,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}
