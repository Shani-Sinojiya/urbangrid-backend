package collections

import (
	"go.mongodb.org/mongo-driver/mongo"
	"urbangrid.com/database"
)

func SMS() *mongo.Collection {
	return database.MongoDBDatabase.Collection("sms-logs")
}
