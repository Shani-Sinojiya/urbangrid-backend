package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"urbangrid.com/functions"
)

var Emergency_clients []Clients

func SetEmergencyVehicleDetectEnable(c *fiber.Ctx) error {
	type Request struct {
		SignalID string `json:"signal_id"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
			"success": false,
		})
	}

	err := functions.SetEmergencyVehicleDetected(req.SignalID, true)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot set emergency vehicle detected",
			"success": false,
		})
	}

	for _, client := range Emergency_clients {
		client.Client.WriteMessage(websocket.TextMessage, []byte("Emergency Vehicle Detect status changed"))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Emergency Vehicle Detected",
		"success": true,
	})
}

func SetEmergencyVehicleDetectDisable(c *fiber.Ctx) error {
	type Request struct {
		SignalID string `json:"signal_id"`
	}

	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
			"success": false,
		})
	}

	err := functions.SetEmergencyVehicleDetected(req.SignalID, false)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot set emergency vehicle detected",
			"success": false,
		})
	}

	for _, client := range Emergency_clients {
		client.Client.WriteMessage(websocket.TextMessage, []byte("Emergency Vehicle Detect status changed"))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Emergency Vehicle Detected",
		"success": true,
	})
}

func GetEmergency(c *fiber.Ctx) error {
	data, err := functions.GetAllEmergencyVehicleDetected()
	if err != nil {
		if err.Error() == "redis: nil" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "No emergency vehicle detected",
				"success": true,
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot get emergency vehicle detected",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Emergency Vehicle Detected",
		"success": true,
		"data":    data,
	})
}

func EmergencySocket(c *websocket.Conn) {
	Emergency_clients = append(Emergency_clients, Clients{
		Client: c,
	})

	defer func() {
		for i, client := range Emergency_clients {
			if client.Client == c {
				Emergency_clients = append(Emergency_clients[:i], Emergency_clients[i+1:]...)
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
