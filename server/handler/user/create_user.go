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
)

func CreateOrSyncUser(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	type UserProfileRequest struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	userReq := new(UserProfileRequest)
	if err := c.BodyParser(userReq); err != nil {
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

	var existingUser models.User
	err = collection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&existingUser)

	if err == nil {
		existingUser.Email = userReq.Email
		existingUser.Username = userReq.Username
		existingUser.FirstName = userReq.FirstName
		existingUser.LastName = userReq.LastName
		existingUser.UpdatedAt = time.Now()

		_, err = collection.ReplaceOne(context.Background(), bson.M{"clerkId": clerkUserID}, existingUser)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "Failed to update user profile"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User profile updated",
			"user": models.UserProfile{
				ID:           existingUser.ID,
				ClerkID:      existingUser.ClerkID,
				Email:        existingUser.Email,
				Username:     existingUser.Username,
				FirstName:    existingUser.FirstName,
				LastName:     existingUser.LastName,
				HasOpenAIKey: existingUser.OpenAIKey != "",
				HasGeminiKey: existingUser.GeminiKey != "",
				CreatedAt:    existingUser.CreatedAt,
				UpdatedAt:    existingUser.UpdatedAt,
				NoteIds:      existingUser.NoteIds,
			},
		})
	}

	newUser := models.User{
		ClerkID:   clerkUserID,
		Email:     userReq.Email,
		Username:  userReq.Username,
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		NoteIds:   []primitive.ObjectID{},
	}

	result, err := collection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to create user profile"})
	}

	newUser.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User profile created",
		"user": models.UserProfile{
			ID:           newUser.ID,
			ClerkID:      newUser.ClerkID,
			Email:        newUser.Email,
			Username:     newUser.Username,
			FirstName:    newUser.FirstName,
			LastName:     newUser.LastName,
			HasOpenAIKey: false,
			HasGeminiKey: false,
			CreatedAt:    newUser.CreatedAt,
			UpdatedAt:    newUser.UpdatedAt,
			NoteIds:      newUser.NoteIds,
		},
	})
}
