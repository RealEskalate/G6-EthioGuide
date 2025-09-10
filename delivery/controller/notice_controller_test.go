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

// MockNoticeUsecase is a mock for the INoticeUseCase interface
type MockNoticeUsecase struct {
	mock.Mock
}

func (m *MockNoticeUsecase) CreateNotice(ctx context.Context, notice *domain.Notice) error {
	args := m.Called(ctx, notice)
	return args.Error(0)
}

func (m *MockNoticeUsecase) GetNoticesByFilter(ctx context.Context, filter *domain.NoticeFilter) ([]*domain.Notice, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Notice), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoticeUsecase) UpdateNotice(ctx context.Context, id string, notice *domain.Notice) error {
	args := m.Called(ctx, id, notice)
	return args.Error(0)
}

func (m *MockNoticeUsecase) DeleteNotice(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// NoticeControllerTestSuite is the test suite for NoticeController
type NoticeControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockNoticeUsecase
	controller  *NoticeController
}

// SetupTest is run before each test in the suite
func (s *NoticeControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.mockUsecase = new(MockNoticeUsecase)
	s.controller = NewNoticeController(s.mockUsecase)

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-org-id")
		c.Next()
	}

	notices := s.router.Group("/notices")
	{
		notices.POST("", authMiddleware, s.controller.CreateNotice)
		notices.GET("", s.controller.GetNoticesByFilter)
		notices.PATCH("/:id", authMiddleware, s.controller.UpdateNotice)
		notices.DELETE("/:id", authMiddleware, s.controller.DeleteNotice)
	}
}

func (s *NoticeControllerTestSuite) TestCreateNotice_Success() {
	// Arrange
	reqBody := NoticeDTO{Title: "New Notice", Content: "Important update"}
	s.mockUsecase.On("CreateNotice", mock.Anything, mock.AnythingOfType("*domain.Notice")).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/notices", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Notice Created successfully.")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestCreateNotice_UsecaseError() {
	// Arrange
	reqBody := NoticeDTO{Title: "New Notice", Content: "Important update"}
	s.mockUsecase.On("CreateNotice", mock.Anything, mock.AnythingOfType("*domain.Notice")).Return(domain.ErrNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/notices", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestGetNoticesByFilter_Success() {
	// Arrange
	expectedNotices := []*domain.Notice{{ID: "notice-1", Title: "Test Notice"}}
	var total int64 = 1
	filter := &domain.NoticeFilter{
		OrganizationID: "org-123",
		Tags:           []string{"urgent", "update"},
		Page:           1,
		Limit:          10,
		SortBy:         "createdAt",
		SortOrder:      domain.SortDesc,
	}
	s.mockUsecase.On("GetNoticesByFilter", mock.Anything, filter).Return(expectedNotices, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	url := "/notices?organizationId=org-123&tags=urgent,update&sortBy=createdAt&sortOrder=DESC"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp NoticeListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp.Total)
	s.Len(resp.Data, 1)
	s.Equal("Test Notice", resp.Data[0].Title)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestGetNoticesByFilter_InvalidPageParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/notices?page=abc", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid 'page' parameter")
}

func (s *NoticeControllerTestSuite) TestUpdateNotice_Success() {
	// Arrange
	noticeID := "notice-123"
	reqBody := NoticeDTO{Title: "Updated Title", Content: "Updated Content"}
	s.mockUsecase.On("UpdateNotice", mock.Anything, noticeID, mock.AnythingOfType("*domain.Notice")).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/notices/%s", noticeID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Notice Updated successfully.")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestUpdateNotice_NotFound() {
	// Arrange
	noticeID := "notice-404"
	reqBody := NoticeDTO{Title: "Updated Title"}
	s.mockUsecase.On("UpdateNotice", mock.Anything, noticeID, mock.AnythingOfType("*domain.Notice")).Return(domain.ErrNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/notices/%s", noticeID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestDeleteNotice_Success() {
	// Arrange
	noticeID := "notice-to-delete"
	s.mockUsecase.On("DeleteNotice", mock.Anything, noticeID).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/notices/%s", noticeID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Notice Deleted successfully.")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *NoticeControllerTestSuite) TestDeleteNotice_NotFound() {
	// Arrange
	noticeID := "notice-404"
	s.mockUsecase.On("DeleteNotice", mock.Anything, noticeID).Return(domain.ErrNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/notices/%s", noticeID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestNoticeControllerTestSuite runs the entire suite
func TestNoticeControllerTestSuite(t *testing.T) {
	suite.Run(t, new(NoticeControllerTestSuite))
}
