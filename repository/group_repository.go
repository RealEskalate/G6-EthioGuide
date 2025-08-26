package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroupModel struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID  `bson:"organization_id"`
	ParentID       *primitive.ObjectID `bson:"parent_id,omitempty"` // Pointer to allow for nil (top-level groups)
	Title          string              `bson:"title"`
}
