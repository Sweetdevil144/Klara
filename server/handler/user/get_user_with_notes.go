package user

import (
	"log"
	"server/database"
	"server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserWithNotes(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User ID is required"})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID format"})
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	userWithNotes, err := utils.GetUserWithNotes(db, userID)
	if err != nil {
		log.Printf("Failed to get user with notes: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve user with notes"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User with notes retrieved successfully",
		"user":    userWithNotes,
	})
}
