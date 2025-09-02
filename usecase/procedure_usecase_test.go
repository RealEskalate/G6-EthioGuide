package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks & Placeholders ---

// MockProcedureRepository mocks IProcedureRepository
type MockProcedureRepository struct {
	mock.Mock
}

func (m *MockProcedureRepository) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	var procedure *domain.Procedure
	if args.Get(0) != nil {
		procedure = args.Get(0).(*domain.Procedure)
	}
	return procedure, args.Error(1)
}

func (m *MockProcedureRepository) Update(ctx context.Context, id string, procedure *domain.Procedure) error {
	return m.Called(ctx, id, procedure).Error(0)
}

func (m *MockProcedureRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// --- Test Suite Definition ---
type ProcedureUsecaseTestSuite struct {
	suite.Suite
	mockProcedureRepo *MockProcedureRepository
	usecase           *ProcedureUsecase
}

func (s *ProcedureUsecaseTestSuite) SetupSuite() {
	setupDomainErrors()
}

func (s *ProcedureUsecaseTestSuite) SetupTest() {
	s.mockProcedureRepo = new(MockProcedureRepository)
	s.usecase = NewProcedureUsecase(s.mockProcedureRepo)
}

func TestProcedureUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureUsecaseTestSuite))
}

// --- GetProcedureByID Tests ---
func (s *ProcedureUsecaseTestSuite) TestGetProcedureByID() {
	procedureID := "procedure-123"
	mockProcedure := &domain.Procedure{
		ID:               procedureID,
		Name:             "Test Procedure",
		OrganizationID:   "org-123",
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(mockProcedure, nil).Once()

		result, err := s.usecase.GetProcedureByID(context.Background(), procedureID)

		s.NoError(err)
		s.Equal(mockProcedure, result)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Procedure Not Found", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, domain.ErrNotFound).Once()

		result, err := s.usecase.GetProcedureByID(context.Background(), procedureID)

		s.ErrorIs(err, domain.ErrNotFound)
		s.Nil(result)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Database Error", func() {
		s.SetupTest()
		dbError := errors.New("database connection failed")
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, dbError).Once()

		result, err := s.usecase.GetProcedureByID(context.Background(), procedureID)

		s.ErrorIs(err, dbError)
		s.Nil(result)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})
}

// --- UpdateProcedure Tests ---
func (s *ProcedureUsecaseTestSuite) TestUpdateProcedure() {
	procedureID := "procedure-123"
	organizationID := "org-123"

	existingProcedure := &domain.Procedure{
		ID:               procedureID,
		Name:             "Original Procedure",
		OrganizationID:   organizationID,
	}

	updatedProcedure := &domain.Procedure{
		ID:               procedureID,
		Name:             "Updated Procedure",
		OrganizationID:   organizationID,
	}

	s.Run("Success", func() {
		s.SetupTest()
		// Mock getting the existing procedure
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		// Mock the update operation
		s.mockProcedureRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(nil).Once()

		err := s.usecase.UpdateProcedure(context.Background(), procedureID, updatedProcedure)

		s.NoError(err)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Procedure Not Found", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, domain.ErrNotFound).Once()

		err := s.usecase.UpdateProcedure(context.Background(), procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrNotFound)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Permission Denied (Different Organization)", func() {
		s.SetupTest()
		// Create a procedure from a different organization
		differentOrgProcedure := &domain.Procedure{
			ID:               procedureID,
			Name:             "Updated Procedure",
			OrganizationID:   "different-org", // Different organization ID
		}

		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.UpdateProcedure(context.Background(), procedureID, differentOrgProcedure)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Database Error on Get", func() {
		s.SetupTest()
		dbError := errors.New("database connection failed")
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, dbError).Once()

		err := s.usecase.UpdateProcedure(context.Background(), procedureID, updatedProcedure)

		s.ErrorIs(err, dbError)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Database Error on Update", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		
		dbError := errors.New("database update failed")
		s.mockProcedureRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(dbError).Once()

		err := s.usecase.UpdateProcedure(context.Background(), procedureID, updatedProcedure)

		s.ErrorIs(err, dbError)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Same Organization Different Fields", func() {
		s.SetupTest()
		// Test that only organization ID matters for permission check
		similarProcedure := &domain.Procedure{
			ID:               procedureID,
			Name:             "Completely Different Name",
			OrganizationID:   organizationID, // Same organization ID
		}

		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockProcedureRepo.On("Update", mock.Anything, procedureID, similarProcedure).Return(nil).Once()

		err := s.usecase.UpdateProcedure(context.Background(), procedureID, similarProcedure)

		s.NoError(err)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})
}