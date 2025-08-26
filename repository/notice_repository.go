package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoticeModel struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	Content        string             `bson:"content"`
	Tags           []string           `bson:"tags,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}
