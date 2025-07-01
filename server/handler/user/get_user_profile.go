package user

import (
	"context"
	"log"
	"server/database"
	"server/middleware"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserProfile(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	collection := db.Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User profile not found",
			})
		}
		log.Printf("Failed to get user: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve user profile"})
	}

	profile := models.UserProfile{
		ID:           user.ID,
		ClerkID:      user.ClerkID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		HasOpenAIKey: user.OpenAIKey != "",
		HasGeminiKey: user.GeminiKey != "",
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		NoteIds:      user.NoteIds,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User profile retrieved successfully",
		"user":    profile,
	})
}
