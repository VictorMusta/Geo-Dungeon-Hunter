package mongodb

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DungeonRepository struct {
	db *mongo.Database
}

func NewDungeonRepository(db *mongo.Database) *DungeonRepository {
	return &DungeonRepository{db: db}
}

func (r *DungeonRepository) Get(params models.QueryParams) ([]models.Dungeon, error) {
	var (
		dungeons []models.Dungeon
		d        models.Dungeon
	)
	collection := r.db.Collection(d.Collection())
	filter := mongodb.SelectConstructeur(params)

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
	return dungeons, nil
}

func (r *DungeonRepository) GetByID(id string) (models.Dungeon, error) {
	var (
		d      models.Dungeon
		params models.QueryParams
	)
	collection := r.db.Collection(d.Collection())
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)

	err := collection.FindOne(context.TODO(), filter).Decode(&d)
	return d, err
}

func (r *DungeonRepository) Create(dungeon *models.Dungeon) error {
	collection := r.db.Collection(dungeon.Collection())
	_, err := collection.InsertOne(context.TODO(), dungeon)
	return err
}

func (r *DungeonRepository) Update(id string, dungeon *models.Dungeon) error {
	var (
		d      models.Dungeon
		params models.QueryParams
	)
	collection := r.db.Collection(d.Collection())
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)

	doc, err := mongodb.ToDoc(dungeon)
	if err != nil {
		return err
	}

	update := bson.M{"$set": doc}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("dungeon not found")
	}
	return nil
}

func (r *DungeonRepository) GetPublished(params models.QueryParams) ([]models.Dungeon, error) {
	params.FilterClause = append(params.FilterClause, "status,published")
	return r.Get(params)
}

func (r *DungeonRepository) GetByMJ(mjId string, params models.QueryParams) ([]models.Dungeon, error) {
	params.FilterClause = append(params.FilterClause, "createdBy,"+mjId)
	return r.Get(params)
}
