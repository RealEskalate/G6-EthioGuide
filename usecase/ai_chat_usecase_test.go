// package usecase_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"

// 	"EthioGuide/domain"
// 	"EthioGuide/usecase"
// )

// type MockEmbeddingService struct {
// 	GenerateEmbeddingFunc func(ctx context.Context, query string) ([]float64, error)
// }

// func (m *MockEmbeddingService) GenerateEmbedding(ctx context.Context, query string) ([]float64, error) {
// 	return m.GenerateEmbeddingFunc(ctx, query)
// }

// type MockProcedureRepository struct {
// 	SearchByEmbeddingFunc func(ctx context.Context, vec []float64, topK int) ([]domain.Procedure, error)
// }

// func (m *MockProcedureRepository) SearchByEmbedding(ctx context.Context, vec []float64, topK int) ([]domain.Procedure, error) {
// 	return m.SearchByEmbeddingFunc(ctx, vec, topK)
// }

// // type MockAIService struct {
// //     mock.Mock
// // }

// // func (m *MockAIService) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
// //     args := m.Called(ctx, prompt)
// //     if rf, ok := args.Get(0).(func(context.Context, string) (string, error)); ok {
// //         return rf(ctx, prompt)
// //     }
// //     return args.String(0), args.Error(1)
// // }

// type ChatUsecaseTestSuite struct {
// 	suite.Suite
// 	usecase   *usecase.AIChatUsecase
// 	embedMock *MockEmbeddingService
// 	procMock  *MockProcedureRepository
// 	llmMock   *MockAIService
// }

// func (s *ChatUsecaseTestSuite) TestAIchat_Success() {
// 	s.embedMock.GenerateEmbeddingFunc = func(ctx context.Context, query string) ([]float64, error) {
// 		return []float64{1.0, 2.0, 3.0}, nil
// 	}
// 	s.procMock.SearchByEmbeddingFunc = func(ctx context.Context, vec []float64, topK int) ([]domain.Procedure, error) {
// 		return []domain.Procedure{{
// 			Name: "Test",
// 			Content: domain.Content{
// 				Prerequisites: []string{"Req1", "Req2"},
// 				Steps:         []string{"Step1", "Step2"},
// 				Result:        []string{"Result1"},
// 			},
// 			Fees: domain.Fees{
// 				Label:    "Service Fee",
// 				Currency: "ETB",
// 				Amount:   100.0,
// 			},
// 			ProcessingTime: domain.ProcessingTime{
// 				MinDays: 1,
// 				MaxDays: 3,
// 			},
// 		}}, nil
// 	}
// 	s.llmMock.On("GenerateCompletion", mock.Anything, mock.Anything).Return("AI answer", nil)

// 	answer, err := s.usecase.AIchat(context.Background(), "How to do X?")
// 	assert.NoError(s.T(), err)
// 	assert.Equal(s.T(), "AI answer", answer)
// 	s.llmMock.AssertExpectations(s.T())
// }

// func (s *ChatUsecaseTestSuite) TestAIchat_EmbeddingError() {
// 	s.embedMock.GenerateEmbeddingFunc = func(ctx context.Context, query string) ([]float64, error) {
// 		return nil, errors.New("embedding error")
// 	}

// 	answer, err := s.usecase.AIchat(context.Background(), "fail embedding")
// 	assert.Error(s.T(), err)
// 	assert.Empty(s.T(), answer)
// }

// func (s *ChatUsecaseTestSuite) TestAIchat_SearchError() {
// 	s.embedMock.GenerateEmbeddingFunc = func(ctx context.Context, query string) ([]float64, error) {
// 		return []float64{1.0}, nil
// 	}
// 	s.procMock.SearchByEmbeddingFunc = func(ctx context.Context, vec []float64, topK int) ([]domain.Procedure, error) {
// 		return nil, errors.New("search error")
// 	}

// 	answer, err := s.usecase.AIchat(context.Background(), "fail search")
// 	assert.Error(s.T(), err)
// 	assert.Empty(s.T(), answer)
// }

// func (s *ChatUsecaseTestSuite) TestAIchat_LLMError() {
// 	s.embedMock.GenerateEmbeddingFunc = func(ctx context.Context, query string) ([]float64, error) {
// 		return []float64{1.0}, nil
// 	}
// 	s.procMock.SearchByEmbeddingFunc = func(ctx context.Context, vec []float64, topK int) ([]domain.Procedure, error) {
// 		return []domain.Procedure{{
// 			Name: "Test",
// 			Content: domain.Content{
// 				Prerequisites: []string{"Req1", "Req2"},
// 				Steps:         []string{"Step1", "Step2"},
// 				Result:        []string{"Result1"},
// 			},
// 			Fees: domain.Fees{
// 				Label:    "Service Fee",
// 				Currency: "ETB",
// 				Amount:   100.0,
// 			},
// 			ProcessingTime: domain.ProcessingTime{
// 				MinDays: 1,
// 				MaxDays: 3,
// 			},
// 		}}, nil
// 	}
// 	s.llmMock.On("GenerateCompletion", mock.Anything, mock.Anything).Return("", errors.New("llm error"))

// 	answer, err := s.usecase.AIchat(context.Background(), "fail llm")
// 	assert.Error(s.T(), err)
// 	assert.Empty(s.T(), answer)
// 	s.llmMock.AssertExpectations(s.T())
// }

// func TestChatUsecaseTestSuite(t *testing.T) {
// 	suite.Run(t, new(ChatUsecaseTestSuite))
// }
