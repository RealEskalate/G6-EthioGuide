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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockCategoryUsecase is a mock for the ICategoryUsecase interface
type MockCategoryUsecase struct {
	mock.Mock
}

func (m *MockCategoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryUsecase) GetCategories(ctx context.Context, options *domain.CategorySearchAndFilter) ([]*domain.Category, int64, error) {
	args := m.Called(ctx, options)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Category), args.Get(1).(int64), args.Error(2)
}

// CategoryControllerTestSuite is the test suite for CategoryController
type CategoryControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockCategoryUsecase
	controller  *CategoryController
}

// SetupTest is run before each test in the suite
func (s *CategoryControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.mockUsecase = new(MockCategoryUsecase)
	s.controller = NewCategoryController(s.mockUsecase)

	// A helper middleware to inject userID into the context for authenticated routes
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-org-id")
		c.Next()
	}

	// Setup routes
	s.router.POST("/categories", authMiddleware, s.controller.CreateCategory)
	s.router.GET("/categories", s.controller.GetCategory)
}

func (s *CategoryControllerTestSuite) TestCreateCategory_Success() {
	// Arrange
	reqBody := CreateCategoryRequest{
		ParentID: "parent-123",
		Title:    "New Test Category",
	}
	// The mock should be called with the userID from the context
	expectedCategoryArg := &domain.Category{
		OrganizationID: "test-org-id",
		ParentID:       "parent-123",
		Title:          "New Test Category",
	}
	s.mockUsecase.On("CreateCategory", mock.Anything, expectedCategoryArg).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/categories", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Category created successfully")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *CategoryControllerTestSuite) TestCreateCategory_InvalidJSON() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(`{"title":`)) // Malformed JSON
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *CategoryControllerTestSuite) TestCreateCategory_UsecaseError() {
	// Arrange
	reqBody := CreateCategoryRequest{
		Title: "Failing Category",
	}
	expectedCategoryArg := &domain.Category{
		OrganizationID: "test-org-id",
		Title:          "Failing Category",
	}
	s.mockUsecase.On("CreateCategory", mock.Anything, expectedCategoryArg).Return(domain.ErrConflict).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/categories", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusConflict, w.Code)
	s.Contains(w.Body.String(), domain.ErrConflict.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *CategoryControllerTestSuite) TestGetCategory_Success_WithFilters() {
	// Arrange
	expectedCategories := []*domain.Category{
		{ID: "cat-1", Title: "Filtered Category", OrganizationID: "org-xyz"},
	}
	var total int64 = 1
	expectedOptions := &domain.CategorySearchAndFilter{
		Page:           2,
		Limit:          5,
		SortBy:         "title",
		SortOrder:      domain.SortAsc,
		ParentID:       "parent-1",
		OrganizationID: "org-xyz",
		Title:          "Filtered",
	}

	s.mockUsecase.On("GetCategories", mock.Anything, expectedOptions).Return(expectedCategories, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	url := "/categories?page=2&limit=5&sortBy=title&sortOrder=asc&parentID=parent-1&organizationID=org-xyz&title=Filtered"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp PaginatedCategoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp.Total)
	s.Equal(int64(2), resp.Page)
	s.Equal(int64(5), resp.Limit)
	s.Len(resp.Data, 1)
	s.Equal("Filtered Category", resp.Data[0].Title)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *CategoryControllerTestSuite) TestGetCategory_Success_Defaults() {
	// Arrange
	expectedCategories := []*domain.Category{}
	var total int64 = 0
	expectedOptions := &domain.CategorySearchAndFilter{
		Page:      1,
		Limit:     10,
		SortOrder: domain.SortDesc,
	}
	s.mockUsecase.On("GetCategories", mock.Anything, expectedOptions).Return(expectedCategories, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *CategoryControllerTestSuite) TestGetCategory_InvalidPageParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/categories?page=invalid", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid 'page' parameter")
}

func (s *CategoryControllerTestSuite) TestGetCategory_InvalidLimitParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/categories?limit=0", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid 'limit' parameter")
}

func (s *CategoryControllerTestSuite) TestGetCategory_UsecaseError() {
	// Arrange
	expectedOptions := &domain.CategorySearchAndFilter{
		Page:      1,
		Limit:     10,
		SortOrder: domain.SortDesc,
	}
	s.mockUsecase.On("GetCategories", mock.Anything, expectedOptions).Return(nil, int64(0), domain.ErrUnableToFetchData).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "An unexpected internal error occurred")
	s.mockUsecase.AssertExpectations(s.T())
}

// TestCategoryControllerTestSuite runs the entire suite
func TestCategoryControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerTestSuite))
}
