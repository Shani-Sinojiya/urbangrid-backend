package functions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"urbangrid.com/collections"
)

func CreateAccident(location string, img primitive.ObjectID) (bool, primitive.ObjectID, error) {
	info, err := collections.Accident().InsertOne(context.Background(), bson.M{
		"location":  location,
		"img":       img,
		"createdAt": primitive.NewDateTimeFromTime(time.Now()),
	})

	if err != nil {
		return false, primitive.NilObjectID, err
	}

	return true, info.InsertedID.(primitive.ObjectID), nil
}

func DeleteAccident(id primitive.ObjectID) (bool, error) {
	_, err := collections.Accident().UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"deletedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

type Accident struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Loaction string `json:"location,omitempty" bson:"location,omitempty"`
	Img      string `json:"img,omitempty" bson:"img,omitempty"`
}

func GetAccidents() (bool, []Accident, error) {
	cursor, err := collections.Accident().Aggregate(context.Background(), bson.A{
		bson.M{
			"$match": bson.M{
				"deletedAt": bson.M{
					"$exists": false,
				},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "files",
				"localField":   "img",
				"foreignField": "_id",
				"as":           "img",
			},
		},
		bson.M{
			"$sort": bson.M{
				"createdAt": -1,
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$img",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$project": bson.M{
				"location": 1,
				"img":      "$img.url",
			},
		},
	})

	if err != nil {
		return false, nil, err
	}

	var data []Accident

	err = cursor.All(context.Background(), &data)

	if err != nil {
		return false, nil, err
	}

	return true, data, nil
}
