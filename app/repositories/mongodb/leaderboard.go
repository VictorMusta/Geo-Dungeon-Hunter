package mongodb

import (
	"context"
	"dungeons/app/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LeaderboardRepository struct {
	db *mongo.Database
}

func NewLeaderboardRepository(db *mongo.Database) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

func (r *LeaderboardRepository) GetByCompletions(limit int) ([]repositories.LeaderboardEntry, error) {
	runCollection := r.db.Collection("run")

	pipeline := bson.A{
		bson.M{"$match": bson.M{"state": "completed"}},
		bson.M{"$group": bson.M{
			"_id":   "$playerId",
			"score": bson.M{"$sum": 1},
		}},
		bson.M{"$sort": bson.M{"score": -1}},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{
			"from":         "player",
			"localField":   "_id",
			"foreignField": "customID",
			"as":           "playerInfo",
		}},
		bson.M{"$unwind": bson.M{"path": "$playerInfo", "preserveNullAndEmptyArrays": true}},
		bson.M{"$addFields": bson.M{
			"displayName": bson.M{"$ifNull": bson.A{"$playerInfo.display_name", "Unknown"}},
		}},
		bson.M{"$project": bson.M{"playerInfo": 0}},
	}

	cursor, err := runCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var entries []repositories.LeaderboardEntry
	if err := cursor.All(context.TODO(), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *LeaderboardRepository) GetByGold(limit int) ([]repositories.LeaderboardEntry, error) {
	playerCollection := r.db.Collection("player")

	opts := options.Find().SetSort(bson.D{{Key: "gold", Value: -1}}).SetLimit(int64(limit))
	cursor, err := playerCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var entries []repositories.LeaderboardEntry
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entry := repositories.LeaderboardEntry{
			Score: 0,
		}
		if id, ok := doc["customID"].(string); ok {
			entry.PlayerID = id
		}
		if name, ok := doc["display_name"].(string); ok {
			entry.DisplayName = name
		}
		switch g := doc["gold"].(type) {
		case int32:
			entry.Score = float64(g)
		case int64:
			entry.Score = float64(g)
		case float64:
			entry.Score = g
		}
		entries = append(entries, entry)
	}
	return entries, cursor.Err()
}

func (r *LeaderboardRepository) GetBySpeed(dungeonID string, limit int) ([]repositories.LeaderboardEntry, error) {
	runCollection := r.db.Collection("run")

	pipeline := bson.A{
		bson.M{"$match": bson.M{"state": "completed", "dungeonId": dungeonID}},
		bson.M{"$addFields": bson.M{
			"duration": bson.M{"$subtract": bson.A{"$endedAt", "$startedAt"}},
		}},
		bson.M{"$sort": bson.M{"duration": 1}},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{
			"from":         "player",
			"localField":   "playerId",
			"foreignField": "customID",
			"as":           "playerInfo",
		}},
		bson.M{"$unwind": bson.M{"path": "$playerInfo", "preserveNullAndEmptyArrays": true}},
		bson.M{"$project": bson.M{
			"playerId":    "$playerId",
			"displayName": bson.M{"$ifNull": bson.A{"$playerInfo.display_name", "Unknown"}},
			"score":       bson.M{"$divide": bson.A{"$duration", 1000}},
		}},
	}

	cursor, err := runCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var entries []repositories.LeaderboardEntry
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		entry := repositories.LeaderboardEntry{}
		if id, ok := doc["playerId"].(string); ok {
			entry.PlayerID = id
		}
		if name, ok := doc["displayName"].(string); ok {
			entry.DisplayName = name
		}
		if score, ok := doc["score"].(float64); ok {
			entry.Score = score
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
