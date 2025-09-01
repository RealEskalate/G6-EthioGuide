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

type AIUserChatRepositoryMongo struct {
    collection *mongo.Collection
}

func NewAIUserChatRepositoryMongo(db *mongo.Database) *AIUserChatRepositoryMongo {
    return &AIUserChatRepositoryMongo{
        collection: db.Collection("ai_chats"),
    }
}

func (r *AIUserChatRepositoryMongo) Save(ctx context.Context, chat *domain.AIChat) error {
	chat.Timestamp = time.Now()
	dto, err := DomainAIChatToMessageModel(chat)
	if err != nil {
		return err
	}
	_, err = r.collection.InsertOne(ctx, dto)
	return err
}

func (r *AIUserChatRepositoryMongo) GetByUser(ctx context.Context, userID string, limit int) ([]*domain.AIChat, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"user_id": userObjID}
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dtos []*AIUserChatMessageModel
	if err := cursor.All(ctx, &dtos); err != nil {
		return nil, err
	}
	var chats []*domain.AIChat
	for _, dto := range dtos {
		chats = append(chats, AIUserChatMessageModelToDomain(dto))
	}
	return chats, nil
}

func (r *AIUserChatRepositoryMongo) DeleteByUser(ctx context.Context, userID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteMany(ctx, bson.M{"user_id": userObjID})
	return err
}



////////-----------------------Dto----------------/////////////////
////////////////==============================///////////////////////


type AIUserChatMessageModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Source    string             `bson:"source,omitempty"`
	Request   string             `bson:"request"`
	Response  string             `bson:"response"`
	Timestamp time.Time          `bson:"timestamp"`
}

// Convert AIUserChatMessageModel (DTO) to domain.AIChat
func AIUserChatMessageModelToDomain(dto *AIUserChatMessageModel) *domain.AIChat {
	return &domain.AIChat{
		ID:        dto.ID.Hex(),
		UserID:    dto.UserID.Hex(),
		Source:    dto.Source,
		Request:   dto.Request,
		Response:  dto.Response,
		Timestamp: dto.Timestamp,
	}
}

// Convert domain.AIChat to AIUserChatMessageModel (DTO)
func DomainAIChatToMessageModel(chat *domain.AIChat) (*AIUserChatMessageModel, error) {
	var id primitive.ObjectID
	var userID primitive.ObjectID
	var err error
	if chat.ID != "" {
		id, err = primitive.ObjectIDFromHex(chat.ID)
		if err != nil {
			return nil, err
		}
	}
	if chat.UserID != "" {
		userID, err = primitive.ObjectIDFromHex(chat.UserID)
		if err != nil {
			return nil, err
		}
	}
	return &AIUserChatMessageModel{
		ID:        id,
		UserID:    userID,
		Source:    chat.Source,
		Request:   chat.Request,
		Response:  chat.Response,
		Timestamp: chat.Timestamp,
	}, nil
}