package chat

import (
	"context"
	"server/database"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func DeleteChatSession(c *fiber.Ctx) error {
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
	messageCollection.DeleteMany(context.Background(), bson.M{
		"sessionId": sessionID,
		"clerkId":   clerkUserID,
	})

	sessionCollection := db.Collection("chat_sessions")
	result, err := sessionCollection.DeleteOne(context.Background(), bson.M{
		"sessionId": sessionID,
		"clerkId":   clerkUserID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to delete chat session"})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Chat session not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Chat session deleted successfully",
	})
}
