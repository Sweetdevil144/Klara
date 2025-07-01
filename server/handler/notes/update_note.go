package notes

import (
	"context"
	"log"
	"server/database"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateNote(c *fiber.Ctx) error {
	noteID := c.Params("id")
	if noteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Note ID is required",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid note ID format",
		})
	}

	type UpdateRequest struct {
		Title   string `json:"title,omitempty"`
		Content string `json:"content,omitempty"`
	}

	updateReq := new(UpdateRequest)
	if err := c.BodyParser(updateReq); err != nil {
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

	collection := db.Collection("notes")

	// Prepare update fields
	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	if updateReq.Title != "" {
		updateFields["title"] = updateReq.Title
	}
	if updateReq.Content != "" {
		updateFields["content"] = updateReq.Content
	}

	// Update the note
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updateFields})
	if err != nil {
		log.Printf("Failed to update note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to update note"})
	}

	// Fetch and return the updated note
	var updatedNote models.Note
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&updatedNote)
	if err != nil {
		log.Printf("Failed to fetch updated note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Note updated but failed to retrieve"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note updated successfully",
		"note":    updatedNote,
	})
}
