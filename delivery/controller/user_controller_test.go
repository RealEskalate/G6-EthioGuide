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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks & Placeholders ---

// MockUserUsecase is a mock implementation of the IUserUsecase interface.
type MockUserUsecase struct {
	mock.Mock
}

// Ensure MockUserUsecase implements IUserUsecase at compile time.
var _ domain.IUserUsecase = (*MockUserUsecase)(nil)

func (m *MockUserUsecase) Register(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockUserUsecase) Login(ctx context.Context, identifier, password string) (*domain.Account, string, string, error) {
	args := m.Called(ctx, identifier, password)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.String(1), args.String(2), args.Error(3)
}

func (m *MockUserUsecase) RefreshTokenForWeb(ctx context.Context, refreshToken string) (string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) RefreshTokenForMobile(ctx context.Context, refreshToken string) (string, string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

// Add missing GetProfile method to satisfy domain.IUserUsecase
func (m *MockUserUsecase) GetProfile(ctx context.Context, userID string) (*domain.Account, error) {
	args := m.Called(ctx, userID)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.Error(1)
}

func (m *MockUserUsecase) UpdatePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	args := m.Called(ctx, userID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserUsecase) LoginWithSocial(ctx context.Context, provider domain.AuthProvider, code string) (*domain.Account, string, string, error) {
	args := m.Called(ctx, provider, code)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.String(1), args.String(2), args.Error(3)
}

func (m *MockUserUsecase) UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) (*domain.Account, error) {
	args := m.Called(ctx, userID, updates)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.Error(1)
}

// --- Test Suite Definition ---

type UserControllerTestSuite struct {
	suite.Suite
	router          *gin.Engine
	mockUsecase     *MockUserUsecase
	controller      *UserController
	recorder        *httptest.ResponseRecorder
	refreshTokenTTL time.Duration
}

