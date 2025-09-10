package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock Repository ---
type MockCategoryRepository struct{ mock.Mock }

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetCategories(c context.Context, options *domain.CategorySearchAndFilter) ([]*domain.Category, int64, error) {
	args := m.Called(c, options)
	var catagories []*domain.Category
	if args.Get(0) != nil {
		catagories = args.Get(0).([]*domain.Category)
	}
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return catagories, total, args.Error(2)
}

// --- Test Suite ---
type CategoryUsecaseTestSuite struct {
	suite.Suite
	mockRepo       *MockCategoryRepository
	usecase        domain.ICategoryUsecase
	contextTimeout time.Duration
}

func (s *CategoryUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockCategoryRepository)
	s.contextTimeout = 2 * time.Second
	s.usecase = NewCategoryUsecase(s.mockRepo, s.contextTimeout)
}

func TestCategoryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryUsecaseTestSuite))
}

// --- Tests ---

func (s *CategoryUsecaseTestSuite) TestCreateCategory() {
	ctx := context.Background()
	category := &domain.Category{
		OrganizationID: "org123",
		ParentID:       "parent456",
		Title:          "Test Category",
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Category) bool {
			return c.OrganizationID == category.OrganizationID && c.Title == category.Title
		})).Return(nil).Once()

		err := s.usecase.CreateCategory(ctx, category)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Repo Error", func() {
		s.SetupTest()
		s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(errors.New("db error")).Once()

		err := s.usecase.CreateCategory(ctx, category)
		s.Error(err)
		s.Contains(err.Error(), "db error")
		s.mockRepo.AssertExpectations(s.T())
	})
}

func (s *CategoryUsecaseTestSuite) TestGetCategories() {
	ctx := context.Background()

	s.Run("Success - Valid Options", func() {
		// Arrange
		s.SetupTest() // Reset mocks for each sub-test
		options := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 20,
		}
		expectedCategories := []*domain.Category{{ID: "cat-1", Title: "Category 1"}}
		var expectedTotal int64 = 1
		s.mockRepo.On("GetCategories", mock.Anything, options).Return(expectedCategories, expectedTotal, nil).Once()

		// Act
		categories, total, err := s.usecase.GetCategories(ctx, options)

		// Assert
		s.NoError(err)
		s.Equal(expectedTotal, total)
		s.Equal(expectedCategories, categories)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Applies Default Limit", func() {
		// Arrange
		s.SetupTest()
		// Input has a limit of 0, should be defaulted to 10
		inputOptions := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 0,
		}
		// We expect the repo to be called with the modified options
		expectedOptionsAfterDefault := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 10,
		}
		s.mockRepo.On("GetCategories", mock.Anything, expectedOptionsAfterDefault).Return([]*domain.Category{}, int64(0), nil).Once()

		// Act
		_, _, err := s.usecase.GetCategories(ctx, inputOptions)

		// Assert
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Enforces Max Limit", func() {
		// Arrange
		s.SetupTest()
		// Input has a limit of 200, should be enforced to 100
		inputOptions := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 200,
		}
		expectedOptionsAfterEnforce := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 100,
		}
		s.mockRepo.On("GetCategories", mock.Anything, expectedOptionsAfterEnforce).Return([]*domain.Category{}, int64(0), nil).Once()

		// Act
		_, _, err := s.usecase.GetCategories(ctx, inputOptions)

		// Assert
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Applies Default Page", func() {
		// Arrange
		s.SetupTest()
		// Input has a page of 0, should be defaulted to 1
		inputOptions := &domain.CategorySearchAndFilter{
			Page:  0,
			Limit: 10,
		}
		expectedOptionsAfterDefault := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 10,
		}
		s.mockRepo.On("GetCategories", mock.Anything, expectedOptionsAfterDefault).Return([]*domain.Category{}, int64(0), nil).Once()

		// Act
		_, _, err := s.usecase.GetCategories(ctx, inputOptions)

		// Assert
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Repo Error", func() {
		// Arrange
		s.SetupTest()
		options := &domain.CategorySearchAndFilter{
			Page:  1,
			Limit: 10,
		}
		expectedError := errors.New("database connection failed")
		s.mockRepo.On("GetCategories", mock.Anything, options).Return(nil, int64(0), expectedError).Once()

		// Act
		categories, total, err := s.usecase.GetCategories(ctx, options)

		// Assert
		s.Error(err)
		s.Equal(expectedError, err)
		s.Nil(categories)
		s.Zero(total)
		s.mockRepo.AssertExpectations(s.T())
	})
}
