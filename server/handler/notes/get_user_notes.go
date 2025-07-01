package notes

import (
	"context"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
	"server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserNotes(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	userCollection := db.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User profile not found",
			})
		}
		log.Printf("Failed to find user: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
	}

	notes, err := utils.GetAllUserNotes(db, user.ID)
	if err != nil {
		log.Printf("Failed to get user notes: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve notes"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notes retrieved successfully",
		"notes":   notes,
		"count":   len(notes),
	})
}

func GetMyNotes(c *fiber.Ctx) error {
	return GetUserNotes(c)
}
