package controller_test

import (
	. "EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockAIChatUsecase is a mock for the IAIChatUsecase interface
type MockAIChatUsecase struct {
	mock.Mock
}

func (m *MockAIChatUsecase) AIchat(ctx context.Context, userId, query string) (*domain.AIChat, error) {
	args := m.Called(ctx, userId, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AIChat), args.Error(1)
}

func (m *MockAIChatUsecase) AIHistory(ctx context.Context, userId string, page, limit int64) ([]*domain.AIChat, int64, error) {
	args := m.Called(ctx, userId, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.AIChat), args.Get(1).(int64), args.Error(2)
}

// AIChatControllerTestSuite is the test suite for AIChatController
type AIChatControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockAIChatUsecase
	controller  *AIChatController
}

// SetupTest is run before each test in the suite
func (s *AIChatControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.mockUsecase = new(MockAIChatUsecase)
	s.controller = NewAIChatController(s.mockUsecase)

	// A helper middleware to inject userID into the context for authenticated routes
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	}

	// Setup routes
	s.router.POST("/ai/guide", authMiddleware, s.controller.AIChatController)
	s.router.GET("/ai/history", authMiddleware, s.controller.AIChatHistoryController)

	// Add routes without auth middleware for testing unauthorized access
	s.router.POST("/unauth/ai/guide", s.controller.AIChatController)
	s.router.GET("/unauth/ai/history", s.controller.AIChatHistoryController)
}

func (s *AIChatControllerTestSuite) TestAIChat_Success() {
	// Arrange
	reqBody := AIChatRequest{Query: "test query"}
	expectedResponse := &domain.AIChat{
		ID:        "chat-id-1",
		UserID:    "test-user-id",
		Request:   "test query",
		Response:  "This is a test response.",
		Timestamp: time.Now(),
	}
	s.mockUsecase.On("AIchat", mock.Anything, "test-user-id", "test query").Return(expectedResponse, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/guide", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp AIConversationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedResponse.ID, resp.ID)
	s.Equal(expectedResponse.UserID, resp.UserID)
	s.Equal(expectedResponse.Response, resp.Response)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *AIChatControllerTestSuite) TestAIChat_Unauthorized() {
	// Arrange
	reqBody := AIChatRequest{Query: "test"}

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/unauth/ai/guide", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "User ID not found in token")
}

func (s *AIChatControllerTestSuite) TestAIChat_InvalidJSON() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/guide", bytes.NewBufferString(`{"query":`)) // Malformed JSON
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid request")
}

func (s *AIChatControllerTestSuite) TestAIChat_UsecaseError() {
	// Arrange
	reqBody := AIChatRequest{Query: "test query"}
	// Use a domain error that maps to a 500 in the HandleError helper
	s.mockUsecase.On("AIchat", mock.Anything, "test-user-id", "test query").Return(nil, domain.ErrUnableToFetchData).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/guide", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusInternalServerError, w.Code)
	// Check for the generic error message produced by the default case in HandleError
	s.Contains(w.Body.String(), "An unexpected internal error occurred")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *AIChatControllerTestSuite) TestAIChatHistory_Success() {
	// Arrange
	expectedHistory := []*domain.AIChat{
		{ID: "hist-1", UserID: "test-user-id", Request: "q1", Response: "a1"},
	}
	var total int64 = 1
	s.mockUsecase.On("AIHistory", mock.Anything, "test-user-id", int64(1), int64(10)).Return(expectedHistory, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ai/history?page=1&limit=10", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp PaginatedAIHisoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp.Pagination.Total)
	s.Equal(int64(1), resp.Pagination.Page)
	s.Equal(int64(10), resp.Pagination.Limit)
	s.Len(resp.Data, 1)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *AIChatControllerTestSuite) TestAIChatHistory_InvalidPageParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ai/history?page=abc", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid 'page' parameter")
}

func (s *AIChatControllerTestSuite) TestAIChatHistory_InvalidLimitParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ai/history?limit=0", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid 'limit' parameter")
}

func (s *AIChatControllerTestSuite) TestAIChatHistory_UsecaseError() {
	// Arrange
	s.mockUsecase.On("AIHistory", mock.Anything, "test-user-id", int64(1), int64(10)).Return(nil, int64(0), domain.ErrNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ai/history", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	// HandleError maps domain.ErrNotFound to a 404 status code
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestAIChatControllerTestSuite runs the entire suite
func TestAIChatControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AIChatControllerTestSuite))
}
