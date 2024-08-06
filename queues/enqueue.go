package queues

import (
	"context"

	"github.com/hibiken/asynq"
)

func EnqueueTask(t *asynq.Task) (*asynq.TaskInfo, error) {
	return Client.EnqueueContext(context.Background(), t)
}
