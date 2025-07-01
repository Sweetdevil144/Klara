package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	ClerkID   string               `json:"clerkId" bson:"clerkId" binding:"required"`
	Email     string               `json:"email,omitempty" bson:"email,omitempty"`
	Username  string               `json:"username,omitempty" bson:"username,omitempty"`
	FirstName string               `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string               `json:"lastName,omitempty" bson:"lastName,omitempty"`
	OpenAIKey string               `json:"openaiKey,omitempty" bson:"openaiKey,omitempty"`
	GeminiKey string               `json:"geminiKey,omitempty" bson:"geminiKey,omitempty"`
	CreatedAt time.Time            `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time            `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	NoteIds   []primitive.ObjectID `json:"noteIds,omitempty" bson:"noteIds,omitempty"`
}

type UserWithNotes struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	ClerkID   string               `json:"clerkId" bson:"clerkId"`
	Email     string               `json:"email,omitempty" bson:"email,omitempty"`
	Username  string               `json:"username,omitempty" bson:"username,omitempty"`
	FirstName string               `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string               `json:"lastName,omitempty" bson:"lastName,omitempty"`
	OpenAIKey string               `json:"openaiKey,omitempty" bson:"openaiKey,omitempty"`
	GeminiKey string               `json:"geminiKey,omitempty" bson:"geminiKey,omitempty"`
	CreatedAt time.Time            `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time            `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	NoteIds   []primitive.ObjectID `json:"noteIds,omitempty" bson:"noteIds,omitempty"`
	Notes     []Note               `json:"notes,omitempty" bson:"notes,omitempty"`
}

type UserProfile struct {
	ID           primitive.ObjectID   `json:"id"`
	ClerkID      string               `json:"clerkId"`
	Email        string               `json:"email,omitempty"`
	Username     string               `json:"username,omitempty"`
	FirstName    string               `json:"firstName,omitempty"`
	LastName     string               `json:"lastName,omitempty"`
	HasOpenAIKey bool                 `json:"hasOpenaiKey"`
	HasGeminiKey bool                 `json:"hasGeminiKey"`
	CreatedAt    time.Time            `json:"createdAt,omitempty"`
	UpdatedAt    time.Time            `json:"updatedAt,omitempty"`
	NoteIds      []primitive.ObjectID `json:"noteIds,omitempty"`
}

type APIKeysUpdate struct {
	OpenAIKey string `json:"openaiKey,omitempty"`
	GeminiKey string `json:"geminiKey,omitempty"`
}
