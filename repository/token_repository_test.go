package repository_test

import (
	"EthioGuide/domain"
	. "EthioGuide/repository" // Dot import for convenience
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.ITokenRepository
	collection *mongo.Collection
}

// SetupSuite runs once for the entire test suite.
func (s *TokenRepositoryTestSuite) SetupSuite() {
	// The testDBClient is already connected from TestMain.
	// We just need to get a database handle.
	s.db = testDBClient.Database(testDBName)
	s.repo = NewTokenRepository(s.db)
	s.collection = s.db.Collection("tokens") // Get a handle for direct checks
}

// TearDownSuite runs once after all tests in the suite have finished.
func (s *TokenRepositoryTestSuite) TearDownSuite() {
	// Clean up the database after the suite is done.
	err := s.db.Drop(context.Background())
	s.Require().NoError(err, "Failed to drop test database")
}

// BeforeTest runs before each individual test function.
func (s *TokenRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	// Clean the collection before each test to ensure isolation.
	_, err := s.collection.DeleteMany(context.Background(), bson.M{})
	s.Require().NoError(err, "Failed to clean up collection before test")
}

// TestTokenRepositoryTestSuite is the entry point for running the suite.
func TestTokenRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(TokenRepositoryTestSuite))
}

// --- Test Cases ---

func (s *TokenRepositoryTestSuite) TestCreateToken() {
	ctx := context.Background()
	token := &domain.Token{
		Id:        "jti-123", // JTI
		Token:     "a.very.long.jwt.string",
		TokenType: domain.RefreshToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	s.Run("Success", func() {
		createdToken, err := s.repo.CreateToken(ctx, token)
		s.NoError(err)
		s.NotNil(createdToken)
		s.Equal(token.Id, createdToken.Id, "JTI should match")

		// Verify directly in the DB
		count, err := s.collection.CountDocuments(ctx, bson.M{"jti": "jti-123"})
		s.NoError(err)
		s.Equal(int64(1), count, "Token should be found in the database")
	})

	s.Run("Failure - Empty JTI", func() {
		invalidToken := &domain.Token{Id: ""} // Empty JTI (Id)
		_, err := s.repo.CreateToken(ctx, invalidToken)
		s.Error(err, "Should return an error for an empty JTI")
	})
}

func (s *TokenRepositoryTestSuite) TestGetToken() {
	ctx := context.Background()
	token := &domain.Token{
		Id:        "jti-456",
		Token:     "another.jwt.string",
		TokenType: domain.RefreshToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	_, err := s.repo.CreateToken(ctx, token)
	s.Require().NoError(err)

	s.Run("Success", func() {
		foundTokenStr, err := s.repo.GetToken(ctx, string(domain.RefreshToken), "jti-456")
		s.NoError(err)
		s.Equal(token.Token, foundTokenStr, "The long JWT string should be returned")
	})

	s.Run("Failure - Not Found", func() {
		_, err := s.repo.GetToken(ctx, string(domain.RefreshToken), "non-existent-jti")
		s.Error(err)
		s.ErrorIs(err, domain.ErrNotFound, "Should return a domain-specific not found error")
	})
}

func (s *TokenRepositoryTestSuite) TestDeleteToken() {
	ctx := context.Background()
	token := &domain.Token{
		Id:        "jti-789",
		Token:     "yet.another.jwt.string",
		TokenType: domain.RefreshToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	_, err := s.repo.CreateToken(ctx, token)
	s.Require().NoError(err)

	s.Run("Success", func() {
		err := s.repo.DeleteToken(ctx, string(domain.RefreshToken), "jti-789")
		s.NoError(err)

		// Verify it's gone
		_, getErr := s.repo.GetToken(ctx, string(domain.RefreshToken), "jti-789")
		s.ErrorIs(getErr, domain.ErrNotFound)
	})

	s.Run("Failure - Not Found", func() {
		// Try to delete a token that doesn't exist
		err := s.repo.DeleteToken(ctx, string(domain.RefreshToken), "non-existent-jti")
		s.Error(err)
		s.ErrorIs(err, domain.ErrNotFound, "Should return not found error when deleting a non-existent token")
	})
}
