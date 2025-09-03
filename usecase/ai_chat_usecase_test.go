package usecase_test

import (
	"context"
	// "errors"
	"testing"

	"EthioGuide/domain"
	"EthioGuide/usecase"

	"github.com/stretchr/testify/suite"

    "github.com/stretchr/testify/mock"
)

// --- Mock for IProcedureRepository ---
type MockProcedureRepository struct {
    mock.Mock
}

func (m *MockProcedureRepository) Create(ctx context.Context, procedure *domain.Procedure) error {
    args := m.Called(ctx, procedure)
    return args.Error(0)
}

func (m *MockProcedureRepository) SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*domain.Procedure, error) {
    args := m.Called(ctx, queryVec, limit)
    if args.Get(0) != nil {
        return args.Get(0).([]*domain.Procedure), args.Error(1)
    }
    return nil, args.Error(1)
}

// --- Mock for AIChatRepository ---
type MockAIChatRepository struct {
    mock.Mock
}

func (m *MockAIChatRepository) Save(ctx context.Context, chat *domain.AIChat) error {
    args := m.Called(ctx, chat)
    return args.Error(0)
}

func (m *MockAIChatRepository) GetByUser(ctx context.Context, userID int, limit int) ([]*domain.AIChat, error) {
    args := m.Called(ctx, userID, limit)
    if args.Get(0) != nil {
        return args.Get(0).([]*domain.AIChat), args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *MockAIChatRepository) DeleteByUser(ctx context.Context, userID int) error {
    args := m.Called(ctx, userID)
    return args.Error(0)
}

type MockEmbeddingService struct {
    mock.Mock
}

func (m *MockEmbeddingService) GenerateEmbedding(ctx context.Context, query string) ([]float64, error) {
    args := m.Called(ctx, query)
    if args.Get(0) != nil {
        return args.Get(0).([]float64), args.Error(1)
    }
    return nil, args.Error(1)
}

type MockLLMService struct {
    mock.Mock
}

func (m *MockLLMService) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
    args := m.Called(ctx, prompt)
    return args.String(0), args.Error(1)
}

type AIChatUsecaseSuite struct {
    suite.Suite
    embedMock *MockEmbeddingService
    procMock  *MockProcedureRepository
    chatRepo  *MockAIChatRepository
    llmMock   *MockLLMService
    usecase   domain.IAIChatUsecase
}

func (s *AIChatUsecaseSuite) SetupTest() {
    s.embedMock = &MockEmbeddingService{}
    s.procMock = &MockProcedureRepository{}
    s.chatRepo = &MockAIChatRepository{}
    s.llmMock = &MockLLMService{}
    s.usecase = &usecase.AIChatUsecase{
        EmbedService:  s.embedMock,
        ProcedureRepo: s.procMock,
        AiChatRepo:    s.chatRepo,
        LLMService:    s.llmMock,
    }
}

func TestAIChatUsecaseSuite(t *testing.T) {
    suite.Run(t, new(AIChatUsecaseSuite))
}
func (s *AIChatUsecaseSuite) TestAIchat_OfficialSource() {
    ctx := context.WithValue(context.Background(), "userID", "user123")

    // 1. LLM: language detection returns "english"
    s.llmMock.On("GenerateCompletion", ctx, mock.AnythingOfType("string")).Return("english", nil).Once()
    // 2. LLM: main answer
    s.llmMock.On("GenerateCompletion", ctx, mock.AnythingOfType("string")).Return("AI answer", nil).Once()

    // Embedding
    s.embedMock.On("GenerateEmbedding", ctx, mock.AnythingOfType("string")).Return([]float64{0.1, 0.2, 0.3}, nil)

    // Vector search returns one procedure
    s.procMock.On("SearchByEmbedding", ctx, []float64{0.1, 0.2, 0.3}, 3).Return([]*domain.Procedure{
        {
            Name: "Procedure1",
            Content: domain.Content{
                Prerequisites: []string{"A"},
                Steps:         []string{"Step1"},
                Result:        []string{"Result1"},
            },
            Fees: domain.Fees{
                Currency: "ETB",
                Amount:   100,
                Label:    "Fee",
            },
            ProcessingTime: domain.ProcessingTime{
                MinDays: 1,
                MaxDays: 2,
            },
        },
    }, nil)

    // Chat repo save
    s.chatRepo.On("Save", ctx, mock.MatchedBy(func(chat *domain.AIChat) bool {
        return chat.UserID == "user123" && chat.Source == "official"
    })).Return(nil)

    answer, err := s.usecase.AIchat(ctx, "How do I get a license?")
    s.NoError(err)
    s.Equal("AI answer", answer)
    s.llmMock.AssertExpectations(s.T())
    s.embedMock.AssertExpectations(s.T())
    s.procMock.AssertExpectations(s.T())
    s.chatRepo.AssertExpectations(s.T())
}

func (s *AIChatUsecaseSuite) TestAIchat_UnofficialSource() {
    ctx := context.WithValue(context.Background(), "userID", "user456")

    // 1. LLM: language detection returns "english"
    s.llmMock.On("GenerateCompletion", ctx, mock.AnythingOfType("string")).Return("english", nil).Once()
    // 2. LLM: main answer
    s.llmMock.On("GenerateCompletion", ctx, mock.AnythingOfType("string")).Return("AI answer", nil).Once()

    // Embedding
    s.embedMock.On("GenerateEmbedding", ctx, mock.AnythingOfType("string")).Return([]float64{0.1, 0.2, 0.3}, nil)

    // Vector search returns no procedures
    s.procMock.On("SearchByEmbedding", ctx, []float64{0.1, 0.2, 0.3}, 3).Return([]*domain.Procedure{}, nil)

    // Chat repo save
    s.chatRepo.On("Save", ctx, mock.MatchedBy(func(chat *domain.AIChat) bool {
        return chat.UserID == "user456" && chat.Source == "unofficial"
    })).Return(nil)

    answer, err := s.usecase.AIchat(ctx, "What is the weather?")
    s.NoError(err)
    s.Equal("AI answer", answer)
    s.llmMock.AssertExpectations(s.T())
    s.embedMock.AssertExpectations(s.T())
    s.procMock.AssertExpectations(s.T())
    s.chatRepo.AssertExpectations(s.T())
}