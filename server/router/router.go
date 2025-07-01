package router

import (
	"github.com/gofiber/fiber/v2"
	"server/handler"
	notes "server/handler/notes"
	users "server/handler/user"
)

func Route(app *fiber.App) {
	app.Get("/", handler.Healthy)

	user := app.Group("/user")
	user.Post("/", users.CreateOrSyncUser)
	user.Put("/:id", users.UpdateUser)
	user.Delete("/:id", users.DeleteUser)
	user.Get("/:id", users.GetUserProfile)
	user.Get("/:id/notes", users.GetUserWithNotes)

	note := app.Group("/note")
	note.Post("/", notes.CreateNote)
	note.Put("/:id", notes.UpdateNote)
	note.Delete("/:id", notes.DeleteNote)
	note.Get("/:id", notes.GetNote)
	note.Get("/user/:userId", notes.GetUserNotes)
}
