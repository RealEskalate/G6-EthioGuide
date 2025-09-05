package repository

import (
	"EthioGuide/domain"
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

func ToDomainUserProcedure(model *UserProcedureModel) *domain.UserProcedure {
	return &domain.UserProcedure{
		ID:          model.ID.Hex(),
		UserID:      model.UserID.Hex(),
		ProcedureID: model.ProcedureID.Hex(),
		Percent:     model.Percent,
		Status:      model.Status,
		UpdatedAt:   model.UpdatedAt,
	}
}
