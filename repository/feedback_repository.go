package repository

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedbackModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"user_id"`
	ProcedureID   primitive.ObjectID `bson:"procedure_id"`
	Content       string             `bson:"content"`
	LikeCount     int                `bson:"like_count"`
	DislikeCount  int                `bson:"dislike_count"`
	ViewCount     int                `bson:"view_count"`
	Type          string             `bson:"type"`                     // "Suggestion", "Bug Report", "Praise"
	Status        string             `bson:"status"`                   // "New", "Pending", "In Progress", "Addressed"
	AdminResponse *string            `bson:"admin_response,omitempty"` // Pointer to allow for nil
	Tags          []string           `bson:"tags,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func toDomainFeedback(f *FeedbackModel) *domain.Feedback {
	return &domain.Feedback{
		ID:           f.ID.Hex(),
		UserID:       f.UserID.Hex(),
		ProcedureID:  f.ProcedureID.Hex(),
		Content:      f.Content,
		LikeCount:    f.LikeCount,
		DislikeCount: f.DislikeCount,
		ViewCount:    f.ViewCount,
		Type:         domain.FeedbackType(f.Type),
		Status:       domain.FeedbackStatus(f.Status),
		AdminResponse: func() string {
			if f.AdminResponse != nil {
				return *f.AdminResponse
			}
			return ""
		}(),
		Tags:      f.Tags,
		CreatedAT: f.CreatedAt,
		UpdatedAT: f.UpdatedAt,
	}
}

func fromDomainFeedback(f *domain.Feedback) *FeedbackModel {
	var adminResponse *string
	if f.AdminResponse != "" {
		adminResponse = &f.AdminResponse
	}
	uID, err := primitive.ObjectIDFromHex(f.UserID)
	if err != nil {
		return nil
	}
	pID, err := primitive.ObjectIDFromHex(f.ProcedureID)
	if err != nil {
		return nil
	}
	return &FeedbackModel{
		ID: func() primitive.ObjectID {
			if f.ID != "" {
				id, _ := primitive.ObjectIDFromHex(f.ID)
				return id
			}
			return primitive.NilObjectID
		}(),
		UserID:        uID,
		ProcedureID:   pID,
		Content:       f.Content,
		LikeCount:     f.LikeCount,
		DislikeCount:  f.DislikeCount,
		ViewCount:     f.ViewCount,
		Type:          string(f.Type),
		Status:        string(f.Status),
		AdminResponse: adminResponse,
		Tags:          f.Tags,
		CreatedAt:     f.CreatedAT,
		UpdatedAt:     f.UpdatedAT,
	}
}

type FeedbackRepository struct {
	collection *mongo.Collection
}

func NewFeedbackRepository(db *mongo.Database) domain.IFeedbackRepository {
	return &FeedbackRepository{
		collection: db.Collection("feedbacks"),
	}
}

func (r *FeedbackRepository) SubmitFeedback(ctx context.Context, feedback *domain.Feedback) error {
	model := fromDomainFeedback(feedback)
	model.LikeCount = 0
	model.DislikeCount = 0
	model.ViewCount = 0
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	model.Status = string(domain.NewFeedback)
	model.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to insert feedback: %w", err)
	}
	feedback.ID = model.ID.Hex()
	feedback.Status = domain.NewFeedback
	feedback.CreatedAT = model.CreatedAt
	feedback.UpdatedAT = model.UpdatedAt
	return nil
}

func (r *FeedbackRepository) GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	pid, err := primitive.ObjectIDFromHex(procedureID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid procedure ID: %w", err)
	}

	f := []bson.M{bson.M{"procedure_id": pid}}
	if filter.Status != nil {
		f = append(f, bson.M{"status": *filter.Status})
	}

	filters := bson.M{"$and": f}

	total, err := r.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}

	opts := options.Find()
	opts.SetSkip((filter.Page - 1) * filter.Limit)
	opts.SetLimit(filter.Limit)
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first

	cursor, err := r.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch feedbacks: %w", err)
	}

	var feedbacks []*FeedbackModel
	if err := cursor.All(ctx, &feedbacks); err != nil {
		return nil, 0, fmt.Errorf("failed to decode feedbacks: %w", err)
	}

	domainFeedbacks := make([]*domain.Feedback, len(feedbacks))
	for i, f := range feedbacks {
		domainFeedbacks[i] = toDomainFeedback(f)
	}

	return domainFeedbacks, total, nil
}

func (r *FeedbackRepository) GetFeedbackByID(ctx context.Context, id string) (*domain.Feedback, error) {
	fid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid feedback ID: %w", err)
	}

	var model FeedbackModel
	err = r.collection.FindOne(ctx, bson.M{"_id": fid}).Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch feedback: %w", err)
	}

	return toDomainFeedback(&model), nil
}

func (r *FeedbackRepository) UpdateFeedbackStatus(ctx context.Context, feedbackID string, newFeedback *domain.Feedback) error {
	model := fromDomainFeedback(newFeedback)

	if model == nil {
		return domain.ErrNotFound
	}

	fid, err := primitive.ObjectIDFromHex(feedbackID)
	if err != nil {
		return fmt.Errorf("invalid feedback ID: %w", err)
	}
	model.ID = fid
	model.UpdatedAt = time.Now()

	filter := bson.M{"_id": model.ID}
	update := bson.M{"$set": model}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *FeedbackRepository) GetAllFeedbacks(ctx context.Context, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	f := []bson.M{}
	if filter.Status != nil {
		f = append(f, bson.M{"status": *filter.Status})
	}
	if filter.ProcedureID != nil {
		pid, err := primitive.ObjectIDFromHex(*filter.ProcedureID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid procedure ID: %w", err)
		}
		f = append(f, bson.M{"procedure_id": pid})
	}

	filters := bson.M{}
	if len(f) > 0 {
		filters = bson.M{"$and": f}
	}

	total, err := r.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}
	opts := options.Find()
	opts.SetSkip((filter.Page - 1) * filter.Limit)
	opts.SetLimit(filter.Limit)
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch feedbacks: %w", err)
	}
	var feedbacks []*FeedbackModel
	if err := cursor.All(ctx, &feedbacks); err != nil {
		return nil, 0, fmt.Errorf("failed to decode feedbacks: %w", err)
	}
	domainFeedbacks := make([]*domain.Feedback, len(feedbacks))
	for i, f := range feedbacks {
		domainFeedbacks[i] = toDomainFeedback(f)
	}
	return domainFeedbacks, total, nil
}
