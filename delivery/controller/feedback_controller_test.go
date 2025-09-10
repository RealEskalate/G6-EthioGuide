package controller_test

import (
	. "EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockFeedbackUsecase is a mock for the IFeedbackUsecase interface
type MockFeedbackUsecase struct {
	mock.Mock
}

func (m *MockFeedbackUsecase) SubmitFeedback(ctx context.Context, feedback *domain.Feedback) error {
	args := m.Called(ctx, feedback)
	return args.Error(0)
}

func (m *MockFeedbackUsecase) GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	args := m.Called(ctx, procedureID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Feedback), args.Get(1).(int64), args.Error(2)
}

func (m *MockFeedbackUsecase) UpdateFeedbackStatus(ctx context.Context, feedbackID, userID string, status domain.FeedbackStatus, adminResponse *string) error {
	args := m.Called(ctx, feedbackID, userID, status, adminResponse)
	return args.Error(0)
}

func (m *MockFeedbackUsecase) GetAllFeedbacks(ctx context.Context, filter *domain.FeedbackFilter) ([]*domain.Feedback, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Feedback), args.Get(1).(int64), args.Error(2)
}

// FeedbackControllerTestSuite is the test suite for FeedbackController
type FeedbackControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockFeedbackUsecase
	controller  *FeedbackController
}

// SetupTest is run before each test in the suite
func (s *FeedbackControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.mockUsecase = new(MockFeedbackUsecase)
	s.controller = NewFeedbackController(s.mockUsecase)

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	}

	// Routes
	s.router.POST("/procedures/:id/feedback", authMiddleware, s.controller.SubmitFeedback)
	s.router.GET("/procedures/:id/feedback", s.controller.GetAllFeedbacksForProcedure)
	s.router.PATCH("/feedback/:id", authMiddleware, s.controller.UpdateFeedbackStatus)
	s.router.GET("/feedback", authMiddleware, s.controller.GetAllFeedbacks)
}

func (s *FeedbackControllerTestSuite) TestSubmitFeedback_Success() {
	// Arrange
	procedureID := "proc-123"
	reqBody := FeedbackCreateRequest{Content: "Great procedure!", Type: "thanks"}
	s.mockUsecase.On("SubmitFeedback", mock.Anything, mock.AnythingOfType("*domain.Feedback")).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/procedures/%s/feedback", procedureID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Feedback submitted successfully")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *FeedbackControllerTestSuite) TestSubmitFeedback_ProcedureNotFound() {
	// Arrange
	procedureID := "proc-404"
	reqBody := FeedbackCreateRequest{Content: "Great procedure!", Type: "thanks"}
	s.mockUsecase.On("SubmitFeedback", mock.Anything, mock.AnythingOfType("*domain.Feedback")).Return(domain.ErrProcedureNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/procedures/%s/feedback", procedureID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrProcedureNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *FeedbackControllerTestSuite) TestGetAllFeedbacksForProcedure_Success() {
	// Arrange
	procedureID := "proc-123"
	expectedFeedbacks := []*domain.Feedback{{ID: "fb-1", Content: "Test content"}}
	var total int64 = 1
	status := "new"
	filter := &domain.FeedbackFilter{Page: 1, Limit: 10, Status: &status}

	s.mockUsecase.On("GetAllFeedbacksForProcedure", mock.Anything, procedureID, filter).Return(expectedFeedbacks, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/procedures/%s/feedback?status=new", procedureID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]FeedbackListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp["feedbacks"].Total)
	s.Len(resp["feedbacks"].Feedbacks, 1)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *FeedbackControllerTestSuite) TestUpdateFeedbackStatus_Success() {
	// Arrange
	feedbackID := "fb-123"
	adminResponse := "Thank you for your feedback."
	reqBody := FeedbackStatePatchRequest{Status: "resolved", AdminResponse: adminResponse}
	s.mockUsecase.On("UpdateFeedbackStatus", mock.Anything, feedbackID, "test-user-id", domain.FeedbackStatus("resolved"), &adminResponse).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/feedback/%s", feedbackID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Feedback status updated successfully")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *FeedbackControllerTestSuite) TestUpdateFeedbackStatus_PermissionDenied() {
	// Arrange
	feedbackID := "fb-123"
	reqBody := FeedbackStatePatchRequest{Status: "resolved"}
	s.mockUsecase.On("UpdateFeedbackStatus", mock.Anything, feedbackID, "test-user-id", domain.FeedbackStatus("resolved"), &reqBody.AdminResponse).Return(domain.ErrPermissionDenied).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/feedback/%s", feedbackID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusForbidden, w.Code)
	s.Contains(w.Body.String(), domain.ErrPermissionDenied.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *FeedbackControllerTestSuite) TestGetAllFeedbacks_Success() {
	// Arrange
	expectedFeedbacks := []*domain.Feedback{{ID: "fb-1", Content: "Test content"}}
	var total int64 = 1
	status := "in_progress"
	procID := "proc-xyz"
	filter := &domain.FeedbackFilter{Page: 2, Limit: 5, Status: &status, ProcedureID: &procID}

	s.mockUsecase.On("GetAllFeedbacks", mock.Anything, filter).Return(expectedFeedbacks, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	url := "/feedback?page=2&limit=5&status=in_progress&procedure_id=proc-xyz"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]FeedbackListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp["feedbacks"].Total)
	s.Len(resp["feedbacks"].Feedbacks, 1)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *FeedbackControllerTestSuite) TestGetAllFeedbacks_InvalidQueryParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/feedback?page=invalid", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), domain.ErrInvalidQueryParam.Error())
}

// TestFeedbackControllerTestSuite runs the entire suite
func TestFeedbackControllerTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackControllerTestSuite))
}
