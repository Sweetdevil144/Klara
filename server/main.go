package main

import (
	"log"
	"server/database"
	"server/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "AI Note Taker API",
		ServerHeader: "AI Note Taker API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	if err := database.Init(); err != nil {
		log.Fatal("unable to connect to client")
	}
	defer database.Disconnect()
	router.Route(app)

	app.Listen(":3000")
}
