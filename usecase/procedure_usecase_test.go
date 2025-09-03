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
		ID:             procedureID,
		Name:           "Test Procedure",
		OrganizationID: "org-123",
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
	organizationID := "org-123" // The ID of the organization that OWNS the procedure.

	// --- Contexts for different user roles and IDs ---
	// Context for a user who IS the owning organization.
	orgOwnerCtx := context.WithValue(context.Background(), "userRole", domain.RoleOrg)
	orgOwnerCtx = context.WithValue(orgOwnerCtx, "userID", organizationID)

	// Context for an admin user.
	adminCtx := context.WithValue(context.Background(), "userRole", domain.RoleAdmin)
	adminCtx = context.WithValue(adminCtx, "userID", "admin-user-id-xyz")

	// Context for a different organization (who should be denied).
	otherOrgCtx := context.WithValue(context.Background(), "userRole", domain.RoleOrg)
	otherOrgCtx = context.WithValue(otherOrgCtx, "userID", "different-org-id-999")

	// Context for a regular user (who should be denied).
	userCtx := context.WithValue(context.Background(), "userRole", domain.RoleUser)
	userCtx = context.WithValue(userCtx, "userID", "regular-user-id-abc")

	// --- Procedure objects for testing ---
	existingProcedure := &domain.Procedure{ID: procedureID, Name: "Original Procedure", OrganizationID: organizationID}
	updatedProcedure := &domain.Procedure{ID: procedureID, Name: "Updated Procedure", OrganizationID: organizationID}

	s.Run("Success - Owning Organization Updates", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockProcedureRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(nil).Once()

		err := s.usecase.UpdateProcedure(orgOwnerCtx, procedureID, updatedProcedure)

		s.NoError(err)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Admin Updates", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockProcedureRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(nil).Once()

		err := s.usecase.UpdateProcedure(adminCtx, procedureID, updatedProcedure)

		s.NoError(err)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Procedure Not Found on Get", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, domain.ErrNotFound).Once()

		err := s.usecase.UpdateProcedure(orgOwnerCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrNotFound)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Permission Denied (Different Organization)", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.UpdateProcedure(otherOrgCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Permission Denied (Regular User)", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.UpdateProcedure(userCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Database Error on Get", func() {
		s.SetupTest()
		dbError := errors.New("database connection failed")
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, dbError).Once()

		err := s.usecase.UpdateProcedure(orgOwnerCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, dbError)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Update")
	})

	s.Run("Failure - Database Error on Update", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		dbError := errors.New("database update failed")
		s.mockProcedureRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(dbError).Once()

		err := s.usecase.UpdateProcedure(orgOwnerCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, dbError)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})
}

// --- DeleteProcedure Tests ---
func (s *ProcedureUsecaseTestSuite) TestDeleteProcedure() {
	procedureID := "procedure-to-delete-123"
	organizationID := "org-abc"

	// --- Contexts for different user roles and IDs ---
	orgOwnerCtx := context.WithValue(context.Background(), "userRole", domain.RoleOrg)
	orgOwnerCtx = context.WithValue(orgOwnerCtx, "userID", organizationID)

	adminCtx := context.WithValue(context.Background(), "userRole", domain.RoleAdmin)
	adminCtx = context.WithValue(adminCtx, "userID", "admin-user-id-xyz")

	otherOrgCtx := context.WithValue(context.Background(), "userRole", domain.RoleOrg)
	otherOrgCtx = context.WithValue(otherOrgCtx, "userID", "different-org-id-999")

	userCtx := context.WithValue(context.Background(), "userRole", domain.RoleUser)
	userCtx = context.WithValue(userCtx, "userID", "regular-user-id-abc")

	// --- Procedure object for testing ---
	existingProcedure := &domain.Procedure{ID: procedureID, Name: "To Be Deleted", OrganizationID: organizationID}

	s.Run("Success - Owning Organization Deletes", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockProcedureRepo.On("Delete", mock.Anything, procedureID).Return(nil).Once()

		err := s.usecase.DeleteProcedure(orgOwnerCtx, procedureID)

		s.NoError(err)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Admin Deletes", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockProcedureRepo.On("Delete", mock.Anything, procedureID).Return(nil).Once()

		err := s.usecase.DeleteProcedure(adminCtx, procedureID)

		s.NoError(err)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Procedure Not Found on Get", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(nil, domain.ErrNotFound).Once()

		err := s.usecase.DeleteProcedure(orgOwnerCtx, procedureID)

		s.ErrorIs(err, domain.ErrNotFound)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Delete")
	})

	s.Run("Failure - Permission Denied (Different Organization)", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.DeleteProcedure(otherOrgCtx, procedureID)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Delete")
	})

	s.Run("Failure - Permission Denied (Regular User)", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.DeleteProcedure(userCtx, procedureID)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockProcedureRepo.AssertExpectations(s.T())
		s.mockProcedureRepo.AssertNotCalled(s.T(), "Delete")
	})

	s.Run("Failure - Database Error on Delete", func() {
		s.SetupTest()
		s.mockProcedureRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		dbError := errors.New("database delete failed")
		s.mockProcedureRepo.On("Delete", mock.Anything, procedureID).Return(dbError).Once()

		err := s.usecase.DeleteProcedure(orgOwnerCtx, procedureID)

		s.ErrorIs(err, dbError)
		s.mockProcedureRepo.AssertExpectations(s.T())
	})
}
