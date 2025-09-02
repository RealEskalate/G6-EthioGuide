package repository

import (
	"EthioGuide/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	Content      string             `bson:"content"`
	Title        string             `bson:"title"`
	// LikeCount    int                `bson:"like_count"`
	// DislikeCount int                `bson:"dislike_count"`
	// ViewCount    int                `bson:"view_count"`
	Tags         []string           `bson:"tags,omitempty"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	// For M-M relationship with Procedure, replacing PostProcedure table
	ProcedureIDs []primitive.ObjectID `bson:"procedure_ids,omitempty"`
}

type IPostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *IPostRepository {
	coll := db.Collection("posts")
	return &IPostRepository{
		collection: coll,
	}
}	

func (r *IPostRepository) CreatePost(ctx context.Context, post *domain.Post) error {
	id, err := primitive.ObjectIDFromHex(post.UserID)
	if err != nil {
		return err
	}

	procedureIDs, err := mapHexToObjectID(post.Procedures)
	if err != nil {
		return err
	}	
	_, err = r.collection.InsertOne(ctx, PostModel{	
		UserID:   id,
		Content: post.Content,
		Title: post.Title,
		Tags : post.Tags,
		ProcedureIDs: procedureIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return domain.ErrUnableToEnterData
	}
	return nil
}


func mapHexToObjectID(hexIDs []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID
	for _, hexID := range hexIDs {
		objID, err := primitive.ObjectIDFromHex(hexID)
		if err != nil {
			return nil, err
		}	
		objectIDs = append(objectIDs, objID)
	}
	return objectIDs, nil
}
