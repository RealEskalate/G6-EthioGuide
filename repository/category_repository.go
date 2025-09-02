package repository

import (
	"EthioGuide/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryModel struct {
	ID             primitive.ObjectID `bson:"_id, omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id"`
	ParentID       primitive.ObjectID `bson:"parent_id,omitempty"`
	Title          string             `bson:"title"`
}

func fromDomainCategory(c *domain.Category) (*CategoryModel, error) {
	orgID, err := primitive.ObjectIDFromHex(c.OrganizationID)
	if err != nil {
		return nil, err
	}

	var parentID primitive.ObjectID
	if c.ParentID != "" {
		parentID, err = primitive.ObjectIDFromHex(c.ParentID)
		if err != nil {
			return nil, err
		}
	}

	return &CategoryModel{
		ID:             primitive.NewObjectID(),
		OrganizationID: orgID,
		ParentID:       parentID,
		Title:          c.Title,
	}, nil
}

func toDomainCategory(m *CategoryModel) *domain.Category {
	return &domain.Category{
		ID:             m.ID.Hex(),
		OrganizationID: m.OrganizationID.Hex(),
		ParentID:       m.ParentID.Hex(),
		Title:          m.Title,
	}
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database, collectionName string) *categoryRepository {
	return &categoryRepository{
		collection: db.Collection("Group"),
	}
}

func (r *categoryRepository) GetCategories(ctx context.Context, opts *domain.CategorySearchAndFilter) ([]*domain.Category, int64, error) {
	// 1. Make the filter
	filter := buildFilter(opts)

	// 2. Get the  total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 3. Find the documents
	findOptions := options.Find()
	findOptions.SetLimit(opts.Limit)
	findOptions.SetSkip((opts.Page - 1) * opts.Limit)

	sortValue := -1 // Default to DESC
	if opts.SortOrder == domain.SortAsc {
		sortValue = 1
	}

	findOptions.SetSort(bson.D{{Key: "title", Value: sortValue}})

	// 4. Execute the find query.
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 5. Decode the results.
	var categories []*domain.Category
	for cursor.Next(ctx) {
		var model CategoryModel
		if err := cursor.Decode(&model); err != nil {
			return nil, 0, err
		}
		categories = append(categories, toDomainCategory(&model))
	}

	return categories, total, cursor.Err()
}

// --- helper ---
func buildFilter(options *domain.CategorySearchAndFilter) bson.M {
	var conditions []bson.M

	if options.Title != "" {
		conditions = append(conditions, bson.M{"title": bson.M{"$regex": options.Title, "$options": "i"}})
	}

	if options.ParentID != "" {
		if parentID, err := primitive.ObjectIDFromHex(options.ParentID); err == nil {
			conditions = append(conditions, bson.M{"parent_id": parentID})
		}
	}

	if options.OrganizationID != "" {
		if orgID, err := primitive.ObjectIDFromHex(options.OrganizationID); err == nil {
			conditions = append(conditions, bson.M{"organization_id": orgID})
		}
	}

	if len(conditions) == 0 {
		return bson.M{}
	}

	return bson.M{"$and": conditions}
}
