package auction

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/auction"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auction struct {
	AuctionService *service.Auction
}

func New(auctionService *service.Auction) *Auction {
	return &Auction{AuctionService: auctionService}
}

func (ctrl *Auction) CreateListing(ctx *gin.Context) {
	var in models.Listing

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "listing.Create.BadRequest", err))
		return
	}

	// Securely set seller from JWT
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.KnownError(http.StatusUnauthorized, "listing.Create.Unauthorized", errors.New("authentication required")))
		return
	}
	in.SellerID = playerID.(string)

	listing, err := ctrl.AuctionService.CreateListing(&in)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "item not found", "seller does not own this item":
			status = http.StatusNotFound
		case "INSUFFICIENT_ITEMS":
			status = http.StatusConflict
		case "item is not tradable":
			status = http.StatusConflict
		}
		common.SendResponse(ctx, status, models.KnownError(status, "listing.Create.Error", err))
		return
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Listing", TotalCount: 1, Count: 1},
		Data: listing,
	}
	common.SendResponse(ctx, http.StatusCreated, response)
}

func (ctrl *Auction) GetListings(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)

	listings, err := ctrl.AuctionService.GetListings(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "listing.Get.Error", err))
		return
	}

	if len(listings) == 0 {
		listings = []models.Listing{}
	}

	response := &models.WSResponse{
		Meta: models.MetaResponse{ObjectName: "Listing", TotalCount: len(listings), Count: len(listings)},
		Data: listings,
	}
	common.SendResponse(ctx, http.StatusOK, response)
}

func (ctrl *Auction) Buy(ctx *gin.Context) {
	listingID := ctx.Param("id")

	var body struct {
		BuyerID string `json:"buyerId"`
		Qty     int    `json:"qty" binding:"min=0"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "listing.Buy.BadRequest", err))
		return
	}

	// Default Qty to 1 if not provided or 0
	if body.Qty <= 0 {
		body.Qty = 1
	}

	// Securely set buyer from JWT
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.KnownError(http.StatusUnauthorized, "listing.Buy.Unauthorized", errors.New("authentication required")))
		return
	}
	body.BuyerID = playerID.(string)

	if err := ctrl.AuctionService.Buy(listingID, body.BuyerID, body.Qty); err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "listing not found", "buyer not found":
			status = http.StatusNotFound
		case "Offre plus active", "LISTING_NOT_ACTIVE":
			status = http.StatusConflict
		case "Or insuffisant pour cet achat", "INSUFFICIENT_GOLD":
			status = http.StatusConflict
		case "Tu ne peux pas acheter ta propre offre":
			status = http.StatusConflict
		case "Quantité demandée supérieure au stock":
			status = http.StatusConflict
		}
		common.SendResponse(ctx, status, models.KnownError(status, "listing.Buy.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "listing.Buy.OK", "purchase successful"))
}

func (ctrl *Auction) Cancel(ctx *gin.Context) {
	listingID := ctx.Param("id")

	var body struct {
		SellerID string `json:"sellerId"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "listing.Cancel.BadRequest", err))
		return
	}

	// Securely set seller from JWT
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.KnownError(http.StatusUnauthorized, "listing.Cancel.Unauthorized", errors.New("authentication required")))
		return
	}
	body.SellerID = playerID.(string)

	if err := ctrl.AuctionService.Cancel(listingID, body.SellerID); err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "listing not found":
			status = http.StatusNotFound
		case "LISTING_NOT_ACTIVE":
			status = http.StatusConflict
		case "only the seller can cancel this listing":
			status = http.StatusForbidden
		}
		common.SendResponse(ctx, status, models.KnownError(status, "listing.Cancel.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "listing.Cancel.OK", "listing cancelled"))
}
