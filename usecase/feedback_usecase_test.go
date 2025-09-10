package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockFeedbackRepository struct {
	mock.Mock
}

func (m *MockFeedbackRepository) SubmitFeedback(ctx context.Context, feedback *domain.Feedback) error {
	args := m.Called(ctx, feedback)
	return args.Error(0)
}
func (m *MockFeedbackRepository) GetFeedbackByID(ctx context.Context, id string) (*domain.Feedback, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Feedback), args.Error(1)
}
func (m *MockFeedbackRepository) GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	args := m.Called(ctx, procedureID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Feedback), args.Get(1).(int64), args.Error(2)
}
func (m *MockFeedbackRepository) UpdateFeedbackStatus(ctx context.Context, feedbackID string, newFeedback *domain.Feedback) error {
	args := m.Called(ctx, feedbackID, newFeedback)
	return args.Error(0)
}
func (m *MockFeedbackRepository) GetAllFeedbacks(ctx context.Context, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Feedback), args.Get(1).(int64), args.Error(2)
}

// --- Test Suite ---
type FeedbackUsecaseTestSuite struct {
	suite.Suite
	mockFeedbackRepo *MockFeedbackRepository
	mockProcRepo     *MockProcedureRepository
	usecase          domain.IFeedbackUsecase
	ctx              context.Context
}

func (s *FeedbackUsecaseTestSuite) SetupTest() {
	s.mockFeedbackRepo = new(MockFeedbackRepository)
	s.mockProcRepo = new(MockProcedureRepository)
	s.usecase = NewFeedbackUsecase(s.mockFeedbackRepo, s.mockProcRepo, 5*time.Second)
	s.ctx = context.Background()
}

func (s *FeedbackUsecaseTestSuite) TestSubmitFeedback_Success() {
	// Arrange
	feedback := &domain.Feedback{ProcedureID: "proc-123"}
	s.mockProcRepo.On("GetByID", mock.Anything, "proc-123").Return(&domain.Procedure{}, nil).Once()
	s.mockFeedbackRepo.On("SubmitFeedback", mock.Anything, feedback).Return(nil).Once()

	// Act
	err := s.usecase.SubmitFeedback(s.ctx, feedback)

	// Assert
	s.NoError(err)
	s.mockProcRepo.AssertExpectations(s.T())
	s.mockFeedbackRepo.AssertExpectations(s.T())
}

func (s *FeedbackUsecaseTestSuite) TestSubmitFeedback_ProcedureNotFound() {
	// Arrange
	feedback := &domain.Feedback{ProcedureID: "proc-404"}
	s.mockProcRepo.On("GetByID", mock.Anything, "proc-404").Return(nil, domain.ErrProcedureNotFound).Once()

	// Act
	err := s.usecase.SubmitFeedback(s.ctx, feedback)

	// Assert
	s.Error(err)
	s.ErrorIs(err, domain.ErrProcedureNotFound)
	s.mockProcRepo.AssertExpectations(s.T())
	// Ensure SubmitFeedback is not called if the procedure doesn't exist
	s.mockFeedbackRepo.AssertNotCalled(s.T(), "SubmitFeedback")
}

func (s *FeedbackUsecaseTestSuite) TestUpdateFeedbackStatus_Success() {
	// Arrange
	feedbackID := "fb-1"
	userID := "org-1" // This is the owner
	status := domain.ResolvedFeedback
	adminResponse := "Resolved."

	mockFeedback := &domain.Feedback{ID: feedbackID, ProcedureID: "proc-1"}
	mockProcedure := &domain.Procedure{ID: "proc-1", OrganizationID: userID}

	s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedbackID).Return(mockFeedback, nil).Once()
	s.mockProcRepo.On("GetByID", mock.Anything, "proc-1").Return(mockProcedure, nil).Once()
	s.mockFeedbackRepo.On("UpdateFeedbackStatus", mock.Anything, feedbackID, mock.AnythingOfType("*domain.Feedback")).Return(nil).Once()

	// Act
	err := s.usecase.UpdateFeedbackStatus(s.ctx, feedbackID, userID, status, &adminResponse)

	// Assert
	s.NoError(err)
	s.mockFeedbackRepo.AssertExpectations(s.T())
	s.mockProcRepo.AssertExpectations(s.T())
}

