package functions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"urbangrid.com/collections"
)

func CreateFile(fileurl, filetype string) (bool, primitive.ObjectID, error) {
	info, err := collections.Files().InsertOne(context.Background(), bson.M{
		"url":  fileurl,
		"type": filetype,

		"createdAt": primitive.NewDateTimeFromTime(time.Now()),
	})

	if err != nil {
		return false, primitive.NilObjectID, err
	}

	return true, info.InsertedID.(primitive.ObjectID), nil
}

func DeleteFile(id primitive.ObjectID) (bool, error) {
	_, err := collections.Files().UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"deletedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
