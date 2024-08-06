package queues

import (
	"github.com/hibiken/asynq"
)

var Client *asynq.Client
var Server *asynq.Server

func InitQueue() {
	r := asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
		DB:   3,
	}

	Client = asynq.NewClient(r)

	Server = asynq.NewServer(r, asynq.Config{
		Concurrency: 10,
	})
}

func DisconnectQueue() error {
	err := Client.Close()
	return err
}
