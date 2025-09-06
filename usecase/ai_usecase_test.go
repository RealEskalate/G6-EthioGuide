package usecase_test

// import (
// 	"EthioGuide/domain"
// 	. "EthioGuide/usecase"
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// // --- Mock Setup ---
// // (MockAIService remains the same)

// // MockAIService is a mock implementation of the IAIService interface.
// type MockAIService struct {
// 	mock.Mock
// }

// // GenerateCompletion is the mock method.
// func (m *MockAIService) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
// 	// m.Called() records the call and its arguments.
// 	args := m.Called(ctx, prompt)

// 	// Check if a 'run' function was provided for this call.
// 	// This is the key to solving the panic.
// 	if rf, ok := args.Get(0).(func(context.Context, string) (string, error)); ok {
// 		// If it was, execute it and return its results.
// 		return rf(ctx, prompt)
// 	}

// 	// Otherwise, fall back to the original behavior of returning static values.
// 	return args.String(0), args.Error(1)
// }

// // --- Test Suite Definition ---

// type GeminiUseCaseTestSuite struct {
// 	suite.Suite
// 	mockAIService *MockAIService
// 	useCase       domain.IGeminiUseCase
// 	ctx           context.Context
// }

// // SetupTest runs before each test in the suite.
// func (s *GeminiUseCaseTestSuite) SetupTest() {
// 	s.mockAIService = new(MockAIService)
// 	s.useCase = NewGeminiUsecase(s.mockAIService, 5*time.Second)
// 	s.ctx = context.Background()
// }

// // TestGeminiUseCaseTestSuite is the entry point for running the test suite.
// func TestGeminiUseCaseTestSuite(t *testing.T) {
// 	suite.Run(t, new(GeminiUseCaseTestSuite))
// }

// // --- Test Cases ---

// func (s *GeminiUseCaseTestSuite) TestTranslateContent_Success() {
// 	// Arrange: Define inputs and expected outputs
// 	content := "  Hello  " // Added whitespace to test trimming
// 	targetLang := "am"
// 	expectedTranslation := "ሰላም"

// 	// Configure the mock. Use mock.Anything for the context because the use case
// 	// creates a derived context with a timeout, which won't be the same object.
// 	s.mockAIService.On("GenerateCompletion", mock.Anything, mock.Anything).
// 		Return("  "+expectedTranslation+"  ", nil).Once() // Return with whitespace to test trimming

// 	// Act: Call the method under test
// 	result, err := s.useCase.TranslateContent(s.ctx, content, targetLang)

// 	// Assert: Check the results
// 	s.NoError(err)
// 	s.Equal(expectedTranslation, result)

// 	// Verify that the mock was called as expected
// 	s.mockAIService.AssertExpectations(s.T())
// }

// func (s *GeminiUseCaseTestSuite) TestTranslateContent_EmptyContent() {
// 	content := "   "
// 	targetLang := "en"

// 	result, err := s.useCase.TranslateContent(s.ctx, content, targetLang)

// 	s.NoError(err)
// 	s.Empty(result)
// 	s.mockAIService.AssertNotCalled(s.T(), "GenerateCompletion", mock.Anything, mock.Anything)
// }

// func (s *GeminiUseCaseTestSuite) TestTranslateContent_UnsupportedLanguage() {
// 	content := "This will fail"
// 	targetLang := "fr"

// 	result, err := s.useCase.TranslateContent(s.ctx, content, targetLang)

// 	s.Empty(result)
// 	s.Error(err)
// 	s.ErrorIs(err, domain.ErrUnsupportedLanguage)
// 	s.Contains(err.Error(), targetLang)
// 	s.mockAIService.AssertNotCalled(s.T(), "GenerateCompletion", mock.Anything, mock.Anything)
// }

// func (s *GeminiUseCaseTestSuite) TestTranslateContent_AIServiceReturnsError() {
// 	// Arrange
// 	content := "API call will fail"
// 	targetLang := "en"
// 	serviceError := errors.New("API rate limit exceeded")
// 	expectedWrappedErrorMsg := "gemini service failed to generate completion"

// 	// Configure the mock to return an error.
// 	s.mockAIService.On("GenerateCompletion", mock.Anything, mock.Anything).Return("", serviceError)

// 	// Act
// 	result, err := s.useCase.TranslateContent(s.ctx, content, targetLang)

// 	// Assert
// 	s.Empty(result)
// 	s.Error(err)
// 	s.ErrorIs(err, serviceError)                     // Check that the original error is still present
// 	s.Contains(err.Error(), expectedWrappedErrorMsg) // Check for our custom wrapping message

// 	// Verify that the mock was called
// 	s.mockAIService.AssertExpectations(s.T())
// }

// func (s *GeminiUseCaseTestSuite) TestTranslateContent_AIReturnsUnknownLanguageString() {
// 	content := "Some content"
// 	targetLang := "en"

// 	s.mockAIService.On("GenerateCompletion", mock.Anything, mock.Anything).Return("unknown language", nil)

// 	result, err := s.useCase.TranslateContent(s.ctx, content, targetLang)

// 	s.Empty(result)
// 	s.Error(err)
// 	s.ErrorIs(err, domain.ErrUnsupportedLanguage)
// 	s.mockAIService.AssertExpectations(s.T())
// }

// func (s *GeminiUseCaseTestSuite) TestTranslateContent_ContextTimeout() {
// 	// Arrange
// 	content := "This will time out"
// 	targetLang := "en"
// 	shortTimeout := 10 * time.Millisecond
// 	longOperation := 20 * time.Millisecond

// 	shortTimeoutUsecase := NewGeminiUsecase(s.mockAIService, shortTimeout)

// 	s.mockAIService.On("GenerateCompletion", mock.Anything, mock.Anything).
// 		Return(func(ctx context.Context, prompt string) (string, error) {
// 			select {
// 			case <-time.After(longOperation):
// 				return "this should not be returned", nil
// 			case <-ctx.Done():
// 				return "", ctx.Err()
// 			}
// 		}).Once()

// 	// Act
// 	result, err := shortTimeoutUsecase.TranslateContent(s.ctx, content, targetLang)

// 	// Assert
// 	s.Empty(result, "The result should be empty on timeout")
// 	s.Error(err, "An error should be returned on timeout")
// 	s.ErrorIs(err, context.DeadlineExceeded, "The wrapped error should be context.DeadlineExceeded")
// 	s.mockAIService.AssertExpectations(s.T())
// }
