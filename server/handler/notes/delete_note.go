package notes

import (
	"context"
	"log"
	"server/database"
	"server/models"
	"server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteNote(c *fiber.Ctx) error {
	note := new(models.Note)
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	collection := db.Collection("notes")
	var existingNote models.Note
	err = collection.FindOne(context.Background(), bson.M{"_id": note.ID}).Decode(&existingNote)
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "Note not found"})
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": note.ID})
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	if err := utils.RemoveNoteFromUser(db, existingNote.UserID, note.ID); err != nil {
		log.Printf("Failed to remove note from user: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Note deleted"})
}
