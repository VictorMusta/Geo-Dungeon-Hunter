package item

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Item struct {
	validate *validator.Validate
}

func New() *Item {
	return &Item{
		validate: validator.New(),
	}
}

func (s *Item) Get(queryParams models.QueryParams) ([]models.ItemDef, error) {
	var (
		items  []models.ItemDef
		item   models.ItemDef
		cursor *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(item.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var i models.ItemDef
		if err := cursor.Decode(&i); err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		items = append(items, i)
	}

	if err := cursor.Err(); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return items, nil
}

func (s *Item) Create(in *models.ItemDef) (*models.ItemDef, error) {
	var item models.ItemDef

	srv := server.GetServer()
	collection := srv.Database.Collection(item.Collection())

	if err := s.validate.Struct(in); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if err := functions.ConvertInputStructToDataStruct(in, &item); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	item.CustomID = functions.NewUUID()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if _, err := collection.InsertOne(context.TODO(), item); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &item, nil
}

func (s *Item) GetByID(id string) (models.ItemDef, error) {
	var (
		item        models.ItemDef
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(item.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err := collection.FindOne(context.TODO(), filter).Decode(&item)
	return item, err
}
