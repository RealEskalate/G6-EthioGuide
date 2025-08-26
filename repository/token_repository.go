package repository

import (
	"EthioGuide/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenModelDTO struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Token     string             `bson:"token"`
	TokenType string             `bson:"tokenType"`
	ExpiresAt time.Time          `bson:"expiresAt"`
}

func ToDomainToken(tokenmodel *TokenModelDTO) *domain.TokenModel {
	return &domain.TokenModel{
		Id:        tokenmodel.Id.Hex(),
		Token:     tokenmodel.Token,
		TokenType: tokenmodel.TokenType,
		ExpiresAt: tokenmodel.ExpiresAt,
	}
}

func ToDTOToken(tokendomain *domain.TokenModel) *TokenModelDTO {
	dtoId, err := primitive.ObjectIDFromHex(tokendomain.Id)
	if err != nil {
		dtoId = primitive.NilObjectID
	}
	return &TokenModelDTO{
		Id:        dtoId,
		Token:     tokendomain.Token,
		TokenType: tokendomain.TokenType,
		ExpiresAt: tokendomain.ExpiresAt,
	}
}

type TokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository(db *mongo.Database) *TokenRepository {
	coll := db.Collection("refresh_token")
	return &TokenRepository{
		collection: coll,
	}
}

func (tr *TokenRepository) CreateToken(ctx context.Context, token *domain.TokenModel) (*domain.TokenModel, error) {
	newDto := ToDTOToken(token)
	if _, err := tr.collection.InsertOne(ctx, newDto); err != nil {
		return &domain.TokenModel{}, err
	}

	return token, nil
}

func (tr *TokenRepository) GetToken(ctx context.Context, tokentype, token string) (string, error) {
	var tokenGet TokenModelDTO
	filter := bson.M{"tokenType": tokentype, "token": token}

	err := tr.collection.FindOne(ctx, filter).Decode(&tokenGet)
	if err != nil {
		return "", err
	}

	return tokenGet.Token, nil
}

func (tr *TokenRepository) DeleteToken(ctx context.Context, tokentype, token string) error {
	filter := bson.M{"tokenType": tokentype, "token": token}
	if _, err := tr.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
