package controller_test

import (
	. "EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks & Placeholders ---

// MockGeminiUseCase is a mock implementation of the IGeminiUseCase interface.
type MockGeminiUseCase struct {
	mock.Mock
}

// Ensure MockGeminiUseCase implements the interface.
var _ domain.IGeminiUseCase = (*MockGeminiUseCase)(nil)

// TranslateContent is the mocked method.
func (m *MockGeminiUseCase) TranslateContent(ctx context.Context, content string, targetLang string) (string, error) {
	args := m.Called(ctx, content, targetLang)
	return args.String(0), args.Error(1)
}

// --- Test Suite Definition ---

type GeminiControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockGeminiUseCase
	controller  *GeminiController
	recorder    *httptest.ResponseRecorder
}

// SetupSuite runs once before the entire suite.
func (s *GeminiControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest runs before each individual test.
func (s *GeminiControllerTestSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUsecase = new(MockGeminiUseCase)
	s.controller = NewGeminiController(s.mockUsecase)

	s.router.POST("/translate", s.controller.Translate)
}

// TestGeminiControllerTestSuite is the entry point for the suite.
func TestGeminiControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GeminiControllerTestSuite))
}

// --- Test Cases ---

func (s *GeminiControllerTestSuite) TestTranslate_Success() {
	// Arrange
	requestBody := TranslateDTO{Content: "Hello"}
	jsonBody, _ := json.Marshal(requestBody)
	targetLang := "am"
	expectedTranslation := "ሰላም"

	// Configure the mock to expect a call with these specific arguments
	s.mockUsecase.On("TranslateContent", mock.Anything, requestBody.Content, targetLang).
		Return(expectedTranslation, nil).Once()

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/translate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("lang", targetLang) // Set the language header
	s.router.ServeHTTP(s.recorder, req)

	// Assert
	s.Equal(http.StatusOK, s.recorder.Code)

	var response map[string]string
	json.Unmarshal(s.recorder.Body.Bytes(), &response)
	s.Equal(expectedTranslation, response["content"])

	s.mockUsecase.AssertExpectations(s.T())
}

func (s *GeminiControllerTestSuite) TestTranslate_Success_DefaultLanguage() {
	// Arrange
	requestBody := TranslateDTO{Content: "Hello"}
	jsonBody, _ := json.Marshal(requestBody)
	defaultLang := "en" // This is the fallback language
	expectedTranslation := "Hello"

	// Configure the mock to expect a call with the default language 'en'
	s.mockUsecase.On("TranslateContent", mock.Anything, requestBody.Content, defaultLang).
		Return(expectedTranslation, nil).Once()

	// Act (No 'lang' header is set)
	req, _ := http.NewRequest(http.MethodPost, "/translate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(s.recorder, req)

	// Assert
	s.Equal(http.StatusOK, s.recorder.Code)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *GeminiControllerTestSuite) TestTranslate_InvalidRequestBody() {
	// Arrange: Malformed JSON
	invalidJsonBody := []byte(`{"content": "missing quote}`)

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/translate", bytes.NewBuffer(invalidJsonBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(s.recorder, req)

	// Assert
	s.Equal(http.StatusBadRequest, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), "Invalid request body")

	// Crucially, assert that the use case was never called
	s.mockUsecase.AssertNotCalled(s.T(), "TranslateContent", mock.Anything, mock.Anything, mock.Anything)
}

func (s *GeminiControllerTestSuite) TestTranslate_UsecaseReturnsError() {
	// Arrange
	requestBody := TranslateDTO{Content: "This will fail"}
	jsonBody, _ := json.Marshal(requestBody)
	targetLang := "klingon"

	// Configure the mock to return a domain-specific error
	s.mockUsecase.On("TranslateContent", mock.Anything, requestBody.Content, targetLang).
		Return("", domain.ErrUnsupportedLanguage).Once()

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/translate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("lang", targetLang)
	s.router.ServeHTTP(s.recorder, req)

	// Assert
	// We check for the status code that our dummy HandleError provides for this error.
	s.Equal(http.StatusBadRequest, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), domain.ErrUnsupportedLanguage.Error())

	s.mockUsecase.AssertExpectations(s.T())
}