func (s *UserControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *UserControllerTestSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUsecase = new(MockUserUsecase)
	s.refreshTokenTTL = 15 * time.Minute

	// Create a new controller for each test, injecting the mock.
	s.controller = NewUserController(s.mockUsecase, s.refreshTokenTTL)

	// Setup routes
	s.router.POST("/register", s.controller.Register)
	s.router.POST("/login", s.controller.Login)
	s.router.POST("/refresh", s.controller.HandleRefreshToken)
	s.router.POST("/social", s.controller.SocialLogin)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

// --- Test Cases ---

func (s *UserControllerTestSuite) TestRegister() {
	s.Run("Success", func() {
		// Arrange
		reqBody := RegisterRequest{Name: "Test User", Email: "test@example.com", Password: "password123", Username: "testuser"}
		jsonBody, _ := json.Marshal(reqBody)

		// We use mock.Anything for the account because the pointer address will differ.
		s.mockUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusCreated, s.recorder.Code)
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Email Exists", func() {
		// Arrange
		s.SetupTest() // Reset recorder and mocks
		reqBody := RegisterRequest{Name: "Test User", Email: "test@example.com", Password: "password123", Username: "testuser"}
		jsonBody, _ := json.Marshal(reqBody)

		s.mockUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(domain.ErrEmailExists).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusConflict, s.recorder.Code) // Assert correct status from HandleError
		s.Contains(s.recorder.Body.String(), domain.ErrEmailExists.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestLogin() {
	reqBody := LoginRequest{Identifier: "test@example.com", Password: "password"}
	jsonBody, _ := json.Marshal(reqBody)
	mockAccount := &domain.Account{
		Email:      "test@example.com",
		UserDetail: &domain.UserDetail{Username: "testuser"},
	}

	s.Run("Success - Web Client", func() {
		// Arrange
		s.SetupTest()
		s.mockUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
			Return(mockAccount, "access_token_123", "refresh_token_abc", nil).Once()

		// Act (No X-Client-Type header)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		// Check for cookie
		cookie := s.recorder.Result().Cookies()[0]
		s.Equal("refresh_token", cookie.Name)
		s.Equal("refresh_token_abc", cookie.Value)
		s.Equal(int(s.refreshTokenTTL.Seconds()), cookie.MaxAge)

		// Check response body (should NOT contain refresh token)
		var resp LoginResponse
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("access_token_123", resp.AccessToken)
		s.Empty(resp.RefreshToken)

		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Success - Mobile Client", func() {
		// Arrange
		s.SetupTest()
		s.mockUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
			Return(mockAccount, "access_token_123", "refresh_token_abc", nil).Once()

		// Act (WITH X-Client-Type header)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Client-Type", "mobile")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		// Check that NO cookie was set
		s.Empty(s.recorder.Result().Cookies())

		// Check response body (SHOULD contain refresh token)
		var resp LoginResponse
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("access_token_123", resp.AccessToken)
		s.Equal("refresh_token_abc", resp.RefreshToken)

		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid Credentials", func() {
		s.SetupTest()
		s.mockUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
			Return(nil, "", "", domain.ErrAuthenticationFailed).Once()

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestHandleRefreshToken() {
	s.Run("Success - Web Client", func() {
		s.SetupTest()
		oldRefreshToken := "old_refresh_from_cookie"
		newAccessToken := "new_access_token_web"

		s.mockUsecase.On("RefreshTokenForWeb", mock.Anything, oldRefreshToken).Return(newAccessToken, nil).Once()

		// Act (Request with a cookie)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: oldRefreshToken})
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		var resp map[string]string
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal(newAccessToken, resp["access_token"])
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Success - Mobile Client", func() {
		s.SetupTest()
		oldRefreshToken := "old_refresh_from_header"
		newAccessToken := "new_access_token_mobile"
		newRefreshToken := "new_refresh_token_mobile"

		s.mockUsecase.On("RefreshTokenForMobile", mock.Anything, oldRefreshToken).
			Return(newAccessToken, newRefreshToken, nil).Once()

		// Act (Request with headers)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		req.Header.Set("Authorization", "Bearer "+oldRefreshToken)
		req.Header.Set("X-Client-Type", "mobile")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		var resp map[string]string
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal(newAccessToken, resp["access_token"])
		s.Equal(newRefreshToken, resp["refresh_token"])
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Web Client No Cookie", func() {
		s.SetupTest()

		// Act (Request with no cookie)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.mockUsecase.AssertNotCalled(s.T(), "RefreshTokenForWeb")
	})

	s.Run("Failure - Mobile Client No Auth Header", func() {
		s.SetupTest()

		// Act (Request with mobile header but no auth header)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		req.Header.Set("X-Client-Type", "mobile")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.mockUsecase.AssertNotCalled(s.T(), "RefreshTokenForMobile")
	})
}

func (s *UserControllerTestSuite) TestGetProfile() {
	// Add the route for testing
	s.router.GET("/profile", func(c *gin.Context) {
		// Simulate middleware setting userID in context
		c.Set("userID", "user123")
		s.controller.GetProfile(c)
	})

	mockAccount := &domain.Account{
		ID:    "user123",
		Email: "test@example.com",
		Name:  "Test User",
		UserDetail: &domain.UserDetail{
			Username:         "testuser",
			SubscriptionPlan: domain.SubscriptionNone,
			IsBanned:         false,
			IsVerified:       true,
		},
	}

	s.Run("Success", func() {
		s.SetupTest()
		// Add the route again after SetupTest
		s.router.GET("/profile", func(c *gin.Context) {
			c.Set("userID", "user123")
			s.controller.GetProfile(c)
		})

		s.mockUsecase.On("GetProfile", mock.Anything, "user123").Return(mockAccount, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusOK, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "testuser")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - User Not Found", func() {
		s.SetupTest()
		s.router.GET("/profile", func(c *gin.Context) {
			c.Set("userID", "user123")
			s.controller.GetProfile(c)
		})

		s.mockUsecase.On("GetProfile", mock.Anything, "user123").Return(nil, domain.ErrUserNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusNotFound, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrUserNotFound.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - No userID in context", func() {
		s.SetupTest()
		// Route without setting userID
		s.router.GET("/profile", func(c *gin.Context) {
			s.controller.GetProfile(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "User ID not found")
	})
}

func (s *UserControllerTestSuite) TestUpdatePassword() {
	// Setup route for PATCH /me/password
	s.router.PATCH("/me/password", func(c *gin.Context) {
		// Simulate middleware setting userID in context
		c.Set("userID", "user-123")
		s.controller.UpdatePassword(c)
	})

	s.Run("Success", func() {
		s.SetupTest()
		s.router.PATCH("/me/password", func(c *gin.Context) {
			c.Set("userID", "user-123")
			s.controller.UpdatePassword(c)
		})

		reqBody := ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "newpass123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		s.mockUsecase.On("UpdatePassword", mock.Anything, "user-123", "oldpass", "newpass123").Return(nil).Once()

		req, _ := http.NewRequest(http.MethodPatch, "/me/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusOK, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Password updated successfully")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - No userID in context", func() {
		s.SetupTest()
		// Route without setting userID
		s.router.PATCH("/me/password", func(c *gin.Context) {
			s.controller.UpdatePassword(c)
		})

		reqBody := ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "newpass123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPatch, "/me/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "User ID not found")
	})

	s.Run("Failure - Invalid request body", func() {
		s.SetupTest()
		s.router.PATCH("/me/password", func(c *gin.Context) {
			c.Set("userID", "user-123")
			s.controller.UpdatePassword(c)
		})

		req, _ := http.NewRequest(http.MethodPatch, "/me/password", bytes.NewBuffer([]byte("not-json")))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusBadRequest, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Invalid request body")
	})

	s.Run("Failure - Usecase error", func() {
		s.SetupTest()
		s.router.PATCH("/me/password", func(c *gin.Context) {
			c.Set("userID", "user-123")
			s.controller.UpdatePassword(c)
		})

		reqBody := ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "newpass123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		s.mockUsecase.On("UpdatePassword", mock.Anything, "user-123", "oldpass", "newpass123").Return(domain.ErrAuthenticationFailed).Once()

		req, _ := http.NewRequest(http.MethodPatch, "/me/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrAuthenticationFailed.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestSocialLogin() {
	reqBody := SocialLoginRequest{Provider: "google", Code: "some_google_auth_code"}
	jsonBody, _ := json.Marshal(reqBody)
	mockAccount := &domain.Account{
		Email:        "social.user@example.com",
		AuthProvider: domain.AuthProviderGoogle,
		UserDetail:   &domain.UserDetail{Username: "testuser"},
	}

	s.Run("Success - Web Client", func() {
		// Arrange
		s.SetupTest()
		s.mockUsecase.On("LoginWithSocial", mock.Anything, reqBody.Provider, reqBody.Code).
			Return(mockAccount, "access_token_social_123", "refresh_token_social_abc", nil).Once()

		// Act (No X-Client-Type header)
		req, _ := http.NewRequest(http.MethodPost, "/social", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		// Check for the auth cookie
		cookie := s.recorder.Result().Cookies()[0]
		s.Equal("refresh_token", cookie.Name)
		s.Equal("refresh_token_social_abc", cookie.Value)
		s.Equal(int(s.refreshTokenTTL.Seconds()), cookie.MaxAge)

		// Check response body (should NOT contain refresh token)
		var resp LoginResponse
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("access_token_social_123", resp.AccessToken)
		s.Empty(resp.RefreshToken)

		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Success - Mobile Client", func() {
		// Arrange
		s.SetupTest()
		s.mockUsecase.On("LoginWithSocial", mock.Anything, reqBody.Provider, reqBody.Code).
			Return(mockAccount, "access_token_social_123", "refresh_token_social_abc", nil).Once()

		// Act (WITH X-Client-Type header)
		req, _ := http.NewRequest(http.MethodPost, "/social", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Client-Type", "mobile")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		// Check that NO cookie was set
		s.Empty(s.recorder.Result().Cookies())

		// Check response body (SHOULD contain refresh token)
		var resp LoginResponse
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("access_token_social_123", resp.AccessToken)
		s.Equal("refresh_token_social_abc", resp.RefreshToken)

		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Usecase returns an error", func() {
		s.SetupTest()
		s.mockUsecase.On("LoginWithSocial", mock.Anything, reqBody.Provider, reqBody.Code).
			Return(nil, "", "", domain.ErrAuthenticationFailed).Once()

		req, _ := http.NewRequest(http.MethodPost, "/social", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrAuthenticationFailed.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid Request Body (missing provider)", func() {
		s.SetupTest()
		invalidBody := `{"code": "some_code"}` // Missing 'provider'

		req, _ := http.NewRequest(http.MethodPost, "/social", bytes.NewBufferString(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusBadRequest, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Invalid request")

		// Ensure the usecase was not called because validation failed first
		s.mockUsecase.AssertNotCalled(s.T(), "LoginWithSocial")
	})
}

func (s *UserControllerTestSuite) TestUpdateProfile() {
	// Setup route for PATCH /me
	s.router.PATCH("/me", func(c *gin.Context) {
		// Simulate middleware setting userID in context
		c.Set("userID", "user-123")
		s.controller.UpdateProfile(c)
	})

	updatedAccount := &domain.Account{
		ID:    "user-123",
		Name:  "Updated Name",
		Email: "updated@example.com",
		UserDetail: &domain.UserDetail{
			Username:         "updateduser",
			SubscriptionPlan: domain.SubscriptionNone,
			IsBanned:         false,
			IsVerified:       true,
		},
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.router.PATCH("/me", func(c *gin.Context) {
			c.Set("userID", "user-123")
			s.controller.UpdateProfile(c)
		})

		updates := map[string]interface{}{
			"name": "Updated Name",
		}
		jsonBody, _ := json.Marshal(updates)

		s.mockUsecase.On("UpdateProfile", mock.Anything, "user-123", updates).Return(updatedAccount, nil).Once()

		req, _ := http.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusOK, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Updated Name")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - No userID in context", func() {
		s.SetupTest()
		// Route without setting userID
		s.router.PATCH("/me", func(c *gin.Context) {
			s.controller.UpdateProfile(c)
		})

		updates := map[string]interface{}{
			"name": "Updated Name",
		}
		jsonBody, _ := json.Marshal(updates)

		req, _ := http.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "User ID not found")
	})

	s.Run("Failure - Invalid form data", func() {
		s.SetupTest()
		s.router.PATCH("/me", func(c *gin.Context) {
			c.Set("userID", "user-123")
			s.controller.UpdateProfile(c)
		})

		req, _ := http.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer([]byte("not-json")))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusBadRequest, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Invalid form data")
	})

	s.Run("Failure - Usecase error", func() {
		s.SetupTest()
		s.router.PATCH("/me", func(c *gin.Context) {
			c.Set("userID", "user-123")
			s.controller.UpdateProfile(c)
		})

		updates := map[string]interface{}{
			"name": "Updated Name",
		}
		jsonBody, _ := json.Marshal(updates)

		s.mockUsecase.On("UpdateProfile", mock.Anything, "user-123", updates).Return(nil, domain.ErrNotFound).Once()

		req, _ := http.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusNotFound, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrNotFound.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})
}
