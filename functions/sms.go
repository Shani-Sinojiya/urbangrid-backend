package functions

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"urbangrid.com/collections"
)

func CreateSMSLog(body string, msgid string) (bool, primitive.ObjectID, error) {
	info, err := collections.SMS().InsertOne(context.Background(), bson.M{
		"body":  body,
		"msgid": msgid,

		"createdAt": primitive.NewDateTimeFromTime(time.Now()),
	})

	if err != nil {
		return false, primitive.NilObjectID, err
	}

	return true, info.InsertedID.(primitive.ObjectID), nil
}

type logs struct {
	ID    string `json:"id,omitempty" bson:"_id,omitempty"`
	Body  string `json:"body,omitempty" bson:"body,omitempty"`
	Msgid string `json:"msgid,omitempty" bson:"msgid,omitempty"`
}

func GetLoga() ([]logs, error) {
	var logs []logs

	cursor, err := collections.SMS().Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{
		"createdAt": -1,
	}).SetProjection(bson.M{
		"body":  1,
		"msgid": 1,
	}))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &logs)

	if err != nil {
		return nil, err
	}

	return logs, nil
}
