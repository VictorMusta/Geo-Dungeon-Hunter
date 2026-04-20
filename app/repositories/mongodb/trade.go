package mongodb

import (
	"context"
	"dungeons/app/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TradeRepository struct {
	db *mongo.Database
}

func NewTradeRepository(db *mongo.Database) *TradeRepository {
	return &TradeRepository{db: db}
}

func (r *TradeRepository) Create(trade *models.Trade) error {
	collection := r.db.Collection(trade.Collection())
	_, err := collection.InsertOne(context.TODO(), trade)
	return err
}

func (r *TradeRepository) GetByPlayer(playerID string) ([]models.Trade, error) {
	var trades []models.Trade
	var t models.Trade
	collection := r.db.Collection(t.Collection())
	
	filter := bson.M{"$or": []bson.M{
		{"buyerId": playerID},
		{"sellerId": playerID},
	}}
	
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var tr models.Trade
		if err := cursor.Decode(&tr); err != nil {
			return nil, err
		}
		trades = append(trades, tr)
	}
	return trades, nil
}
