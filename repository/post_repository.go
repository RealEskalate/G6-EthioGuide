package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Content      string             `bson:"content"`
	LikeCount    int                `bson:"like_count"`
	DislikeCount int                `bson:"dislike_count"`
	ViewCount    int                `bson:"view_count"`
	Tags         []string           `bson:"tags,omitempty"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	// For M-M relationship with Procedure, replacing PostProcedure table
	ProcedureIDs []primitive.ObjectID `bson:"procedure_ids,omitempty"`
}
