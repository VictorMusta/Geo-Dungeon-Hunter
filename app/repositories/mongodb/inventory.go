package mongodb

import (
	"context"
	"dungeons/app/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InventoryRepository struct {
	db *mongo.Database
}

func NewInventoryRepository(db *mongo.Database) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetByPlayerID(playerID string) ([]models.InventoryEntry, error) {
	var entries []models.InventoryEntry
	var i models.InventoryEntry
	collection := r.db.Collection(i.Collection())
	
	cursor, err := collection.Find(context.TODO(), bson.M{"playerId": playerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var entry models.InventoryEntry
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *InventoryRepository) GetByItem(playerID, itemID string) (models.InventoryEntry, error) {
	var i models.InventoryEntry
	collection := r.db.Collection(i.Collection())
	err := collection.FindOne(context.TODO(), bson.M{"playerId": playerID, "itemId": itemID}).Decode(&i)
	return i, err
}

func (r *InventoryRepository) Update(playerID string, itemID string, qtyDelta int64) error {
	var i models.InventoryEntry
	collection := r.db.Collection(i.Collection())
	
	filter := bson.M{"playerId": playerID, "itemId": itemID}
	update := bson.M{
		"$inc": bson.M{"qty": qtyDelta},
		"$set": bson.M{"updatedAt": time.Now()},
		"$setOnInsert": bson.M{
			"playerId": playerID,
			"itemId":   itemID,
		},
	}
	opts := options.UpdateOne().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}
