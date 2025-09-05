package repository

import (
	"EthioGuide/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PreferencesRepository struct {
    collection *mongo.Collection
}

func NewPreferencesRepository(db *mongo.Database) domain.IPreferencesRepository{
    return &PreferencesRepository{
        collection: db.Collection("preferences"),
    }
}

// DTO for MongoDB
type PreferencesDTO struct {
    ID                primitive.ObjectID `bson:"_id,omitempty"`
    UserID            string             `bson:"user_id"`
    PreferredLang     string             `bson:"preferredLang"`
    PushNotification  bool               `bson:"pushNotification"`
    EmailNotification bool               `bson:"emailNotification"`
}

func dtoFromDomain(p *domain.Preferences) *PreferencesDTO {
    var oid primitive.ObjectID
    if p.ID != "" {
        oid, _ = primitive.ObjectIDFromHex(p.ID)
    }
    return &PreferencesDTO{
        ID:                oid,
        UserID:            p.UserID,
        PreferredLang:     string(p.PreferredLang),
        PushNotification:  p.PushNotification,
        EmailNotification: p.EmailNotification,
    }
}

func (dto *PreferencesDTO) ToDomain() *domain.Preferences {
    return &domain.Preferences{
        ID:                dto.ID.Hex(),
        UserID:            dto.UserID,
        PreferredLang:     domain.Lang(dto.PreferredLang),
        PushNotification:  dto.PushNotification,
        EmailNotification: dto.EmailNotification,
    }
}

func (r *PreferencesRepository) Create(ctx context.Context, preferences *domain.Preferences) error {
    dto := dtoFromDomain(preferences)
    res, err := r.collection.InsertOne(ctx, dto)
    if err != nil {
        return err
    }
    if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
        preferences.ID = oid.Hex()
    }
    return nil
}

func (r *PreferencesRepository) GetByUserID(ctx context.Context, userID string) (*domain.Preferences, error) {
    filter := bson.M{"user_id": userID}
    var dto PreferencesDTO
    err := r.collection.FindOne(ctx, filter).Decode(&dto)
    if err != nil {
        return nil, err
    }
    return dto.ToDomain(), nil
}

func (r *PreferencesRepository) UpdateByUserID(ctx context.Context, userID string, preferences *domain.Preferences) error {
    filter := bson.M{"user_id": userID}
    update := bson.M{
        "$set": dtoFromDomain(preferences),
    }
    _, err := r.collection.UpdateOne(ctx, filter, update)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrUserNotFound
	}
    return err
}