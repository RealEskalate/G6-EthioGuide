package usecase_test

import (
	"EthioGuide/domain"
	"EthioGuide/usecase"
	"context"
	"errors"
	"testing"
	"time"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks ---

type MockSearchRepository struct {
	mock.Mock
}

// Ensure MockSearchRepository implements ISearchRepository at compile time.
var _ domain.ISearchRepository = (*MockSearchRepository)(nil)

func (m *MockSearchRepository) Search(ctx context.Context, filter domain.SearchFilterRequest) (*domain.SearchResult, error) {
	args := m.Called(ctx, filter)
	var res *domain.SearchResult
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.SearchResult)
	}
	return res, args.Error(1)
}

// The other methods of the interface are not used by this use case, so we can give them empty implementations.
func (m *MockSearchRepository) FindProcedures(ctx context.Context, filter domain.SearchFilterRequest) ([]*domain.Procedure, error) {
	return nil, nil
}

func (m *MockSearchRepository) FindOrganizations(ctx context.Context, filter domain.SearchFilterRequest) ([]*domain.AccountOrgSearch, error) {
	return nil, nil
}

// --- Test Suite ---

type SearchUsecaseTestSuite struct {
	suite.Suite
	mockRepo    *MockSearchRepository
	searchUsecase domain.ISearchUseCase
}

func (s *SearchUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockSearchRepository)
	// Use a 5-second timeout for tests
	s.searchUsecase = usecase.NewSearchUsecase(s.mockRepo, 5*time.Second)
}

func TestSearchUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(SearchUsecaseTestSuite))
}

// --- Test Cases ---

func (s *SearchUsecaseTestSuite) TestSearch_Success() {
	// Arrange
	filter := domain.SearchFilterRequest{
		Query: "test",
		Page:  1,
		Limit: 10,
	}

	expectedResult := &domain.SearchResult{
		Procedures:    []*domain.Procedure{{Name: "Test Procedure"}},
		Organizations: []*domain.AccountOrgSearch{{Name: "Test Org"}},
	}

	s.mockRepo.On("Search", mock.Anything, filter).Return(expectedResult, nil).Once()

	// Act
	result, err := s.searchUsecase.Search(context.Background(), filter)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(expectedResult, result)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *SearchUsecaseTestSuite) TestSearch_RepositoryError() {
	// Arrange
	filter := domain.SearchFilterRequest{
		Query: "test",
		Page:  1,
		Limit: 10,
	}

	expectedError := errors.New("database error")

	s.mockRepo.On("Search", mock.Anything, filter).Return(nil, expectedError).Once()

	// Act
	result, err := s.searchUsecase.Search(context.Background(), filter)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedError, err)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *SearchUsecaseTestSuite) TestSearch_PaginationDefaults() {
	// Arrange
	// Input filter with no pagination details
	inputFilter := domain.SearchFilterRequest{
		Query: "test",
		Page:  0,  // Invalid page
		Limit: -5, // Invalid limit
	}

	// Expected filter with default pagination applied by the use case
	expectedFilter := domain.SearchFilterRequest{
		Query: "test",
		Page:  1,  // Default page
		Limit: 20, // Default limit
	}

	expectedResult := &domain.SearchResult{}

	// We assert that the mock is called with the *expected* filter, not the input one.
	s.mockRepo.On("Search", mock.Anything, expectedFilter).Return(expectedResult, nil).Once()

	// Act
	_, err := s.searchUsecase.Search(context.Background(), inputFilter)

	// Assert
	s.NoError(err)
	s.mockRepo.AssertExpectations(s.T())
}
