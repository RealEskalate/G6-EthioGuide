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

// --- Mocks ---
type MockEmbeddingService struct {
	mock.Mock
}

func (m *MockEmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float64), args.Error(1)
}

// MockProcedureRepository is a complete mock for the IProcedureRepository interface
type MockProcedureRepository struct {
	mock.Mock
}

func (m *MockProcedureRepository) Create(ctx context.Context, procedure *domain.Procedure) error {
	args := m.Called(ctx, procedure)
	return args.Error(0)
}

func (m *MockProcedureRepository) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Procedure), args.Error(1)
}

func (m *MockProcedureRepository) Update(ctx context.Context, id string, procedure *domain.Procedure) error {
	args := m.Called(ctx, id, procedure)
	return args.Error(0)
}

func (m *MockProcedureRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProcedureRepository) SearchAndFilter(ctx context.Context, options domain.ProcedureSearchFilterOptions) ([]*domain.Procedure, int64, error) {
	args := m.Called(ctx, options)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Procedure), args.Get(1).(int64), args.Error(2)
}

func (m *MockProcedureRepository) SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*domain.Procedure, error) {
	args := m.Called(ctx, queryVec, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Procedure), args.Error(1)
}

type MockAIChatRepository struct {
	mock.Mock
}

func (m *MockAIChatRepository) Save(ctx context.Context, chat *domain.AIChat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}
func (m *MockAIChatRepository) GetByUser(ctx context.Context, userID string, page, limit int64) ([]*domain.AIChat, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.AIChat), args.Get(1).(int64), args.Error(2)
}
func (m *MockAIChatRepository) DeleteByUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockLLMService struct {
	mock.Mock
}

func (m *MockLLMService) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

// --- Test Suite ---
type AIChatUsecaseTestSuite struct {
	suite.Suite
	mockEmbedSvc *MockEmbeddingService
	mockProcRepo *MockProcedureRepository
	mockChatRepo *MockAIChatRepository
	mockLLMSvc   *MockLLMService
	usecase      domain.IAIChatUsecase
}

func (s *AIChatUsecaseTestSuite) SetupTest() {
	s.mockEmbedSvc = new(MockEmbeddingService)
	s.mockProcRepo = new(MockProcedureRepository)
	s.mockChatRepo = new(MockAIChatRepository)
	s.mockLLMSvc = new(MockLLMService)
	s.usecase = NewChatUsecase(s.mockEmbedSvc, s.mockProcRepo, s.mockChatRepo, s.mockLLMSvc, 5*time.Second)
}

func (s *AIChatUsecaseTestSuite) TestAIchat_Success_English() {
	// Arrange
	userID, query := "user-1", "How to get passport?"
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("procedure", nil).Once() // Classifier
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("english", nil).Once()   // Language detector
	s.mockEmbedSvc.On("GenerateEmbedding", mock.Anything, query).Return([]float64{0.1}, nil).Once()

	// Use mock.Anything for the vector, since it's hard to predict
	s.mockProcRepo.On("SearchByEmbedding", mock.Anything, mock.AnythingOfType("[]float64"), 3).Return([]*domain.Procedure{{ID: "proc-1", Name: "Passport Application"}}, nil).Once()

	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("Final Answer", nil).Once() // Final answer
	s.mockChatRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.AIChat")).Return(nil).Once()

	// Act
	chat, err := s.usecase.AIchat(context.Background(), userID, query)

	// Assert
	s.NoError(err)
	s.NotNil(chat)
	s.Equal("Final Answer", chat.Response)
	s.Equal("official", chat.Source)
	s.Len(chat.RelatedProcedures, 1)
	s.mockLLMSvc.AssertExpectations(s.T())
	s.mockEmbedSvc.AssertExpectations(s.T())
	s.mockChatRepo.AssertExpectations(s.T())
	s.mockProcRepo.AssertExpectations(s.T())
}

