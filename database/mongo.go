package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client
var MongoDBDatabase *mongo.Database

func InitMongo() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	MongoDBClient = client
	MongoDBDatabase = client.Database("urban-grid")

	return nil
}

func DisconnectMongo() error {
	err := MongoDBClient.Disconnect(context.Background())
	if err != nil {
		return err
	}

	return nil
}
