package inventory

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/server"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Inventory struct{}

func New() *Inventory {
	return &Inventory{}
}

func (s *Inventory) GetByPlayerID(playerID string) (*models.InventoryResponse, error) {
	srv := server.GetServer()
	invCollection := srv.Database.Collection("inventory")

	cursor, err := invCollection.Find(context.TODO(), bson.M{"playerId": playerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var items []models.InventoryItemDTO
	for cursor.Next(context.TODO()) {
		var entry models.InventoryEntry
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		items = append(items, models.InventoryItemDTO{
			ItemID: entry.ItemID,
			Qty:    entry.Qty,
		})
	}

	if items == nil {
		items = []models.InventoryItemDTO{}
	}

	return &models.InventoryResponse{
		PlayerID: playerID,
		Items:    items,
	}, cursor.Err()
}