func (s *AIChatUsecaseTestSuite) TestAIchat_Success_Amharic_Translation() {
	// Arrange
	userID, amharicQuery := "user-1", "ፓስፖርት እንዴት ማውጣት እችላለሁ?"
	englishQuery := "How can I get a passport?"
	amharicAnswer := "መልሱ ይኸውና"

	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("procedure", nil).Once()  // Classifier
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("amharic", nil).Once()    // Language detector
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return(englishQuery, nil).Once() // Query translation
	s.mockEmbedSvc.On("GenerateEmbedding", mock.Anything, englishQuery).Return([]float64{0.1}, nil).Once()
	s.mockProcRepo.On("SearchByEmbedding", mock.Anything, mock.Anything, 3).Return([]*domain.Procedure{}, nil).Once()
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("English Answer", nil).Once() // Final answer
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return(amharicAnswer, nil).Once()    // Answer translation
	s.mockChatRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.AIChat")).Return(nil).Once()

	// Act
	chat, err := s.usecase.AIchat(context.Background(), userID, amharicQuery)

	// Assert
	s.NoError(err)
	s.NotNil(chat)
	s.Equal(amharicAnswer, chat.Response)
	s.Equal("unofficial", chat.Source) // No docs found
	s.mockLLMSvc.AssertExpectations(s.T())
}

func (s *AIChatUsecaseTestSuite) TestAIchat_OffensiveQuery() {
	// Arrange
	userID, query := "user-1", "some offensive words"
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("offensive", nil).Once()

	// Act
	chat, err := s.usecase.AIchat(context.Background(), userID, query)

	// Assert
	s.Error(err)
	s.NotNil(chat)
	s.Contains(err.Error(), "offensive content")
	s.Contains(chat.Response, "violate our content policy")
	// Ensure no other calls were made
	s.mockEmbedSvc.AssertNotCalled(s.T(), "GenerateEmbedding")
}

func (s *AIChatUsecaseTestSuite) TestAIchat_IrrelevantQuery() {
	// Arrange
	userID, query := "user-1", "what is 2+2?"
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("irrelevant", nil).Once()

	// Act
	chat, err := s.usecase.AIchat(context.Background(), userID, query)

	// Assert
	s.NoError(err)
	s.NotNil(chat)
	s.Contains(chat.Response, "I can only answer questions about Ethiopian government procedures")
	// Ensure no other calls were made
	s.mockEmbedSvc.AssertNotCalled(s.T(), "GenerateEmbedding")
}

func (s *AIChatUsecaseTestSuite) TestAIchat_EmbeddingError() {
	// Arrange
	userID, query := "user-1", "How to get passport?"
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("procedure", nil).Once()
	s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.AnythingOfType("string")).Return("english", nil).Once()
	s.mockEmbedSvc.On("GenerateEmbedding", mock.Anything, query).Return(nil, errors.New("embedding failed")).Once()

	// Act
	chat, err := s.usecase.AIchat(context.Background(), userID, query)

	// Assert
	s.Error(err)
	s.NotNil(chat)
	s.Contains(chat.Response, "I'm sorry, but I encountered a technical issue")
	s.EqualError(err, "embedding failed")
}

func (s *AIChatUsecaseTestSuite) TestAIHistory_Success() {
	// Arrange
	userID := "user-1"
	var page, limit int64 = 1, 10
	expectedHistory := []*domain.AIChat{{ID: "chat-1"}}
	var total int64 = 1
	s.mockChatRepo.On("GetByUser", mock.Anything, userID, page, limit).Return(expectedHistory, total, nil).Once()

	// Act
	history, resTotal, err := s.usecase.AIHistory(context.Background(), userID, page, limit)

	// Assert
	s.NoError(err)
	s.Equal(total, resTotal)
	s.Equal(expectedHistory, history)
	s.mockChatRepo.AssertExpectations(s.T())
}

func TestAIChatUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AIChatUsecaseTestSuite))
}
