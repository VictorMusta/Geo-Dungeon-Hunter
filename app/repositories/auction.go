package repositories

import (
	"context"
	"dungeons/app/models"
)

type AuctionRepository interface {
	GetByStatus(status string) ([]models.Listing, error)
	GetByID(id string) (models.Listing, error)
	Create(listing *models.Listing) error
	Update(id string, listing *models.Listing) error

	// Transactional operations
	BuyListing(ctx context.Context, buyerID, listingID string, priceTotal int64, sellerID string, itemID string, qty int) error
}
