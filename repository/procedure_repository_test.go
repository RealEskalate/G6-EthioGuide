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

type ProcedureRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repo       domain.IProcedureRepository
	collection *mongo.Collection
}

func (s *ProcedureRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewProcedureRepository(s.db)
	s.collection = s.db.Collection("procedures")
}

func (s *ProcedureRepositoryTestSuite) TearDownSuite() {
	s.db.Drop(context.Background())
}

func (s *ProcedureRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	s.collection.DeleteMany(context.Background(), bson.M{})
}

func TestProcedureRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(ProcedureRepositoryTestSuite))
}

func (s *ProcedureRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID()
	proc := &domain.Procedure{
		OrganizationID: orgID.Hex(),
		Name:           "Test Procedure",
	}

	err := s.repo.Create(ctx, proc)
	s.NoError(err)
	s.NotEmpty(proc.ID)

	// Verify in DB
	var model ProcedureModel
	objID, _ := primitive.ObjectIDFromHex(proc.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	s.NoError(err)
	s.Equal("Test Procedure", model.Name)
}

func (s *ProcedureRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	// Arrange
	model := &ProcedureModel{ID: primitive.NewObjectID(), OrganizationID: primitive.NewObjectID(), Name: "Find Me"}
	_, err := s.collection.InsertOne(ctx, model)
	s.Require().NoError(err)

	s.Run("Success", func() {
		proc, err := s.repo.GetByID(ctx, model.ID.Hex())
		s.NoError(err)
		s.NotNil(proc)
		s.Equal("Find Me", proc.Name)
	})

	s.Run("Failure - Not Found", func() {
		_, err := s.repo.GetByID(ctx, primitive.NewObjectID().Hex())
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *ProcedureRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	// Arrange
	model := &ProcedureModel{ID: primitive.NewObjectID(), OrganizationID: primitive.NewObjectID(), Name: "Original"}
	_, err := s.collection.InsertOne(ctx, model)
	s.Require().NoError(err)

	s.Run("Success", func() {
		updateDomain := &domain.Procedure{Name: "Updated"}
		err := s.repo.Update(ctx, model.ID.Hex(), updateDomain)
		s.NoError(err)

		var updatedModel ProcedureModel
		s.collection.FindOne(ctx, bson.M{"_id": model.ID}).Decode(&updatedModel)
		s.Equal("Updated", updatedModel.Name)
	})

	s.Run("Failure - Not Found", func() {
		err := s.repo.Update(ctx, primitive.NewObjectID().Hex(), &domain.Procedure{})
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *ProcedureRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	// Arrange
	model := &ProcedureModel{ID: primitive.NewObjectID()}
	_, err := s.collection.InsertOne(ctx, model)
	s.Require().NoError(err)

	s.Run("Success", func() {
		err := s.repo.Delete(ctx, model.ID.Hex())
		s.NoError(err)
		count, _ := s.collection.CountDocuments(ctx, bson.M{"_id": model.ID})
		s.Equal(int64(0), count)
	})

	s.Run("Failure - Not Found", func() {
		err := s.repo.Delete(ctx, primitive.NewObjectID().Hex())
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *ProcedureRepositoryTestSuite) TestSearchAndFilter() {
	ctx := context.Background()
	orgID1 := primitive.NewObjectID()
	orgID2 := primitive.NewObjectID()
	groupID1 := primitive.NewObjectID()

	// Arrange: Create procedures with different properties
	procs := []interface{}{
		ProcedureModel{Name: "Passport Renewal", OrganizationID: orgID1, GroupID: &groupID1, Fees: ProcedureFeeModel{Amount: 100}, ProcessingTime: ProcessingTimeModel{MinDays: 10}, CreatedAt: time.Now().Add(-2 * time.Hour)},
		ProcedureModel{Name: "Driver's License", OrganizationID: orgID2, Fees: ProcedureFeeModel{Amount: 200}, ProcessingTime: ProcessingTimeModel{MinDays: 20}, CreatedAt: time.Now().Add(-1 * time.Hour)},
		ProcedureModel{Name: "Passport Application", OrganizationID: orgID1, Fees: ProcedureFeeModel{Amount: 150}, ProcessingTime: ProcessingTimeModel{MinDays: 15}, CreatedAt: time.Now()},
	}
	_, err := s.collection.InsertMany(ctx, procs)
	s.Require().NoError(err)

	s.Run("Filter by Name", func() {
		name := "Passport"
		opts := domain.ProcedureSearchFilterOptions{Name: &name, Page: 1, Limit: 10}
		results, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total)
		s.Len(results, 2)
	})

	s.Run("Filter by OrganizationID", func() {
		orgIDHex := orgID2.Hex()
		opts := domain.ProcedureSearchFilterOptions{OrganizationID: &orgIDHex, Page: 1, Limit: 10}
		results, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total)
		s.Equal("Driver's License", results[0].Name)
	})

	s.Run("Filter by GroupID", func() {
		groupIDHex := groupID1.Hex()
		opts := domain.ProcedureSearchFilterOptions{GroupID: &groupIDHex, Page: 1, Limit: 10}
		results, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(1), total)
		s.Equal("Passport Renewal", results[0].Name)
	})

	s.Run("Filter by Fee Range", func() {
		minFee := 120.0
		maxFee := 210.0
		opts := domain.ProcedureSearchFilterOptions{MinFee: &minFee, MaxFee: &maxFee, Page: 1, Limit: 10}
		_, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total) // Should find 150 and 200
	})

	s.Run("Filter by Processing Time", func() {
		minDays := 12
		opts := domain.ProcedureSearchFilterOptions{MinProcessingDays: &minDays, Page: 1, Limit: 10}
		_, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total) // Should find 15 and 20
	})

	s.Run("Filter by Date Range", func() {
		startDate := time.Now().Add(-90 * time.Minute) // 1.5 hours ago
		opts := domain.ProcedureSearchFilterOptions{StartDate: &startDate, Page: 1, Limit: 10}
		_, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(2), total) // Should find the ones created now and 1 hour ago
	})

	s.Run("Pagination and Sorting", func() {
		opts := domain.ProcedureSearchFilterOptions{Page: 2, Limit: 1, SortBy: "created_at", SortOrder: domain.SortAsc}
		results, total, err := s.repo.SearchAndFilter(ctx, opts)
		s.NoError(err)
		s.Equal(int64(3), total)
		s.Len(results, 1)
		s.Equal("Driver's License", results[0].Name) // Second oldest
	})
}
