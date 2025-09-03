package repository

import (
	"EthioGuide/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcedureContentModel struct {
	Prerequisites []string `bson:"prerequisites,omitempty"`
	Steps         []string `bson:"steps,omitempty"`
	Result        []string `bson:"result,omitempty"`
}

// ProcedureFee is a nested struct for the Procedure.fees field.
type ProcedureFeeModel struct {
	Label    string  `bson:"label,omitempty"`
	Currency string  `bson:"currency,omitempty"`
	Amount   float64 `bson:"amount,omitempty"`
}

// ProcessingTime is a nested struct for the Procedure.processing_time field.
type ProcessingTimeModel struct {
	MinDays int `bson:"min_days,omitempty"`
	MaxDays int `bson:"max_days,omitempty"`
}

type ProcedureModel struct {
	ID             primitive.ObjectID    `bson:"_id,omitempty"`
	GroupID        *primitive.ObjectID   `bson:"group_id,omitempty"` // Pointer to allow for nil
	OrganizationID primitive.ObjectID    `bson:"organization_id"`
	Name           string                `bson:"name"`
	Content        ProcedureContentModel `bson:"content,omitempty"`
	Fees           ProcedureFeeModel     `bson:"fees,omitempty"`
	ProcessingTime ProcessingTimeModel   `bson:"processing_time,omitempty"`
	CreatedAt      time.Time             `bson:"created_at"`
	// For M-M relationship with Notice
	NoticeIDs []primitive.ObjectID `bson:"notice_ids,omitempty"`
}

func ToDomainProcedure(model *ProcedureModel) *domain.Procedure {
	var groupID string
	if model.GroupID != nil {
		groupID = model.GroupID.Hex()
	}

	return &domain.Procedure{
		ID:             model.ID.Hex(),
		GroupID:        groupID,
		OrganizationID: model.OrganizationID.Hex(),
		Name:           model.Name,
		Content: domain.Content{
			Prerequisites: model.Content.Prerequisites,
			Steps:         model.Content.Steps,
			Result:        model.Content.Result,
		},
		Fees: domain.Fees{
			Label:    model.Fees.Label,
			Currency: model.Fees.Currency,
			Amount:   model.Fees.Amount,
		},
		ProcessingTime: domain.ProcessingTime{
			MinDays: model.ProcessingTime.MinDays,
			MaxDays: model.ProcessingTime.MaxDays,
		},
		CreatedAt: model.CreatedAt,
		UpdatedAt: time.Now(), // You can adjust this if you track updates separately
	}
}
