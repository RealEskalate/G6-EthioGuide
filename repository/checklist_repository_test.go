package repository_test

import (
	. "EthioGuide/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChecklistRepositoryTestSuite struct {
	suite.Suite
	db                      *mongo.Database
	repo                    *ChecklistRepository // Use concrete type to call the repo directly
	proceduresCollection    *mongo.Collection
	userProcedureCollection *mongo.Collection
	checklistCollection     *mongo.Collection
}

func (s *ChecklistRepositoryTestSuite) SetupSuite() {
	s.db = testDBClient.Database(testDBName)
	// We use the concrete type here to have access to the collections for assertions
	s.repo = NewChecklistRepository(s.db)
	s.proceduresCollection = s.db.Collection("procedures")
	s.userProcedureCollection = s.db.Collection("user_procedure")
	s.checklistCollection = s.db.Collection("checklist")
}

func (s *ChecklistRepositoryTestSuite) TearDownSuite() {
	s.db.Drop(context.Background())
}

func (s *ChecklistRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	// Clean all related collections before each test
	s.proceduresCollection.DeleteMany(context.Background(), bson.M{})
	s.userProcedureCollection.DeleteMany(context.Background(), bson.M{})
	s.checklistCollection.DeleteMany(context.Background(), bson.M{})
}

func TestChecklistRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(ChecklistRepositoryTestSuite))
}

func (s *ChecklistRepositoryTestSuite) TestCreateChecklist() {
	ctx := context.Background()
	userID := primitive.NewObjectID()

	// Arrange: Create a master procedure in the database for the checklist to be based on.
	masterProcedure := ProcedureModel{
		ID:   primitive.NewObjectID(),
		Name: "Get Passport",
		Content: ProcedureContentModel{
			Prerequisites: []string{"Birth Certificate", "ID Card"},
			Steps:         map[int]string{1: "Fill form", 2: "Submit"},
		},
	}
	_, err := s.proceduresCollection.InsertOne(ctx, masterProcedure)
	s.Require().NoError(err)

	s.Run("Success", func() {
		// Act
		userProcedure, err := s.repo.CreateChecklist(ctx, userID.Hex(), masterProcedure.ID.Hex())

		// Assert
		s.NoError(err)
		s.NotNil(userProcedure)
		s.Equal("Not Started", userProcedure.Status)
		s.Equal(0, userProcedure.Percent)
		s.Equal(userID.Hex(), userProcedure.UserID)

		// Verify directly in the DB
		// Check that the user_procedure document was created
		count, err := s.userProcedureCollection.CountDocuments(ctx, bson.M{"user_id": userID, "procedure_id": masterProcedure.ID})
		s.NoError(err)
		s.Equal(int64(1), count)

		// Check that the correct number of checklist items were created
		userProcObjID, _ := primitive.ObjectIDFromHex(userProcedure.ID)
		checklistCount, err := s.checklistCollection.CountDocuments(ctx, bson.M{"user_procedure_id": userProcObjID})
		s.NoError(err)
		s.Equal(int64(4), checklistCount) // 2 prerequisites + 2 steps
	})

	s.Run("Failure - Procedure Not Found", func() {
		nonExistentProcID := primitive.NewObjectID().Hex()
		_, err := s.repo.CreateChecklist(ctx, userID.Hex(), nonExistentProcID)
		s.Error(err)
		s.ErrorIs(err, mongo.ErrNoDocuments)
	})
}

func (s *ChecklistRepositoryTestSuite) TestGetProcedures() {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	procID := primitive.NewObjectID()

	// Arrange
	_, err := s.userProcedureCollection.InsertOne(ctx, bson.M{"user_id": userID, "procedure_id": procID})
	s.Require().NoError(err)
	_, err = s.userProcedureCollection.InsertOne(ctx, bson.M{"user_id": userID, "procedure_id": primitive.NewObjectID()})
	s.Require().NoError(err)
	// Add one for another user
	_, err = s.userProcedureCollection.InsertOne(ctx, bson.M{"user_id": primitive.NewObjectID(), "procedure_id": procID})
	s.Require().NoError(err)

	// Act
	procedures, err := s.repo.GetProcedures(ctx, userID.Hex())

	// Assert
	s.NoError(err)
	s.Len(procedures, 2, "Should only find procedures for the specified user")
}

func (s *ChecklistRepositoryTestSuite) TestGetChecklistByUserProcedureID() {
	ctx := context.Background()
	userProcID := primitive.NewObjectID()

	// Arrange
	_, err := s.checklistCollection.InsertOne(ctx, bson.M{"user_procedure_id": userProcID, "content": "Item 1"})
	s.Require().NoError(err)
	_, err = s.checklistCollection.InsertOne(ctx, bson.M{"user_procedure_id": userProcID, "content": "Item 2"})
	s.Require().NoError(err)
	// Add an item for another procedure
	_, err = s.checklistCollection.InsertOne(ctx, bson.M{"user_procedure_id": primitive.NewObjectID(), "content": "Other Item"})
	s.Require().NoError(err)

	// Act
	checklists, err := s.repo.GetChecklistByUserProcedureID(ctx, userProcID.Hex())

	// Assert
	s.NoError(err)
	s.Len(checklists, 2, "Should only find checklist items for the specified user procedure")
}
