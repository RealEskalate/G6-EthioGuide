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

type FeedbackRepositoryTestSuite struct {
    suite.Suite
    db         *mongo.Database
    repo       domain.IFeedbackRepository
    collection *mongo.Collection
}

func (s *FeedbackRepositoryTestSuite) SetupSuite() {
    s.db = testDBClient.Database(testDBName)
    s.repo = NewFeedbackRepository(s.db)
    s.collection = s.db.Collection("feedbacks")
}

func (s *FeedbackRepositoryTestSuite) TearDownSuite() {
    err := s.db.Drop(context.Background())
    s.Require().NoError(err)
}

func (s *FeedbackRepositoryTestSuite) BeforeTest(suiteName, testName string) {
    _, err := s.collection.DeleteMany(context.Background(), bson.M{})
    s.Require().NoError(err)
}

func TestFeedbackRepositoryTestSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode.")
    }
    suite.Run(t, new(FeedbackRepositoryTestSuite))
}

func (s *FeedbackRepositoryTestSuite) TestSubmitFeedback() {
    ctx := context.Background()
    feedback := &domain.Feedback{
        UserID:      primitive.NewObjectID().Hex(),
        ProcedureID: primitive.NewObjectID().Hex(),
        Content:     "Test feedback",
        Type:        domain.ThanksFeedback,
        Tags:        []string{"tag1"},
    }

    err := s.repo.SubmitFeedback(ctx, feedback)
    s.NoError(err)
    s.NotEmpty(feedback.ID, "Domain feedback ID should be back-filled after creation")

    // Verify directly in the DB
    var createdModel FeedbackModel
    objID, _ := primitive.ObjectIDFromHex(feedback.ID)
    err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdModel)
    s.NoError(err)
    s.Equal("Test feedback", createdModel.Content)
    s.Equal("thanks", createdModel.Type)
}

func (s *FeedbackRepositoryTestSuite) TestGetFeedbackByID() {
    ctx := context.Background()
    feedback := &domain.Feedback{
        UserID:      primitive.NewObjectID().Hex(),
        ProcedureID: primitive.NewObjectID().Hex(),
        Content:     "Find me",
        Type:        domain.InaccuracyFeedback,
        Tags:        []string{"tag2"},
    }
    err := s.repo.SubmitFeedback(ctx, feedback)
    s.Require().NoError(err)

    s.Run("Success", func() {
        found, err := s.repo.GetFeedbackByID(ctx, feedback.ID)
        s.NoError(err)
        s.NotNil(found)
        s.Equal(feedback.ID, found.ID)
        s.Equal("Find me", found.Content)
    })

    s.Run("Failure - Not Found", func() {
        nonExistentID := primitive.NewObjectID().Hex()
        found, err := s.repo.GetFeedbackByID(ctx, nonExistentID)
        s.Error(err)
        s.Nil(found)
        s.ErrorIs(err, domain.ErrNotFound)
    })

    s.Run("Failure - Invalid ID Format", func() {
        found, err := s.repo.GetFeedbackByID(ctx, "not-an-objectid")
        s.Error(err)
        s.Nil(found)
    })
}

func (s *FeedbackRepositoryTestSuite) TestGetAllFeedbacksForProcedure() {
    ctx := context.Background()
    procedureID := primitive.NewObjectID().Hex()
    feedback := &domain.Feedback{
        UserID:      primitive.NewObjectID().Hex(),
        ProcedureID: procedureID,
        Content:     "List me",
        Type:        domain.OutdatedFeedback,
        Tags:        []string{"tag3"},
    }
    err := s.repo.SubmitFeedback(ctx, feedback)
    s.Require().NoError(err)

    filter := &domain.FeedbackFilter{Page: 1, Limit: 10}
    s.Run("Success", func() {
        list, total, err := s.repo.GetAllFeedbacksForProcedure(ctx, procedureID, filter)
        s.NoError(err)
        s.Equal(int64(1), total)
        s.Len(list, 1)
        s.Equal("List me", list[0].Content)
    })

    s.Run("Failure - Invalid ProcedureID", func() {
        list, total, err := s.repo.GetAllFeedbacksForProcedure(ctx, "not-an-objectid", filter)
        s.Error(err)
        s.Nil(list)
        s.Equal(int64(0), total)
    })
}