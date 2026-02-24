package Player

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Player struct {
	validate *validator.Validate
}

func New() *Player {
	return &Player{
		validate: validator.New(),
	}
}

// Get services to get list of Player on db
func (p *Player) Get(queryParams models.QueryParams) ([]models.Player, error) {
	var (
		err     error
		players []models.Player
		player  models.Player
		cursor  *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var player models.Player
		err = cursor.Decode(&player)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		players = append(players, player)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return players, err
}

// Create new player on db
func (p *Player) Create(in *models.Player) (*models.Player, error) {
	var player models.Player

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

	// Check input fields
	err := p.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &player)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	player.ID = functions.NewUUID()
	player.CreatedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), player)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &player, nil
}

// GetByID controller to get one Player by ID
func (p *Player) GetByID(id string) (models.Player, error) {
	var (
		err         error
		player      models.Player
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err = collection.FindOne(context.TODO(), filter).Decode(&player)
	if err == nil {
		if err == mongo.ErrNoDocuments {
			log.Error().Err(err).Msg("")
			return player, err
		}

	}
	return player, err
}

// Update controller to update a Player
func (p *Player) Update(id string, in *models.Player) error {
	var (
		doc         interface{}
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
		Player      models.Player
	)

	srv := server.GetServer()

	// Check input fields
	err = p.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	Player, err = p.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	err = functions.ConvertInputStructToDataStruct(in, &Player)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	Player.UpdatedAt = time.Now()
	collection := srv.Database.Collection(Player.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	if doc, err = mongodb.ToDoc(Player); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	update := bson.M{"$set": doc}
	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("Player to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("Player could not be updated")
	}
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

// Suspend or Delete controller to suspend a Player
func (s *Player) Suspend(id string) error {
	var (
		err         error
		queryParams models.QueryParams
		Player      models.Player
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(Player.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "suspended", Value: true},
		}},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)

	return err
}

// GetByIds controller to get list of Player by Ids
func (s *Player) GetByIds(ids []string) ([]models.Player, error) {
	var Players []models.Player
	for _, id := range ids {
		Player, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
		Players = append(Players, Player)
	}
	return Players, nil
}
