package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"urbangrid.com/constants"
	"urbangrid.com/functions"
	"urbangrid.com/queues"
)

var Accident_clients []Clients

func UploadAccident(c *fiber.Ctx) error {
	// form in get file and location
	img, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid body",
			"success": false,
		})
	}

	file, err := img.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid body",
			"success": false,
		})
	}

	defer file.Close()

	ctype := img.Header.Get("Content-Type")

	var dest string

	if ctype == "image/png" {
		dest = fmt.Sprintf("./uploads/%s.png", primitive.NewObjectIDFromTimestamp(time.Now()).Hex())
	} else if ctype == "image/jpg" {
		dest = fmt.Sprintf("./uploads/%s.jpg", primitive.NewObjectIDFromTimestamp(time.Now()).Hex())
	} else if ctype == "image/jpeg" {
		dest = fmt.Sprintf("./uploads/%s.jpeg", primitive.NewObjectIDFromTimestamp(time.Now()).Hex())
	} else {
		dest = fmt.Sprintf("./uploads/%s.jpeg", primitive.NewObjectIDFromTimestamp(time.Now()).Hex())
	}

	err = c.SaveFile(img, dest)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"success": false,
		})
	}

	_, fid, err := functions.CreateFile(strings.ReplaceAll(dest, "./uploads", "/accidents"), ctype)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"success": false,
		})
	}

	locations := c.FormValue("locations")

	_, _, err = functions.CreateAccident(locations, fid)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"success": false,
		})
	}

	if _, err := queues.EnqueueTask(asynq.NewTask(constants.ACCIDENT_SMS, []byte(locations))); err != nil {
		var err error
		for err != nil {
			_, err = queues.EnqueueTask(asynq.NewTask(constants.ACCIDENT_SMS, []byte(locations)))
		}
	}

	for _, client := range Accident_clients {
		client.Client.WriteMessage(websocket.TextMessage, []byte("new accident detected"))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "accident uploaded successfully",
		"success": true,
	})
}

func GetAccidents(c *fiber.Ctx) error {
	_, data, err := functions.GetAccidents()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "accident found successfully",
		"success": true,
		"data":    data,
	})
}

func AccidentSocket(c *websocket.Conn) {
	Accident_clients = append(Accident_clients, Clients{
		Client: c,
	})

	defer func() {
		for i, client := range Accident_clients {
			if client.Client == c {
				Accident_clients = append(Accident_clients[:i], Accident_clients[i+1:]...)
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
