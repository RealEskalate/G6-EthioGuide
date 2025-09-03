package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProcedureContentModel struct {
	Prerequisites []string       `bson:"prerequisites,omitempty"`
	Steps         map[int]string `bson:"steps"`
	Result        string         `bson:"result,omitempty"`
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
	GroupID        *primitive.ObjectID   `bson:"group_id,omitempty"`
	OrganizationID primitive.ObjectID    `bson:"organization_id"`
	Name           string                `bson:"name"`
	Content        ProcedureContentModel `bson:"content"`
	Fees           ProcedureFeeModel     `bson:"fees,omitempty"`
	ProcessingTime ProcessingTimeModel   `bson:"processing_time,omitempty"`
	CreatedAt      time.Time             `bson:"created_at"`
	// For M-M relationship with Notice
	NoticeIDs []primitive.ObjectID `bson:"notice_ids,omitempty"`
}

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
		Content: domain.ProcedureContent{
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
	var id, orgid primitive.ObjectID
	id, _ = primitive.ObjectIDFromHex(proc.ID)
	orgid, _ = primitive.ObjectIDFromHex(proc.OrganizationID)
	var groupID *primitive.ObjectID
	if proc.GroupID != nil {
		gID, err := primitive.ObjectIDFromHex(*proc.GroupID)
		if err == nil {
			groupID = &gID
		}
	}
	noticeIDs := make([]primitive.ObjectID, 0, len(proc.NoticeIDs))
	if proc.NoticeIDs != nil {
		for _, idStr := range proc.NoticeIDs {
			objID, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				noticeIDs = append(noticeIDs, objID)
			}
		}
	}
	return &ProcedureModel{
		ID:             id,
		GroupID:        groupID,
		OrganizationID: orgid,
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
		NoticeIDs: noticeIDs,
	}
}

func ToUpdateBSON(proc *domain.Procedure) bson.M {
	update := bson.M{
		"name": proc.Name,
		"content": ProcedureContentModel{
			Prerequisites: proc.Content.Prerequisites,
			Steps:         proc.Content.Steps,
			Result:        proc.Content.Result,
		},
		"fees": ProcedureFeeModel{
			Label:    proc.Fees.Label,
			Currency: proc.Fees.Currency,
			Amount:   proc.Fees.Amount,
		},
		"processing_time": ProcessingTimeModel{
			MinDays: proc.ProcessingTime.MinDays,
			MaxDays: proc.ProcessingTime.MaxDays,
		},
	}

	// Handle optional GroupID
	if proc.GroupID != nil {
		groupID, err := primitive.ObjectIDFromHex(*proc.GroupID)
		if err == nil {
			update["group_id"] = groupID
		}

	} else {
		update["group_id"] = nil
	}

	// Handle optional NoticeIDs
	if proc.NoticeIDs != nil {
		noticeIDs := make([]primitive.ObjectID, len(proc.NoticeIDs))
		for i, idStr := range proc.NoticeIDs {
			objID, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				noticeIDs[i] = objID
			}
		}
		update["notice_ids"] = noticeIDs
	}

	return update
}

type ProcedureRepository struct {
	db *mongo.Collection
}

func NewProcedureRepository(db *mongo.Collection) *ProcedureRepository {
	return &ProcedureRepository{
		db: db,
	}
}

func (pr *ProcedureRepository) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var procedure ProcedureModel
	err = pr.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&procedure)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return procedure.ToDomain(), nil
}

func (pr *ProcedureRepository) Update(ctx context.Context, id string, procedure *domain.Procedure) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotFound
	}

	updateData := ToUpdateBSON(procedure)

	result, err := pr.db.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (pr *ProcedureRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotFound
	}

	result, err := pr.db.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
