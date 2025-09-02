package controller_test

import (
	"EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockPostUseCase struct {
	mock.Mock
}

func (m *MockPostUseCase) CreatePost(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

type PostControllerTestSuite struct {
	suite.Suite
	mockUseCase *MockPostUseCase
	router      *gin.Engine
}

func (s *PostControllerTestSuite) SetupTest() {
	s.mockUseCase = new(MockPostUseCase)
	pc := controller.NewPostController(s.mockUseCase)

	s.router = gin.Default()
	s.router.Use(func(c *gin.Context) {
		c.Set("userID", "123")
		c.Next()
	})
	s.router.POST("/discussions", pc.CreatePost)
}

func (s *PostControllerTestSuite) TestCreatePost_Success() {
	body := `{"title": "Hello", "content": "World", "procedures": [], "tags": []}`
	req := httptest.NewRequest(http.MethodPost, "/discussions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// inject userID in context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Define the expected Post object
	expectedPost := &domain.Post{
		UserID:     "123",
		Title:      "Hello",
		Content:    "World",
		Procedures: []string{},
		Tags:       []string{},
	}

	// Use mock.MatchedBy to compare the actual Post object with the expected one
	s.mockUseCase.On("CreatePost", mock.Anything, mock.MatchedBy(func(post *domain.Post) bool {
		// Custom matching logic: compare relevant fields
		return post.UserID == expectedPost.UserID &&
			post.Title == expectedPost.Title &&
			post.Content == expectedPost.Content
		// Add more fields to compare as needed
	})).Return(nil)

	// call handler
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Post created successfully")

	// Assert that CreatePost was called with the expected Post object
	s.mockUseCase.AssertCalled(s.T(), "CreatePost", mock.Anything, mock.MatchedBy(func(post *domain.Post) bool {
		return post.UserID == expectedPost.UserID &&
			post.Title == expectedPost.Title &&
			post.Content == expectedPost.Content
	}))
}

func TestPostControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PostControllerTestSuite))
}
