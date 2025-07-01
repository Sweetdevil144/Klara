package chat

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
)

func GetChatSessions(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	sessionCollection := db.Collection("chat_sessions")
	findOptions := options.Find().SetSort(bson.M{"lastActivity": -1})
	cursor, err := sessionCollection.Find(
		context.Background(),
		bson.M{"clerkId": clerkUserID},
		findOptions,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve chat sessions"})
	}
	defer cursor.Close(context.Background())

	var sessions []models.ChatSession
	if err = cursor.All(context.Background(), &sessions); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to decode chat sessions"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Chat sessions retrieved successfully",
		"sessions": sessions,
		"count":    len(sessions),
	})
}
