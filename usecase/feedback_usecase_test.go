package usecase_test

import (
	"EthioGuide/domain"
	"EthioGuide/usecase"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks ---
type MockFeedbackRepo struct {
    mock.Mock
}

func (m *MockFeedbackRepo) SubmitFeedback(ctx context.Context, feedback *domain.Feedback) error {
    args := m.Called(ctx, feedback)
    return args.Error(0)
}
func (m *MockFeedbackRepo) GetFeedbackByID(ctx context.Context, id string) (*domain.Feedback, error) {
    args := m.Called(ctx, id)
    var fb *domain.Feedback
    if args.Get(0) != nil {
        fb = args.Get(0).(*domain.Feedback)
    }
    return fb, args.Error(1)
}
func (m *MockFeedbackRepo) GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
    args := m.Called(ctx, procedureID, filter)
    var feedbacks []*domain.Feedback
    if args.Get(0) != nil {
        feedbacks = args.Get(0).([]*domain.Feedback)
    }
    return feedbacks, args.Get(1).(int64), args.Error(2)
}
func (m *MockFeedbackRepo) UpdateFeedbackStatus(ctx context.Context, feedbackID string, newFeedback *domain.Feedback) error {
    args := m.Called(ctx, feedbackID, newFeedback)
    return args.Error(0)
}

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

// --- Test Suite ---
type FeedbackUsecaseTestSuite struct {
    suite.Suite
    mockFeedbackRepo  *MockFeedbackRepo
    mockProcedureRepo *MockProcedureRepo
    usecase           domain.IFeedbackUsecase
    timeout           time.Duration
}

func (s *FeedbackUsecaseTestSuite) SetupTest() {
    s.mockFeedbackRepo = new(MockFeedbackRepo)
    s.mockProcedureRepo = new(MockProcedureRepo)
    s.timeout = 2 * time.Second
    s.usecase = usecase.NewFeedbackUsecase(s.mockFeedbackRepo, s.mockProcedureRepo, s.timeout)
}

func TestFeedbackUsecaseTestSuite(t *testing.T) {
    suite.Run(t, new(FeedbackUsecaseTestSuite))
}

// --- Tests ---

func (s *FeedbackUsecaseTestSuite) TestSubmitFeedback() {
    feedback := &domain.Feedback{
        UserID:      "user1",
        ProcedureID: "proc1",
        Content:     "Nice!",
        Type:        domain.ThanksFeedback,
        Tags:        []string{"tag1"},
    }
    procedure := &domain.Procedure{ID: "proc1"}

    s.Run("Success", func() {
        s.SetupTest()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(procedure, nil).Once()
        s.mockFeedbackRepo.On("SubmitFeedback", mock.Anything, feedback).Return(nil).Once()

        err := s.usecase.SubmitFeedback(context.Background(), feedback)
        s.NoError(err)
        s.mockProcedureRepo.AssertExpectations(s.T())
        s.mockFeedbackRepo.AssertExpectations(s.T())
    })

    s.Run("Failure - Procedure Not Found", func() {
        s.SetupTest()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(nil, domain.ErrNotFound).Once()

        err := s.usecase.SubmitFeedback(context.Background(), feedback)
        s.ErrorIs(err, domain.ErrNotFound)
        s.mockProcedureRepo.AssertExpectations(s.T())
    })

    s.Run("Failure - Repo Error", func() {
        s.SetupTest()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(procedure, nil).Once()
        s.mockFeedbackRepo.On("SubmitFeedback", mock.Anything, feedback).Return(errors.New("db error")).Once()

        err := s.usecase.SubmitFeedback(context.Background(), feedback)
        s.Error(err)
        s.Contains(err.Error(), "db error")
        s.mockFeedbackRepo.AssertExpectations(s.T())
    })
}

