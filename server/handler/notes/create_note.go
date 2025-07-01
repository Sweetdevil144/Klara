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

func CreateNote(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	type NoteRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	noteReq := new(NoteRequest)
	if err := c.BodyParser(noteReq); err != nil {
		log.Printf("Body parsing error: %v", err)
		log.Printf("Request body: %s", c.Body())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// More flexible validation - allow empty content but require title
	if noteReq.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Title is required",
		})
	}

	if noteReq.Content == "" {
		noteReq.Content = " "
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
		} else {
			log.Printf("Failed to find user: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
		}
	}

	note := models.Note{
		Title:     noteReq.Title,
		Content:   noteReq.Content,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := db.Collection("notes")
	result, err := collection.InsertOne(context.Background(), note)
	if err != nil {
		log.Printf("Failed to create note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to create note"})
	}

	noteID := result.InsertedID.(primitive.ObjectID)

	if err := utils.AddNoteToUser(db, user.ID, noteID); err != nil {
		log.Printf("Failed to add note to user: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Note created successfully",
		"noteId":  noteID,
		"note": fiber.Map{
			"id":        noteID,
			"title":     note.Title,
			"content":   note.Content,
			"userId":    note.UserID,
			"createdAt": note.CreatedAt,
			"updatedAt": note.UpdatedAt,
		},
	})
}
