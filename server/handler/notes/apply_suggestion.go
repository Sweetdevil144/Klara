package notes

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

type SuggestionRequest struct {
	NewTitle   string `json:"newTitle,omitempty"`
	NewContent string `json:"newContent,omitempty"`
}

func ApplySuggestion(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	noteID := c.Params("id")
	if noteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Note ID is required",
		})
	}

	var suggestionReq SuggestionRequest
	if err := c.BodyParser(&suggestionReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// At least one field should be provided
	if suggestionReq.NewTitle == "" && suggestionReq.NewContent == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one field (newTitle or newContent) must be provided",
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	// Get user to verify ownership
	userCollection := db.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User profile not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
	}

	// Convert noteID to ObjectID
	noteObjID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid note ID",
		})
	}

	// Update the note
	notesCollection := db.Collection("notes")

	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	if suggestionReq.NewTitle != "" {
		updateFields["title"] = suggestionReq.NewTitle
	}

	if suggestionReq.NewContent != "" {
		updateFields["content"] = suggestionReq.NewContent
	}

	filter := bson.M{
		"_id":    noteObjID,
		"userId": user.ID,
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := notesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed to update note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to update note"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Note not found or you don't have permission to update it",
		})
	}

	// Get the updated note
	var updatedNote models.Note
	err = notesCollection.FindOne(context.Background(), filter).Decode(&updatedNote)
	if err != nil {
		log.Printf("Failed to fetch updated note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Note updated but failed to fetch updated version"})
	}

	return c.Status(fiber.StatusOK).JSON(updatedNote)
}
