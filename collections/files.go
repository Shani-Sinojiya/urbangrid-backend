package collections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"urbangrid.com/database"
)

func Files() *mongo.Collection {
	return database.MongoDBDatabase.Collection("files")
}
