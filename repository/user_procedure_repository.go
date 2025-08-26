package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProcedureModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	ProcedureID primitive.ObjectID `bson:"procedure_id"`
	Percent     int                `bson:"percent"`
	Status      string             `bson:"status"` // "Not Started", "In Progress", "Completed"
	UpdatedAt   time.Time          `bson:"updated_at"`
}
