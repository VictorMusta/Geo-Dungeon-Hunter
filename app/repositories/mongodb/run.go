package mongodb

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RunRepository struct {
	db *mongo.Database
}

func NewRunRepository(db *mongo.Database) *RunRepository {
	return &RunRepository{db: db}
}

func (r *RunRepository) GetByPlayerID(playerID string) ([]models.Run, error) {
	var runs []models.Run
	var run models.Run
	collection := r.db.Collection(run.Collection())
	
	opts := options.Find().SetSort(bson.D{{Key: "startedAt", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.M{"playerId": playerID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var rd models.Run
		if err := cursor.Decode(&rd); err != nil {
			return nil, err
		}
		runs = append(runs, rd)
	}
	return runs, nil
}

func (r *RunRepository) GetByID(id string) (models.Run, error) {
	var run models.Run
	var params models.QueryParams
	collection := r.db.Collection(run.Collection())
	
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)
	err := collection.FindOne(context.TODO(), filter).Decode(&run)
	return run, err
}

func (r *RunRepository) GetActiveRun(playerID, dungeonId string) (models.Run, error) {
	var run models.Run
	collection := r.db.Collection(run.Collection())
	err := collection.FindOne(context.TODO(), bson.M{
		"playerId":  playerID,
		"dungeonId": dungeonId,
		"state":     "active",
	}).Decode(&run)
	return run, err
}

func (r *RunRepository) Create(run *models.Run) error {
	collection := r.db.Collection(run.Collection())
	_, err := collection.InsertOne(context.TODO(), run)
	return err
}

func (r *RunRepository) Update(id string, run *models.Run) error {
	var ru models.Run
	collection := r.db.Collection(ru.Collection())
	
	doc, err := mongodb.ToDoc(run)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"customID": id}, bson.M{"$set": doc})
	return err
}

func (r *RunRepository) ExecuteBossAttempt(ctx context.Context, runID string, playerID string, rewards models.RewardsGiven, isCompleted bool, killedStep models.KilledStep) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		rColl := r.db.Collection("run")
		pColl := r.db.Collection("player")
		iColl := r.db.Collection("inventory")

		// 1. Update Run
		runUpdate := bson.M{
			"$push": bson.M{"killedSteps": killedStep},
			"$inc":  bson.M{"currentStep": 1},
		}
		if isCompleted {
			now := time.Now()
			runUpdate["$set"] = bson.M{"state": "completed", "endedAt": now}
		}
		_, err := rColl.UpdateOne(sessCtx, bson.M{"customID": runID}, runUpdate)
		if err != nil {
			return nil, err
		}

		// 2. Update Player Gold
		_, err = pColl.UpdateOne(sessCtx, bson.M{"customID": playerID}, bson.M{
			"$inc": bson.M{"gold": rewards.Gold},
		})
		if err != nil {
			return nil, err
		}

		// 3. Update Inventory Items
		for _, item := range rewards.Items {
			filter := bson.M{"playerId": playerID, "itemId": item.ItemID}
			update := bson.M{
				"$inc": bson.M{"qty": int64(item.Qty)},
				"$set": bson.M{"updatedAt": time.Now()},
				"$setOnInsert": bson.M{
					"playerId": playerID,
					"itemId":   item.ItemID,
				},
			}
			opts := options.UpdateOne().SetUpsert(true)
			_, err := iColl.UpdateOne(sessCtx, filter, update, opts)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}
