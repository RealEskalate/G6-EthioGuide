package repository

import (
	"EthioGuide/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoticeModel struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	Title          string             `bson:"title"`
	Content        string             `bson:"content"`
	Tags           []string           `bson:"tags,omitempty"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

func (nm *NoticeModel) ToDomain() *domain.Notice {
	return &domain.Notice{
		ID:             nm.ID.Hex(),
		OrganizationID: nm.OrganizationID.Hex(),
		Title:          nm.Title,
		Content:        nm.Content,
		Tags:           nm.Tags,
		CreatedAt:      nm.CreatedAt,
		UpdatedAt:      nm.UpdatedAt,
	}
}

func FromDomain(notice *domain.Notice) *NoticeModel {
	var id primitive.ObjectID
	if notice.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(notice.ID)
		if err == nil {
			id = objectID
		}
	}
	var orgID primitive.ObjectID
	if notice.OrganizationID != "" {
		objectID, err := primitive.ObjectIDFromHex(notice.OrganizationID)
		if err == nil {
			orgID = objectID
		}
	}
	return &NoticeModel{
		ID:             id,
		OrganizationID: orgID,
		Title:          notice.Title,
		Content:        notice.Content,
		Tags:           notice.Tags,
		CreatedAt:      notice.CreatedAt,
		UpdatedAt:      notice.UpdatedAt,
	}
}

type NoticeRepository struct {
	collection *mongo.Collection
}

func NewNoticeRepository(db *mongo.Database) *NoticeRepository {
	return &NoticeRepository{
		collection: db.Collection("notices"),
	}
}

func (nr *NoticeRepository) Create(ctx context.Context, notice *domain.Notice) error {
	model := FromDomain(notice)
	model.CreatedAt = time.Now()
	model.ID = primitive.NewObjectID()

	_, err := nr.collection.InsertOne(ctx, model)

	notice.ID = model.ID.Hex()

	return err
}

func (nr *NoticeRepository) GetByFilter(ctx context.Context, filter *domain.NoticeFilter) ([]*domain.Notice, error) {
	// Build filter
	var conditions []bson.M

	if filter.OrganizationID != "" {
		if oid, err := primitive.ObjectIDFromHex(filter.OrganizationID); err == nil {
			conditions = append(conditions, bson.M{"organization_id": oid})
		}
	}

	if len(filter.Tags) > 0 {
		// AND logic: all tags must be present
		conditions = append(conditions, bson.M{"tags": bson.M{"$all": filter.Tags}})
	}

	var mongoFilter bson.M
	if len(conditions) == 0 {
		mongoFilter = bson.M{}
	} else {
		mongoFilter = bson.M{"$and": conditions}
	}

	// Pagination defaults and bounds
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}

	// Sorting: only created_at supported via sortBy=created_At
	sortValue := -1 // DESC default
	if filter.SortOrder == "ASC" {
		sortValue = 1
	}
	// Map sortBy to DB field (only created_At supported)
	sortField := "created_at"
	// If provided but not "created_At", we still default to created_at

	findOpts := options.Find().
		SetLimit(limit).
		SetSkip((page - 1) * limit).
		SetSort(bson.D{{Key: sortField, Value: sortValue}})

	cur, err := nr.collection.Find(ctx, mongoFilter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var notices []*domain.Notice
	for cur.Next(ctx) {
		var m NoticeModel
		if err := cur.Decode(&m); err != nil {
			return nil, err
		}
		notices = append(notices, m.ToDomain())
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return notices, nil
}

func (nr *NoticeRepository) Update(ctx context.Context, id string, notice *domain.Notice) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	model := FromDomain(notice)
	_, err = nr.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": model})
	return err
}

func (nr *NoticeRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = nr.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// CountByFilter returns the total number of notices matching the filter (ignores pagination)
func (nr *NoticeRepository) CountByFilter(ctx context.Context, filter *domain.NoticeFilter) (int64, error) {
	var conditions []bson.M

	if filter.OrganizationID != "" {
		if oid, err := primitive.ObjectIDFromHex(filter.OrganizationID); err == nil {
			conditions = append(conditions, bson.M{"organization_id": oid})
		}
	}

	if len(filter.Tags) > 0 {
		conditions = append(conditions, bson.M{"tags": bson.M{"$all": filter.Tags}})
	}

	var mongoFilter bson.M
	if len(conditions) == 0 {
		mongoFilter = bson.M{}
	} else {
		mongoFilter = bson.M{"$and": conditions}
	}

	count, err := nr.collection.CountDocuments(ctx, mongoFilter)
	return count, err
}
