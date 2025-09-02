package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"EthioGuide/delivery/controller"
	"EthioGuide/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock Usecase ---
type MockCategoryUsecase struct {
	mock.Mock
}

func (m *MockCategoryUsecase) CreateCategory(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

// --- Test Suite ---
type CategoryControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	controller  *controller.CategoryController
	mockUsecase *MockCategoryUsecase
	recorder    *httptest.ResponseRecorder
}

func (s *CategoryControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *CategoryControllerTestSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUsecase = new(MockCategoryUsecase)
	s.controller = controller.NewCategoryController(s.mockUsecase)
	s.router.POST("/category", s.controller.CreateCategory)
}

func TestCategoryControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerTestSuite))
}

// --- Tests ---

func (s *CategoryControllerTestSuite) TestCreateCategory() {
	reqBody := controller.CreateCategoryRequest{
		OrganizationID: "org123",
		ParentID:       "parent456",
		Title:          "New Category",
	}
	jsonBody, _ := json.Marshal(reqBody)

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUsecase.On("CreateCategory", mock.Anything, mock.MatchedBy(func(c *domain.Category) bool {
			return c.OrganizationID == reqBody.OrganizationID && c.Title == reqBody.Title
		})).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusCreated, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Category created successfully")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid Request", func() {
		s.SetupTest()
		req, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewBuffer([]byte("not-json")))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusBadRequest, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "error")
	})

	s.Run("Failure - Usecase Error", func() {
		s.SetupTest()
		s.mockUsecase.On("CreateCategory", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(domain.ErrConflict).Once()

		req, _ := http.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusConflict, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrConflict.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})
}
