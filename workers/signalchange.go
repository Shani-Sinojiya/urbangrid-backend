package workers

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/hibiken/asynq"
	"urbangrid.com/handlers"
)

func SignalChangeNotification(ctx context.Context, t *asynq.Task) error {
	Signal_clients := handlers.Signal_clients

	for _, client := range Signal_clients {
		client.Client.WriteMessage(websocket.TextMessage, []byte("signal changed"))
	}

	return nil
}
