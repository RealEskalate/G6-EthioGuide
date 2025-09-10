package infrastructure_test

import (
	. "EthioGuide/infrastructure"
	"context"
	"encoding/json"
	"errors"
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

// TranslationMiddlewareTestSuite is the test suite
type TranslationMiddlewareTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockGeminiUseCase
}

// SetupTest is run before each test in the suite
func (s *TranslationMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	s.mockUsecase = new(MockGeminiUseCase)

	// Apply the middleware to the router
	s.router.Use(NewTranslationMiddleware(s.mockUsecase))
}

func (s *TranslationMiddlewareTestSuite) TestTranslation_Success() {
	// Arrange
	originalResponse := gin.H{"message": "Hello"}
	translatedResponse := map[string]interface{}{"message": "Hola"}
	targetLang := "es"

	s.mockUsecase.On("TranslateJSON", mock.Anything, map[string]interface{}{"message": "Hello"}, targetLang).Return(translatedResponse, nil).Once()

	s.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, originalResponse)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("lang", targetLang)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("Hola", resp["message"]) // Check for translated content
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *TranslationMiddlewareTestSuite) TestTranslation_SkipsWhenLangIsEmpty() {
	// Arrange
	originalResponse := gin.H{"message": "Hello"}
	s.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, originalResponse)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	// No 'lang' header
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Hello") // Should contain original content
	s.mockUsecase.AssertNotCalled(s.T(), "TranslateJSON")
}

func (s *TranslationMiddlewareTestSuite) TestTranslation_SkipsWhenLangIsEnglish() {
	// Arrange
	originalResponse := gin.H{"message": "Hello"}
	s.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, originalResponse)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("lang", "en")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Hello")
	s.mockUsecase.AssertNotCalled(s.T(), "TranslateJSON")
}

func (s *TranslationMiddlewareTestSuite) TestTranslation_SkipsForNon200Status() {
	// Arrange
	s.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("lang", "es")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "Not Found") // Original body
	s.mockUsecase.AssertNotCalled(s.T(), "TranslateJSON")
}

func (s *TranslationMiddlewareTestSuite) TestTranslation_SkipsForNonJSONResponse() {
	// Arrange
	s.router.GET("/test", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<h1>Hello</h1>"))
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("lang", "es")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal("<h1>Hello</h1>", w.Body.String()) // Original body
	s.mockUsecase.AssertNotCalled(s.T(), "TranslateJSON")
}

func (s *TranslationMiddlewareTestSuite) TestTranslation_ReturnsOriginalBodyOnUsecaseError() {
	// Arrange
	originalResponse := gin.H{"message": "Hello"}
	targetLang := "es"
	s.mockUsecase.On("TranslateJSON", mock.Anything, mock.Anything, targetLang).Return(nil, errors.New("translation service down")).Once()

	s.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, originalResponse)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("lang", targetLang)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Hello") // Check for original content
	s.mockUsecase.AssertExpectations(s.T())
}

// TestTranslationMiddlewareTestSuite runs the entire suite
func TestTranslationMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(TranslationMiddlewareTestSuite))
}
