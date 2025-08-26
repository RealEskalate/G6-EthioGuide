package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIUserChatMessageModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Source    string             `bson:"source,omitempty"`
	Request   string             `bson:"request"`
	Response  string             `bson:"response"`
	Timestamp time.Time          `bson:"timestamp"`
}