func (s *FeedbackUsecaseTestSuite) TestGetAllFeedbacksForProcedure() {
    filter := &domain.FeedbackFilter{Page: 1, Limit: 10}
    feedbacks := []*domain.Feedback{{ID: "fb1"}}

    s.Run("Success", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetAllFeedbacksForProcedure", mock.Anything, "proc1", filter).Return(feedbacks, int64(1), nil).Once()

        res, total, err := s.usecase.GetAllFeedbacksForProcedure(context.Background(), "proc1", filter)
        s.NoError(err)
        s.Equal(int64(1), total)
        s.Equal(feedbacks, res)
        s.mockFeedbackRepo.AssertExpectations(s.T())
    })

    s.Run("Failure - Repo Error", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetAllFeedbacksForProcedure", mock.Anything, "proc1", filter).Return(nil, int64(0), errors.New("db error")).Once()

        res, total, err := s.usecase.GetAllFeedbacksForProcedure(context.Background(), "proc1", filter)
        s.Error(err)
        s.Nil(res)
        s.Equal(int64(0), total)
        s.mockFeedbackRepo.AssertExpectations(s.T())
    })
}

func (s *FeedbackUsecaseTestSuite) TestUpdateFeedbackStatus() {
    feedback := &domain.Feedback{
        ID:          "fb1",
        ProcedureID: "proc1",
        Status:      domain.NewFeedback,
    }
    procedure := &domain.Procedure{
        ID:             "proc1",
        OrganizationID: "org1",
    }
    adminResponse := "Done"

    s.Run("Success", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedback.ID).Return(feedback, nil).Once()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(procedure, nil).Once()
        s.mockFeedbackRepo.On("UpdateFeedbackStatus", mock.Anything, feedback.ID, mock.AnythingOfType("*domain.Feedback")).Return(nil).Once()

        err := s.usecase.UpdateFeedbackStatus(context.Background(), feedback.ID, "org1", domain.ResolvedFeedback, &adminResponse)
        s.NoError(err)
        s.mockFeedbackRepo.AssertExpectations(s.T())
        s.mockProcedureRepo.AssertExpectations(s.T())
    })

    s.Run("Failure - Feedback Not Found", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedback.ID).Return(nil, domain.ErrNotFound).Once()

        err := s.usecase.UpdateFeedbackStatus(context.Background(), feedback.ID, "org1", domain.ResolvedFeedback, &adminResponse)
        s.ErrorIs(err, domain.ErrNotFound)
        s.mockFeedbackRepo.AssertExpectations(s.T())
    })

    s.Run("Failure - Procedure Not Found", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedback.ID).Return(feedback, nil).Once()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(nil, domain.ErrNotFound).Once()

        err := s.usecase.UpdateFeedbackStatus(context.Background(), feedback.ID, "org1", domain.ResolvedFeedback, &adminResponse)
        s.ErrorIs(err, domain.ErrNotFound)
        s.mockProcedureRepo.AssertExpectations(s.T())
    })

    s.Run("Failure - Permission Denied", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedback.ID).Return(feedback, nil).Once()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(procedure, nil).Once()

        err := s.usecase.UpdateFeedbackStatus(context.Background(), feedback.ID, "wrong_org", domain.ResolvedFeedback, &adminResponse)
        s.ErrorIs(err, domain.ErrPermissionDenied)
    })

    s.Run("Failure - Missing Admin Response", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedback.ID).Return(feedback, nil).Once()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(procedure, nil).Once()

        err := s.usecase.UpdateFeedbackStatus(context.Background(), feedback.ID, "org1", domain.ResolvedFeedback, nil)
        s.Error(err)
        s.Contains(err.Error(), "admin response is required")
    })

    s.Run("Failure - Repo Update Error", func() {
        s.SetupTest()
        s.mockFeedbackRepo.On("GetFeedbackByID", mock.Anything, feedback.ID).Return(feedback, nil).Once()
        s.mockProcedureRepo.On("GetByID", mock.Anything, feedback.ProcedureID).Return(procedure, nil).Once()
        s.mockFeedbackRepo.On("UpdateFeedbackStatus", mock.Anything, feedback.ID, mock.AnythingOfType("*domain.Feedback")).Return(errors.New("db error")).Once()

        err := s.usecase.UpdateFeedbackStatus(context.Background(), feedback.ID, "org1", domain.ResolvedFeedback, &adminResponse)
        s.Error(err)
        s.Contains(err.Error(), "db error")
    })
}