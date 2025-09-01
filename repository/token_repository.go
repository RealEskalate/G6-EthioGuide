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

// This JTI is the unique identifier of the token itself.
type TokenModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	JTI       string             `bson:"jti"` // JTI stands for JWT ID
	Token     string             `bson:"token"`
	TokenType domain.TokenType   `bson:"tokenType"`
	ExpiresAt time.Time          `bson:"expiresAt"`
}

// Mapper from the database model to the domain model.
func toDomainToken(model *TokenModel) *domain.Token {
	return &domain.Token{
		// Note: We are setting the domain 'Id' to the token's JTI, not the DB _id.
		// This is correct as the use case logic relies on the JTI.
		Id:        model.JTI,
		Token:     model.Token,
		TokenType: model.TokenType,
		ExpiresAt: model.ExpiresAt,
	}
}

// Mapper from the domain model to the database model.
func fromDomainToken(domainToken *domain.Token) (*TokenModel, error) {
	// The domainToken.Id is the JTI. We don't try to parse it as an ObjectID.
	// The database _id will be generated automatically by MongoDB upon insertion.
	if domainToken.Id == "" {
		return nil, fmt.Errorf("domain token JTI (Id) cannot be empty")
	}
	return &TokenModel{
		JTI:       domainToken.Id,
		Token:     domainToken.Token,
		TokenType: domainToken.TokenType,
		ExpiresAt: domainToken.ExpiresAt,
	}, nil
}

type TokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository(db *mongo.Database) domain.ITokenRepository {
	coll := db.Collection("tokens")
	return &TokenRepository{
		collection: coll,
	}
}

func (tr *TokenRepository) CreateToken(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	model, err := fromDomainToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to map domain token to model: %w", err)
	}

	res, err := tr.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("failed to insert token: %w", err)
	}

	// IMPROVEMENT: Return a complete domain object mapped from the newly created record.
	// This isn't strictly necessary if you don't use the DB _id, but it's a robust pattern.
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to get inserted ID from database")
	}
	model.ID = insertedID

	return toDomainToken(model), nil
}

func (tr *TokenRepository) GetToken(ctx context.Context, tokentype, tokenID string) (string, error) {
	var tokenGet TokenModel
	filter := bson.M{"tokenType": tokentype, "jti": tokenID}

	err := tr.collection.FindOne(ctx, filter).Decode(&tokenGet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", domain.ErrNotFound // Return a domain-specific error
		}
		return "", err
	}

	return tokenGet.Token, nil
}

func (tr *TokenRepository) DeleteToken(ctx context.Context, tokentype, tokenID string) error {
	filter := bson.M{"tokenType": tokentype, "jti": tokenID}

	res, err := tr.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Check if a document was actually deleted. If not, the token didn't exist.
	if res.DeletedCount == 0 {
		return domain.ErrNotFound // Return a domain-specific error
	}

	return nil
}
