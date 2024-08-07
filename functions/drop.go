package functions

import (
	"context"

	"urbangrid.com/database"
)

func DropMongodb() error {
	return database.MongoDBDatabase.Drop(context.Background())
}
