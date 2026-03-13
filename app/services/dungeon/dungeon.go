package dungeon

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
)

type Dungeon struct {
	validate *validator.Validate
}

func New() *Dungeon {
	return &Dungeon{validate: validator.New()}
}

func (s *Dungeon) Create(in *models.Dungeon) (*models.Dungeon, error) {
	var d models.Dungeon

	srv := server.GetServer()
	collection := srv.Database.Collection(d.Collection())

	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	if err := functions.ConvertInputStructToDataStruct(in, &d); err != nil {
		return nil, err
	}

	d.CustomID = functions.NewUUID()
	d.Status = "draft"
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	if _, err := collection.InsertOne(context.TODO(), d); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &d, nil
}

func (s *Dungeon) GetByID(id string) (models.Dungeon, error) {
	var (
		d           models.Dungeon
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(d.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err := collection.FindOne(context.TODO(), filter).Decode(&d)
	return d, err
}

func (s *Dungeon) Update(id string, in *models.Dungeon) error {
	d, err := s.GetByID(id)
	if err != nil {
		return errors.New("dungeon not found")
	}

	if d.Status != "draft" {
		return errors.New("only draft dungeons can be modified")
	}

	srv := server.GetServer()
	collection := srv.Database.Collection(d.Collection())

	if in.Title != "" {
		d.Title = in.Title
	}
	if in.Description != "" {
		d.Description = in.Description
	}
	if in.Area.Name != "" {
		d.Area = in.Area
	}
	d.UpdatedAt = time.Now()

	var queryParams models.QueryParams
	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	doc, err := mongodb.ToDoc(d)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": doc})
	return err
}

func (s *Dungeon) Publish(id string) error {
	d, err := s.GetByID(id)
	if err != nil {
		return errors.New("dungeon not found")
	}

	if d.Status != "draft" {
		return errors.New("only draft dungeons can be published")
	}

	srv := server.GetServer()
	var bs models.BossStep
	bsCollection := srv.Database.Collection(bs.Collection())

	count, err := bsCollection.CountDocuments(context.TODO(), bson.M{"dungeonId": id})
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("dungeon must have at least one boss step to publish")
	}

	collection := srv.Database.Collection(d.Collection())
	var queryParams models.QueryParams
	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{
		"$set": bson.M{"status": "published", "updatedAt": time.Now()},
	})
	return err
}

func (s *Dungeon) GetPublished(queryParams models.QueryParams) ([]models.Dungeon, error) {
	var dungeons []models.Dungeon
	var d models.Dungeon

	srv := server.GetServer()
	collection := srv.Database.Collection(d.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "status,published")
	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var dg models.Dungeon
		if err := cursor.Decode(&dg); err != nil {
			return nil, err
		}
		dungeons = append(dungeons, dg)
	}
	return dungeons, cursor.Err()
}
