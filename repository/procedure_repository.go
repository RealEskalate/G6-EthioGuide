// CreateEmbeddingVectorIndex creates a vector index on the 'embedding' field for vector search.
package repository

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	embedding      []float64             `bson:"embedding,omitempty"`
	Fees           ProcedureFeeModel     `bson:"fees,omitempty"`
	ProcessingTime ProcessingTimeModel   `bson:"processing_time,omitempty"`
	CreatedAt      time.Time             `bson:"created_at"`
	// For M-M relationship with Notice
	NoticeIDs []primitive.ObjectID `bson:"notice_ids,omitempty"`
}

// --- mappers ---

func fromDomainProcedure(p *domain.Procedure) (*ProcedureModel, error) {
	groupID, err := primitive.ObjectIDFromHex(p.GroupID)
	if err != nil && p.GroupID != "" {
		return nil, err
	}
	orgID, err := primitive.ObjectIDFromHex(p.OrganizationID)
	if err != nil {
		return nil, err
	}
	return &ProcedureModel{
		Name:           p.Name,
		GroupID:        &groupID,
		OrganizationID: orgID,
		Content: ProcedureContentModel{
			Prerequisites: p.Content.Prerequisites,
			Steps:         p.Content.Steps,
			Result:        p.Content.Result,
		},
		Fees: ProcedureFeeModel{
			Label:    p.Fees.Label,
			Currency: p.Fees.Currency,
			Amount:   p.Fees.Amount,
		},
		ProcessingTime: ProcessingTimeModel{
			MinDays: p.ProcessingTime.MinDays,
			MaxDays: p.ProcessingTime.MaxDays,
		},
	}, nil
}

// --- implementation ---

type procedureRepository struct {
	collection *mongo.Collection
}

func NewProcedureRepository(db *mongo.Database) domain.IProcedureRepository {
	return &procedureRepository{
		collection: db.Collection("procedures"),
	}
}

func (r *procedureRepository) Create(ctx context.Context, procedure *domain.Procedure) error {
	model, err := fromDomainProcedure(procedure)
	if err != nil {
		return fmt.Errorf("failed to map domain procedure to model: %w", err)
	}

	model.CreatedAt = time.Now()
	model.ID = primitive.NewObjectID()

	_, err = r.collection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	procedure.ID = model.ID.Hex()
	return nil
}

func (r *procedureRepository) SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*domain.Procedure, error) {
	// Vector search stage
	searchStage := bson.D{
		{"$vectorSearch", bson.D{
			{"index", "vector_index"}, // must match your Atlas vector index name
			{"path", "embedding"},     // the field that stores vectors
			{"queryVector", queryVec}, // the query embedding
			{"numCandidates", 100},    // candidate pool before top-k filtering
			{"limit", limit},          // top-k results
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

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Procedure
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
