package auction

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Auction struct {
	validate *validator.Validate
}

func New() *Auction {
	return &Auction{validate: validator.New()}
}

func (s *Auction) CreateListing(in *models.Listing) (*models.Listing, error) {
	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	srv := server.GetServer()

	// Verify item exists and is tradable
	var item models.ItemDef
	itemCollection := srv.Database.Collection(item.Collection())
	err := itemCollection.FindOne(context.TODO(), bson.M{"customID": in.ItemID}).Decode(&item)
	if err != nil {
		return nil, errors.New("item not found")
	}
	if !item.Tradable {
		return nil, errors.New("item is not tradable")
	}

	// Verify seller has enough items
	invCollection := srv.Database.Collection("inventory")
	var inv models.InventoryEntry
	err = invCollection.FindOne(context.TODO(), bson.M{
		"playerId": in.SellerID,
		"itemId":   in.ItemID,
	}).Decode(&inv)
	if err != nil {
		return nil, errors.New("seller does not own this item")
	}
	if inv.Qty < int64(in.Qty) {
		return nil, errors.New("INSUFFICIENT_ITEMS")
	}

	// Deduct items from seller inventory
	newQty := inv.Qty - int64(in.Qty)
	if newQty <= 0 {
		_, err = invCollection.DeleteOne(context.TODO(), bson.M{"playerId": in.SellerID, "itemId": in.ItemID})
	} else {
		_, err = invCollection.UpdateOne(context.TODO(), bson.M{"playerId": in.SellerID, "itemId": in.ItemID},
			bson.M{"$inc": bson.M{"qty": -int64(in.Qty)}, "$set": bson.M{"updatedAt": time.Now()}})
	}
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

	listingCollection := srv.Database.Collection(listing.Collection())
	if _, err := listingCollection.InsertOne(context.TODO(), listing); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &listing, nil
}

func (s *Auction) GetListings(queryParams models.QueryParams) ([]models.Listing, error) {
	var listings []models.Listing
	var l models.Listing

	srv := server.GetServer()
	collection := srv.Database.Collection(l.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "status,active")
	filter := mongodb.SelectConstructeur(queryParams)

	opts := options.Find().SetSort(bson.D{{Key: "pricePerUnit", Value: 1}})
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var listing models.Listing
		if err := cursor.Decode(&listing); err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}
	return listings, cursor.Err()
}

func (s *Auction) Buy(listingID, buyerID string, qty int) error {
	srv := server.GetServer()

	// Get listing
	var listing models.Listing
	listingCollection := srv.Database.Collection(listing.Collection())
	err := listingCollection.FindOne(context.TODO(), bson.M{"customID": listingID}).Decode(&listing)
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
	var buyer models.Player
	pCollection := srv.Database.Collection(buyer.Collection())
	err = pCollection.FindOne(context.TODO(), bson.M{"customID": buyerID}).Decode(&buyer)
	if err != nil {
		return errors.New("buyer not found")
	}
	if buyer.Gold < totalPrice {
		return errors.New("INSUFFICIENT_GOLD")
	}

	if buyerID == listing.SellerID {
		return errors.New("cannot buy your own listing")
	}

	// Atomic transaction
	session, err := srv.Database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx context.Context) (interface{}, error) {
		// Debit buyer
		_, err := pCollection.UpdateOne(ctx, bson.M{"customID": buyerID}, bson.M{
			"$inc": bson.M{"gold": -totalPrice},
		})
		if err != nil {
			return nil, err
		}

		// Credit seller
		_, err = pCollection.UpdateOne(ctx, bson.M{"customID": listing.SellerID}, bson.M{
			"$inc": bson.M{"gold": totalPrice},
		})
		if err != nil {
			return nil, err
		}

		// Transfer items to buyer inventory
		invCollection := srv.Database.Collection("inventory")
		filter := bson.M{"playerId": buyerID, "itemId": listing.ItemID}
		update := bson.M{
			"$inc": bson.M{"qty": int64(qty)},
			"$set": bson.M{"updatedAt": time.Now()},
			"$setOnInsert": bson.M{
				"playerId": buyerID,
				"itemId":   listing.ItemID,
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err = invCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return nil, err
		}

		// Update listing
		remainingQty := listing.Qty - qty
		listingUpdate := bson.M{}
		if remainingQty <= 0 {
			listingUpdate["$set"] = bson.M{"status": "sold", "qty": 0}
		} else {
			listingUpdate["$set"] = bson.M{"qty": remainingQty}
		}
		_, err = listingCollection.UpdateOne(ctx, bson.M{"customID": listingID}, listingUpdate)
		if err != nil {
			return nil, err
		}

		// Create trade record
		trade := models.Trade{
			CustomID:   functions.NewUUID(),
			BuyerID:    buyerID,
			SellerID:   listing.SellerID,
			ListingID:  listingID,
			Qty:        qty,
			TotalPrice: totalPrice,
			CreatedAt:  time.Now(),
		}
		tradeCollection := srv.Database.Collection(trade.Collection())
		_, err = tradeCollection.InsertOne(ctx, trade)
		return nil, err
	})

	return err
}

func (s *Auction) Cancel(listingID, sellerID string) error {
	srv := server.GetServer()

	var listing models.Listing
	listingCollection := srv.Database.Collection(listing.Collection())
	err := listingCollection.FindOne(context.TODO(), bson.M{"customID": listingID}).Decode(&listing)
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
	invCollection := srv.Database.Collection("inventory")
	filter := bson.M{"playerId": sellerID, "itemId": listing.ItemID}
	update := bson.M{
		"$inc": bson.M{"qty": int64(listing.Qty)},
		"$set": bson.M{"updatedAt": time.Now()},
		"$setOnInsert": bson.M{
			"playerId": sellerID,
			"itemId":   listing.ItemID,
		},
	}
		opts := options.UpdateOne().SetUpsert(true)
		_, err = invCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	// Mark listing as cancelled
	_, err = listingCollection.UpdateOne(context.TODO(), bson.M{"customID": listingID}, bson.M{
		"$set": bson.M{"status": "cancelled"},
	})
	return err
}
