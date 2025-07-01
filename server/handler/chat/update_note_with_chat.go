package chat

import (
	"context"
	"fmt"
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

func UpdateNoteWithChat(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	var updateReq models.UpdateNoteRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if updateReq.Model != "openai" && updateReq.Model != "gemini" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Model must be 'openai' or 'gemini'",
		})
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
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
	}

	var apiKey string
	switch updateReq.Model {
	case "openai":
		apiKey = user.OpenAIKey
	case "gemini":
		apiKey = user.GeminiKey
	}

	if apiKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No %s API key found", updateReq.Model),
		})
	}

	noteID, err := primitive.ObjectIDFromHex(updateReq.NoteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid note ID format",
		})
	}

	noteCollection := db.Collection("notes")
	var note models.Note
	err = noteCollection.FindOne(context.Background(), bson.M{
		"_id":    noteID,
		"userId": user.ID,
	}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Note not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Failed to retrieve note"})
	}

	updatedContent, err := aiService.UpdateNoteWithAI(
		user.ID.Hex(),
		clerkUserID,
		updateReq.SessionID,
		updateReq.NoteID,
		note.Content,
		updateReq.Model,
		apiKey,
		updateReq.Prompt,
	)
	if err != nil {
		log.Printf("AI update error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update note with AI",
			"error":   err.Error(),
		})
	}

	note.Content = updatedContent
	note.UpdatedAt = time.Now()

	_, err = noteCollection.ReplaceOne(context.Background(), bson.M{"_id": noteID}, note)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save updated note"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note updated successfully with AI assistance",
		"note": fiber.Map{
			"id":        note.ID,
			"title":     note.Title,
			"content":   note.Content,
			"updatedAt": note.UpdatedAt,
		},
	})
}
