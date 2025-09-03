package usecase_test

import (
	"EthioGuide/domain"
	"EthioGuide/usecase"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockProcedureRepo struct {
	mock.Mock
}

func (m *MockProcedureRepo) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	var proc *domain.Procedure
	if args.Get(0) != nil {
		proc = args.Get(0).(*domain.Procedure)
	}
	return proc, args.Error(1)
}
func (m *MockProcedureRepo) Update(ctx context.Context, id string, procedure *domain.Procedure) error {
	args := m.Called(ctx, id, procedure)
	return args.Error(0)
}
func (m *MockProcedureRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProcedureRepo) Create(ctx context.Context, procedure *domain.Procedure) error {
	args := m.Called(ctx, procedure)
	return args.Error(0)
}

// --- Test Suite ---
type ProcedureUsecaseTestSuite struct {
	suite.Suite
	mockRepo       *MockProcedureRepo
	usecase        domain.IProcedureUsecase
	contextTimeout time.Duration
}

func (s *ProcedureUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockProcedureRepo)
	s.contextTimeout = 2 * time.Second
	s.usecase = usecase.NewProcedureUsecase(s.mockRepo, s.contextTimeout)
}

func TestProcedureUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureUsecaseTestSuite))
}

// --- Tests ---
func (s *ProcedureUsecaseTestSuite) TestCreateProcedure() {
	validProc := &domain.Procedure{
		Name:           "Test Procedure",
		OrganizationID: "org456",
		Content: domain.ProcedureContent{
			Prerequisites: []string{"A", "B"},
			Steps:         map[int]string{1: "Step1", 2: "Step2"},
			Result:        "Result1",
		},
		Fees: domain.ProcedureFee{
			Label:    "FeeLabel",
			Currency: "ETB",
			Amount:   100.0,
		},
		ProcessingTime: domain.ProcessingTime{
			MinDays: 1,
			MaxDays: 5,
		},
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Procedure")).Return(nil).Once()

		err := s.usecase.CreateProcedure(context.Background(), validProc)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})
}

// --- UpdateProcedure Tests ---
func (s *ProcedureUsecaseTestSuite) TestUpdateProcedure() {
	procedureID := "proc-123"
	organizationID := "org-abc"

	// Contexts for different user roles
	orgOwnerCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleOrg), "userID", organizationID)
	adminCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleAdmin), "userID", "admin-user-id")
	otherOrgCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleOrg), "userID", "different-org-id")
	userCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleUser), "userID", "regular-user-id")

	existingProcedure := &domain.Procedure{ID: procedureID, OrganizationID: organizationID}
	updatedProcedure := &domain.Procedure{ID: procedureID, Name: "Updated Name", OrganizationID: organizationID}

	s.Run("Success - Owning Organization Updates", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(nil).Once()

		err := s.usecase.UpdateProcedure(orgOwnerCtx, procedureID, updatedProcedure)

		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Admin Updates", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockRepo.On("Update", mock.Anything, procedureID, updatedProcedure).Return(nil).Once()

		err := s.usecase.UpdateProcedure(adminCtx, procedureID, updatedProcedure)

		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Permission Denied (Different Organization)", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.UpdateProcedure(otherOrgCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	s.Run("Failure - Permission Denied (Regular User)", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.UpdateProcedure(userCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
	})

	s.Run("Failure - Procedure Not Found", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(nil, domain.ErrNotFound).Once()

		err := s.usecase.UpdateProcedure(orgOwnerCtx, procedureID, updatedProcedure)

		s.ErrorIs(err, domain.ErrNotFound)
		s.mockRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
	})
}

// --- DeleteProcedure Tests ---
func (s *ProcedureUsecaseTestSuite) TestDeleteProcedure() {
	procedureID := "proc-456"
	organizationID := "org-xyz"

	orgOwnerCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleOrg), "userID", organizationID)
	adminCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleAdmin), "userID", "admin-user-id")
	otherOrgCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleOrg), "userID", "different-org-id")
	userCtx := context.WithValue(context.WithValue(context.Background(), "userRole", domain.RoleUser), "userID", "regular-user-id")

	existingProcedure := &domain.Procedure{ID: procedureID, OrganizationID: organizationID}

	s.Run("Success - Owning Organization Deletes", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockRepo.On("Delete", mock.Anything, procedureID).Return(nil).Once()

		err := s.usecase.DeleteProcedure(orgOwnerCtx, procedureID)

		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Admin Deletes", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()
		s.mockRepo.On("Delete", mock.Anything, procedureID).Return(nil).Once()

		err := s.usecase.DeleteProcedure(adminCtx, procedureID)

		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Permission Denied (Different Organization)", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.DeleteProcedure(otherOrgCtx, procedureID)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
	})

	s.Run("Failure - Permission Denied (Regular User)", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(existingProcedure, nil).Once()

		err := s.usecase.DeleteProcedure(userCtx, procedureID)

		s.ErrorIs(err, domain.ErrPermissionDenied)
		s.mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
	})

	s.Run("Failure - Procedure Not Found", func() {
		s.SetupTest()
		s.mockRepo.On("GetByID", mock.Anything, procedureID).Return(nil, domain.ErrNotFound).Once()

		err := s.usecase.DeleteProcedure(orgOwnerCtx, procedureID)

		s.ErrorIs(err, domain.ErrNotFound)
		s.mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
	})
}
