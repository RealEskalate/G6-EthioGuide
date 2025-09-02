package repository

import (
	"EthioGuide/domain"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database, collectionName string) *categoryRepository {
	return &categoryRepository{
		collection: db.Collection("Group"),
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	model, err := fromDomainCategory(category)
	if err != nil {
		return fmt.Errorf("failed to map domain category to model: %w", err)
	}

	model.ID = primitive.NewObjectID()

	_, err = r.collection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}

	category.ID = model.ID.Hex()
	return nil
}