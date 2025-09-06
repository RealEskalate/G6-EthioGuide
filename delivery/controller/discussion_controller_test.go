package controller_test

// import (
// 	"EthioGuide/delivery/controller"
// 	"EthioGuide/domain"
// 	"context"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type MockPostUseCase struct {
// 	mock.Mock
// }

// func (m *MockPostUseCase) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
// 	args := m.Called(ctx, post)
// 	return args.Get(0).(*domain.Post), args.Error(1)
// }
// func (m *MockPostUseCase) GetPosts(ctx context.Context, opts domain.PostFilters) ([]*domain.Post, int64, error) {
// 	args := m.Called(ctx, opts)
// 	return args.Get(0).([]*domain.Post), args.Get(1).(int64), args.Error(2)
// }
// func (m *MockPostUseCase) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(*domain.Post), args.Error(1)
// }
// func (m *MockPostUseCase) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
// 	args := m.Called(ctx, post)
// 	return args.Get(0).(*domain.Post), args.Error(1)
// }
// func (m *MockPostUseCase) DeletePost(ctx context.Context, id, userID, role string) error {
// 	args := m.Called(ctx, id, userID, role)
// 	return args.Error(0)
// }

// type PostControllerTestSuite struct {
// 	suite.Suite
// 	mockUseCase *MockPostUseCase
// 	router      *gin.Engine
// }

// func (s *PostControllerTestSuite) SetupTest() {
// 	s.mockUseCase = new(MockPostUseCase)
// 	pc := controller.NewPostController(s.mockUseCase)

// 	gin.SetMode(gin.TestMode)
// 	s.router = gin.Default()

// 	// Set both userID and userRole in the context
// 	s.router.Use(func(c *gin.Context) {
// 		c.Set("userID", "123")
// 		c.Set("userRole", domain.RoleAdmin)
// 		c.Next()
// 	})

// 	s.router.POST("/discussions", pc.CreatePost)
// 	s.router.GET("/discussions", pc.GetPosts)
// 	s.router.GET("/discussions/:id", pc.GetPostByID)
// 	s.router.PUT("/discussions/:id", pc.UpdatePost)
// 	s.router.DELETE("/discussions/:id", pc.DeletePost)
// }

// func (s *PostControllerTestSuite) TestCreatePost_Success() {
// 	body := `{"title": "Hello", "content": "World", "procedures": [], "tags": []}`
// 	req := httptest.NewRequest(http.MethodPost, "/discussions", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	expectedPost := &domain.Post{
// 		UserID:     "123",
// 		Title:      "Hello",
// 		Content:    "World",
// 		Procedures: []string{},
// 		Tags:       []string{},
// 	}
// 	s.mockUseCase.On("CreatePost", mock.Anything, mock.MatchedBy(func(post *domain.Post) bool {
// 		return post.UserID == expectedPost.UserID &&
// 			post.Title == expectedPost.Title &&
// 			post.Content == expectedPost.Content
// 	})).Return(expectedPost, nil)

// 	s.router.ServeHTTP(w, req)
// 	s.Equal(http.StatusCreated, w.Code)
// 	s.Contains(w.Body.String(), "post")
// }

// func (s *PostControllerTestSuite) TestGetPosts_Success() {
// 	expectedPosts := []*domain.Post{
// 		{ID: "1", Title: "Hello", Content: "World"},
// 	}
// 	s.mockUseCase.On("GetPosts", mock.Anything, mock.AnythingOfType("domain.PostFilters")).Return(expectedPosts, int64(1), nil)

// 	req := httptest.NewRequest(http.MethodGet, "/discussions?title=Hello", nil)
// 	w := httptest.NewRecorder()
// 	s.router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// 	s.Contains(w.Body.String(), "posts")
// }

// func (s *PostControllerTestSuite) TestGetPostByID_Success() {
// 	expectedPost := &domain.Post{ID: "1", Title: "Hello", Content: "World"}
// 	s.mockUseCase.On("GetPostByID", mock.Anything, "1").Return(expectedPost, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/discussions/1", nil)
// 	w := httptest.NewRecorder()
// 	s.router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// 	s.Contains(w.Body.String(), "post")
// }

// func (s *PostControllerTestSuite) TestUpdatePost_Success() {
// 	body := `{"title": "Updated", "content": "World", "procedures": [], "tags": []}`
// 	expectedPost := &domain.Post{
// 		ID:         "1",
// 		UserID:     "123",
// 		Title:      "Updated",
// 		Content:    "World",
// 		Procedures: []string{},
// 		Tags:       []string{},
// 	}
// 	s.mockUseCase.On("UpdatePost", mock.Anything, mock.MatchedBy(func(post *domain.Post) bool {
// 		return post.ID == expectedPost.ID &&
// 			post.UserID == expectedPost.UserID &&
// 			post.Title == expectedPost.Title &&
// 			post.Content == expectedPost.Content
// 	})).Return(expectedPost, nil)

// 	req := httptest.NewRequest(http.MethodPut, "/discussions/1", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	s.router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// 	s.Contains(w.Body.String(), "post")
// }

// func (s *PostControllerTestSuite) TestDeletePost_Success() {
// 	// Set up expected behavior
// 	s.mockUseCase.On("DeletePost", mock.Anything, "1", "123", string(domain.RoleAdmin)).Return(nil)

// 	// Create request
// 	req := httptest.NewRequest(http.MethodDelete, "/discussions/1", nil)
// 	w := httptest.NewRecorder()

// 	// Inject userRole into context
// 	s.router.Use(func(c *gin.Context) {
// 		c.Set("userRole", domain.RoleAdmin)
// 		c.Next()
// 	})

// 	// Serve request
// 	s.router.ServeHTTP(w, req)

// 	// Assert
// 	s.Equal(http.StatusNoContent, w.Code)
// 	s.mockUseCase.AssertExpectations(s.T())
// }

// func (s *PostControllerTestSuite) TestDeletePost_Failure() {
// 	s.mockUseCase.On("DeletePost", mock.Anything, "1", "123", string(domain.RoleAdmin)).Return(errors.New("delete failed"))

// 	req := httptest.NewRequest(http.MethodDelete, "/discussions/1", nil)
// 	w := httptest.NewRecorder()

// 	// Create a test context and set the userRole
// 	c, _ := gin.CreateTestContext(w)
// 	c.Request = req
// 	c.Set("userRole", domain.RoleAdmin) // Set the userRole in the Gin context

// 	s.router.ServeHTTP(w, req)

// 	s.Equal(http.StatusInternalServerError, w.Code)
// 	s.mockUseCase.AssertExpectations(s.T())
// }

// func TestPostControllerTestSuite(t *testing.T) {
// 	suite.Run(t, new(PostControllerTestSuite))
// }
