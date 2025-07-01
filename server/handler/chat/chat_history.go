package chat

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/middleware"
	"server/models"
	"server/database"
)

func GetChatHistory(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	sessionID := c.Params("sessionId")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Session ID is required",
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	messageCollection := db.Collection("chat_messages")
	findOptions := options.Find().SetSort(bson.M{"createdAt": 1})
	cursor, err := messageCollection.Find(
		context.Background(),
		bson.M{
			"sessionId": sessionID,
			"clerkId":   clerkUserID,
		},
		findOptions,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve chat history"})
	}
	defer cursor.Close(context.Background())

	var messages []models.ChatMessage
	if err = cursor.All(context.Background(), &messages); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to decode chat messages"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Chat history retrieved successfully",
		"sessionId": sessionID,
		"messages":  messages,
		"count":     len(messages),
	})
}
