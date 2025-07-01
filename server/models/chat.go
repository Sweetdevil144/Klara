package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID string             `json:"sessionId" bson:"sessionId"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	ClerkID   string             `json:"clerkId" bson:"clerkId"`
	Role      string             `json:"role" bson:"role"`
	Content   string             `json:"content" bson:"content"`
	Model     string             `json:"model" bson:"model"`
	MemoryIds []string           `json:"memoryIds,omitempty" bson:"memoryIds,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type ChatSession struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID    string             `json:"sessionId" bson:"sessionId"`
	UserID       primitive.ObjectID `json:"userId" bson:"userId"`
	ClerkID      string             `json:"clerkId" bson:"clerkId"`
	Title        string             `json:"title" bson:"title"`
	Model        string             `json:"model" bson:"model"`
	MessageCount int                `json:"messageCount" bson:"messageCount"`
	LastActivity time.Time          `json:"lastActivity" bson:"lastActivity"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Memory struct {
	ID         string                 `json:"id"`
	Memory     string                 `json:"memory"`
	UserID     string                 `json:"user_id"`
	AgentID    string                 `json:"agent_id,omitempty"`
	AppID      string                 `json:"app_id,omitempty"`
	RunID      string                 `json:"run_id,omitempty"`
	Hash       string                 `json:"hash,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Categories []string               `json:"categories,omitempty"`
	Immutable  bool                   `json:"immutable"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type ChatRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message" binding:"required"`
	Model     string `json:"model" binding:"required"`
}

type ChatResponse struct {
	SessionID string    `json:"sessionId"`
	Message   string    `json:"message"`
	Role      string    `json:"role"`
	Model     string    `json:"model"`
	Memories  []Memory  `json:"memories,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateNoteRequest struct {
	NoteID    string `json:"noteId" binding:"required"`
	SessionID string `json:"sessionId" binding:"required"`
	Model     string `json:"model" binding:"required"`
	Prompt    string `json:"prompt,omitempty"`
}

type Mem0AddRequest struct {
	Messages         []map[string]string    `json:"messages"`
	AgentID          string                 `json:"agent_id,omitempty"`
	UserID           string                 `json:"user_id"`
	AppID            string                 `json:"app_id,omitempty"`
	RunID            string                 `json:"run_id,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Includes         string                 `json:"includes,omitempty"`
	Excludes         string                 `json:"excludes,omitempty"`
	Infer            bool                   `json:"infer"`
	CustomCategories map[string]interface{} `json:"custom_categories,omitempty"`
	OrgID            string                 `json:"org_id,omitempty"`
	ProjectID        string                 `json:"project_id,omitempty"`
	Version          string                 `json:"version,omitempty"`
}

type Mem0SearchRequest struct {
	Query     string                 `json:"query"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	TopK      int                    `json:"top_k,omitempty"`
	Fields    []string               `json:"fields,omitempty"`
	Rerank    bool                   `json:"rerank,omitempty"`
	OrgID     string                 `json:"org_id,omitempty"`
	ProjectID string                 `json:"project_id,omitempty"`
}

type Mem0GetRequest struct {
	Filters map[string]interface{} `json:"filters,omitempty"`
}

type Mem0BatchDeleteRequest struct {
	MemoryIds []string `json:"memory_ids"`
}
