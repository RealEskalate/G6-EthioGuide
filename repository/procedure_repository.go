package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Embedding      []float64             `bson:"embedding,omitempty"`
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
		Embedding: pm.Embedding,
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
		Embedding: proc.Embedding,
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
		"embedding": proc.Embedding,
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
		for _, idStr := range proc.NoticeIDs {
			objID, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				noticeIDs = append(noticeIDs, objID)
			}
		}
		update["notice_ids"] = noticeIDs
	}

	return update
}

type ProcedureRepository struct {
	col *mongo.Collection
}

func NewProcedureRepository(db *mongo.Database) domain.IProcedureRepository {
	return &ProcedureRepository{
		col: db.Collection("procedures"),
	}
}

func (r *ProcedureRepository) Create(ctx context.Context, procedure *domain.Procedure) error {
	model := ToDTO(procedure)

	model.CreatedAt = time.Now()
	model.ID = primitive.NewObjectID()

	_, err := r.col.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	procedure.ID = model.ID.Hex()
	return nil
}

func (pr *ProcedureRepository) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var procedure ProcedureModel
	err = pr.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&procedure)
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

	result, err := pr.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
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

	result, err := pr.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// SearchAndFilter searches procedures based on the options and returns paginated results.
func (r *ProcedureRepository) SearchAndFilter(ctx context.Context, opts domain.ProcedureSearchFilterOptions) ([]*domain.Procedure, int64, error) {
	filter := bson.M{}

	// --- Search by Name (case-insensitive regex) ---
	if opts.Name != nil && *opts.Name != "" {
		filter["name"] = bson.M{"$regex": *opts.Name, "$options": "i"}
	}

	// --- OrganizationID ---
	if opts.OrganizationID != nil {
		if *opts.OrganizationID == "" {
			filter["organization_id"] = bson.M{"$exists": false}
		} else {
			oid, err := primitive.ObjectIDFromHex(*opts.OrganizationID)
			if err == nil {
				filter["organization_id"] = oid
			}
		}
	}

	// --- GroupID ---
	if opts.GroupID != nil {
		if *opts.GroupID == "" {
			filter["group_id"] = bson.M{"$exists": false}
		} else {
			gid, err := primitive.ObjectIDFromHex(*opts.GroupID)
			if err == nil {
				filter["group_id"] = gid
			}
		}
	}

	// --- Fee range ---
	if opts.MinFee != nil || opts.MaxFee != nil {
		feeFilter := bson.M{}
		if opts.MinFee != nil {
			feeFilter["$gte"] = *opts.MinFee
		}
		if opts.MaxFee != nil {
			feeFilter["$lte"] = *opts.MaxFee
		}
		filter["fees.amount"] = feeFilter
	}

	// --- Processing time ---
	if opts.MinProcessingDays != nil || opts.MaxProcessingDays != nil {
		timeFilter := bson.M{}
		if opts.MinProcessingDays != nil {
			timeFilter["$gte"] = *opts.MinProcessingDays
		}
		if opts.MaxProcessingDays != nil {
			timeFilter["$lte"] = *opts.MaxProcessingDays
		}
		filter["processing_time.min_days"] = timeFilter
	}

	// --- Date range ---
	if opts.StartDate != nil || opts.EndDate != nil {
		dateFilter := bson.M{}
		if opts.StartDate != nil {
			dateFilter["$gte"] = *opts.StartDate
		}
		if opts.EndDate != nil {
			dateFilter["$lte"] = *opts.EndDate
		}
		filter["created_at"] = dateFilter
	}

	// --- Count total first ---
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// --- Sorting ---
	findOpts := options.Find()
	if opts.SortBy != "" {
		order := -1 // DESC by default
		if opts.SortOrder == domain.SortAsc {
			order = 1
		}
		findOpts.SetSort(bson.D{{Key: opts.SortBy, Value: order}})
	}

	// --- Pagination ---
	skip := (opts.Page - 1) * opts.Limit
	findOpts.SetSkip(skip)
	findOpts.SetLimit(opts.Limit)

	// --- Query ---
	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var models []*ProcedureModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, 0, err
	}

	// --- Convert to domain ---
	results := make([]*domain.Procedure, len(models))
	for i, m := range models {
		results[i] = m.ToDomain()
	}

	return results, total, nil
}

func (r *ProcedureRepository) SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*domain.Procedure, error) {
	// Vector search stage
	searchStage := bson.D{
		{Key: "$vectorSearch", Value: bson.D{
			{Key: "index", Value: "vector_index"}, // must match your Atlas vector index name
			{Key: "path", Value: "embedding"},     // the field that stores vectors
			{Key: "queryVector", Value: queryVec}, // the query embedding
			{Key: "numCandidates", Value: 100},    // candidate pool before top-k filtering
			{Key: "limit", Value: limit},          // top-k results
		}},
	}

	// // You can add a $project stage if you want only specific fields
	// projectStage := bson.D{
	//     {"$project", bson.D{
	//         {"title", 1},
	//         {"requirements", 1},
	//         {"steps", 1},
	//         {"fees", 1},
	//         {"score", bson.D{{"$meta", "vectorSearchScore"}}}, // optional score
	//     }},
	// }

	// pipeline := mongo.Pipeline{searchStage, projectStage}
	pipeline := mongo.Pipeline{searchStage}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []*ProcedureModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	// --- Convert to domain ---
	results := make([]*domain.Procedure, len(models))
	for i, m := range models {
		results[i] = m.ToDomain()
	}
	return results, nil
}
