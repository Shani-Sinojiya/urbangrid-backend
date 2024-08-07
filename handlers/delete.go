package handlers

import (
	"github.com/gofiber/fiber/v2"
	"urbangrid.com/functions"
)

func DeleteEverything(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthoraized",
			"success": false,
		})
	}

	if body.Email != "admin@urbangrid.com" || body.Password != "admin@1234" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthoraized",
			"success": false,
		})
	}

	err = functions.DropMongodb()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while deleting",
			"success": false,
		})
	}

	err = functions.DropRedis()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while deleting",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Everything is deleted",
		"success": true,
	})
}
