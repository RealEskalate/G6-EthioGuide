package controller_test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"EthioGuide/delivery/controller"
// 	"EthioGuide/domain"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// // --- Mock Usecase ---
// type MockCategoryUsecase struct {
// 	mock.Mock
// }

// func (m *MockCategoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) error {
// 	args := m.Called(ctx, category)
// 	return args.Error(0)
// }

// func (m *MockCategoryUsecase) GetCategories(c context.Context, options *domain.CategorySearchAndFilter) ([]*domain.Category, int64, error) {
// 	args := m.Called(c, options)
// 	var catagories []*domain.Category
// 	if args.Get(0) != nil {
// 		catagories = args.Get(0).([]*domain.Category)
// 	}
// 	var total int64
// 	if args.Get(1) != nil {
// 		total = args.Get(1).(int64)
// 	}
// 	return catagories, total, args.Error(2)
// }

// // --- Test Suite ---
// type CategoryControllerTestSuite struct {
// 	suite.Suite
// 	router      *gin.Engine
// 	controller  *controller.CategoryController
// 	mockUsecase *MockCategoryUsecase
// 	recorder    *httptest.ResponseRecorder
// }

// // TestCategoryControllerTestSuite runs the test suite
// func TestCategoryControllerTestSuite(t *testing.T) {
// 	suite.Run(t, new(CategoryControllerTestSuite))
// }

// func (s *CategoryControllerTestSuite) SetupSuite() {
// 	gin.SetMode(gin.TestMode)
// }

// func (s *CategoryControllerTestSuite) SetupTest() {
// 	s.recorder = httptest.NewRecorder()
// 	s.router = gin.Default()
// 	s.mockUsecase = new(MockCategoryUsecase)
// 	s.controller = controller.NewCategoryController(s.mockUsecase)

// 	// Add middleware to set userID in context
// 	s.router.Use(func(c *gin.Context) {
// 		c.Set("userID", "org123")
// 		c.Next()
// 	})

// 	s.router.POST("/category", s.controller.CreateCategory)
// 	s.router.GET("/category", s.controller.GetCategory)
// }

// func (s *CategoryControllerTestSuite) TestCreateCategory() {
// 	reqBody := controller.CreateCategoryRequest{
// 		// Remove OrganizationID as it should come from context
// 		ParentID: "parent456",
// 		Title:    "New Category",
// 	}
// 	jsonBody, _ := json.Marshal(reqBody)

// 	s.Run("Success", func() {
// 		s.SetupTest()
// 		s.mockUsecase.On("CreateCategory", mock.Anything, mock.MatchedBy(func(c *domain.Category) bool {
// 			return c.OrganizationID == "org123" && c.Title == reqBody.Title
// 		})).Return(nil).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusCreated, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), "Category created successfully")
// 		s.mockUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Invalid Request", func() {
// 		s.SetupTest()
// 		req, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewBuffer([]byte("not-json")))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusBadRequest, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), "error")
// 	})

// 	s.Run("Failure - Usecase Error", func() {
// 		s.SetupTest()
// 		s.mockUsecase.On("CreateCategory", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(domain.ErrConflict).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusConflict, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), domain.ErrConflict.Error())
// 		s.mockUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *CategoryControllerTestSuite) TestGetCategory() {
// 	s.Run("Success - Default Parameters", func() {
// 		s.SetupTest()

// 		categories := []*domain.Category{
// 			{
// 				ID:             "cat1",
// 				OrganizationID: "org123",
// 				Title:          "Category 1",
// 			},
// 			{
// 				ID:             "cat2",
// 				OrganizationID: "org123",
// 				ParentID:       "cat1",
// 				Title:          "Category 2",
// 			},
// 		}

// 		s.mockUsecase.On("GetCategories", mock.Anything, mock.AnythingOfType("*domain.CategorySearchAndFilter")).Return(categories, int64(2), nil).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/category", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)

// 		var response controller.PaginatedCategoryResponse
// 		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
// 		s.NoError(err)
// 		s.Equal(int64(2), response.Total)
// 		s.Equal(int64(1), response.Page)
// 		s.Equal(int64(10), response.Limit)
// 		s.Len(response.Data, 2)
// 		s.Equal("cat1", response.Data[0].ID)
// 		s.Equal("Category 1", response.Data[0].Title)
// 		s.Equal("cat2", response.Data[1].ID)
// 		s.Equal("Category 2", response.Data[1].Title)
// 		s.Equal("cat1", response.Data[1].ParentID)

// 		s.mockUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Success - Custom Parameters", func() {
// 		s.SetupTest()

// 		categories := []*domain.Category{
// 			{
// 				ID:             "cat3",
// 				OrganizationID: "org456",
// 				Title:          "Category 3",
// 			},
// 		}

// 		s.mockUsecase.On("GetCategories", mock.Anything, mock.AnythingOfType("*domain.CategorySearchAndFilter")).Return(categories, int64(1), nil).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/category?page=2&limit=5&sortBy=title&sortOrder=asc&parentID=parent123&organizationID=org456&title=Category", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)

// 		var response controller.PaginatedCategoryResponse
// 		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
// 		s.NoError(err)
// 		s.Equal(int64(1), response.Total)
// 		s.Equal(int64(2), response.Page)
// 		s.Equal(int64(5), response.Limit)
// 		s.Len(response.Data, 1)
// 		s.Equal("cat3", response.Data[0].ID)
// 		s.Equal("Category 3", response.Data[0].Title)

// 		s.mockUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Invalid Page Parameter", func() {
// 		s.SetupTest()

// 		req, _ := http.NewRequest(http.MethodGet, "/category?page=invalid", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusBadRequest, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), "Invalid 'page' parameter")
// 	})

// 	s.Run("Failure - Invalid Limit Parameter", func() {
// 		s.SetupTest()

// 		req, _ := http.NewRequest(http.MethodGet, "/category?limit=invalid", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusBadRequest, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), "Invalid 'limit' parameter")
// 	})

// 	s.Run("Failure - Usecase Error", func() {
// 		s.SetupTest()

// 		s.mockUsecase.On("GetCategories", mock.Anything, mock.AnythingOfType("*domain.CategorySearchAndFilter")).Return(nil, int64(0), domain.ErrNotFound).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/category", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusNotFound, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), domain.ErrNotFound.Error())

// 		s.mockUsecase.AssertExpectations(s.T())
// 	})
// }
