package main

import (
	"urbangrid.com/database"
	"urbangrid.com/queues"
)

func init() {
	err := database.InitMongo()

	if err != nil {
		panic(err)
	}

	err = database.InitRedis()

	if err != nil {
		panic(err)
	}

	queues.InitQueue()
}
