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

type NoticeRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       *NoticeRepository // Use concrete type for direct access
	collection *mongo.Collection
}

func (s *NoticeRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewNoticeRepository(s.db)
	s.collection = s.db.Collection("notices")
}

func (s *NoticeRepositoryTestSuite) TearDownSuite() {
	s.db.Drop(context.Background())
}

func (s *NoticeRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	s.collection.DeleteMany(context.Background(), bson.M{})
}

func TestNoticeRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(NoticeRepositoryTestSuite))
}

func (s *NoticeRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID()
	notice := &domain.Notice{
		OrganizationID: orgID.Hex(),
		Title:          "New Announcement",
		Content:        "Details here.",
	}

	err := s.repo.Create(ctx, notice)
	s.NoError(err)
	s.NotEmpty(notice.ID)

	// Verify in DB
	var model NoticeModel
	objID, _ := primitive.ObjectIDFromHex(notice.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	s.NoError(err)
	s.Equal("New Announcement", model.Title)
	s.Equal(orgID, model.OrganizationID)
}

func (s *NoticeRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID()
	// Arrange: Create a notice to update
	originalModel := &NoticeModel{
		ID:             primitive.NewObjectID(),
		OrganizationID: orgID,
		Title:          "Original Title",
		Content:        "Original Content",
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		UpdatedAt:      time.Now().Add(-1 * time.Hour),
	}
	_, err := s.collection.InsertOne(ctx, originalModel)
	s.Require().NoError(err)

	s.Run("Success", func() {
		updateDomain := &domain.Notice{
			Title:   "Updated Title",
			Content: "Updated Content",
		}

		err := s.repo.Update(ctx, originalModel.ID.Hex(), updateDomain)
		s.NoError(err)

		// Verify in DB
		var updatedModel NoticeModel
		err = s.collection.FindOne(ctx, bson.M{"_id": originalModel.ID}).Decode(&updatedModel)
		s.NoError(err)
		s.Equal("Updated Title", updatedModel.Title)
		s.Equal("Updated Content", updatedModel.Content)
		s.True(updatedModel.UpdatedAt.After(originalModel.UpdatedAt))
		s.Equal(originalModel.CreatedAt.Unix(), updatedModel.CreatedAt.Unix()) // CreatedAt should not change
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		err := s.repo.Update(ctx, nonExistentID, &domain.Notice{})
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *NoticeRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	// Arrange
	model := &NoticeModel{ID: primitive.NewObjectID()}
	_, err := s.collection.InsertOne(ctx, model)
	s.Require().NoError(err)

	// Act
	err = s.repo.Delete(ctx, model.ID.Hex())
	s.NoError(err)

	// Verify
	count, err := s.collection.CountDocuments(ctx, bson.M{"_id": model.ID})
	s.NoError(err)
	s.Equal(int64(0), count)
}

func (s *NoticeRepositoryTestSuite) TestGetAndCountByFilter() {
	ctx := context.Background()
	orgID1 := primitive.NewObjectID()
	orgID2 := primitive.NewObjectID()

	// Arrange: Create notices
	notices := []interface{}{
		NoticeModel{OrganizationID: orgID1, Title: "A", Tags: []string{"urgent", "update"}, CreatedAt: time.Now().Add(-3 * time.Hour)},
		NoticeModel{OrganizationID: orgID1, Title: "B", Tags: []string{"update"}, CreatedAt: time.Now().Add(-2 * time.Hour)},
		NoticeModel{OrganizationID: orgID2, Title: "C", Tags: []string{"general", "update"}, CreatedAt: time.Now().Add(-1 * time.Hour)},
	}
	_, err := s.collection.InsertMany(ctx, notices)
	s.Require().NoError(err)

	s.Run("Filter by OrganizationID", func() {
		filter := &domain.NoticeFilter{OrganizationID: orgID1.Hex()}
		results, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(results, 2)

		count, err := s.repo.CountByFilter(ctx, filter)
		s.NoError(err)
		s.Equal(int64(2), count)
	})

	s.Run("Filter by single Tag", func() {
		filter := &domain.NoticeFilter{Tags: []string{"urgent"}}
		results, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(results, 1)
		s.Equal("A", results[0].Title)
	})

	s.Run("Filter by multiple Tags ($all)", func() {
		filter := &domain.NoticeFilter{Tags: []string{"urgent", "update"}}
		results, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(results, 1)
		s.Equal("A", results[0].Title)
	})

	s.Run("Filter by common Tag", func() {
		filter := &domain.NoticeFilter{Tags: []string{"update"}}
		results, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(results, 3)

		count, err := s.repo.CountByFilter(ctx, filter)
		s.NoError(err)
		s.Equal(int64(3), count)
	})

	s.Run("Combined Filter", func() {
		filter := &domain.NoticeFilter{OrganizationID: orgID2.Hex(), Tags: []string{"update"}}
		results, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(results, 1)
		s.Equal("C", results[0].Title)
	})

	s.Run("Pagination and Sorting", func() {
		filter := &domain.NoticeFilter{Page: 2, Limit: 1, SortOrder: "ASC"} // Get 2nd item when sorted oldest first
		results, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(results, 1)
		s.Equal("B", results[0].Title)
	})
}
