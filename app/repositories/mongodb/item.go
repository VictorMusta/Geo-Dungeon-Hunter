package mongodb

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ItemRepository struct {
	db *mongo.Database
}

func NewItemRepository(db *mongo.Database) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Get(params models.QueryParams) ([]models.ItemDef, error) {
	var items []models.ItemDef
	var i models.ItemDef
	collection := r.db.Collection(i.Collection())
	filter := mongodb.SelectConstructeur(params)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var item models.ItemDef
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) GetByID(id string) (models.ItemDef, error) {
	var i models.ItemDef
	var params models.QueryParams
	collection := r.db.Collection(i.Collection())
	params.FilterClause = append(params.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(params)

	err := collection.FindOne(context.TODO(), filter).Decode(&i)
	return i, err
}

func (r *ItemRepository) Create(item *models.ItemDef) error {
	collection := r.db.Collection(item.Collection())
	_, err := collection.InsertOne(context.TODO(), item)
	return err
}

func (r *ItemRepository) Update(id string, item *models.ItemDef) error {
	var i models.ItemDef
	collection := r.db.Collection(i.Collection())

	doc, err := mongodb.ToDoc(item)
	if err != nil {
		return err
	}

	result, err := collection.UpdateOne(context.TODO(), bson.M{"customID": id}, bson.M{"$set": doc})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("item not found")
	}
	return nil
}

func (r *ItemRepository) Delete(id string) error {
	var i models.ItemDef
	collection := r.db.Collection(i.Collection())

	result, err := collection.DeleteOne(context.TODO(), bson.M{"customID": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("item not found")
	}
	return nil
}
