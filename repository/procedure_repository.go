package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcedureContent struct {
	Prerequisites []string `bson:"prerequisites,omitempty"`
	Steps         []string `bson:"steps,omitempty"`
	Result        string   `bson:"result,omitempty"`
}

// ProcedureFee is a nested struct for the Procedure.fees field.
type ProcedureFee struct {
	Label    string  `bson:"label,omitempty"`
	Currency string  `bson:"currency,omitempty"`
	Amount   float64 `bson:"amount,omitempty"`
}

// ProcessingTime is a nested struct for the Procedure.processing_time field.
type ProcessingTime struct {
	MinDays int `bson:"min_days,omitempty"`
	MaxDays int `bson:"max_days,omitempty"`
}

type ProcedureModel struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	GroupID        *primitive.ObjectID `bson:"group_id,omitempty"` // Pointer to allow for nil
	OrganizationID primitive.ObjectID  `bson:"organization_id"`
	Name           string              `bson:"name"`
	Content        ProcedureContent    `bson:"content,omitempty"`
	Fees           ProcedureFee        `bson:"fees,omitempty"`
	ProcessingTime ProcessingTime      `bson:"processing_time,omitempty"`
	CreatedAt      time.Time           `bson:"created_at"`
	// For M-M relationship with Notice
	NoticeIDs []primitive.ObjectID `bson:"notice_ids,omitempty"`
}
