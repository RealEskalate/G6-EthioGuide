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

// MockPreferencesUsecase is a mock for the IPreferencesUsecase interface
type MockPreferencesUsecase struct {
	mock.Mock
}

func (m *MockPreferencesUsecase) CreateUserPreferences(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPreferencesUsecase) GetUserPreferences(ctx context.Context, userId string) (*domain.Preferences, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Preferences), args.Error(1)
}

func (m *MockPreferencesUsecase) UpdateUserPreferences(ctx context.Context, preferences *domain.Preferences) error {
	args := m.Called(ctx, preferences)
	return args.Error(0)
}

// PreferencesControllerTestSuite is the test suite for PreferencesController
type PreferencesControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockPreferencesUsecase
	controller  *PreferencesController
}

// SetupTest is run before each test in the suite
func (s *PreferencesControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.mockUsecase = new(MockPreferencesUsecase)
	s.controller = NewPreferencesController(s.mockUsecase)

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	}

	authGroup := s.router.Group("/auth/me/preferences")
	authGroup.Use(authMiddleware)
	{
		authGroup.GET("", s.controller.GetUserPreferences)
		authGroup.PATCH("", s.controller.UpdateUserPreferences)
	}
}

func (s *PreferencesControllerTestSuite) TestGetUserPreferences_Success() {
	// Arrange
	expectedPrefs := &domain.Preferences{
		UserID:            "test-user-id",
		PreferredLang:     domain.Amharic,
		PushNotification:  true,
		EmailNotification: false,
	}
	s.mockUsecase.On("GetUserPreferences", mock.Anything, "test-user-id").Return(expectedPrefs, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/me/preferences", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp PreferencesDTO
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(string(expectedPrefs.PreferredLang), resp.PreferredLang)
	s.Equal(expectedPrefs.PushNotification, resp.PushNotification)
	s.Equal(expectedPrefs.EmailNotification, resp.EmailNotification)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PreferencesControllerTestSuite) TestGetUserPreferences_NotFound() {
	// Arrange
	s.mockUsecase.On("GetUserPreferences", mock.Anything, "test-user-id").Return(nil, domain.ErrNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/me/preferences", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "Preferences not found")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PreferencesControllerTestSuite) TestUpdateUserPreferences_Success() {
	// Arrange
	reqBody := PreferencesDTO{
		PreferredLang:     "en",
		PushNotification:  false,
		EmailNotification: true,
	}
	expectedPrefsArg := &domain.Preferences{
		UserID:            "test-user-id",
		PreferredLang:     domain.English,
		PushNotification:  false,
		EmailNotification: true,
	}
	s.mockUsecase.On("UpdateUserPreferences", mock.Anything, expectedPrefsArg).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/auth/me/preferences", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp PreferencesDTO
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(reqBody.PreferredLang, resp.PreferredLang)
	s.Equal(reqBody.PushNotification, resp.PushNotification)
	s.Equal(reqBody.EmailNotification, resp.EmailNotification)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *PreferencesControllerTestSuite) TestUpdateUserPreferences_InvalidJSON() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/auth/me/preferences", bytes.NewBufferString(`{"preferredLang":`))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid input")
}

func (s *PreferencesControllerTestSuite) TestUpdateUserPreferences_UsecaseError() {
	// Arrange
	reqBody := PreferencesDTO{PreferredLang: "xx"} // Invalid lang
	expectedPrefsArg := &domain.Preferences{
		UserID:        "test-user-id",
		PreferredLang: "xx",
	}
	s.mockUsecase.On("UpdateUserPreferences", mock.Anything, expectedPrefsArg).Return(domain.ErrUnsupportedLanguage).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/auth/me/preferences", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), domain.ErrUnsupportedLanguage.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

// TestPreferencesControllerTestSuite runs the entire suite
func TestPreferencesControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesControllerTestSuite))
}
