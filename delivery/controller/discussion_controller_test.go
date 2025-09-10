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

// MockPostUsecase is a mock for the IPostUseCase interface
type MockPostUsecase struct {
	mock.Mock
}

func (m *MockPostUsecase) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostUsecase) GetPosts(ctx context.Context, opts domain.PostFilters) ([]*domain.Post, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostUsecase) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostUsecase) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostUsecase) DeletePost(ctx context.Context, id, userID, role string) error {
	args := m.Called(ctx, id, userID, role)
	return args.Error(0)
}

// PostControllerTestSuite is the test suite for PostController
type PostControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockPostUsecase
	controller  *PostController
}

// SetupTest is run before each test in the suite
func (s *PostControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.mockUsecase = new(MockPostUsecase)
	s.controller = NewPostController(s.mockUsecase)

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Set("userRole", domain.RoleUser)
		c.Next()
	}

	discussions := s.router.Group("/discussions")
	{
		discussions.POST("", authMiddleware, s.controller.CreatePost)
		discussions.GET("", s.controller.GetPosts)
		discussions.GET("/:id", s.controller.GetPostByID)
		discussions.PATCH("/:id", authMiddleware, s.controller.UpdatePost)
		discussions.DELETE("/:id", authMiddleware, s.controller.DeletePost)
	}
}

func (s *PostControllerTestSuite) TestCreatePost_Success() {
	// Arrange
	reqBody := CreatePostDTO{Title: "New Post", Content: "Some content"}
	expectedPost := &domain.Post{ID: "post-123", UserID: "test-user-id", Title: "New Post", Content: "Some content"}

	s.mockUsecase.On("CreatePost", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(expectedPost, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/discussions", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	var resp map[string]*domain.Post
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedPost.ID, resp["post"].ID)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestGetPosts_Success() {
	// Arrange
	expectedPosts := []*domain.Post{{ID: "post-1", Title: "Test Post"}}
	var total int64 = 1
	var page int64 = 0
	var limit int64 = 10
	title := "Test"
	userId := "user-1"
	expectedFilters := domain.PostFilters{
		Title:       &title,
		UserId:      &userId,
		ProcedureID: []string{"proc-1"},
		Tags:        []string{"tag1"},
		Page:        page,
		Limit:       limit,
		SortBy:      "created_at",
		SortOrder:   domain.SortOrder("desc"),
	}

	s.mockUsecase.On("GetPosts", mock.Anything, expectedFilters).Return(expectedPosts, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	url := "/discussions?title=Test&userId=user-1&procedure_ids=proc-1&tags=tag1"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]PaginatedPostsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp["Posts"].Total)
	s.Len(resp["Posts"].Posts, 1)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestGetPostByID_Success() {
	// Arrange
	postID := "post-123"
	expectedPost := &domain.Post{ID: postID, Title: "A Post"}
	s.mockUsecase.On("GetPostByID", mock.Anything, postID).Return(expectedPost, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/discussions/%s", postID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]*domain.Post
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedPost.ID, resp["post"].ID)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestGetPostByID_NotFound() {
	// Arrange
	postID := "non-existent-id"
	s.mockUsecase.On("GetPostByID", mock.Anything, postID).Return(nil, domain.ErrPostNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/discussions/%s", postID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code) // HandleError maps ErrNotFound to 404
	s.Contains(w.Body.String(), domain.ErrPostNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestUpdatePost_Success() {
	// Arrange
	postID := "post-to-update"
	reqBody := UpdatePostDTO{Title: "Updated Title"}
	updatedPost := &domain.Post{ID: postID, Title: "Updated Title", UserID: "test-user-id"}
	s.mockUsecase.On("UpdatePost", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(updatedPost, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/discussions/%s", postID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]*domain.Post
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(updatedPost.Title, resp["post"].Title)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestUpdatePost_Forbidden() {
	// Arrange
	postID := "post-to-update"
	reqBody := UpdatePostDTO{Title: "Updated Title"}
	s.mockUsecase.On("UpdatePost", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(nil, domain.ErrPermissionDenied).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/discussions/%s", postID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusForbidden, w.Code)
	s.Contains(w.Body.String(), domain.ErrPermissionDenied.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestDeletePost_Success() {
	// Arrange
	postID := "post-to-delete"
	s.mockUsecase.On("DeletePost", mock.Anything, postID, "test-user-id", string(domain.RoleUser)).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/discussions/%s", postID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNoContent, w.Code)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PostControllerTestSuite) TestDeletePost_NotFound() {
	// Arrange
	postID := "post-to-delete"
	s.mockUsecase.On("DeletePost", mock.Anything, postID, "test-user-id", string(domain.RoleUser)).Return(domain.ErrPostNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/discussions/%s", postID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrPostNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestPostControllerTestSuite runs the entire suite
func TestPostControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PostControllerTestSuite))
}
