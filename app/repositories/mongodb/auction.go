package mongodb

import (
	"context"
	"dungeons/app/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuctionRepository struct {
	db *mongo.Database
}

func NewAuctionRepository(db *mongo.Database) *AuctionRepository {
	return &AuctionRepository{db: db}
}

func (r *AuctionRepository) GetByStatus(status string) ([]models.Listing, error) {
	var listings []models.Listing
	var l models.Listing
	collection := r.db.Collection(l.Collection())

	cursor, err := collection.Find(context.TODO(), bson.M{"status": status})
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
	return listings, nil
}

func (r *AuctionRepository) GetByID(id string) (models.Listing, error) {
	var l models.Listing
	collection := r.db.Collection(l.Collection())
	err := collection.FindOne(context.TODO(), bson.M{"customID": id}).Decode(&l)
	return l, err
}

func (r *AuctionRepository) Create(listing *models.Listing) error {
	collection := r.db.Collection(listing.Collection())
	_, err := collection.InsertOne(context.TODO(), listing)
	return err
}

func (r *AuctionRepository) Update(id string, listing *models.Listing) error {
	var l models.Listing
	collection := r.db.Collection(l.Collection())

	_, err := collection.UpdateOne(context.TODO(), bson.M{"customID": id}, bson.M{"$set": listing})
	return err
}

func (r *AuctionRepository) BuyListing(ctx context.Context, buyerID, listingID string, priceTotal int64, sellerID string, itemID string, qty int, trade *models.Trade) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		lColl := r.db.Collection("listing")
		pColl := r.db.Collection("player")
		vColl := r.db.Collection("inventory")
		tColl := r.db.Collection("trade")

		// 1. Update Listing Status
		res, err := lColl.UpdateOne(sessCtx,
			bson.M{"customID": listingID, "status": "active"},
			bson.M{"$set": bson.M{"status": "sold", "updatedAt": time.Now()}})
		if err != nil || res.MatchedCount == 0 {
			return nil, errors.New("listing no longer available")
		}

		// 2. Subtract Gold from Buyer
		res, err = pColl.UpdateOne(sessCtx,
			bson.M{"customID": buyerID, "gold": bson.M{"$gte": priceTotal}},
			bson.M{"$inc": bson.M{"gold": -priceTotal}})
		if err != nil || res.MatchedCount == 0 {
			return nil, errors.New("insufficient gold")
		}

		// 3. Add Gold to Seller
		_, err = pColl.UpdateOne(sessCtx,
			bson.M{"customID": sellerID},
			bson.M{"$inc": bson.M{"gold": priceTotal}})
		if err != nil {
			return nil, err
		}

		// 4. Update Buyer Inventory
		_, err = vColl.UpdateOne(sessCtx,
			bson.M{"playerId": buyerID, "itemId": itemID},
			bson.M{
				"$inc":         bson.M{"qty": int64(qty)},
				"$set":         bson.M{"updatedAt": time.Now()},
				"$setOnInsert": bson.M{"playerId": buyerID, "itemId": itemID},
			},
			options.UpdateOne().SetUpsert(true))
		if err != nil {
			return nil, err
		}

		// 5. Create Trade Record
		_, err = tColl.InsertOne(sessCtx, trade)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func BoolPointer(b bool) *bool {
	return &b
}
