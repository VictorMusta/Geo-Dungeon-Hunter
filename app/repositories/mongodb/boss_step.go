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

type BossStepRepository struct {
	db *mongo.Database
}

func NewBossStepRepository(db *mongo.Database) *BossStepRepository {
	return &BossStepRepository{db: db}
}

func (r *BossStepRepository) GetByDungeon(dungeonId string) ([]models.BossStep, error) {
	var steps []models.BossStep
	var b models.BossStep
	collection := r.db.Collection(b.Collection())
	
	cursor, err := collection.Find(context.TODO(), bson.M{"dungeonId": dungeonId})
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
	return steps, nil
}

func (r *BossStepRepository) GetByDungeonOrdered(dungeonId string) ([]models.BossStep, error) {
	var steps []models.BossStep
	var b models.BossStep
	collection := r.db.Collection(b.Collection())
	
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	cursor, err := collection.Find(context.TODO(), bson.M{"dungeonId": dungeonId}, opts)
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
	return steps, nil
}

func (r *BossStepRepository) GetByID(dungeonId, id string) (models.BossStep, error) {
	var b models.BossStep
	collection := r.db.Collection(b.Collection())
	
	err := collection.FindOne(context.TODO(), bson.M{"dungeonId": dungeonId, "customID": id}).Decode(&b)
	return b, err
}

func (r *BossStepRepository) Create(step *models.BossStep) error {
	collection := r.db.Collection(step.Collection())
	_, err := collection.InsertOne(context.TODO(), step)
	return err
}

func (r *BossStepRepository) Update(id string, step *models.BossStep) error {
	var b models.BossStep
	collection := r.db.Collection(b.Collection())
	
	doc, err := mongodb.ToDoc(step)
	if err != nil {
		return err
	}

	result, err := collection.UpdateOne(context.TODO(), bson.M{"customID": id}, bson.M{"$set": doc})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("boss step not found")
	}
	return nil
}

func (r *BossStepRepository) Delete(dungeonId, id string) error {
	var b models.BossStep
	collection := r.db.Collection(b.Collection())
	_, err := collection.DeleteOne(context.TODO(), bson.M{"dungeonId": dungeonId, "customID": id})
	return err
}

func (r *BossStepRepository) CountByDungeon(dungeonId string) (int64, error) {
	var b models.BossStep
	collection := r.db.Collection(b.Collection())
	return collection.CountDocuments(context.TODO(), bson.M{"dungeonId": dungeonId})
}
