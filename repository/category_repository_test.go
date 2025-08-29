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