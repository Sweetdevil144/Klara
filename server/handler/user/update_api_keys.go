package user

import (
	"context"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateAPIKeys(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	apiKeys := new(models.APIKeysUpdate)
	if err := c.BodyParser(apiKeys); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	collection := db.Collection("users")

	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	if apiKeys.OpenAIKey != "" {
		updateFields["openaiKey"] = apiKeys.OpenAIKey
	}
	if apiKeys.GeminiKey != "" {
		updateFields["geminiKey"] = apiKeys.GeminiKey
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"clerkId": clerkUserID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		log.Printf("Failed to update API keys: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to update API keys"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User profile not found",
		})
	}

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&user)
	if err != nil {
		log.Printf("Failed to get updated user: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "API keys updated but failed to retrieve user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "API keys updated successfully",
		"apiKeyStatus": fiber.Map{
			"hasOpenaiKey": user.OpenAIKey != "",
			"hasGeminiKey": user.GeminiKey != "",
		},
	})
}

func DeleteAPIKey(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	keyType := c.Params("keyType")
	if keyType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Key type parameter is required",
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	collection := db.Collection("users")

	var updateField string
	switch keyType {
	case "openai":
		updateField = "openaiKey"
	case "gemini":
		updateField = "geminiKey"
	case "claude":
		updateField = "claudeKey"
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid key type. Must be 'openai', 'gemini', or 'claude'",
		})
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"clerkId": clerkUserID},
		bson.M{"$set": bson.M{updateField: "", "updatedAt": time.Now()}},
	)

	if err != nil {
		log.Printf("Failed to delete API key: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to delete API key"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User profile not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": keyType + " API key deleted successfully",
	})
}
