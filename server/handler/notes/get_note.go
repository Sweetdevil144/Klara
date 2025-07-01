package notes

import (
	"context"
	"log"
	"server/database"
	"server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNote(c *fiber.Ctx) error {
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
	err = collection.FindOne(context.Background(), bson.M{"_id": note.ID}).Decode(note)
	if err != nil {
		log.Fatal(err.Error())
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Note retrieved", "note": note})
}
