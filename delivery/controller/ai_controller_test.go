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

// MockGeminiUseCase is a mock for the IGeminiUseCase interface
type MockGeminiUseCase struct {
	mock.Mock
}

func (m *MockGeminiUseCase) TranslateJSON(ctx context.Context, data map[string]interface{}, targetLang string) (map[string]interface{}, error) {
	args := m.Called(ctx, data, targetLang)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GeminiControllerTestSuite is the test suite for GeminiController
type GeminiControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockGeminiUseCase
	controller  *GeminiController
}

// SetupTest is run before each test in the suite
func (s *GeminiControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.mockUsecase = new(MockGeminiUseCase)
	s.controller = NewGeminiController(s.mockUsecase)

	// Setup routes
	s.router.POST("/ai/translate", s.controller.Translate)
}

func (s *GeminiControllerTestSuite) TestTranslate_Success_WithLangHeader() {
	// Arrange
	requestBody := TranslateDTO{
		Content: map[string]interface{}{"title": "Hello", "body": "World"},
	}
	expectedResponse := map[string]interface{}{"title": "Hola", "body": "Mundo"}
	targetLang := "es"

	s.mockUsecase.On("TranslateJSON", mock.Anything, requestBody.Content, targetLang).Return(expectedResponse, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/translate", toJSON(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("lang", targetLang)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedResponse, resp["content"])
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *GeminiControllerTestSuite) TestTranslate_Success_DefaultLang() {
	// Arrange
	requestBody := TranslateDTO{
		Content: map[string]interface{}{"title": "Hello"},
	}
	expectedResponse := map[string]interface{}{"title": "Hello"} // Assuming default 'en' doesn't change it
	defaultLang := "en"

	s.mockUsecase.On("TranslateJSON", mock.Anything, requestBody.Content, defaultLang).Return(expectedResponse, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/translate", toJSON(requestBody))
	req.Header.Set("Content-Type", "application/json")
	// No 'lang' header is set
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedResponse, resp["content"])
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *GeminiControllerTestSuite) TestTranslate_InvalidJSON() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/translate", bytes.NewBufferString(`{"content":`)) // Malformed JSON
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid request body")
}

func (s *GeminiControllerTestSuite) TestTranslate_MissingContentField() {
	// Arrange
	invalidBody := map[string]string{"message": "This is not the right structure"}

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/translate", toJSON(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid request body")
}

func (s *GeminiControllerTestSuite) TestTranslate_UsecaseError() {
	// Arrange
	requestBody := TranslateDTO{
		Content: map[string]interface{}{"title": "Hello"},
	}
	targetLang := "fr"

	// Use a specific domain error that HandleError can map
	s.mockUsecase.On("TranslateJSON", mock.Anything, requestBody.Content, targetLang).Return(nil, domain.ErrTranslationMismatch).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/ai/translate", toJSON(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("lang", targetLang)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadGateway, w.Code) // HandleError maps this to 502
	s.Contains(w.Body.String(), domain.ErrTranslationMismatch.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestGeminiControllerTestSuite runs the entire suite
func TestGeminiControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GeminiControllerTestSuite))
}
