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
	repo       *ProcedureRepository
	collection *mongo.Collection
}

// SetupSuite connects to the database initialized in TestMain.
func (s *ProcedureRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewProcedureRepository(s.db)
	s.collection = s.db.Collection("procedures")
}

// TearDownSuite drops the database after all tests in this suite are done.
func (s *ProcedureRepositoryTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())
	s.Require().NoError(err)
}

// BeforeTest cleans the collection to ensure test isolation.
func (s *ProcedureRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{})
	s.Require().NoError(err)
}

// TestProcedureRepositoryTestSuite is the entry point for the suite.
func TestProcedureRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(ProcedureRepositoryTestSuite))
}

// --- Helper to create a procedure for tests ---
func (s *ProcedureRepositoryTestSuite) createTestProcedure() *domain.Procedure {
	ctx := context.Background()
	orgID := primitive.NewObjectID()
	procedure := &domain.Procedure{
		OrganizationID: orgID.Hex(),
		Name:           "Apply for Passport",
		Content: domain.ProcedureContent{
			Prerequisites: []string{"Birth Certificate", "ID Card"},
			Steps:         map[int]string{1: "Fill form", 2: "Submit documents"},
			Result:        "Receive Passport",
		},
		Fees: domain.ProcedureFee{
			Label:    "Standard Fee",
			Currency: "ETB",
			Amount:   600,
		},
	}
	err := s.repo.Create(ctx, procedure)
	s.Require().NoError(err, "Setup: failed to create test procedure")
	s.Require().NotEmpty(procedure.ID)
	return procedure
}

// --- Test Cases ---

func (s *ProcedureRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	orgID := primitive.NewObjectID()
	procedure := &domain.Procedure{
		OrganizationID: orgID.Hex(),
		Name:           "Test Procedure Creation",
		Content: domain.ProcedureContent{
			Steps: map[int]string{1: "Step 1"},
		},
	}

	err := s.repo.Create(ctx, procedure)

	s.NoError(err)
	s.NotEmpty(procedure.ID, "Domain procedure ID should be back-filled after creation")

	// Verify directly in the DB
	var createdModel ProcedureModel
	objID, _ := primitive.ObjectIDFromHex(procedure.ID)
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdModel)
	s.NoError(err)
	s.Equal("Test Procedure Creation", createdModel.Name)
	s.Equal(orgID, createdModel.OrganizationID)
	s.WithinDuration(time.Now(), createdModel.CreatedAt, 5*time.Second)
}

func (s *ProcedureRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	// Arrange: Create a procedure to fetch
	createdProcedure := s.createTestProcedure()

	s.Run("Success", func() {
		foundProcedure, err := s.repo.GetByID(ctx, createdProcedure.ID)
		s.NoError(err)
		s.NotNil(foundProcedure)
		s.Equal(createdProcedure.ID, foundProcedure.ID)
		s.Equal("Apply for Passport", foundProcedure.Name)
		s.Equal(600.0, foundProcedure.Fees.Amount)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		foundProcedure, err := s.repo.GetByID(ctx, nonExistentID)
		s.Error(err)
		s.Nil(foundProcedure)
		s.ErrorIs(err, domain.ErrNotFound)
	})

	s.Run("Failure - Invalid ID Format", func() {
		foundProcedure, err := s.repo.GetByID(ctx, "this-is-not-an-object-id")
		s.Error(err)
		s.Nil(foundProcedure)
	})
}

func (s *ProcedureRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	// Arrange: Create a procedure to update
	originalProcedure := s.createTestProcedure()

	s.Run("Success", func() {
		// Create an updated version of the domain object
		updatedProcedure := *originalProcedure
		updatedProcedure.Name = "Updated Passport Application"
		updatedProcedure.Fees.Amount = 750.50
		updatedProcedure.Content.Steps = map[int]string{1: "New Step 1", 2: "New Step 2"}

		err := s.repo.Update(ctx, originalProcedure.ID, &updatedProcedure)
		s.NoError(err)

		// Verify by fetching from the DB
		var modelInDB ProcedureModel
		objID, _ := primitive.ObjectIDFromHex(originalProcedure.ID)
		err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&modelInDB)
		s.NoError(err)
		s.Equal("Updated Passport Application", modelInDB.Name)
		s.Equal(750.50, modelInDB.Fees.Amount)
		s.Equal(map[int]string{1: "New Step 1", 2: "New Step 2"}, modelInDB.Content.Steps)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		procedureToUpdate := domain.Procedure{Name: "Won't be saved"}
		err := s.repo.Update(ctx, nonExistentID, &procedureToUpdate)
		s.Error(err) // The driver might not return an error if no doc is matched, depends on config
		// A better verification is to check the update count, but for simplicity,
		// we can fetch to ensure nothing was created or changed.
		count, err := s.collection.CountDocuments(ctx, bson.M{"name": "Won't be saved"})
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *ProcedureRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	// Arrange: Create a procedure to delete
	procedureToDelete := s.createTestProcedure()
	objID, _ := primitive.ObjectIDFromHex(procedureToDelete.ID)

	// Ensure it exists before deletion
	countBefore, err := s.collection.CountDocuments(ctx, bson.M{"_id": objID})
	s.Require().NoError(err)
	s.Require().Equal(int64(1), countBefore)

	s.Run("Success", func() {
		err := s.repo.Delete(ctx, procedureToDelete.ID)
		s.NoError(err)

		// Verify it's gone from the DB
		countAfter, err := s.collection.CountDocuments(ctx, bson.M{"_id": objID})
		s.NoError(err)
		s.Equal(int64(0), countAfter)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		err := s.repo.Delete(ctx, nonExistentID)
		s.ErrorIs(err, domain.ErrNotFound)
	})
}
