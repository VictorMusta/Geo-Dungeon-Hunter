package mongodb

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PlayerRepository struct {
	db *mongo.Database
}

func NewPlayerRepository(db *mongo.Database) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Get(params models.QueryParams) ([]models.Player, error) {
	var (
		players []models.Player
		player  models.Player
	)
	collection := r.db.Collection(player.Collection())
	filter := mongodb.SelectConstructeur(params)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var p models.Player
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (r *PlayerRepository) GetByID(id string) (models.Player, error) {
	var (
		player models.Player
		params models.QueryParams
	)
	collection := r.db.Collection(player.Collection())
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)

	err := collection.FindOne(context.TODO(), filter).Decode(&player)
	return player, err
}

func (r *PlayerRepository) Create(player *models.Player) error {
	collection := r.db.Collection(player.Collection())
	_, err := collection.InsertOne(context.TODO(), player)
	return err
}

func (r *PlayerRepository) Update(id string, player *models.Player) error {
	var (
		params models.QueryParams
		p      models.Player
	)
	collection := r.db.Collection(p.Collection())
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)

	doc, err := mongodb.ToDoc(player)
	if err != nil {
		return err
	}

	update := bson.M{"$set": doc}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("player not found")
	}
	return nil
}

func (r *PlayerRepository) Suspend(id string) error {
	var (
		p      models.Player
		params models.QueryParams
	)
	collection := r.db.Collection(p.Collection())
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)

	update := bson.M{"$set": bson.M{"suspended": true}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *PlayerRepository) FindByDisplayName(name string) (models.Player, error) {
	var p models.Player
	collection := r.db.Collection(p.Collection())
	err := collection.FindOne(context.TODO(), bson.M{"display_name": name}).Decode(&p)
	return p, err
}

func (r *PlayerRepository) GetByIDs(ids []string) ([]models.Player, error) {
	var players []models.Player
	for _, id := range ids {
		p, err := r.GetByID(id)
		if err == nil {
			players = append(players, p)
		}
	}
	return players, nil
}

func (r *PlayerRepository) EnsureIndexes() error {
	collection := r.db.Collection("player")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"display_name": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}
