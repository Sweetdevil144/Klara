package handler

import "github.com/gofiber/fiber/v2"

func Healthy(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Server is running",
	})
}
