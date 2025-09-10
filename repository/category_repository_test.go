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

type CategoryRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.ICategoryRepository
	collection *mongo.Collection
}

func (s *CategoryRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewCategoryRepository(s.db, "categories")
	s.collection = s.db.Collection("Group")
}

func (s *CategoryRepositoryTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())
	s.Require().NoError(err)
}

func (s *CategoryRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{})
	s.Require().NoError(err)
}

func TestCategoryRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

func (s *CategoryRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()
	category := &domain.Category{
		OrganizationID: orgID,
		ParentID:       parentID,
		Title:          "Test Category",
	}

	err := s.repo.Create(ctx, category)
	s.NoError(err)
	s.NotEmpty(category.ID, "Category ID should be back-filled after creation")

	// Verify directly in the DB
	var createdModel CategoryModel
	objID, _ := primitive.ObjectIDFromHex(category.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdModel)
	s.NoError(err)
	s.Equal("Test Category", createdModel.Title)
	s.Equal(category.OrganizationID, createdModel.OrganizationID.Hex())
	s.Equal(category.ParentID, createdModel.ParentID.Hex())
}

func (s *CategoryRepositoryTestSuite) TestCreate_InvalidOrgID() {
	ctx := context.Background()
	category := &domain.Category{
		OrganizationID: "badid",
		ParentID:       primitive.NewObjectID().Hex(),
		Title:          "Should Fail",
	}

	err := s.repo.Create(ctx, category)
	s.Error(err)
	s.Contains(err.Error(), "failed to map domain category to model")
}

func (s *CategoryRepositoryTestSuite) TestCreate_InvalidParentID() {
	ctx := context.Background()
	category := &domain.Category{
		OrganizationID: primitive.NewObjectID().Hex(),
		ParentID:       "badid",
		Title:          "Should Fail",
	}

	err := s.repo.Create(ctx, category)
	s.Error(err)
	s.Contains(err.Error(), "failed to map domain category to model")
}

func (s *CategoryRepositoryTestSuite) TestGetCategories() {
	ctx := context.Background()
	orgID1 := primitive.NewObjectID()
	orgID2 := primitive.NewObjectID()
	parentID1 := primitive.NewObjectID()

	// Arrange: Create a set of categories to filter through
	categories := []interface{}{
		CategoryModel{ID: primitive.NewObjectID(), OrganizationID: orgID1, ParentID: primitive.NilObjectID, Title: "Top Level Alpha"},
		CategoryModel{ID: primitive.NewObjectID(), OrganizationID: orgID1, ParentID: parentID1, Title: "Child of Alpha"},
		CategoryModel{ID: primitive.NewObjectID(), OrganizationID: orgID2, ParentID: primitive.NilObjectID, Title: "Top Level Beta"},
		CategoryModel{ID: primitive.NewObjectID(), OrganizationID: orgID2, ParentID: parentID1, Title: "Another Child of Alpha"},
	}
	_, err := s.collection.InsertMany(ctx, categories)
	s.Require().NoError(err)

	s.Run("Success - Filter by Title (finds only top-level)", func() {
		// Logic: No ParentID provided, so it defaults to top-level.
		// It should find "Top Level Alpha" but not "Child of Alpha".
		opts := &domain.CategorySearchAndFilter{Title: "alpha"}
		results, total, err := s.repo.GetCategories(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total) // CORRECTED EXPECTATION
		s.Len(results, 1)
		s.Equal("Top Level Alpha", results[0].Title)
	})

	s.Run("Success - Filter by specific ParentID", func() {
		opts := &domain.CategorySearchAndFilter{ParentID: parentID1.Hex()}
		results, total, err := s.repo.GetCategories(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total)
		s.Len(results, 2)
	})

	s.Run("Success - Filter for Top-Level categories (empty ParentID)", func() {
		opts := &domain.CategorySearchAndFilter{ParentID: ""}
		results, total, err := s.repo.GetCategories(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total)
		s.Len(results, 2)
	})

	s.Run("Success - Filter by OrganizationID (finds only top-level)", func() {
		// Logic: No ParentID provided, so it defaults to top-level.
		// It should find "Top Level Alpha" but not "Child of Alpha".
		opts := &domain.CategorySearchAndFilter{OrganizationID: orgID1.Hex()}
		results, total, err := s.repo.GetCategories(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total) // CORRECTED EXPECTATION
		s.Len(results, 1)
	})

	s.Run("Success - Combined Filters (Parent and Org)", func() {
		opts := &domain.CategorySearchAndFilter{
			ParentID:       parentID1.Hex(),
			OrganizationID: orgID2.Hex(),
		}
		results, total, err := s.repo.GetCategories(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total)
		s.Len(results, 1)
		s.Equal("Another Child of Alpha", results[0].Title)
	})

	s.Run("Success - Pagination and Sorting (finds only top-level)", func() {
		// Logic: No ParentID, so it finds the 2 top-level categories.
		opts := &domain.CategorySearchAndFilter{
			Page:      1,
			Limit:     4,
			SortOrder: domain.SortAsc,
		}
		results, total, err := s.repo.GetCategories(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total) // CORRECTED EXPECTATION
		s.Len(results, 2)        // CORRECTED EXPECTATION
		s.Equal("Top Level Alpha", results[0].Title, "Should be first alphabetically")
		s.Equal("Top Level Beta", results[1].Title, "Should be last alphabetically")
	})
}
