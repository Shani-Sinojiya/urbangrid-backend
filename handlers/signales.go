package handlers

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"urbangrid.com/functions"
)

type socketData struct {
	Id    string `json:"id,omitempty"`
	Count int64  `json:"vehicle_count,omitempty"`
}

func SetSignal(c *fiber.Ctx) error {
	var req socketData

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid body",
			"success": true,
		})
	}

	if err := functions.UpdateCount(req.Id, req.Count); err != nil {
		// log.Print(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to set signal",
			"success": false,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Updated Successfully",
			"success": true,
		})
	}
}

func GetSignalUpdate(c *fiber.Ctx) error {
	data, err := functions.GetSignalData()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get signal",
			"success": false,
		})
	}

	timer, err := functions.GetSignalTimer()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get signal",
			"success": false,
		})
	}

	parsedTime, _ := time.Parse(time.RFC3339, timer.Timer)
	seconds := int(time.Since(parsedTime).Seconds())

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Signal data",
		"success":    true,
		"data":       data,
		"greentimer": seconds,
	})
}

type Clients struct {
	Client *websocket.Conn
}

var Signal_clients []Clients

func SignalSocket(c *websocket.Conn) {
	Signal_clients = append(Signal_clients, Clients{
		Client: c,
	})

	defer func() {
		for i, client := range Signal_clients {
			if client.Client == c {
				Signal_clients = append(Signal_clients[:i], Signal_clients[i+1:]...)
			}
		}
	}()

	for {
		_, _, err := c.NextReader()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				break
			}
		}
	}
}
