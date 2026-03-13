package run

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"errors"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Run struct {
	validate *validator.Validate
}

func New() *Run {
	return &Run{validate: validator.New()}
}

func (s *Run) Create(in *models.Run) (*models.Run, error) {
	srv := server.GetServer()

	// Verify dungeon exists and is published
	var d models.Dungeon
	dCollection := srv.Database.Collection(d.Collection())
	err := dCollection.FindOne(context.TODO(), bson.M{"customID": in.DungeonID}).Decode(&d)
	if err != nil {
		return nil, errors.New("dungeon not found")
	}
	if d.Status != "published" {
		return nil, errors.New("dungeon is not published")
	}

	// Verify player exists
	var p models.Player
	pCollection := srv.Database.Collection(p.Collection())
	err = pCollection.FindOne(context.TODO(), bson.M{"customID": in.PlayerID}).Decode(&p)
	if err != nil {
		return nil, errors.New("player not found")
	}

	// Check no active run for this player+dungeon
	var r models.Run
	rCollection := srv.Database.Collection(r.Collection())
	err = rCollection.FindOne(context.TODO(), bson.M{
		"dungeonId": in.DungeonID,
		"playerId":  in.PlayerID,
		"state":     "active",
	}).Decode(&r)
	if err == nil {
		return nil, errors.New("player already has an active run for this dungeon")
	}

	run := models.Run{
		CustomID:    functions.NewUUID(),
		DungeonID:   in.DungeonID,
		PlayerID:    in.PlayerID,
		State:       "active",
		CurrentStep: 1,
		KilledSteps: []models.KilledStep{},
		StartedAt:   time.Now(),
	}

	if _, err := rCollection.InsertOne(context.TODO(), run); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &run, nil
}

func (s *Run) GetByID(id string) (models.Run, error) {
	var (
		r           models.Run
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(r.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err := collection.FindOne(context.TODO(), filter).Decode(&r)
	return r, err
}

func (s *Run) GetByPlayerID(playerID string) ([]models.Run, error) {
	var runs []models.Run
	var r models.Run

	srv := server.GetServer()
	collection := srv.Database.Collection(r.Collection())

	opts := options.Find().SetSort(bson.D{{Key: "startedAt", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.M{"playerId": playerID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var run models.Run
		if err := cursor.Decode(&run); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, cursor.Err()
}

func (s *Run) Abandon(id string) error {
	r, err := s.GetByID(id)
	if err != nil {
		return errors.New("run not found")
	}
	if r.State != "active" {
		return errors.New("only active runs can be abandoned")
	}

	srv := server.GetServer()
	collection := srv.Database.Collection(r.Collection())

	now := time.Now()
	_, err = collection.UpdateOne(context.TODO(), bson.M{"customID": id}, bson.M{
		"$set": bson.M{"state": "abandoned", "endedAt": now, "updatedAt": now},
	})
	return err
}

type AttemptResult struct {
	Success bool              `json:"success"`
	Rewards models.RewardsGiven `json:"rewards"`
	RunCompleted bool         `json:"runCompleted"`
}

func (s *Run) AttemptBoss(runID, stepID string, lat, lon float64) (*AttemptResult, error) {
	r, err := s.GetByID(runID)
	if err != nil {
		return nil, errors.New("run not found")
	}

	if r.State != "active" {
		return nil, errors.New("run is not active")
	}

	// Idempotency: check if boss already killed
	for _, ks := range r.KilledSteps {
		if ks.BossStepID == stepID {
			return &AttemptResult{
				Success:      true,
				Rewards:      ks.RewardsGiven,
				RunCompleted: r.State == "completed",
			}, nil
		}
	}

	srv := server.GetServer()

	// Get the boss step
	var bs models.BossStep
	bsCollection := srv.Database.Collection(bs.Collection())
	err = bsCollection.FindOne(context.TODO(), bson.M{"customID": stepID}).Decode(&bs)
	if err != nil {
		return nil, errors.New("boss step not found")
	}

	// Validate step order
	if bs.Order != r.CurrentStep {
		return nil, errors.New("WRONG_STEP_ORDER")
	}

	// Validate location
	distance := functions.HaversineDistance(lat, lon, bs.Location.Lat, bs.Location.Lon)
	if distance > bs.Location.RadiusMeters {
		return nil, errors.New("NOT_IN_RANGE")
	}

	// Roll loot
	rewards := models.RewardsGiven{Gold: bs.GoldReward}
	for _, loot := range bs.LootTable {
		if rand.Float64() <= loot.DropRate {
			qty := loot.MinQty
			if loot.MaxQty > loot.MinQty {
				qty = loot.MinQty + rand.Intn(loot.MaxQty-loot.MinQty+1)
			}
			rewards.Items = append(rewards.Items, models.RewardItem{ItemID: loot.ItemID, Qty: qty})
		}
	}

	// Count total steps to check completion
	totalSteps, err := bsCollection.CountDocuments(context.TODO(), bson.M{"dungeonId": r.DungeonID})
	if err != nil {
		return nil, err
	}

	runCompleted := int64(r.CurrentStep) >= totalSteps

	// Atomic transaction: update run + player gold + inventory
	session, err := srv.Database.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx context.Context) (interface{}, error) {
		rCollection := srv.Database.Collection(r.Collection())

		killedStep := models.KilledStep{
			BossStepID:   stepID,
			KilledAt:     time.Now(),
			RewardsGiven: rewards,
		}

		runUpdate := bson.M{
			"$push": bson.M{"killedSteps": killedStep},
			"$inc":  bson.M{"currentStep": 1},
		}
		if runCompleted {
			now := time.Now()
			runUpdate["$set"] = bson.M{"state": "completed", "endedAt": now}
		}

		_, err := rCollection.UpdateOne(ctx, bson.M{"customID": runID}, runUpdate)
		if err != nil {
			return nil, err
		}

		// Credit player gold
		pCollection := srv.Database.Collection("player")
		_, err = pCollection.UpdateOne(ctx, bson.M{"customID": r.PlayerID}, bson.M{
			"$inc": bson.M{"gold": rewards.Gold},
		})
		if err != nil {
			return nil, err
		}

		// Add items to inventory
		invCollection := srv.Database.Collection("inventory")
		for _, item := range rewards.Items {
			filter := bson.M{"playerId": r.PlayerID, "itemId": item.ItemID}
			update := bson.M{
				"$inc": bson.M{"qty": int64(item.Qty)},
				"$set": bson.M{"updatedAt": time.Now()},
				"$setOnInsert": bson.M{
					"playerId": r.PlayerID,
					"itemId":   item.ItemID,
				},
			}
			opts := options.UpdateOne().SetUpsert(true)
			_, err := invCollection.UpdateOne(ctx, filter, update, opts)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &AttemptResult{
		Success:      true,
		Rewards:      rewards,
		RunCompleted: runCompleted,
	}, nil
}
