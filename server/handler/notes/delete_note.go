package notes

import (
	"context"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
	"server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteNote(c *fiber.Ctx) error {
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

	objectID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid note ID format",
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	// First, get the user to verify ownership
	userCollection := db.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&user)
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
	}

	// Find the note and verify it belongs to the user
	collection := db.Collection("notes")
	var existingNote models.Note
	err = collection.FindOne(context.Background(), bson.M{
		"_id":    objectID,
		"userId": user.ID,
	}).Decode(&existingNote)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Note not found",
			})
		}
		log.Printf("Failed to find note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find note"})
	}

	// Delete the note
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Failed to delete note: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Failed to delete note"})
	}

	// Remove note reference from user
	if err := utils.RemoveNoteFromUser(db, existingNote.UserID, objectID); err != nil {
		log.Printf("Failed to remove note from user: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}
