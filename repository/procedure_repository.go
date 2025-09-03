package repository

import (
	"EthioGuide/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProcedureContentModel struct {
	Prerequisites []string `bson:"prerequisites,omitempty"`
	Steps         []string `bson:"steps,omitempty"`
	Result        string   `bson:"result,omitempty"`
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

// package repository

// import (
// 	"EthioGuide/domain"
// 	"context"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type ProcedureContentModel struct {
// 	Prerequisites []string `bson:"prerequisites,omitempty"`
// 	Steps         map[int]string `bson:"steps"`
// 	Result        string   `bson:"result,omitempty"`
// }

// // ProcedureFee is a nested struct for the Procedure.fees field.
// type ProcedureFeeModel struct {
// 	Label    string  `bson:"label,omitempty"`
// 	Currency string  `bson:"currency,omitempty"`
// 	Amount   float64 `bson:"amount,omitempty"`
// }

// // ProcessingTime is a nested struct for the Procedure.processing_time field.
// type ProcessingTimeModel struct {
// 	MinDays int `bson:"min_days,omitempty"`
// 	MaxDays int `bson:"max_days,omitempty"`
// }

// type ProcedureModel struct {
// 	ID             primitive.ObjectID    `bson:"_id,omitempty"`
// 	GroupID        *primitive.ObjectID   `bson:"group_id,omitempty"` 
// 	OrganizationID primitive.ObjectID    `bson:"organization_id"`
// 	Name           string                `bson:"name"`
// 	Content        ProcedureContentModel `bson:"content"`
// 	Fees           ProcedureFeeModel     `bson:"fees,omitempty"`
// 	ProcessingTime ProcessingTimeModel   `bson:"processing_time,omitempty"`
// 	CreatedAt      time.Time             `bson:"created_at"`
// 	// For M-M relationship with Notice
// 	NoticeIDs []primitive.ObjectID 		 `bson:"notice_ids,omitempty"`
// }


func (pm *ProcedureModel) ToDomain() *domain.Procedure {
	var groupID *string
	if pm.GroupID != nil {
		id := pm.GroupID.Hex()
		groupID = &id
	}
	// Convert []primitive.ObjectID to []string
	noticeIDs := make([]string, len(pm.NoticeIDs))
	for i, id := range pm.NoticeIDs {
		noticeIDs[i] = id.Hex()
	}
	return &domain.Procedure{
		ID:             pm.ID.Hex(),
		GroupID:        groupID,
		OrganizationID: pm.OrganizationID.Hex(),
		Name:           pm.Name,
		Content:        domain.ProcedureContent{
			Prerequisites: pm.Content.Prerequisites,
			Steps:         pm.Content.Steps,
			Result:        pm.Content.Result,
		},
		Fees: domain.ProcedureFee{
			Label:    pm.Fees.Label,
			Currency: pm.Fees.Currency,
			Amount:   pm.Fees.Amount,
		},
		ProcessingTime: domain.ProcessingTime{
			MinDays: pm.ProcessingTime.MinDays,
			MaxDays: pm.ProcessingTime.MaxDays,
		},
		CreatedAt: pm.CreatedAt,
		NoticeIDs: noticeIDs,
	}
}

func ToDTO(proc *domain.Procedure) *ProcedureModel {
	return &ProcedureModel{
		ID:             primitive.NewObjectID(),
		GroupID:        nil,
		OrganizationID: primitive.NewObjectID(),
		Name:           proc.Name,
		Content: ProcedureContentModel{
			Prerequisites: proc.Content.Prerequisites,
			Steps:         proc.Content.Steps,
			Result:        proc.Content.Result,
		},
		Fees: ProcedureFeeModel{
			Label:    proc.Fees.Label,
			Currency: proc.Fees.Currency,
			Amount:   proc.Fees.Amount,
		},
		ProcessingTime: ProcessingTimeModel{
			MinDays: proc.ProcessingTime.MinDays,
			MaxDays: proc.ProcessingTime.MaxDays,
		},
		CreatedAt: proc.CreatedAt,
		NoticeIDs: nil,
	}
}


type ProcedureRepository struct {
	db *mongo.Collection
}

func NewProcedureRepository(db *mongo.Database) *ProcedureRepository {
	return &ProcedureRepository{
		db: db.Collection("procedures"),
	}
}

func (pr *ProcedureRepository) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	var procedure ProcedureModel
	err := pr.db.FindOne(ctx, bson.M{"_id": id}).Decode(&procedure)
	if err != nil {
		return nil, err
	}
	return procedure.ToDomain(), nil
}

func (pr *ProcedureRepository) Update(ctx context.Context, id string, procedure *domain.Procedure) error {
	_, err := pr.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": ToDTO(procedure)})
	return err
}

func (pr *ProcedureRepository) Delete(ctx context.Context, id string) error {
	_, err := pr.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}