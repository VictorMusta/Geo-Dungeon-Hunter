package auction

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Auction struct {
	repo          repositories.AuctionRepository
	itemRepo      repositories.ItemRepository
	inventoryRepo repositories.InventoryRepository
	playerRepo    repositories.PlayerRepository
	validate      *validator.Validate
}

func New(
	repo repositories.AuctionRepository,
	itemRepo repositories.ItemRepository,
	inventoryRepo repositories.InventoryRepository,
	playerRepo repositories.PlayerRepository,
) *Auction {
	return &Auction{
		repo:          repo,
		itemRepo:      itemRepo,
		inventoryRepo: inventoryRepo,
		playerRepo:    playerRepo,
		validate:      validator.New(),
	}
}

func (s *Auction) CreateListing(in *models.Listing) (*models.Listing, error) {
	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	// Verify item exists and is tradable
	item, err := s.itemRepo.GetByID(in.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}
	if !item.Tradable {
		return nil, errors.New("item is not tradable")
	}

	// Verify seller has enough items
	inv, err := s.inventoryRepo.GetByItem(in.SellerID, in.ItemID)
	if err != nil {
		return nil, errors.New("seller does not own this item")
	}
	if inv.Qty < int64(in.Qty) {
		return nil, errors.New("INSUFFICIENT_ITEMS")
	}

	// Deduct items from seller inventory
	err = s.inventoryRepo.Update(in.SellerID, in.ItemID, -int64(in.Qty))
	if err != nil {
		return nil, err
	}

	listing := models.Listing{
		CustomID:     functions.NewUUID(),
		SellerID:     in.SellerID,
		ItemID:       in.ItemID,
		Qty:          in.Qty,
		PricePerUnit: in.PricePerUnit,
		Status:       "active",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(72 * time.Hour),
	}

	if err := s.repo.Create(&listing); err != nil {
		return nil, err
	}

	return &listing, nil
}

func (s *Auction) GetListings(queryParams models.QueryParams) ([]models.Listing, error) {
	return s.repo.GetByStatus("active")
}

func (s *Auction) Buy(listingID, buyerID string, qty int) error {
	// Get listing
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return errors.New("listing not found")
	}
	if listing.Status != "active" {
		return errors.New("LISTING_NOT_ACTIVE")
	}
	if qty > listing.Qty {
		return errors.New("requested quantity exceeds listing quantity")
	}

	totalPrice := int64(qty) * listing.PricePerUnit

	// Verify buyer has enough gold
	buyer, err := s.playerRepo.GetByID(buyerID)
	if err != nil {
		return errors.New("buyer not found")
	}
	if buyer.Gold < totalPrice {
		return errors.New("INSUFFICIENT_GOLD")
	}

	if buyerID == listing.SellerID {
		return errors.New("cannot buy your own listing")
	}

	trade := &models.Trade{
		CustomID:   functions.NewUUID(),
		BuyerID:    buyerID,
		SellerID:   listing.SellerID,
		ListingID:  listingID,
		Qty:        qty,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
	}

	// Atomic transaction in repository
	return s.repo.BuyListing(context.TODO(), buyerID, listingID, totalPrice, listing.SellerID, listing.ItemID, qty, trade)
}

func (s *Auction) Cancel(listingID, sellerID string) error {
	listing, err := s.repo.GetByID(listingID)
	if err != nil {
		return errors.New("listing not found")
	}
	if listing.Status != "active" {
		return errors.New("LISTING_NOT_ACTIVE")
	}
	if listing.SellerID != sellerID {
		return errors.New("only the seller can cancel this listing")
	}

	// Return items to seller inventory
	err = s.inventoryRepo.Update(sellerID, listing.ItemID, int64(listing.Qty))
	if err != nil {
		return err
	}

	// Mark listing as cancelled
	listing.Status = "cancelled"
	return s.repo.Update(listingID, &listing)
}
