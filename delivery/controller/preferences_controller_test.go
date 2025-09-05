package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"EthioGuide/delivery/controller"
	"EthioGuide/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock Usecase ---
type MockPreferencesUsecase struct {
	mock.Mock
}

func (m *MockPreferencesUsecase) CreateUserPreferences(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockPreferencesUsecase) GetUserPreferences(ctx context.Context, userID string) (*domain.Preferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Preferences), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPreferencesUsecase) UpdateUserPreferences(ctx context.Context, pref *domain.Preferences) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

// --- Test Suite ---
type PreferencesControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockPreferencesUsecase
	controller  *controller.PreferencesController
	recorder    *httptest.ResponseRecorder
}

func (s *PreferencesControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUsecase = new(MockPreferencesUsecase)
	s.controller = controller.NewPreferencesController(s.mockUsecase)
	s.router.GET("/users/me/preferences", func(c *gin.Context) {
		c.Set("userID", "user123")
		s.controller.GetUserPreferences(c)
	})
	s.router.PATCH("/users/me/preferences", func(c *gin.Context) {
		c.Set("userID", "user123")
		s.controller.UpdateUserPreferences(c)
	})
}

func TestPreferencesControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesControllerTestSuite))
}

// --- Tests ---

func (s *PreferencesControllerTestSuite) TestGetUserPreferences_Success() {
	expected := &domain.Preferences{
		UserID:            "user123",
		PreferredLang:     "en",
		PushNotification:  true,
		EmailNotification: false,
	}
	s.mockUsecase.On("GetUserPreferences", mock.Anything, "user123").Return(expected, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/users/me/preferences", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusOK, s.recorder.Code)
	var resp map[string]interface{}
	json.Unmarshal(s.recorder.Body.Bytes(), &resp)
	s.Equal("en", resp["preferredLang"])
	s.Equal(true, resp["pushNotification"])
	s.Equal(false, resp["emailNotification"])
}

func (s *PreferencesControllerTestSuite) TestUpdateUserPreferences_Success() {
	dto := controller.PreferencesDTO{
		PreferredLang:     "am",
		PushNotification:  false,
		EmailNotification: true,
	}
	body, _ := json.Marshal(dto)
	expected := &domain.Preferences{
		UserID:            "user123",
		PreferredLang:     "am",
		PushNotification:  false,
		EmailNotification: true,
	}
	s.mockUsecase.On("UpdateUserPreferences", mock.Anything, expected).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodPatch, "/users/me/preferences", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(s.recorder, req)

	s.Equal(http.StatusOK, s.recorder.Code)
	var resp map[string]interface{}
	json.Unmarshal(s.recorder.Body.Bytes(), &resp)
	s.Equal("am", resp["preferredLang"])
	s.Equal(false, resp["pushNotification"])
	s.Equal(true, resp["emailNotification"])
}
