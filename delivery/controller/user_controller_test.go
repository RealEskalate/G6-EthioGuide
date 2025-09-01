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

func (m *MockUserUsecase) ForgetPassword(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	args := m.Called(ctx, resetToken, newPassword)
	return args.Error(0)
}

func (m *MockUserUsecase) VerifyAccount(ctx context.Context, activationToken string) error {
	args := m.Called(ctx, activationToken)
	return args.Error(0)
}

// MockSearchUsecase is a mock implementation of the ISearchUseCase interface.
type MockSearchUsecase struct {
	mock.Mock
}

var _ domain.ISearchUseCase = (*MockSearchUsecase)(nil)

func (m *MockSearchUsecase) Search(ctx context.Context, filter domain.SearchFilterRequest) (*domain.SearchResult, error) {
	args := m.Called(ctx, filter)
	var res *domain.SearchResult
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.SearchResult)
	}
	return res, args.Error(1)
}

// --- Test Suite Definition ---

type UserControllerTestSuite struct {
	suite.Suite
	router             *gin.Engine
	mockUserUsecase    *MockUserUsecase
	mockSearchUsecase  *MockSearchUsecase // Added search mock
	controller         *UserController
	recorder           *httptest.ResponseRecorder
	refreshTokenTTL    time.Duration
}

func (s *UserControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *UserControllerTestSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUserUsecase = new(MockUserUsecase)
	s.mockSearchUsecase = new(MockSearchUsecase) // Instantiate search mock
	s.refreshTokenTTL = 15 * time.Minute

	// CRITICAL FIX: Correctly instantiate the controller with all dependencies.
	s.controller = NewUserController(s.mockUserUsecase, s.mockSearchUsecase, s.refreshTokenTTL)

	// Setup routes
	s.router.POST("/register", s.controller.Register)
	s.router.POST("/login", s.controller.Login)
	s.router.POST("/refresh", s.controller.HandleRefreshToken)
	s.router.POST("/forgot-password", s.controller.HandleForgot)
	s.router.POST("/reset-password", s.controller.HandleReset)
	s.router.POST("/verify", s.controller.HandleVerify)
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

		s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusCreated, s.recorder.Code)
		s.mockUserUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Email Exists", func() {
		// Arrange
		s.SetupTest() // Reset recorder and mocks
		reqBody := RegisterRequest{Name: "Test User", Email: "test@example.com", Password: "password123", Username: "testuser"}
		jsonBody, _ := json.Marshal(reqBody)

		s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(domain.ErrEmailExists).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusConflict, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrEmailExists.Error())
		s.mockUserUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestLogin() {
	reqBody := LoginRequest{Identifier: "test@example.com", Password: "password"}
	jsonBody, _ := json.Marshal(reqBody)
	mockAccount := &domain.Account{
		ID:    "user-123",
		Name:  "Test User",
		Email: "test@example.com",
		UserDetail: &domain.UserDetail{
			Username: "testuser",
		},
	}

	expectedUserResponse := ToUserResponse(mockAccount)

	s.Run("Success - Web Client", func() {
		// Arrange
		s.SetupTest()
		s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
			Return(mockAccount, "access_token_123", "refresh_token_abc", nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		cookie := s.recorder.Result().Cookies()[0]
		s.Equal("refresh_token", cookie.Name)
		s.Equal("refresh_token_abc", cookie.Value)

		var resp LoginResponse
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("access_token_123", resp.AccessToken)
		s.Empty(resp.RefreshToken)

		s.Equal(expectedUserResponse, resp.User)

		s.mockUserUsecase.AssertExpectations(s.T())
	})

	s.Run("Success - Mobile Client", func() {
		// Arrange
		s.SetupTest()
		s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
			Return(mockAccount, "access_token_123", "refresh_token_abc", nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Client-Type", "mobile")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		s.Empty(s.recorder.Result().Cookies())

		var resp LoginResponse
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("access_token_123", resp.AccessToken)
		s.Equal("refresh_token_abc", resp.RefreshToken)

		s.Equal(expectedUserResponse, resp.User)

		s.mockUserUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid Credentials", func() {
		s.SetupTest()
		s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
			Return(nil, "", "", domain.ErrAuthenticationFailed).Once()

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.mockUserUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestHandleRefreshToken() {
	s.Run("Success - Web Client", func() {
		s.SetupTest()
		oldRefreshToken := "old_refresh_from_cookie"
		newAccessToken := "new_access_token_web"

		s.mockUserUsecase.On("RefreshTokenForWeb", mock.Anything, oldRefreshToken).Return(newAccessToken, nil).Once()

		// Act (Request with a cookie)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: oldRefreshToken})
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		var resp map[string]string
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal(newAccessToken, resp["access_token"])
		s.mockUserUsecase.AssertExpectations(s.T())
	})

	s.Run("Success - Mobile Client", func() {
		s.SetupTest()
		oldRefreshToken := "old_refresh_from_header"
		newAccessToken := "new_access_token_mobile"
		newRefreshToken := "new_refresh_token_mobile"

		s.mockUserUsecase.On("RefreshTokenForMobile", mock.Anything, oldRefreshToken).
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
		s.mockUserUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Web Client No Cookie", func() {
		s.SetupTest()

		// Act (Request with no cookie)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.mockUserUsecase.AssertNotCalled(s.T(), "RefreshTokenForWeb")
	})

	s.Run("Failure - Mobile Client No Auth Header", func() {
		s.SetupTest()

		// Act (Request with mobile header but no auth header)
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		req.Header.Set("X-Client-Type", "mobile")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusUnauthorized, s.recorder.Code)
		s.mockUserUsecase.AssertNotCalled(s.T(), "RefreshTokenForMobile")
	})
}

func (s *UserControllerTestSuite) TestHandleForgot() {
	s.Run("Success", func() {
		// Arrange
		s.SetupTest()
		reqBody := ForgotDTO{Email: "test@example.com"}
		jsonBody, _ := json.Marshal(reqBody)
		resetToken := "reset_token_123"

		s.mockUserUsecase.On("ForgetPassword", mock.Anything, reqBody.Email).Return(resetToken, nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/forgot-password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		var resp map[string]string
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal(resetToken, resp["resetToken"])
		s.mockUserUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestHandleReset() {
	s.Run("Success", func() {
		// Arrange
		s.SetupTest()
		reqBody := ResetDTO{ResetToken: "reset_token_123", NewPassword: "new_password"}
		jsonBody, _ := json.Marshal(reqBody)

		s.mockUserUsecase.On("ResetPassword", mock.Anything, reqBody.ResetToken, reqBody.NewPassword).Return(nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/reset-password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		var resp map[string]string
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("Password Updated Successfully", resp["message"])
		s.mockUserUsecase.AssertExpectations(s.T())
	})
}

func (s *UserControllerTestSuite) TestHandleVerify() {
	s.Run("Success", func() {
		// Arrange
		s.SetupTest()
		reqBody := ActivateDTO{ActivateToken: "activation_token_123"}
		jsonBody, _ := json.Marshal(reqBody)

		s.mockUserUsecase.On("VerifyAccount", mock.Anything, reqBody.ActivateToken).Return(nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		var resp map[string]string
		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
		s.Equal("User Activated Successfully", resp["message"])
		s.mockUserUsecase.AssertExpectations(s.T())
	})
}
