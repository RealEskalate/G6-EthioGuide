// filepath: usecase/discussions_usecase_test.go
package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockPostRepository is a mock implementation of domain.IPostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostRepository) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostRepository) GetPosts(ctx context.Context, opts domain.PostFilters) ([]*domain.Post, int64, error) {
	args := m.Called(ctx, opts)
	posts, _ := args.Get(0).([]*domain.Post)
	total, _ := args.Get(1).(int64)
	return posts, total, args.Error(2)
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

// In MockPostRepository
func (m *MockPostRepository) DeletePost(ctx context.Context, id, userID, role string) error {
	args := m.Called(ctx, id, userID, role)
	return args.Error(0)
}




// PostUseCaseTestSuite defines the test suite for PostUseCase
type PostUseCaseTestSuite struct {
	suite.Suite
	mockPostRepo *MockPostRepository
	postUseCase  domain.IPostUseCase
}

// SetupTest initializes the test suite before each test
func (suite *PostUseCaseTestSuite) SetupTest() {
	suite.mockPostRepo = new(MockPostRepository)
	suite.postUseCase = NewPostUseCase(suite.mockPostRepo, time.Second*5) // Provide a timeout
}

func (suite *PostUseCaseTestSuite) TestCreatePost() {
	testCases := []struct {
		name          string
		post          *domain.Post
		expectedPost  *domain.Post
		expectedError error
		mockSetup     func(post *domain.Post, expectedPost *domain.Post, expectedError error)
	}{
		{
			name: "Success",
			post: &domain.Post{
				UserID:     "user123",
				Title:      "Test Post",
				Content:    "This is a test post.",
				Procedures: []string{"Step 1", "Step 2"},
				Tags:       []string{"test", "post"},
			},
			expectedPost: &domain.Post{
				ID:         "1", // Mocked ID
				UserID:     "user123",
				Title:      "Test Post",
				Content:    "This is a test post.",
				Procedures: []string{"Step 1", "Step 2"},
				Tags:       []string{"test", "post"},
			},
			expectedError: nil,
			mockSetup: func(post *domain.Post, expectedPost *domain.Post, expectedError error) {
				suite.mockPostRepo.On("CreatePost", mock.Anything, mock.MatchedBy(func(p *domain.Post) bool {
					return reflect.DeepEqual(p, post)
				})).Return(expectedPost, expectedError).Once()
			},
		},
		{
			name: "Failure",
			post: &domain.Post{
				UserID:     "user123",
				Title:      "Test Post",
				Content:    "This is a test post.",
				Procedures: []string{"Step 1", "Step 2"},
				Tags:       []string{"test", "post"},
			},
			expectedPost:  nil,
			expectedError: errors.New("creation failed"),
			mockSetup: func(post *domain.Post, expectedPost *domain.Post, expectedError error) {
				suite.mockPostRepo.On("CreatePost", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(nil, expectedError).Once()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup(tc.post, tc.expectedPost, tc.expectedError)

			createdPost, err := suite.postUseCase.CreatePost(context.Background(), tc.post)

			suite.Equal(tc.expectedError, err)
			suite.Equal(tc.expectedPost, createdPost)
			suite.mockPostRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *PostUseCaseTestSuite) TestGetPosts() {
	testCases := []struct {
		name          string
		opts          domain.PostFilters
		expectedPosts []*domain.Post
		expectedTotal int64
		expectedError error
		mockSetup     func(opts domain.PostFilters, expectedPosts []*domain.Post, expectedTotal int64, expectedError error)
	}{
		{
			name: "Success",
			opts: domain.PostFilters{
				Title:     stringPtr("Test"),
				SortBy:    "title",
				SortOrder: "asc",
				Page:      1,
				Limit:     10,
			},
			expectedPosts: []*domain.Post{
				{ID: "1", Title: "Test Post 1"},
				{ID: "2", Title: "Test Post 2"},
			},
			expectedTotal: 2,
			expectedError: nil,
			mockSetup: func(opts domain.PostFilters, expectedPosts []*domain.Post, expectedTotal int64, expectedError error) {
				suite.mockPostRepo.On("GetPosts", mock.Anything, opts).Return(expectedPosts, expectedTotal, expectedError).Once()
			},
		},
		{
			name: "No Posts Found",
			opts: domain.PostFilters{
				Title: stringPtr("NonExistent"),
			},
			expectedPosts: []*domain.Post{},
			expectedTotal: 0,
			expectedError: nil,
			mockSetup: func(opts domain.PostFilters, expectedPosts []*domain.Post, expectedTotal int64, expectedError error) {
				suite.mockPostRepo.On("GetPosts", mock.Anything, opts).Return(expectedPosts, expectedTotal, expectedError).Once()
			},
		},
		{
			name: "Error from Repository",
			opts: domain.PostFilters{
				Title: stringPtr("Error"),
			},
			expectedPosts: []*domain.Post{},
			expectedTotal: 0,
			expectedError: errors.New("db error"),
			mockSetup: func(opts domain.PostFilters, expectedPosts []*domain.Post, expectedTotal int64, expectedError error) {
				suite.mockPostRepo.On("GetPosts", mock.Anything, opts).Return(expectedPosts, expectedTotal, expectedError).Once()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup(tc.opts, tc.expectedPosts, tc.expectedTotal, tc.expectedError)

			posts, total, err := suite.postUseCase.GetPosts(context.Background(), tc.opts)

			suite.Equal(tc.expectedError, err)
			suite.Equal(tc.expectedPosts, posts)
			suite.Equal(tc.expectedTotal, total)
			suite.mockPostRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *PostUseCaseTestSuite) TestGetPostByID() {
	testCases := []struct {
		name          string
		id            string
		expectedPost  *domain.Post
		expectedError error
		mockSetup     func(id string, expectedPost *domain.Post, expectedError error)
	}{
		{
			name: "Success",
			id:   "1",
			expectedPost: &domain.Post{
				ID:    "1",
				Title: "Test Post",
			},
			expectedError: nil,
			mockSetup: func(id string, expectedPost *domain.Post, expectedError error) {
				suite.mockPostRepo.On("GetPostByID", mock.Anything, id).Return(expectedPost, expectedError).Once()
			},
		},
		{
			name:          "Not Found",
			id:            "999",
			expectedPost:  nil,
			expectedError: errors.New("not found"),
			mockSetup: func(id string, expectedPost *domain.Post, expectedError error) {
				suite.mockPostRepo.On("GetPostByID", mock.Anything, id).Return(nil, expectedError).Once()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup(tc.id, tc.expectedPost, tc.expectedError)

			post, err := suite.postUseCase.GetPostByID(context.Background(), tc.id)

			suite.Equal(tc.expectedError, err)
			suite.Equal(tc.expectedPost, post)
			suite.mockPostRepo.AssertExpectations(suite.T())
		})
	}
}

func (suite *PostUseCaseTestSuite) TestUpdatePost() {
	testCases := []struct {
		name          string
		post          *domain.Post
		expectedPost  *domain.Post
		expectedError error
		mockSetup     func(post *domain.Post, expectedPost *domain.Post, expectedError error)
	}{
		{
			name: "Success",
			post: &domain.Post{
				ID:      "1",
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			expectedPost: &domain.Post{
				ID:      "1",
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			expectedError: nil,
			mockSetup: func(post *domain.Post, expectedPost *domain.Post, expectedError error) {
				suite.mockPostRepo.On("UpdatePost", mock.Anything, mock.MatchedBy(func(p *domain.Post) bool {
					return p.ID == post.ID
				})).Return(expectedPost, expectedError).Once()
			},
		},
		{
			name: "Failure",
			post: &domain.Post{
				ID:      "1",
				Title:   "Updated Title",
				Content: "Updated Content",
			},
			expectedPost:  nil,
			expectedError: errors.New("update failed"),
			mockSetup: func(post *domain.Post, expectedPost *domain.Post, expectedError error) {
				suite.mockPostRepo.On("UpdatePost", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(nil, expectedError).Once()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup(tc.post, tc.expectedPost, tc.expectedError)

			updatedPost, err := suite.postUseCase.UpdatePost(context.Background(), tc.post)

			suite.Equal(tc.expectedError, err)
			suite.Equal(tc.expectedPost, updatedPost)
			suite.mockPostRepo.AssertExpectations(suite.T())
		})
	}
}

// In TestDeletePost
func (suite *PostUseCaseTestSuite) TestDeletePost() {
	testCases := []struct {
		name          string
		id            string
		userID        string
		role          string
		expectedError error
		mockSetup     func(id, userID, role string, expectedError error)
	}{
		{
			name:          "Success",
			id:            "1",
			userID:        "user123",
			role:          "admin",
			expectedError: nil,
			mockSetup: func(id, userID, role string, expectedError error) {
				suite.mockPostRepo.On("DeletePost", mock.Anything, id, userID, role).Return(expectedError).Once()
			},
		},
		{
			name:          "Failure",
			id:            "1",
			userID:        "user123",
			role:          "user",
			expectedError: errors.New("deletion failed"),
			mockSetup: func(id, userID, role string, expectedError error) {
				suite.mockPostRepo.On("DeletePost", mock.Anything, id, userID, role).Return(expectedError).Once()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup(tc.id, tc.userID, tc.role, tc.expectedError)

			err := suite.postUseCase.DeletePost(context.Background(), tc.id, tc.userID, tc.role)

			suite.Equal(tc.expectedError, err)
			suite.mockPostRepo.AssertExpectations(suite.T())
		})
	}
}

func TestPostUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(PostUseCaseTestSuite))
}

func stringPtr(s string) *string {
	return &s
}
