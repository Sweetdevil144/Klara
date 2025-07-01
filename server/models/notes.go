package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty" binding:"required"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty" binding:"required"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId,omitempty" binding:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty" binding:"required"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty" binding:"required"`
}
