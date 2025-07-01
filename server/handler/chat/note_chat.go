package chat

import (
	"context"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
	"server/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteChatRequest struct {
	Message  string `json:"message" binding:"required"`
	Model    string `json:"model" binding:"required"`    // Specific model ID like "gpt-4o-mini" or "gemini-1.5-flash"
	Provider string `json:"provider" binding:"required"` // "openai" or "gemini"
}

type NoteChatResponse struct {
	Message     string `json:"message"`
	Model       string `json:"model"`
	NoteContext string `json:"noteContext"`
	Suggestion  string `json:"suggestion,omitempty"`
}

func ChatWithNote(c *fiber.Ctx) error {
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

	var chatReq NoteChatRequest
	if err := c.BodyParser(&chatReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if chatReq.Provider != "openai" && chatReq.Provider != "gemini" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Provider must be 'openai' or 'gemini'",
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

	// Get user to access API keys
	userCollection := db.Collection("users")
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"clerkId": clerkUserID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User profile not found. Please create your profile first.",
			})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find user"})
	}

	// Get note
	noteObjID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid note ID",
		})
	}

	notesCollection := db.Collection("notes")
	var note models.Note
	err = notesCollection.FindOne(context.Background(), bson.M{
		"_id":    noteObjID,
		"userId": user.ID,
	}).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Note not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Failed to find note"})
	}

	// Get API key for the selected model
	var apiKey string
	switch chatReq.Provider {
	case "openai":
		apiKey = user.OpenAIKey
	case "gemini":
		apiKey = user.GeminiKey
	}

	if apiKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No " + chatReq.Provider + " API key found. Please add your API key in profile settings.",
		})
	}

	// Create context-aware prompt
	contextPrompt := createNoteContextPrompt(note, chatReq.Message)

	// Get AI response using the existing ChatWithAI method with specific model ID
	aiService := services.NewAIService()
	sessionID := "note-" + noteID // Create a unique session ID for note-specific chats
	response, err := aiService.ChatWithAI(
		user.ID.Hex(),
		clerkUserID,
		sessionID,
		contextPrompt,
		chatReq.Model,    // Pass the specific model ID
		chatReq.Provider, // Pass the provider
		apiKey,
	)
	if err != nil {
		log.Printf("AI service error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get AI response",
			"error":   err.Error(),
		})
	}

	chatResponse := NoteChatResponse{
		Message:     response.Message,
		Model:       chatReq.Provider, // Keep using provider for backward compatibility
		NoteContext: note.Title + ": " + note.Content,
		Suggestion:  response.Message, // The AI response can be applied as a suggestion
	}

	return c.Status(fiber.StatusOK).JSON(chatResponse)
}

func createNoteContextPrompt(note models.Note, userMessage string) string {
	return "You are an AI assistant helping with note-taking. Here's the current note:\n\n" +
		"Title: " + note.Title + "\n" +
		"Content: " + note.Content + "\n\n" +
		"User's request: " + userMessage + "\n\n" +
		"Please provide a helpful response. If the user is asking for improvements, suggestions, or modifications to the note content, " +
		"provide your response in a way that could be directly applied to enhance the note. " +
		"Focus on being concise and actionable. If you're suggesting content changes, provide the improved version that can be used to update the note."
}
