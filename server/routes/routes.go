package routes

import (
	"server/handler/notes"
	"server/handler/user"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	public := api.Group("/public")
	public.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "AI Notes App API is running",
		})
	})

	protected := api.Group("/", middleware.ClerkMiddleware())

	userRoutes := protected.Group("/user")
	userRoutes.Post("/profile", user.CreateOrSyncUser)
	userRoutes.Get("/profile", user.GetUserProfile)
	userRoutes.Put("/api-keys", user.UpdateAPIKeys)
	userRoutes.Delete("/api-keys/:keyType", user.DeleteAPIKey)
	userRoutes.Get("/with-notes", user.GetUserWithNotes)

	notesRoutes := protected.Group("/notes")
	notesRoutes.Post("/", notes.CreateNote)
	notesRoutes.Get("/", notes.GetMyNotes)
	notesRoutes.Get("/:id", notes.GetNote)
	notesRoutes.Put("/:id", notes.UpdateNote)
	notesRoutes.Delete("/:id", notes.DeleteNote)

	api.Get("/my-notes", middleware.ClerkMiddleware(), notes.GetMyNotes)
}
