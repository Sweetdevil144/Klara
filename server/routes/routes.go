package routes

import (
	"server/handler/chat"
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
			"message": "Klara API is running",
		})
	})

	public.Get("/debug/auth", middleware.OptionalClerkMiddleware(), func(c *fiber.Ctx) error {
		userID, err := middleware.GetClerkUserIDFromContext(c)
		authHeader := c.Get("Authorization")

		return c.JSON(fiber.Map{
			"authenticated": err == nil,
			"userID":        userID,
			"hasAuthHeader": authHeader != "",
			"authHeader":    authHeader,
			"error": func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
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

	notesRoutes.Post("/:id/chat", chat.ChatWithNote)
	notesRoutes.Post("/:id/apply-suggestion", notes.ApplySuggestion)

	chatRoutes := protected.Group("/chat")
	chatRoutes.Post("/", chat.StartChat)
	chatRoutes.Get("/sessions", chat.GetChatSessions)
	chatRoutes.Get("/sessions/:sessionId", chat.GetChatHistory)
	chatRoutes.Delete("/sessions/:sessionId", chat.DeleteChatSession)
	chatRoutes.Post("/update-note", chat.UpdateNoteWithChat)
}
