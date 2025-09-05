package infrastructure_test

import (
	"EthioGuide/config"
	"EthioGuide/domain"
	. "EthioGuide/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GeminiAIServiceSuite struct {
	suite.Suite
	apiKey    string
	modelName string
	service   domain.IAIService
}

func (s *GeminiAIServiceSuite) SetupSuite() {
	cfg := config.LoadForTest()
	key := cfg.GeminiAPIKey
	if key == "" {
		s.T().Skip("Skipping GeminiAIService integration tests: OPENAI_API_KEY not set.")
	}

	s.apiKey = key
	s.modelName = cfg.GeminiModel
	if s.modelName == "" {
		s.modelName = "gemini-2.5-pro"
	}
}

func (s *GeminiAIServiceSuite) SetupTest() {
	service, err := NewGeminiAIService(s.apiKey, s.modelName)
	s.Require().NoError(err, "Failed to create a new Gemini AI Service instance")
	s.service = service
}

// TestGeminiAIServiceSuite is the entry point for the test suite.
func TestGeminiAIServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(GeminiAIServiceSuite))
}

// --- Tests ---

func (s *GeminiAIServiceSuite) TestGenerateCompletion() {
	ctx := context.Background()

	s.Run("Success - Simple Factual Prompt", func() {
		// Arrange: A simple prompt with a predictable answer.
		prompt := "What is the capital of France?"

		// Act: Call the real Gemini API.
		response, err := s.service.GenerateCompletion(ctx, prompt)

		// Assert
		s.NoError(err, "The API call should not return an error")
		s.NotEmpty(response, "The response from the API should not be empty")
		s.Contains(strings.ToLower(response), "paris", "The response should contain the correct answer")
	})

	s.Run("Success - Structured JSON Prompt", func() {
		// Arrange
		prompt := `Return a single JSON object with one key "country" and the value "Japan". Do not include any other text or markdown formatting.`

		// Act
		response, err := s.service.GenerateCompletion(ctx, prompt)

		// Assert
		s.NoError(err)
		s.NotEmpty(response)

		// The most important assertion: can we unmarshal the response?
		var data struct {
			Country string `json:"country"`
		}
		// Clean the response just in case the AI adds markdown fences.
		cleanedResponse := strings.Trim(response, " \n\t`json")
		unmarshalErr := json.Unmarshal([]byte(cleanedResponse), &data)

		s.NoError(unmarshalErr, "The AI response should be valid JSON that can be unmarshalled")
		s.Equal("Japan", data.Country, "The unmarshalled data should contain the correct value")
	})

	s.Run("Success - Translation Prompt Logic", func() {
		s.T().Helper() // Marks this function as a test helper

		// This helper function creates the exact prompt from your use case
		createTranslationPrompt := func(targetLang, content string) string {
			return fmt.Sprintf(`
	Translate the following content to %s.
	Ensure the translation is accurate, maintains the original meaning, and uses natural, fluent language appropriate for the context. Do not add or omit any information. 
	Do not include any other text or markdown formatting.
	If you don't know the language I requested, send only the text "unknown language".
	Here is the content:
	%s
	`, targetLang, content)
		}

		s.Run("Valid Language", func() {
			// Arrange
			content := "Hello, how are you?"
			prompt := createTranslationPrompt("Spanish", content)

			// Act
			response, err := s.service.GenerateCompletion(ctx, prompt)

			// Assert
			s.NoError(err, "API call for a valid translation should succeed")
			s.NotEmpty(response)
			// Use Contains because LLM output can have slight variations.
			// This is more robust than asserting for exact equality.
			s.Contains(strings.ToLower(response), "hola", "The Spanish translation should contain 'hola'")
		})

		s.Run("Unknown Language", func() {
			// Arrange
			content := "This should not be translated."
			// Use a fictional language to test the specific instruction in the prompt.
			prompt := createTranslationPrompt("abracadabra", content)

			// Act
			response, err := s.service.GenerateCompletion(ctx, prompt)

			// Assert
			s.NoError(err, "API call should succeed even if the language is unknown")
			s.NotEmpty(response)
			// Here we expect an exact string match, as requested in the prompt.
			s.Equal("unknown language", strings.ToLower(strings.TrimSpace(response)), "The AI should return the exact fallback string")
		})
	})
}

// TestNewGeminiAIService_Failures tests the constructor's error handling.
func (s *GeminiAIServiceSuite) TestNewGeminiAIService_Failures() {
	s.Run("Failure - Empty API Key", func() {
		// Act
		service, err := NewGeminiAIService("", "")

		// Assert
		s.Error(err, "Constructor should return an error for an empty API key")
		s.Nil(service, "Service instance should be nil on failure")
	})
}
