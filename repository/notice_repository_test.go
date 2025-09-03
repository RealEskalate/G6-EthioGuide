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
	collection *mongo.Collection
	repo       *NoticeRepository
}

// SetupSuite connects to the database initialized in TestMain (same harness used by other repo tests).
func (s *NoticeRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.collection = s.db.Collection("notices")
	s.repo = NewNoticeController(s.collection)
}

// TearDownSuite drops the database after all tests in this suite are done.
func (s *NoticeRepositoryTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())
	s.Require().NoError(err)
}

// Before each test, clean the collection to ensure isolation.
func (s *NoticeRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{})
	s.Require().NoError(err)
}

func TestNoticeRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(NoticeRepositoryTestSuite))
}

func (s *NoticeRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID().Hex()
	n := &domain.Notice{
		OrganizationID: orgID,
		Title:          "Create Test",
		Content:        "Body",
		Tags:           []string{"a", "b"},
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		UpdatedAt:      time.Now().Add(-1 * time.Hour),
	}

	err := s.repo.Create(ctx, n)
	s.NoError(err)
	s.NotEmpty(n.ID, "ID should be set after creation")

	// Verify in DB
	var model NoticeModel
	objID, _ := primitive.ObjectIDFromHex(n.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	s.NoError(err)
	s.Equal("Create Test", model.Title)
	s.ElementsMatch([]string{"a", "b"}, model.Tags)
}

func (s *NoticeRepositoryTestSuite) TestGetByFilter() {
	ctx := context.Background()
	org1 := primitive.NewObjectID().Hex()
	org2 := primitive.NewObjectID().Hex()

	now := time.Now()
	toCreate := []*domain.Notice{
		{OrganizationID: org1, Title: "A", Tags: []string{"x", "y"}, CreatedAt: now.Add(-120 * time.Hour), UpdatedAt: now.Add(-120 * time.Hour)},
		{OrganizationID: org1, Title: "B", Tags: []string{"x", "z"}, CreatedAt: now.Add(-96 * time.Hour), UpdatedAt: now.Add(-96 * time.Hour)},
		{OrganizationID: org2, Title: "C", Tags: []string{"x", "y", "z"}, CreatedAt: now.Add(-72 * time.Hour), UpdatedAt: now.Add(-72 * time.Hour)},
		{OrganizationID: org2, Title: "D", Tags: []string{"y", "z"}, CreatedAt: now.Add(-48 * time.Hour), UpdatedAt: now.Add(-48 * time.Hour)},
		{OrganizationID: org1, Title: "E", Tags: []string{"x", "y"}, CreatedAt: now.Add(-24 * time.Hour), UpdatedAt: now.Add(-24 * time.Hour)},
	}
	for _, n := range toCreate {
		s.Require().NoError(s.repo.Create(ctx, n))
	}

	s.Run("Tags AND logic, default DESC sort by created_at", func() {
		filter := &domain.NoticeFilter{
			Tags:      []string{"x", "y"},
			Page:      1,
			Limit:     10,
			SortBy:    "created_At",
			SortOrder: "DESC",
		}
		got, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		// Should include A, C, E
		s.Len(got, 3)
		titles := []string{got[0].Title, got[1].Title, got[2].Title}
		s.ElementsMatch([]string{"A", "C", "E"}, titles)

		// Verify non-increasing order by CreatedAt
		for i := 1; i < len(got); i++ {
			s.True(got[i-1].CreatedAt.After(got[i].CreatedAt) || got[i-1].CreatedAt.Equal(got[i].CreatedAt))
		}
	})

	s.Run("ASC sort by created_at + pagination page=2 limit=2", func() {
		filter := &domain.NoticeFilter{
			Page:      2,
			Limit:     2,
			SortBy:    "created_At",
			SortOrder: "ASC",
		}
		got, err := s.repo.GetByFilter(ctx, filter)
		s.NoError(err)
		s.Len(got, 2)
		// Global ASC order by created_at among all: [A,B,C,D,E], so page2,limit2 -> [C,D]
		s.Equal("C", got[0].Title)
		s.Equal("D", got[1].Title)
	})
}

func (s *NoticeRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID().Hex()
	n := &domain.Notice{
		OrganizationID: orgID,
		Title:          "Old Title",
		Content:        "Old Content",
		Tags:           []string{"one"},
		CreatedAt:      time.Now().Add(-2 * time.Hour),
		UpdatedAt:      time.Now().Add(-2 * time.Hour),
	}
	s.Require().NoError(s.repo.Create(ctx, n))
	s.NotEmpty(n.ID)

	// Update fields
	n.Title = "New Title"
	n.Content = "New Content"
	n.Tags = []string{"one", "two"}
	n.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, n.ID, n)
	s.NoError(err)

	// Verify in DB
	var model NoticeModel
	objID, _ := primitive.ObjectIDFromHex(n.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	s.NoError(err)
	s.Equal("New Title", model.Title)
	s.Equal("New Content", model.Content)
	s.ElementsMatch([]string{"one", "two"}, model.Tags)
}

func (s *NoticeRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	orgID := primitive.NewObjectID().Hex()
	n := &domain.Notice{
		OrganizationID: orgID,
		Title:          "To Delete",
		Content:        "X",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	s.Require().NoError(s.repo.Create(ctx, n))
	s.NotEmpty(n.ID)

	err := s.repo.Delete(ctx, n.ID)
	s.NoError(err)

	// Should not be found
	var model NoticeModel
	objID, _ := primitive.ObjectIDFromHex(n.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	s.Error(err)
}
