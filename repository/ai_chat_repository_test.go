package repository_test

import (
	"EthioGuide/domain"
	. "EthioGuide/repository"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AIChatRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.IAIChatRepository
	collection *mongo.Collection
}

func (s *AIChatRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewAIChatRepository(s.db)
	s.collection = s.db.Collection("ai_chats")
}

func (s *AIChatRepositoryTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())
	s.Require().NoError(err)
}

func (s *AIChatRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{})
	s.Require().NoError(err)
}

func TestAIChatRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(AIChatRepositoryTestSuite))
}

func (s *AIChatRepositoryTestSuite) TestSave() {
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	procID := primitive.NewObjectID().Hex()

	chat := &domain.AIChat{
		UserID:   userID,
		Request:  "Test Request",
		Response: "Test Response",
		RelatedProcedures: []*domain.AIProcedure{
			{Id: procID, Name: "Test Proc"},
		},
	}

	err := s.repo.Save(ctx, chat)

	s.NoError(err)
	s.NotEmpty(chat.ID, "Domain chat ID should be back-filled after saving")

	// Verify directly in the DB
	var createdModel AIUserChatMessageModel
	objID, _ := primitive.ObjectIDFromHex(chat.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdModel)
	s.NoError(err)
	s.Equal(chat.Request, createdModel.Request)
	s.Equal(chat.Response, createdModel.Response)
	s.Len(createdModel.RelatedProcedures, 1)
	s.Equal("Test Proc", createdModel.RelatedProcedures[0].Name)
}

func (s *AIChatRepositoryTestSuite) TestGetByUser() {
	ctx := context.Background()
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Arrange: Create some chat messages for two different users
	chatsUser1 := []interface{}{
		AIUserChatMessageModel{UserID: userID1, Request: "Q1", Timestamp: time.Now().Add(-2 * time.Minute)},
		AIUserChatMessageModel{UserID: userID1, Request: "Q2", Timestamp: time.Now().Add(-1 * time.Minute)},
		AIUserChatMessageModel{UserID: userID1, Request: "Q3", Timestamp: time.Now()},
	}
	chatsUser2 := []interface{}{
		AIUserChatMessageModel{UserID: userID2, Request: "Other user question"},
	}
	_, err := s.collection.InsertMany(ctx, chatsUser1)
	s.Require().NoError(err)
	_, err = s.collection.InsertMany(ctx, chatsUser2)
	s.Require().NoError(err)

	s.Run("Success - Get all for user with pagination", func() {
		// Get first page (2 items)
		chats, total, err := s.repo.GetByUser(ctx, userID1.Hex(), 1, 2)
		s.NoError(err)
		s.Equal(int64(3), total)
		s.Len(chats, 2)
		// Check sorting (newest first)
		s.Equal("Q3", chats[0].Request)
		s.Equal("Q2", chats[1].Request)

		// Get second page (1 item)
		chats, total, err = s.repo.GetByUser(ctx, userID1.Hex(), 2, 2)
		s.NoError(err)
		s.Equal(int64(3), total)
		s.Len(chats, 1)
		s.Equal("Q1", chats[0].Request)
	})

	s.Run("Success - No chats for user", func() {
		nonExistentUserID := primitive.NewObjectID().Hex()
		chats, total, err := s.repo.GetByUser(ctx, nonExistentUserID, 1, 10)
		s.NoError(err)
		s.Equal(int64(0), total)
		s.Empty(chats)
	})

	s.Run("Failure - Invalid User ID", func() {
		_, _, err := s.repo.GetByUser(ctx, "invalid-id", 1, 10)
		s.Error(err)
	})
}

func (s *AIChatRepositoryTestSuite) TestDeleteByUser() {
	ctx := context.Background()
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Arrange
	chats := []interface{}{
		AIUserChatMessageModel{UserID: userID1, Request: "Q1"},
		AIUserChatMessageModel{UserID: userID1, Request: "Q2"},
		AIUserChatMessageModel{UserID: userID2, Request: "Keep this one"},
	}
	_, err := s.collection.InsertMany(ctx, chats)
	s.Require().NoError(err)

	s.Run("Success - Deletes all for a user", func() {
		err := s.repo.DeleteByUser(ctx, userID1.Hex())
		s.NoError(err)

		// Verify in DB
		countUser1, err := s.collection.CountDocuments(ctx, bson.M{"user_id": userID1})
		s.NoError(err)
		s.Equal(int64(0), countUser1)

		countUser2, err := s.collection.CountDocuments(ctx, bson.M{"user_id": userID2})
		s.NoError(err)
		s.Equal(int64(1), countUser2, "Should not delete other users' chats")
	})

	s.Run("Failure - Invalid User ID", func() {
		err := s.repo.DeleteByUser(ctx, "invalid-id")
		s.Error(err)
	})
}
