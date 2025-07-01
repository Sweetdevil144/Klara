package main

import (
	"log"
	"os"
	"server/database"
	"server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Klara API",
		ServerHeader: "Klara API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			log.Printf("Error: %v", err)
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	app.Use(cors.New())

	if err := database.Init(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Disconnect()

	routes.SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🚀 Klara API is running!",
			"status":  "healthy",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("🚀 Klara API starting on port %s", port)
	log.Printf("📝 API available at: http://localhost:%s", port)
	log.Printf("🩺 Health check: http://localhost:%s/api/v1/public/health", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
