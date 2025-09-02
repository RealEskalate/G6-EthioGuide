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
