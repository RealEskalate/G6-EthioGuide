package repository_test

import (
	"EthioGuide/domain"
	. "EthioGuide/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PreferencesRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.IPreferencesRepository
	collection *mongo.Collection
}

func (s *PreferencesRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewPreferencesRepository(s.db)
	s.collection = s.db.Collection("preferences")
}

func (s *PreferencesRepositoryTestSuite) TearDownSuite() {
	s.db.Drop(context.Background())
}

func (s *PreferencesRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	s.collection.DeleteMany(context.Background(), bson.M{})
}

func TestPreferencesRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(PreferencesRepositoryTestSuite))
}

func (s *PreferencesRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	preferences := &domain.Preferences{
		UserID:            userID,
		PreferredLang:     domain.English,
		PushNotification:  true,
		EmailNotification: false,
	}

	err := s.repo.Create(ctx, preferences)
	s.NoError(err)
	s.NotEmpty(preferences.ID, "Domain preferences ID should be back-filled after creation")

	// Verify directly in the DB
	var createdModel PreferencesDTO
	objID, _ := primitive.ObjectIDFromHex(preferences.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdModel)
	s.NoError(err)
	s.Equal(userID, createdModel.UserID)
	s.Equal("en", createdModel.PreferredLang)
	s.True(createdModel.PushNotification)
}

func (s *PreferencesRepositoryTestSuite) TestGetByUserID() {
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	// Arrange: Create a document to fetch
	prefToCreate := &domain.Preferences{UserID: userID, PreferredLang: domain.Amharic}
	err := s.repo.Create(ctx, prefToCreate)
	s.Require().NoError(err)

	s.Run("Success", func() {
		found, err := s.repo.GetByUserID(ctx, userID)
		s.NoError(err)
		s.NotNil(found)
		s.Equal(userID, found.UserID)
		s.Equal(domain.Amharic, found.PreferredLang)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentUserID := primitive.NewObjectID().Hex()
		_, err := s.repo.GetByUserID(ctx, nonExistentUserID)
		s.Error(err)
		s.ErrorIs(err, mongo.ErrNoDocuments)
	})
}

func (s *PreferencesRepositoryTestSuite) TestUpdateByUserID() {
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	// Arrange: Create a document to update
	prefToCreate := &domain.Preferences{
		UserID:            userID,
		PreferredLang:     domain.English,
		PushNotification:  false,
		EmailNotification: false,
	}
	err := s.repo.Create(ctx, prefToCreate)
	s.Require().NoError(err)

	s.Run("Success", func() {
		// Prepare the updated domain object
		updatedPrefs := &domain.Preferences{
			UserID:            userID,
			PreferredLang:     domain.Amharic,
			PushNotification:  true,
			EmailNotification: true,
		}

		err := s.repo.UpdateByUserID(ctx, userID, updatedPrefs)
		s.NoError(err)

		// Verify directly in the DB
		var result PreferencesDTO
		err = s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
		s.NoError(err)
		s.Equal("am", result.PreferredLang)
		s.True(result.PushNotification)
		s.True(result.EmailNotification)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentUserID := primitive.NewObjectID().Hex()
		err := s.repo.UpdateByUserID(ctx, nonExistentUserID, &domain.Preferences{UserID: nonExistentUserID})
		s.ErrorIs(err, domain.ErrUserNotFound)
	})
}
