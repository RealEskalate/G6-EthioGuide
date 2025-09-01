package repository_test

import (
	"EthioGuide/domain"
	. "EthioGuide/repository" // Dot import for convenience
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

// SetupSuite connects to the database initialized in TestMain.
func (s *ProcedureRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	s.repo = NewProcedureRepository(s.db.Collection("procedures"))
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

// --- Test Cases ---

func (s *ProcedureRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	
	// Create a test procedure first
	procedure := &domain.Procedure{
		Name:           "Test Procedure",
		OrganizationID: primitive.NewObjectID().Hex(),
		Content: domain.ProcedureContent{
			Prerequisites: []string{"prereq1", "prereq2"},
			Steps:         map[int]string{1: "Step 1", 2: "Step 2"},
			Result:        "Expected result",
		},
		Fees: domain.ProcedureFee{
			Label:    "Service Fee",
			Currency: "USD",
			Amount:   100.50,
		},
		ProcessingTime: domain.ProcessingTime{
			MinDays: 1,
			MaxDays: 3,
		},
		CreatedAt: time.Now(),
		NoticeIDs: []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()},
	}

	// Insert directly to get the MongoDB ID
	procedureModel := ToDTO(procedure)
	result, err := s.collection.InsertOne(ctx, procedureModel)
	s.Require().NoError(err, "Setup: failed to insert procedure")
	
	// Get the inserted ID
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	procedure.ID = insertedID

	s.Run("Success", func() {
		foundProcedure, err := s.repo.GetByID(ctx, insertedID)

		s.NoError(err)
		s.NotNil(foundProcedure)
		s.Equal(insertedID, foundProcedure.ID)
		s.Equal("Test Procedure", foundProcedure.Name)
		s.Equal(procedure.OrganizationID, foundProcedure.OrganizationID)
		s.Equal([]string{"prereq1", "prereq2"}, foundProcedure.Content.Prerequisites)
		s.Equal("Service Fee", foundProcedure.Fees.Label)
		s.Equal(100.50, foundProcedure.Fees.Amount)
		s.Equal(1, foundProcedure.ProcessingTime.MinDays)
		s.Equal(3, foundProcedure.ProcessingTime.MaxDays)
		s.Len(foundProcedure.NoticeIDs, 2)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		foundProcedure, err := s.repo.GetByID(ctx, nonExistentID)

		s.Error(err)
		s.Nil(foundProcedure)
		s.ErrorIs(err, domain.ErrNotFound, "Should return a domain-specific not found error")
	})

	s.Run("Failure - Invalid ID Format", func() {
		foundProcedure, err := s.repo.GetByID(ctx, "this-is-not-an-object-id")

		s.Error(err)
		s.Nil(foundProcedure)
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *ProcedureRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	
	// Create a test procedure first
	procedure := &domain.Procedure{
		Name:           "Original Procedure",
		OrganizationID: primitive.NewObjectID().Hex(),
		Content: domain.ProcedureContent{
			Prerequisites: []string{"original prereq"},
			Steps:         map[int]string{1: "Original step"},
			Result:        "Original result",
		},
		Fees: domain.ProcedureFee{
			Label:    "Original Fee",
			Currency: "USD",
			Amount:   50.00,
		},
		ProcessingTime: domain.ProcessingTime{
			MinDays: 1,
			MaxDays: 2,
		},
		CreatedAt: time.Now(),
	}

	// Insert directly to get the MongoDB ID
	procedureModel := ToDTO(procedure)
	result, err := s.collection.InsertOne(ctx, procedureModel)
	s.Require().NoError(err, "Setup: failed to insert procedure")
	
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	s.Run("Success", func() {
		// Prepare updated procedure
		updatedProcedure := &domain.Procedure{
			ID:             insertedID,
			Name:           "Updated Procedure",
			OrganizationID: procedure.OrganizationID,
			Content: domain.ProcedureContent{
				Prerequisites: []string{"updated prereq1", "updated prereq2"},
				Steps:         map[int]string{1: "Updated step 1", 2: "Updated step 2"},
				Result:        "Updated result",
			},
			Fees: domain.ProcedureFee{
				Label:    "Updated Fee",
				Currency: "ETB",
				Amount:   75.25,
			},
			ProcessingTime: domain.ProcessingTime{
				MinDays: 2,
				MaxDays: 5,
			},
			CreatedAt: procedure.CreatedAt,
			NoticeIDs: []string{primitive.NewObjectID().Hex()},
		}

		err := s.repo.Update(ctx, insertedID, updatedProcedure)
		s.NoError(err)

		// Verify the update in the database
		var updatedModel ProcedureModel
		objID, _ := primitive.ObjectIDFromHex(insertedID)
		err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedModel)
		s.NoError(err)

		s.Equal("Updated Procedure", updatedModel.Name)
		s.Equal([]string{"updated prereq1", "updated prereq2"}, updatedModel.Content.Prerequisites)
		s.Equal("Updated Fee", updatedModel.Fees.Label)
		s.Equal("ETB", updatedModel.Fees.Currency)
		s.Equal(75.25, updatedModel.Fees.Amount)
		s.Equal(2, updatedModel.ProcessingTime.MinDays)
		s.Equal(5, updatedModel.ProcessingTime.MaxDays)
	})

	s.Run("Failure - Not Found", func() {
		nonExistentID := primitive.NewObjectID().Hex()
		updatedProcedure := &domain.Procedure{
			ID:   nonExistentID,
			Name: "Should Not Update",
		}

		err := s.repo.Update(ctx, nonExistentID, updatedProcedure)
		s.Error(err)
		s.ErrorIs(err, domain.ErrNotFound)
	})

	s.Run("Failure - Invalid ID Format", func() {
		updatedProcedure := &domain.Procedure{
			ID:   "invalid-id",
			Name: "Should Not Update",
		}

		err := s.repo.Update(ctx, "invalid-id", updatedProcedure)
		s.Error(err)
		s.ErrorIs(err, domain.ErrNotFound)
	})
}

func (s *ProcedureRepositoryTestSuite) TestMappingLogic() {
	ctx := context.Background()

	s.Run("Correctly maps procedure with all fields", func() {
		procedure := &domain.Procedure{
			Name:           "Complete Procedure",
			OrganizationID: primitive.NewObjectID().Hex(),
			Content: domain.ProcedureContent{
				Prerequisites: []string{"doc1", "doc2"},
				Steps: map[int]string{
					1: "Fill form",
					2: "Submit documents",
					3: "Wait for approval",
				},
				Result: "Certificate issued",
			},
			Fees: domain.ProcedureFee{
				Label:    "Processing Fee",
				Currency: "ETB",
				Amount:   150.75,
			},
			ProcessingTime: domain.ProcessingTime{
				MinDays: 3,
				MaxDays: 7,
			},
			CreatedAt: time.Now(),
			NoticeIDs: []string{
				primitive.NewObjectID().Hex(),
				primitive.NewObjectID().Hex(),
			},
		}

		// Insert via repository (if Create method exists) or directly
		procedureModel := ToDTO(procedure)
		result, err := s.collection.InsertOne(ctx, procedureModel)
		s.Require().NoError(err)

		insertedID := result.InsertedID.(primitive.ObjectID).Hex()

		// Retrieve via repository
		foundProcedure, err := s.repo.GetByID(ctx, insertedID)
		s.NoError(err)

		// Verify all fields are correctly mapped
		s.Equal("Complete Procedure", foundProcedure.Name)
		s.Equal(procedure.OrganizationID, foundProcedure.OrganizationID)
		s.Equal([]string{"doc1", "doc2"}, foundProcedure.Content.Prerequisites)
		s.Equal(map[int]string{1: "Fill form", 2: "Submit documents", 3: "Wait for approval"}, foundProcedure.Content.Steps)
		s.Equal("Certificate issued", foundProcedure.Content.Result)
		s.Equal("Processing Fee", foundProcedure.Fees.Label)
		s.Equal("ETB", foundProcedure.Fees.Currency)
		s.Equal(150.75, foundProcedure.Fees.Amount)
		s.Equal(3, foundProcedure.ProcessingTime.MinDays)
		s.Equal(7, foundProcedure.ProcessingTime.MaxDays)
		s.Len(foundProcedure.NoticeIDs, 2)
		s.WithinDuration(procedure.CreatedAt, foundProcedure.CreatedAt, time.Second)
	})

	s.Run("Correctly maps procedure with optional fields empty", func() {
		procedure := &domain.Procedure{
			Name:           "Minimal Procedure",
			OrganizationID: primitive.NewObjectID().Hex(),
			Content: domain.ProcedureContent{
				Steps: map[int]string{1: "Single step"},
			},
			CreatedAt: time.Now(),
		}

		procedureModel := ToDTO(procedure)
		result, err := s.collection.InsertOne(ctx, procedureModel)
		s.Require().NoError(err)

		insertedID := result.InsertedID.(primitive.ObjectID).Hex()

		foundProcedure, err := s.repo.GetByID(ctx, insertedID)
		s.NoError(err)

		s.Equal("Minimal Procedure", foundProcedure.Name)
		s.Empty(foundProcedure.Content.Prerequisites)
		s.Equal(map[int]string{1: "Single step"}, foundProcedure.Content.Steps)
		s.Empty(foundProcedure.Content.Result)
		s.Empty(foundProcedure.Fees.Label)
		s.Empty(foundProcedure.Fees.Currency)
		s.Equal(0.0, foundProcedure.Fees.Amount)
		s.Equal(0, foundProcedure.ProcessingTime.MinDays)
		s.Equal(0, foundProcedure.ProcessingTime.MaxDays)
		s.Empty(foundProcedure.NoticeIDs)
	})
}