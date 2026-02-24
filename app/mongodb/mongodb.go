package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var mongoDatabase *mongo.Database

// func OpenMongoDB(dbhost string) (*mongo.Client, error) {
// 	var (
// 		mc  *mongo.Client
// 		err error
// 	)
// 	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
// 	defer cancel()
// 	mc, err = mongo.Connect(ctx, options.Client().ApplyURI(dbhost))
// 	if err == nil {
// 		err = mc.Ping(ctx, nil)
// 	}

// 	return mc, err
// }

// OpenMongoDB to open MongoDB connection
func OpenMongoDB(dbhost string) (*mongo.Client, error) {
	var (
		mc  *mongo.Client
		err error
	)

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	fmt.Println("dbhost: %s", dbhost)
	opts := options.Client().ApplyURI(dbhost).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	mc, err = mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		mc.Disconnect(context.TODO())
	}()

	// Send a ping to confirm a successful connection
	if err = mc.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return mc, err
}

// SetInstance init mongo database
func SetInstance(d *mongo.Database) {
	mongoDatabase = d
}

// GetInstance ...
func GetInstance() *mongo.Database {
	return mongoDatabase
}

// ToDoc ...
func ToDoc(v interface{}) (interface{}, error) {
	var doc interface{}

	data, err := bson.Marshal(v)
	if err == nil {
		err = bson.Unmarshal(data, &doc)
	}
	return doc, err
}
