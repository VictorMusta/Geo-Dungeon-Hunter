package bossstep

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BossStep struct {
	validate *validator.Validate
}

func New() *BossStep {
	return &BossStep{validate: validator.New()}
}

func (s *BossStep) Create(in *models.BossStep) (*models.BossStep, error) {
	var bs models.BossStep

	srv := server.GetServer()
	collection := srv.Database.Collection(bs.Collection())

	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	if err := functions.ConvertInputStructToDataStruct(in, &bs); err != nil {
		return nil, err
	}

	bs.CustomID = functions.NewUUID()
	bs.CreatedAt = time.Now()
	bs.UpdatedAt = time.Now()

	nextOrder, err := s.nextOrder(bs.DungeonID)
	if err != nil {
		return nil, err
	}
	bs.Order = nextOrder

	if _, err := collection.InsertOne(context.TODO(), bs); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &bs, nil
}

func (s *BossStep) nextOrder(dungeonID string) (int, error) {
	var bs models.BossStep

	srv := server.GetServer()
	collection := srv.Database.Collection(bs.Collection())

	opts := options.FindOne().SetSort(bson.D{{Key: "order", Value: -1}})
	err := collection.FindOne(context.TODO(), bson.M{"dungeonId": dungeonID}, opts).Decode(&bs)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}
		return 0, err
	}
	return bs.Order + 1, nil
}

func (s *BossStep) GetByID(id string) (models.BossStep, error) {
	var (
		bs          models.BossStep
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(bs.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err := collection.FindOne(context.TODO(), filter).Decode(&bs)
	return bs, err
}

func (s *BossStep) Update(dungeonID, stepID string, in *models.BossStep) error {
	bs, err := s.GetByID(stepID)
	if err != nil {
		return err
	}

	if bs.DungeonID != dungeonID {
		return mongo.ErrNoDocuments
	}

	srv := server.GetServer()
	collection := srv.Database.Collection(bs.Collection())

	if in.Name != "" {
		bs.Name = in.Name
	}
	if in.Location.RadiusMeters > 0 {
		bs.Location = in.Location
	}
	if in.ZoneDescription != "" {
		bs.ZoneDescription = in.ZoneDescription
	}
	if in.Difficulty > 0 {
		bs.Difficulty = in.Difficulty
	}
	if in.GoldReward > 0 {
		bs.GoldReward = in.GoldReward
	}
	if in.LootTable != nil {
		bs.LootTable = in.LootTable
	}
	bs.UpdatedAt = time.Now()

	var queryParams models.QueryParams
	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+stepID)
	filter := mongodb.SelectConstructeur(queryParams)

	doc, err := mongodb.ToDoc(bs)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": doc})
	return err
}

func (s *BossStep) GetByDungeonID(dungeonID string) ([]models.BossStep, error) {
	var steps []models.BossStep
	var bs models.BossStep

	srv := server.GetServer()
	collection := srv.Database.Collection(bs.Collection())

	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	cursor, err := collection.Find(context.TODO(), bson.M{"dungeonId": dungeonID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var step models.BossStep
		if err := cursor.Decode(&step); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, cursor.Err()
}

func (s *BossStep) Reorder(dungeonID string, orderedIDs []string) error {
	var bs models.BossStep

	srv := server.GetServer()
	collection := srv.Database.Collection(bs.Collection())

	for i, id := range orderedIDs {
		filter := bson.M{"customID": id, "dungeonId": dungeonID}
		update := bson.M{"$set": bson.M{"order": i + 1, "updatedAt": time.Now()}}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}
