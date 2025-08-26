package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChecklistItemModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserProcedureID primitive.ObjectID `bson:"user_procedure_id"`
	Type            string             `bson:"type"` // "Requirement" or "Step"
	Content         string             `bson:"content"`
	IsChecked       bool               `bson:"is_checked"` // Renamed from 'status' for clarity
}
