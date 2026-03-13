package item

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/item"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ItemService *service.Item
}

func New(itemService *service.Item) *Item {
	return &Item{ItemService: itemService}
}

func (ctrl *Item) Create(ctx *gin.Context) {
	var in models.ItemDef

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "item.Create.BadRequest", err))
		return
	}

	item, err := ctrl.ItemService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "item.Create.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Item", TotalCount: 1, Count: 1},
		Data: item,
	}
	common.SendResponse(ctx, http.StatusCreated, response)
}

func (ctrl *Item) Get(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)

	items, err := ctrl.ItemService.Get(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "item.Search.Error", err))
		return
	}

	if len(items) == 0 {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "item.Search.NotFound", errors.New("no items found")))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Item", TotalCount: len(items), Count: len(items)},
		Data: items,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Item) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	item, err := ctrl.ItemService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "item.Get.NotFound", errors.New("item not found")))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Item", TotalCount: 1, Count: 1},
		Data: item,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}
