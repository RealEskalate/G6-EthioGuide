package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"user_id"`
	ProcedureID   primitive.ObjectID `bson:"procedure_id"`
	Content       string             `bson:"content"`
	LikeCount     int                `bson:"like_count"`
	DislikeCount  int                `bson:"dislike_count"`
	ViewCount     int                `bson:"view_count"`
	Type          string             `bson:"type"`                     // "Suggestion", "Bug Report", "Praise"
	Status        string             `bson:"status"`                   // "New", "Pending", "In Progress", "Addressed"
	AdminResponse *string            `bson:"admin_response,omitempty"` // Pointer to allow for nil
	Tags          []string           `bson:"tags,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}
