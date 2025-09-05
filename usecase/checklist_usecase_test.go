package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock Checklist Repository ---

type MockChecklistRepository struct {
	mock.Mock
}

var _ domain.IChecklistRepository = (*MockChecklistRepository)(nil)

func (m *MockChecklistRepository) CreateChecklist(ctx context.Context, userID, procedureID string) (*domain.UserProcedure, error) {
	args := m.Called(ctx, userID, procedureID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserProcedure), args.Error(1)
}

func (m *MockChecklistRepository) GetProcedures(ctx context.Context, userID string) ([]*domain.UserProcedure, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserProcedure), args.Error(1)
}

func (m *MockChecklistRepository) GetChecklistByUserProcedureID(ctx context.Context, userProcedureID string) ([]*domain.Checklist, error) {
	args := m.Called(ctx, userProcedureID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Checklist), args.Error(1)
}

func (m *MockChecklistRepository) FindCheck(ctx context.Context, checklistID string) (*domain.Checklist, error) {
	args := m.Called(ctx, checklistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Checklist), args.Error(1)
}

func (m *MockChecklistRepository) ToggleCheck(ctx context.Context, checklistID string) error {
	args := m.Called(ctx, checklistID)
	return args.Error(0)
}

func (m *MockChecklistRepository) CountDocumentsChecklist(ctx context.Context, filter interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockChecklistRepository) UpdateUserProcedure(ctx context.Context, filter interface{}, updatefields map[string]interface{}) error {
	args := m.Called(ctx, filter, updatefields)
	return args.Error(0)
}

// --- Test Suite Definition ---

type ChecklistUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockChecklistRepository
	usecase  domain.IChecklistUsecase
	ctx      context.Context
}

func (s *ChecklistUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockChecklistRepository)
	s.usecase = NewChecklistUsecase(s.mockRepo)
	s.ctx = context.Background()
}

func TestChecklistUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ChecklistUsecaseTestSuite))
}

// --- Test Cases ---

