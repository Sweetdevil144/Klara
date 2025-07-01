package notes

import (
	"context"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
	"server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			log.Printf("User not found, creating profile for clerkID: %s", clerkUserID)
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

			result, createErr := userCollection.InsertOne(context.Background(), user)
			if createErr != nil {
				log.Printf("Failed to create user profile: %v", createErr)
				return c.Status(500).JSON(fiber.Map{"message": "Failed to create user profile"})
			}
			user.ID = result.InsertedID.(primitive.ObjectID)
			log.Printf("Created user profile with ID: %s", user.ID.Hex())

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Notes retrieved successfully",
				"notes":   []models.Note{},
				"count":   0,
			})
		} else {
			log.Printf("Failed to find user: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
		}
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