func (s *FeedbackUsecaseTestSuite) TestUpdateFeedbackStatus_PermissionDenied() {
	// Arrange
	feedbackID := "fb-1"
	userID := "not-the-owner"
	status := domain.ResolvedFeedback
	adminResponse := "Should not work."

	mockFeedback := &domain.Feedback{ID: feedbackID, ProcedureID: "proc-1"}
	mockProcedure := &domain.Procedure{ID: "proc-1", OrganizationID: "org-owner"} // Different owner

	s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedbackID).Return(mockFeedback, nil).Once()
	s.mockProcRepo.On("GetByID", mock.Anything, "proc-1").Return(mockProcedure, nil).Once()

	// Act
	err := s.usecase.UpdateFeedbackStatus(s.ctx, feedbackID, userID, status, &adminResponse)

	// Assert
	s.Error(err)
	s.ErrorIs(err, domain.ErrPermissionDenied)
	s.mockFeedbackRepo.AssertExpectations(s.T())
	s.mockProcRepo.AssertExpectations(s.T())
	s.mockFeedbackRepo.AssertNotCalled(s.T(), "UpdateFeedbackStatus") // Should not be called on permission failure
}

func (s *FeedbackUsecaseTestSuite) TestUpdateFeedbackStatus_MissingAdminResponse() {
	// Arrange
	feedbackID, userID := "fb-1", "org-1"
	mockFeedback := &domain.Feedback{ID: feedbackID, ProcedureID: "proc-1"}
	mockProcedure := &domain.Procedure{ID: "proc-1", OrganizationID: userID}

	s.Run("Failure when Resolving without admin response", func() {
		// Arrange: Reset mocks and set expectations for this specific sub-test.
		s.SetupTest()
		s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedbackID).Return(mockFeedback, nil).Once()
		s.mockProcRepo.On("GetByID", mock.Anything, "proc-1").Return(mockProcedure, nil).Once()

		// Act
		err := s.usecase.UpdateFeedbackStatus(s.ctx, feedbackID, userID, domain.ResolvedFeedback, nil)

		// Assert
		s.Error(err)
		s.Contains(err.Error(), "admin response is required")
		s.mockFeedbackRepo.AssertExpectations(s.T())
		s.mockProcRepo.AssertExpectations(s.T())
	})

	s.Run("Failure when Declining without admin response", func() {
		// Arrange: Reset mocks and set expectations for this specific sub-test.
		s.SetupTest()
		s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedbackID).Return(mockFeedback, nil).Once()
		s.mockProcRepo.On("GetByID", mock.Anything, "proc-1").Return(mockProcedure, nil).Once()

		// Act
		err := s.usecase.UpdateFeedbackStatus(s.ctx, feedbackID, userID, domain.DeclinedFeedback, nil)

		// Assert
		s.Error(err)
		s.Contains(err.Error(), "admin response is required")
		s.mockFeedbackRepo.AssertExpectations(s.T())
		s.mockProcRepo.AssertExpectations(s.T())
	})
}

func (s *FeedbackUsecaseTestSuite) TestGetAllFeedbacks_PassThrough() {
	// Arrange
	filter := &domain.FeedbackFilter{Page: 1, Limit: 10}
	s.mockFeedbackRepo.On("GetAllFeedbacks", mock.Anything, filter).Return([]*domain.Feedback{}, int64(0), nil).Once()

	// Act
	_, _, err := s.usecase.GetAllFeedbacks(s.ctx, filter)

	// Assert
	s.NoError(err)
	s.mockFeedbackRepo.AssertExpectations(s.T())
}

func (s *FeedbackUsecaseTestSuite) TestGetAllFeedbacksForProcedure_PassThrough() {
	// Arrange
	procedureID := "proc-1"
	filter := &domain.FeedbackFilter{Page: 1, Limit: 10}
	s.mockFeedbackRepo.On("GetAllFeedbacksForProcedure", mock.Anything, procedureID, filter).Return([]*domain.Feedback{}, int64(0), errors.New("repo error")).Once()

	// Act
	_, _, err := s.usecase.GetAllFeedbacksForProcedure(s.ctx, procedureID, filter)

	// Assert
	s.Error(err)
	s.mockFeedbackRepo.AssertExpectations(s.T())
}

func TestFeedbackUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackUsecaseTestSuite))
}