func (s *ChecklistUsecaseTestSuite) TestCreateChecklist() {
	userID := uuid.New().String()
	procedureID := uuid.New().String()
	expectedUserProcedure := &domain.UserProcedure{ID: uuid.New().String(), UserID: userID, ProcedureID: procedureID}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("CreateChecklist", mock.Anything, userID, procedureID).Return(expectedUserProcedure, nil).Once()

		userProcedure, err := s.usecase.CreateChecklist(s.ctx, userID, procedureID)

		s.NoError(err)
		s.Equal(expectedUserProcedure, userProcedure)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid ID", func() {
		s.SetupTest()
		_, err := s.usecase.CreateChecklist(s.ctx, "", procedureID)
		s.ErrorIs(err, domain.ErrInvalidID)

		_, err = s.usecase.CreateChecklist(s.ctx, userID, "")
		s.ErrorIs(err, domain.ErrInvalidID)
	})

	s.Run("Failure - Repository Error", func() {
		s.SetupTest()
		repoErr := errors.New("database error")
		s.mockRepo.On("CreateChecklist", mock.Anything, userID, procedureID).Return(nil, repoErr).Once()

		_, err := s.usecase.CreateChecklist(s.ctx, userID, procedureID)

		s.Error(err)
		s.ErrorIs(err, repoErr)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *ChecklistUsecaseTestSuite) TestGetProcedures() {
	userID := uuid.New().String()
	expectedProcedures := []*domain.UserProcedure{{ID: uuid.New().String()}, {ID: uuid.New().String()}}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("GetProcedures", mock.Anything, userID).Return(expectedProcedures, nil).Once()

		procedures, err := s.usecase.GetProcedures(s.ctx, userID)

		s.NoError(err)
		s.Equal(expectedProcedures, procedures)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid ID", func() {
		s.SetupTest()
		_, err := s.usecase.GetProcedures(s.ctx, "")
		s.ErrorIs(err, domain.ErrInvalidID)
	})

	s.Run("Failure - Repository Error", func() {
		s.SetupTest()
		repoErr := errors.New("database error")
		s.mockRepo.On("GetProcedures", mock.Anything, userID).Return(nil, repoErr).Once()

		_, err := s.usecase.GetProcedures(s.ctx, userID)

		s.Error(err)
		s.ErrorIs(err, repoErr)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *ChecklistUsecaseTestSuite) TestGetChecklistByUserProcedureID() {
	userProcedureID := uuid.New().String()
	expectedChecklists := []*domain.Checklist{{ID: uuid.New().String()}, {ID: uuid.New().String()}}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("GetChecklistByUserProcedureID", mock.Anything, userProcedureID).Return(expectedChecklists, nil).Once()

		checklists, err := s.usecase.GetChecklistByUserProcedureID(s.ctx, userProcedureID)

		s.NoError(err)
		s.Equal(expectedChecklists, checklists)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid ID", func() {
		s.SetupTest()
		_, err := s.usecase.GetChecklistByUserProcedureID(s.ctx, "")
		s.ErrorIs(err, domain.ErrInvalidID)
	})

	s.Run("Failure - Repository Error", func() {
		s.SetupTest()
		repoErr := errors.New("database error")
		s.mockRepo.On("GetChecklistByUserProcedureID", mock.Anything, userProcedureID).Return(nil, repoErr).Once()

		_, err := s.usecase.GetChecklistByUserProcedureID(s.ctx, userProcedureID)

		s.Error(err)
		s.ErrorIs(err, repoErr)
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *ChecklistUsecaseTestSuite) TestUpdateChecklist() {
	checklistID := uuid.New().String()
	userProcedureID := uuid.New().String()
	mockChecklist := &domain.Checklist{ID: checklistID, UserProcedureID: userProcedureID, IsChecked: false}

	s.Run("Success - Not Started to In Progress", func() {
		s.SetupTest()
		s.mockRepo.On("FindCheck", mock.Anything, checklistID).Return(mockChecklist, nil).Once()
		s.mockRepo.On("ToggleCheck", mock.Anything, checklistID).Return(nil).Once()
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(10, nil).Once() // Total
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(5, nil).Once()  // Checked
		s.mockRepo.On("UpdateUserProcedure", mock.Anything, mock.Anything, mock.MatchedBy(func(update map[string]interface{}) bool {
			return update["status"] == "In Progress"
		})).Return(nil).Once()

		checklist, err := s.usecase.UpdateChecklist(s.ctx, checklistID)

		s.NoError(err)
		s.True(checklist.IsChecked)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - In Progress to Completed", func() {
		s.SetupTest()
		mockChecklist.IsChecked = false
		s.mockRepo.On("FindCheck", mock.Anything, checklistID).Return(mockChecklist, nil).Once()
		s.mockRepo.On("ToggleCheck", mock.Anything, checklistID).Return(nil).Once()
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(10, nil).Once() // Total
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(10, nil).Once() // Checked
		s.mockRepo.On("UpdateUserProcedure", mock.Anything, mock.Anything, mock.MatchedBy(func(update map[string]interface{}) bool {
			return update["status"] == "Completed"
		})).Return(nil).Once()

		checklist, err := s.usecase.UpdateChecklist(s.ctx, checklistID)

		s.NoError(err)
		s.True(checklist.IsChecked)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Completed to In Progress", func() {
		s.SetupTest()
		mockChecklist.IsChecked = true // Start as checked
		s.mockRepo.On("FindCheck", mock.Anything, checklistID).Return(mockChecklist, nil).Once()
		s.mockRepo.On("ToggleCheck", mock.Anything, checklistID).Return(nil).Once()
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(10, nil).Once() // Total
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(9, nil).Once()   // Checked (after un-checking one)
		s.mockRepo.On("UpdateUserProcedure", mock.Anything, mock.Anything, mock.MatchedBy(func(update map[string]interface{}) bool {
			return update["status"] == "In Progress"
		})).Return(nil).Once()

		checklist, err := s.usecase.UpdateChecklist(s.ctx, checklistID)

		s.NoError(err)
		s.False(checklist.IsChecked) // It should be toggled to false
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - In Progress to Not Started", func() {
		s.SetupTest()
		mockChecklist.IsChecked = true // Start as checked
		s.mockRepo.On("FindCheck", mock.Anything, checklistID).Return(mockChecklist, nil).Once()
		s.mockRepo.On("ToggleCheck", mock.Anything, checklistID).Return(nil).Once()
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(10, nil).Once() // Total
		s.mockRepo.On("CountDocumentsChecklist", mock.Anything, mock.AnythingOfType("map[string]interface {}")).Return(0, nil).Once()   // Checked (after un-checking the last one)
		s.mockRepo.On("UpdateUserProcedure", mock.Anything, mock.Anything, mock.MatchedBy(func(update map[string]interface{}) bool {
			return update["status"] == "Not Started"
		})).Return(nil).Once()

		checklist, err := s.usecase.UpdateChecklist(s.ctx, checklistID)

		s.NoError(err)
		s.False(checklist.IsChecked)
		s.mockRepo.AssertExpectations(s.T())
	})

		s.Run("Failure - Invalid ID", func() {
		s.SetupTest()
		_, err := s.usecase.UpdateChecklist(s.ctx, "")
		s.ErrorIs(err, domain.ErrInvalidID)
	})

	s.Run("Failure - FindCheck Error", func() {
		s.SetupTest()
		repoErr := errors.New("find error")
		s.mockRepo.On("FindCheck", mock.Anything, checklistID).Return(nil, repoErr).Once()

		_, err := s.usecase.UpdateChecklist(s.ctx, checklistID)

		s.Error(err)
		s.ErrorIs(err, repoErr)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - ToggleCheck Error", func() {
		s.SetupTest()
		repoErr := errors.New("toggle error")
		s.mockRepo.On("FindCheck", mock.Anything, checklistID).Return(mockChecklist, nil).Once()
		s.mockRepo.On("ToggleCheck", mock.Anything, checklistID).Return(repoErr).Once()

		_, err := s.usecase.UpdateChecklist(s.ctx, checklistID)

		s.Error(err)
		s.ErrorIs(err, repoErr)
		s.mockRepo.AssertExpectations(s.T())
	})
}
