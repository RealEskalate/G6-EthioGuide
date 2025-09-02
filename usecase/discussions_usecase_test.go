// filepath: usecase/discussions_usecase_test.go
package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockPostRepository is a mock implementation of domain.IPostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) GetPostByID(ctx context.Context, id int) (*domain.Post, error) {
	args := m.Called(ctx, id)
	post, ok := args.Get(0).(*domain.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return post, args.Error(1)
}

func (m *MockPostRepository) GetPosts(ctx context.Context) ([]*domain.Post, error) {
	args := m.Called(ctx)
	posts, ok := args.Get(0).([]*domain.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return posts, args.Error(1)
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) DeletePost(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
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
	suite.postUseCase = NewPostUseCase(suite.mockPostRepo)
}

func (suite *PostUseCaseTestSuite) TestCreatePost() {
	testCases := []struct {
		name          string
		post          *domain.Post
		expectedError error
		mockSetup     func(post *domain.Post, expectedError error)
	}{
		{
			name: "Success",
			post: &domain.Post{
				UserID:    "user123",
				Title:     "Test Post",
				Content:   "This is a test post.",
				Procedures: []string{"Step 1", "Step 2"},
				Tags:      []string{"test", "post"},
			},
			expectedError: nil,
			mockSetup: func(post *domain.Post, expectedError error) {
				suite.mockPostRepo.On("CreatePost", mock.Anything, post).Return(expectedError).Once()
			},
		},
		{
			name: "Failure",
			post: &domain.Post{
				UserID:    "user123",
				Title:     "Test Post",
				Content:   "This is a test post.",
				Procedures: []string{"Step 1", "Step 2"},
				Tags:      []string{"test", "post"},
			},
			expectedError: errors.New("creation failed"),
			mockSetup: func(post *domain.Post, expectedError error) {
				suite.mockPostRepo.On("CreatePost", mock.Anything, post).Return(expectedError).Once()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.mockSetup(tc.post, tc.expectedError)

			err := suite.postUseCase.CreatePost(context.Background(), tc.post)

			suite.Equal(tc.expectedError, err)
			suite.mockPostRepo.AssertExpectations(suite.T())
		})
	}
}

func TestPostUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(PostUseCaseTestSuite))
}