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

type PostModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  primitive.ObjectID `bson:"user_id"`
	Content string             `bson:"content"`
	Title   string             `bson:"title"`
	// LikeCount    int                `bson:"like_count"`
	// DislikeCount int                `bson:"dislike_count"`
	// ViewCount    int                `bson:"view_count"`
	Tags      []string  `bson:"tags,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
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

func (r *IPostRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	id, err := primitive.ObjectIDFromHex(post.UserID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	procedureIDs, err := mapHexToObjectID(post.Procedures)
	if err != nil {
		return nil, err
	}
	res, err := r.collection.InsertOne(ctx, PostModel{
		UserID:       id,
		Content:      post.Content,
		Title:        post.Title,
		Tags:         post.Tags,
		ProcedureIDs: procedureIDs,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return nil, domain.ErrUnableToEnterData
	}
	var pm PostModel
	_ = r.collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&pm)
	Newpost := &domain.Post{
		ID:         pm.ID.Hex(),
		UserID:     pm.UserID.Hex(),
		Title:      pm.Title,
		Content:    pm.Content,
		Procedures: mapObjectIDtoHex(pm.ProcedureIDs),
		Tags:       pm.Tags,
		CreatedAt:  pm.CreatedAt,
		UpdatedAt:  pm.UpdatedAt,
	}
	return Newpost, nil
}
func buildPostFilter(opts domain.PostFilters) (bson.M, error) {
	filter := bson.M{}

	if opts.Title != nil && *opts.Title != "" {
		filter["title"] = bson.M{"$regex": *opts.Title, "$options": "i"}
	}

	if opts.UserId != nil && *opts.UserId != "" {
		userIdObj, err := primitive.ObjectIDFromHex(*opts.UserId)
		if err == nil {
			filter["user_id"] = userIdObj
		}
	}

	if len(opts.ProcedureID) > 0 {
		procedureIDs, err := mapHexToObjectID(opts.ProcedureID)
		if err != nil {
			return nil, err
		}
		filter["procedure_ids"] = bson.M{"$all": procedureIDs}
	}

	if len(opts.Tags) > 0 {
		filter["tags"] = bson.M{"$all": opts.Tags}
	}

	return filter, nil
}

func (r *IPostRepository) GetPosts(ctx context.Context, opts domain.PostFilters) ([]*domain.Post, int64, error) {

	filter, err := buildPostFilter(opts)
	if err != nil {
		return nil, 0, err
	}
	skip := opts.Page * opts.Limit
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(opts.Limit)

	if opts.SortBy != "" {
		sortOrder := 1 // ASC
		if opts.SortOrder == domain.SortDesc {
			sortOrder = -1
		}
		findOptions.Sort = bson.M{opts.SortBy: sortOrder}
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, domain.ErrUnableToFetchData
	}
	defer cursor.Close(ctx)

	var postModels []PostModel
	if err := cursor.All(ctx, &postModels); err != nil {
		return nil, 0, domain.ErrUnableToFetchData
	}

	var posts []*domain.Post
	for _, pm := range postModels {
		post := &domain.Post{
			ID:         pm.ID.Hex(),
			UserID:     pm.UserID.Hex(),
			Title:      pm.Title,
			Content:    pm.Content,
			Procedures: mapObjectIDtoHex(pm.ProcedureIDs),
			Tags:       pm.Tags,
			CreatedAt:  pm.CreatedAt,
			UpdatedAt:  pm.UpdatedAt,
		}
		posts = append(posts, post)
	}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, domain.ErrUnableToFetchData
	}

	return posts, totalCount, nil
}

func (r *IPostRepository) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidIDFormat
	}
	var pm PostModel
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&pm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrPostNotFound
		}
		return nil, domain.ErrUnableToFetchData
	}
	post := &domain.Post{
		UserID:     pm.UserID.Hex(),
		Title:      pm.Title,
		Content:    pm.Content,
		Procedures: mapObjectIDtoHex(pm.ProcedureIDs),
		Tags:       pm.Tags,
		CreatedAt:  pm.CreatedAt,
		UpdatedAt:  pm.UpdatedAt,
	}
	return post, nil
}

func (r *IPostRepository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	oid, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return nil, domain.ErrInvalidIDFormat
	}

	filter := bson.M{"_id": oid}
	var pm PostModel
	err = r.collection.FindOne(ctx, filter).Decode(&pm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrPostNotFound
		}
		return nil, domain.ErrUnableToFetchData
	}
	//check if user is the author of the post
	if pm.UserID.Hex() != post.UserID {
		return nil, domain.ErrPermissionDenied
	}

	updates := bson.M{
		"title":        post.Title,
		"content":      post.Content,
		"procedureIDs": post.Procedures,
		"tags":         post.Tags,
		"updatedAt":    time.Now(),
	}

	_, err = r.collection.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return nil, domain.ErrUnableToUpdateData
	}

	err = r.collection.FindOne(ctx, filter).Decode(&pm)
	if err != nil {
		return nil, domain.ErrUnableToFetchData
	}
	updatedPost := &domain.Post{
		ID:         pm.ID.Hex(),
		UserID:     pm.UserID.Hex(),
		Title:      pm.Title,
		Content:    pm.Content,
		Procedures: mapObjectIDtoHex(pm.ProcedureIDs),
		Tags:       pm.Tags,
		CreatedAt:  pm.CreatedAt,
		UpdatedAt:  pm.UpdatedAt,
	}
	return updatedPost, nil
}

func (r *IPostRepository) DeletePost(ctx context.Context, id, userID, role string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidIDFormat
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrInvalidIDFormat
	}
	var pm PostModel
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&pm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ErrPostNotFound
		}
		return domain.ErrUnableToFetchData
	}
	if pm.UserID != uid && role != "admin" {
		return domain.ErrPermissionDenied
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return domain.ErrUnableToDeleteData
	}
	if res.DeletedCount == 0 {
		return domain.ErrPostNotFound
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

func mapObjectIDtoHex(objectIDs []primitive.ObjectID) []string {
	var hexIDs []string
	for _, objID := range objectIDs {
		hexIDs = append(hexIDs, objID.Hex())
	}
	return hexIDs
}
