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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			// Automatically create user profile if it doesn't exist
			log.Printf("User profile not found, creating for clerkID: %s", clerkUserID)
			user = models.User{
				ClerkID:   clerkUserID,
				Email:     "",
				Username:  "",
				FirstName: "",
				LastName:  "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				NoteIds:   []primitive.ObjectID{},
			}

			result, createErr := collection.InsertOne(context.Background(), user)
			if createErr != nil {
				log.Printf("Failed to create user profile: %v", createErr)
				return c.Status(500).JSON(fiber.Map{"message": "Failed to create user profile"})
			}
			user.ID = result.InsertedID.(primitive.ObjectID)
			log.Printf("Created user profile with ID: %s", user.ID.Hex())
		} else {
			log.Printf("Failed to get user: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve user profile"})
		}
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
