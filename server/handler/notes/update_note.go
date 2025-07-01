package notes

import (
	"context"
	"log"
	"server/database"
	"server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateNote(c *fiber.Ctx) error {
	note := new(models.Note)
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	note.UpdatedAt = time.Now()
	collection := db.Collection("notes")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": note.ID}, bson.M{"$set": note})
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Note updated"})
}
