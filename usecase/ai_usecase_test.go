package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockGeminiLLMService struct {
	mock.Mock
}

func (m *MockGeminiLLMService) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)

	ret0 := args.Get(0)
	var r0 string
	if rf, ok := ret0.(func(context.Context, string) string); ok {
		// If it's a function, execute it to get the string result
		r0 = rf(ctx, prompt)
	} else {
		// Otherwise, treat it as a regular string
		r0 = ret0.(string)
	}

	return r0, args.Error(1)
}

// GeminiUsecaseTestSuite is the test suite for geminiUseCase
type GeminiUsecaseTestSuite struct {
	suite.Suite
	mockLLMSvc *MockGeminiLLMService
	usecase    domain.IGeminiUseCase
}

func (s *GeminiUsecaseTestSuite) SetupTest() {
	s.mockLLMSvc = new(MockGeminiLLMService)
	s.usecase = NewGeminiUsecase(s.mockLLMSvc, 5*time.Second)
}

func TestGeminiUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(GeminiUsecaseTestSuite))
}

func (s *GeminiUsecaseTestSuite) TestTranslateJSON() {
	// Arrange
	inputData := map[string]interface{}{
		"title":       "Welcome",
		"description": "This is a test.",
		"id":          "user-123", // This key should be ignored
		"details": []interface{}{
			"First item",
			map[string]interface{}{
				"step": "Second item",
				"role": "admin", // This key should be ignored
			},
		},
	}
	targetLang := "am"
	separator := "<!--EthioGuideTranslationSeparator-->"

	s.Run("Success", func() {
		s.SetupTest()

		// Arrange
		// Because we now sort the originals, the order sent to the LLM is predictable.
		// Order will be: "First item", "Second item", "This is a test.", "Welcome"
		sortedOriginals := []string{"First item", "Second item", "This is a test.", "Welcome"}
		contentToTranslate := strings.Join(sortedOriginals, separator)

		// Our mock response must now also be in the same sorted order.
		sortedTranslations := []string{"የመጀመሪያው ንጥል", "ሁለተኛ ንጥል", "ይህ ፈተና ነው።", "እንኳን ደህና መጣህ"}
		translatedBlock := strings.Join(sortedTranslations, separator)

		s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.MatchedBy(func(prompt string) bool {
			return strings.Contains(prompt, contentToTranslate)
		})).Return(translatedBlock, nil).Once()

		// Act
		result, err := s.usecase.TranslateJSON(context.Background(), inputData, targetLang)

		// Assert
		s.NoError(err)
		s.NotNil(result)

		// Assert the final, re-assembled object
		s.Equal("እንኳን ደህና መጣህ", result["title"])
		s.Equal("ይህ ፈተና ነው።", result["description"])
		s.Equal("user-123", result["id"]) // Untranslated
		detailsSlice := result["details"].([]interface{})
		s.Equal("የመጀመሪያው ንጥል", detailsSlice[0])
		stepMap := detailsSlice[1].(map[string]interface{})
		s.Equal("ሁለተኛ ንጥል", stepMap["step"])
		s.Equal("admin", stepMap["role"]) // Untranslated

		s.mockLLMSvc.AssertExpectations(s.T())
	})

	s.Run("Unsupported Language", func() {
		// --- ISOLATION ---
		s.SetupTest()
		_, err := s.usecase.TranslateJSON(context.Background(), inputData, "fr")
		s.ErrorIs(err, domain.ErrUnsupportedLanguage)
		s.mockLLMSvc.AssertNotCalled(s.T(), "GenerateCompletion")
	})

	s.Run("No Translatable Strings", func() {
		// --- ISOLATION ---
		s.SetupTest()
		nonTranslatableData := map[string]interface{}{"id": "123", "role": "admin"}
		result, err := s.usecase.TranslateJSON(context.Background(), nonTranslatableData, "am")
		s.NoError(err)
		s.Equal(nonTranslatableData, result)
		s.mockLLMSvc.AssertNotCalled(s.T(), "GenerateCompletion")
	})

	s.Run("Translation Mismatch Error", func() {
		// --- ISOLATION ---
		s.SetupTest()
		s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.Anything).Return("Only one part", nil).Once()
		_, err := s.usecase.TranslateJSON(context.Background(), inputData, "am")
		s.ErrorIs(err, domain.ErrTranslationMismatch)
		s.mockLLMSvc.AssertExpectations(s.T())
	})

	s.Run("LLM Service Error", func() {
		// --- ISOLATION ---
		s.SetupTest()
		expectedError := errors.New("API limit reached")
		s.mockLLMSvc.On("GenerateCompletion", mock.Anything, mock.Anything).Return("", expectedError).Once()
		_, err := s.usecase.TranslateJSON(context.Background(), inputData, "am")
		s.Error(err)
		s.Contains(err.Error(), expectedError.Error())
		s.mockLLMSvc.AssertExpectations(s.T())
	})
}
