package chat

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"server/database"
	"server/middleware"
	"server/models"
	"server/services"
	"time"
)

var (
	aiService = services.NewAIService()
)

func StartChat(c *fiber.Ctx) error {
	clerkUserID, err := middleware.GetClerkUserIDFromContext(c)
	if err != nil {
		return err
	}

	var chatReq models.ChatRequest
	if err := c.BodyParser(&chatReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if chatReq.Model != "openai" && chatReq.Model != "gemini" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Model must be 'openai' or 'gemini'",
		})
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "Database connection failed"})
	}

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

	var apiKey string
	switch chatReq.Model {
	case "openai":
		apiKey = user.OpenAIKey
	case "gemini":
		apiKey = user.GeminiKey
	}

	if apiKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No %s API key found. Please add your API key in profile settings.", chatReq.Model),
		})
	}

	if chatReq.SessionID == "" {
		chatReq.SessionID = uuid.New().String()
	}

	sessionCollection := db.Collection("chat_sessions")
	session := models.ChatSession{
		SessionID:    chatReq.SessionID,
		UserID:       user.ID,
		ClerkID:      clerkUserID,
		Title:        generateSessionTitle(chatReq.Message),
		Model:        chatReq.Model,
		MessageCount: 1,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	var existingSession models.ChatSession
	err = sessionCollection.FindOne(context.Background(), bson.M{"sessionId": chatReq.SessionID}).Decode(&existingSession)
	if err == nil {
		existingSession.MessageCount++
		existingSession.LastActivity = time.Now()
		existingSession.UpdatedAt = time.Now()
		sessionCollection.ReplaceOne(context.Background(), bson.M{"sessionId": chatReq.SessionID}, existingSession)
	} else {
		sessionCollection.InsertOne(context.Background(), session)
	}

	messageCollection := db.Collection("chat_messages")
	userMessage := models.ChatMessage{
		SessionID: chatReq.SessionID,
		UserID:    user.ID,
		ClerkID:   clerkUserID,
		Role:      "user",
		Content:   chatReq.Message,
		Model:     chatReq.Model,
		CreatedAt: time.Now(),
	}
	messageCollection.InsertOne(context.Background(), userMessage)

	response, err := aiService.ChatWithAI(
		user.ID.Hex(),
		clerkUserID,
		chatReq.SessionID,
		chatReq.Message,
		chatReq.Model,
		apiKey,
	)
	if err != nil {
		log.Printf("AI service error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get AI response",
			"error":   err.Error(),
		})
	}

	aiMessage := models.ChatMessage{
		SessionID: chatReq.SessionID,
		UserID:    user.ID,
		ClerkID:   clerkUserID,
		Role:      "assistant",
		Content:   response.Message,
		Model:     chatReq.Model,
		CreatedAt: time.Now(),
	}
	messageCollection.InsertOne(context.Background(), aiMessage)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Chat response generated successfully",
		"data":    response,
	})
}

func generateSessionTitle(message string) string {
	if len(message) > 50 {
		return message[:47] + "..."
	}
	return message
}
