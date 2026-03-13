package inventory

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/inventory"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Inventory struct {
	InventoryService *service.Inventory
}

func New(inventoryService *service.Inventory) *Inventory {
	return &Inventory{InventoryService: inventoryService}
}

func (ctrl *Inventory) Get(ctx *gin.Context) {
	playerID := ctx.Query("playerId")
	if playerID == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "inventory.Get.BadRequest", errors.New("playerId query parameter is required")))
		return
	}

	inv, err := ctrl.InventoryService.GetByPlayerID(playerID)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "inventory.Get.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Inventory", TotalCount: len(inv.Items), Count: len(inv.Items)},
		Data: inv,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}
