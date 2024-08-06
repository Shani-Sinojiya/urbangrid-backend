package collections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"urbangrid.com/database"
)

func Accident() *mongo.Collection {
	return database.MongoDBDatabase.Collection("accidents")
}
