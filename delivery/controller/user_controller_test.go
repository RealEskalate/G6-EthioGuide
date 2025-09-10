package controller_test

import (
	. "EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks for all dependencies ---

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, user *domain.Account) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserUsecase) Login(ctx context.Context, identifier, password string) (*domain.Account, string, string, error) {
	args := m.Called(ctx, identifier, password)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*domain.Account), args.String(1), args.String(2), args.Error(3)
}
func (m *MockUserUsecase) VerifyAccount(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockUserUsecase) RefreshTokenForWeb(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}
func (m *MockUserUsecase) RefreshTokenForMobile(ctx context.Context, token string) (string, string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.String(1), args.Error(2)
}
func (m *MockUserUsecase) ForgetPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}
func (m *MockUserUsecase) ResetPassword(ctx context.Context, token, password string) error {
	args := m.Called(ctx, token, password)
	return args.Error(0)
}
func (m *MockUserUsecase) GetProfile(ctx context.Context, userID string) (*domain.Account, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}
func (m *MockUserUsecase) UpdatePassword(ctx context.Context, userID, old, new string) error {
	args := m.Called(ctx, userID, old, new)
	return args.Error(0)
}
func (m *MockUserUsecase) LoginWithSocial(ctx context.Context, provider domain.AuthProvider, code string) (*domain.Account, string, string, error) {
	args := m.Called(ctx, provider, code)
	if args.Get(0) == nil {
		return nil, "", "", args.Error(3)
	}
	return args.Get(0).(*domain.Account), args.String(1), args.String(2), args.Error(3)
}
func (m *MockUserUsecase) Logout(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockUserUsecase) UpdateProfile(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	args := m.Called(ctx, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}
func (m *MockUserUsecase) RegisterOrg(ctx context.Context, name, email, orgType string) error {
	args := m.Called(ctx, name, email, orgType)
	return args.Error(0)
}
func (m *MockUserUsecase) GetOrgs(ctx context.Context, filter domain.GetOrgsFilter) ([]*domain.Account, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Account), args.Get(1).(int64), args.Error(2)
}
func (m *MockUserUsecase) GetOrgById(ctx context.Context, orgId string) (*domain.Account, error) {
	args := m.Called(ctx, orgId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}
func (m *MockUserUsecase) UpdateOrgFields(ctx context.Context, orgId string, update map[string]interface{}) error {
	args := m.Called(ctx, orgId, update)
	return args.Error(0)
}

type MockSearchUsecase struct {
	mock.Mock
}

func (m *MockSearchUsecase) Search(ctx context.Context, filter domain.SearchFilterRequest) (*domain.SearchResult, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SearchResult), args.Error(1)
}

type MockChecklistUsecase struct {
	mock.Mock
}

func (m *MockChecklistUsecase) CreateChecklist(ctx context.Context, userid, procedureID string) (*domain.UserProcedure, error) {
	args := m.Called(ctx, userid, procedureID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserProcedure), args.Error(1)
}
func (m *MockChecklistUsecase) GetProcedures(ctx context.Context, userid string) ([]*domain.UserProcedure, error) {
	args := m.Called(ctx, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserProcedure), args.Error(1)
}
func (m *MockChecklistUsecase) GetChecklistByUserProcedureID(ctx context.Context, userprocedureID string) ([]*domain.Checklist, error) {
	args := m.Called(ctx, userprocedureID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Checklist), args.Error(1)
}
func (m *MockChecklistUsecase) UpdateChecklist(ctx context.Context, checklistID string) (*domain.Checklist, error) {
	args := m.Called(ctx, checklistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Checklist), args.Error(1)
}

// UserControllerTestSuite is the test suite
type UserControllerTestSuite struct {
	suite.Suite
	router               *gin.Engine
	mockUserUsecase      *MockUserUsecase
	mockSearchUsecase    *MockSearchUsecase
	mockChecklistUsecase *MockChecklistUsecase
	controller           *UserController
}

func (s *UserControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.mockUserUsecase = new(MockUserUsecase)
	s.mockSearchUsecase = new(MockSearchUsecase)
	s.mockChecklistUsecase = new(MockChecklistUsecase)
	s.controller = NewUserController(s.mockUserUsecase, s.mockSearchUsecase, s.mockChecklistUsecase, time.Hour)

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Set("userRole", domain.RoleUser)
		c.Next()
	}

	authGroup := s.router.Group("/auth")
	{
		authGroup.POST("/register", s.controller.Register)
		authGroup.POST("/login", s.controller.Login)
		authGroup.POST("/refresh", s.controller.HandleRefreshToken)
		authGroup.POST("/social", s.controller.SocialLogin)
		authGroup.POST("/forgot", s.controller.HandleForgot)
		authGroup.POST("/reset", s.controller.HandleReset)
		authGroup.POST("/verify", s.controller.HandleVerify)
		authGroup.POST("/logout", authMiddleware, s.controller.Logout)
		meGroup := authGroup.Group("/me", authMiddleware)
		{
			meGroup.GET("", s.controller.GetProfile)
			meGroup.PATCH("", s.controller.UpdateProfile)
			meGroup.PATCH("/password", s.controller.UpdatePassword)
		}
	}

	orgsGroup := s.router.Group("/orgs")
	{
		orgsGroup.POST("", authMiddleware, s.controller.HandleCreateOrg)
		orgsGroup.GET("", s.controller.HandleGetOrgs)
		orgsGroup.GET("/:id", s.controller.HandleGetOrgById)
		orgsGroup.PATCH("/:id", authMiddleware, s.controller.HandleUpdateOrgs)
	}
	s.router.GET("/search", s.controller.HandleSearch)
	checklistsGroup := s.router.Group("/checklists", authMiddleware)
	{
		checklistsGroup.POST("", s.controller.HandleCreateChecklist)
		checklistsGroup.GET("/myProcedures", s.controller.HandleGetProcedures)
		checklistsGroup.GET("/:userProcedureId", s.controller.HandleGetChecklistById)
		checklistsGroup.PATCH("/:checklistID", s.controller.HandleUpdateChecklist)
	}
}

// --- Registration Tests ---
func (s *UserControllerTestSuite) TestRegister_Success() {
	// Arrange
	reqBody := RegisterRequest{Name: "Test", Username: "tester", Email: "test@example.com", Password: "password123"}
	s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "User registered successfully")
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestRegister_EmailExistsError() {
	// Arrange
	reqBody := RegisterRequest{Name: "Test", Username: "tester", Email: "test@example.com", Password: "password123"}
	s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(domain.ErrEmailExists).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusConflict, w.Code)
	s.Contains(w.Body.String(), domain.ErrEmailExists.Error())
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Login Tests ---
func (s *UserControllerTestSuite) TestLogin_Success_Web() {
	// Arrange
	reqBody := LoginRequest{Identifier: "test@example.com", Password: "password"}
	expectedAccount := &domain.Account{ID: "user-1", UserDetail: &domain.UserDetail{}}
	s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).Return(expectedAccount, "access-token", "refresh-token", nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Header().Get("Set-Cookie"), "refresh_token=refresh-token")
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestLogin_Success_Mobile() {
	// Arrange
	reqBody := LoginRequest{Identifier: "test@example.com", Password: "password"}
	expectedAccount := &domain.Account{ID: "user-1", UserDetail: &domain.UserDetail{}}
	s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).Return(expectedAccount, "access-token", "refresh-token", nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Type", "mobile")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.NotContains(w.Header().Get("Set-Cookie"), "refresh_token")
	var resp LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("refresh-token", resp.RefreshToken)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestLogin_AuthFailed() {
	// Arrange
	reqBody := LoginRequest{Identifier: "test@example.com", Password: "wrong-password"}
	s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).Return(nil, "", "", domain.ErrAuthenticationFailed).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), domain.ErrAuthenticationFailed.Error())
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Profile Tests ---
func (s *UserControllerTestSuite) TestGetProfile_Success() {
	// Arrange
	expectedAccount := &domain.Account{ID: "test-user-id", UserDetail: &domain.UserDetail{Username: "tester"}}
	s.mockUserUsecase.On("GetProfile", mock.Anything, "test-user-id").Return(expectedAccount, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/me", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedAccount.ID, resp.ID)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestUpdateProfile_Success() {
	// Arrange
	nameUpdate := "New Name"
	reqBody := UserUpdateRequest{Name: &nameUpdate}
	existingAccount := &domain.Account{ID: "test-user-id", Name: "Old Name", UserDetail: &domain.UserDetail{}}
	updatedAccount := &domain.Account{ID: "test-user-id", Name: "New Name", UserDetail: &domain.UserDetail{}}

	s.mockUserUsecase.On("GetProfile", mock.Anything, "test-user-id").Return(existingAccount, nil).Once()
	s.mockUserUsecase.On("UpdateProfile", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(updatedAccount, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/auth/me", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(updatedAccount.Name, resp.Name)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Password & Token Tests ---
func (s *UserControllerTestSuite) TestUpdatePassword_Success() {
	// Arrange
	reqBody := ChangePasswordRequest{OldPassword: "old", NewPassword: "new"}
	s.mockUserUsecase.On("UpdatePassword", mock.Anything, "test-user-id", "old", "new").Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/auth/me/password", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Password updated successfully")
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleRefreshToken_Web_Success() {
	// Arrange
	s.mockUserUsecase.On("RefreshTokenForWeb", mock.Anything, "old-refresh-token").Return("new-access-token", nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "old-refresh-token"})
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "new-access-token")
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleRefreshToken_Mobile_Success() {
	// Arrange
	s.mockUserUsecase.On("RefreshTokenForMobile", mock.Anything, "old-refresh-token").Return("new-access-token", "new-refresh-token", nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer old-refresh-token")
	req.Header.Set("X-Client-Type", "mobile")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "new-access-token")
	s.Contains(w.Body.String(), "new-refresh-token")
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleForgot_Success() {
	// Arrange
	reqBody := ForgotDTO{Email: "test@example.com"}
	s.mockUserUsecase.On("ForgetPassword", mock.Anything, reqBody.Email).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/forgot", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Logout Tests ---
func (s *UserControllerTestSuite) TestLogout_Web() {
	// Arrange
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	response := w.Result()
	cookies := response.Cookies()
	var refreshTokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			refreshTokenCookie = c
			break
		}
	}
	s.NotNil(refreshTokenCookie)
	s.Less(refreshTokenCookie.MaxAge, 0)
}

func (s *UserControllerTestSuite) TestLogout_Mobile() {
	// Arrange
	s.mockUserUsecase.On("Logout", mock.Anything, "test-user-id").Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("X-Client-Type", "mobile")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Org & Search & Checklist Tests ---
func (s *UserControllerTestSuite) TestHandleCreateOrg_Success() {
	// Arrange
	reqBody := OrgCreateRequest{Name: "New Org", Email: "org@example.com", OrgType: "gov"}
	s.mockUserUsecase.On("RegisterOrg", mock.Anything, reqBody.Name, reqBody.Email, reqBody.OrgType).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orgs", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleSearch_Success() {
	// Arrange
	expectedResult := &domain.SearchResult{Procedures: []*domain.Procedure{{ID: "proc-1"}}}
	filter := domain.SearchFilterRequest{Query: "passport", Page: 1, Limit: 10}
	s.mockSearchUsecase.On("Search", mock.Anything, filter).Return(expectedResult, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?q=passport&page=1&limit=10", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp SearchResultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Len(resp.Procedures, 1)
	s.mockSearchUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleCreateChecklist_Success() {
	// Arrange
	reqBody := CreateChecklistRequest{ProcedureID: "proc-1"}
	expectedResult := &domain.UserProcedure{ID: "up-1"}
	s.mockChecklistUsecase.On("CreateChecklist", mock.Anything, "test-user-id", "proc-1").Return(expectedResult, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/checklists", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.mockChecklistUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleGetProcedures_Success() {
	// Arrange
	expectedResult := []*domain.UserProcedure{{ID: "up-1", UserID: "test-user-id"}}
	s.mockChecklistUsecase.On("GetProcedures", mock.Anything, "test-user-id").Return(expectedResult, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/checklists/myProcedures", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string][]UserProcedureResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Len(resp["message"], 1)
	s.mockChecklistUsecase.AssertExpectations(s.T())
}

// --- Profile Tests (Sad Paths) ---
func (s *UserControllerTestSuite) TestGetProfile_NotFound() {
	// Arrange
	s.mockUserUsecase.On("GetProfile", mock.Anything, "test-user-id").Return(nil, domain.ErrUserNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/me", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrUserNotFound.Error())
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestUpdateProfile_UsecaseError() {
	// Arrange
	nameUpdate := "New Name"
	reqBody := UserUpdateRequest{Name: &nameUpdate}
	existingAccount := &domain.Account{ID: "test-user-id", Name: "Old Name", UserDetail: &domain.UserDetail{}}

	s.mockUserUsecase.On("GetProfile", mock.Anything, "test-user-id").Return(existingAccount, nil).Once()
	s.mockUserUsecase.On("UpdateProfile", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil, domain.ErrConflict).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/auth/me", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusConflict, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Social Login Tests ---
func (s *UserControllerTestSuite) TestSocialLogin_Success() {
	// Arrange
	reqBody := SocialLoginRequest{Provider: domain.AuthProviderGoogle, Code: "auth_code"}
	expectedAccount := &domain.Account{ID: "social-user-1", UserDetail: &domain.UserDetail{}}
	s.mockUserUsecase.On("LoginWithSocial", mock.Anything, reqBody.Provider, reqBody.Code).Return(expectedAccount, "access", "refresh", nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/social", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Password Management Tests ---

func (s *UserControllerTestSuite) TestHandleReset_Success() {
	// Arrange
	reqBody := ResetDTO{ResetToken: "valid-token", NewPassword: "newPassword123"}
	s.mockUserUsecase.On("ResetPassword", mock.Anything, reqBody.ResetToken, reqBody.NewPassword).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/reset", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleVerify_Success() {
	// Arrange
	reqBody := ActivateDTO{ActivateToken: "valid-token"}
	s.mockUserUsecase.On("VerifyAccount", mock.Anything, reqBody.ActivateToken).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Organization Tests ---
func (s *UserControllerTestSuite) TestHandleGetOrgs_Success() {
	// Arrange
	expectedOrgs := []*domain.Account{{ID: "org-1", Name: "Test Org"}}
	var total int64 = 1
	filter := domain.GetOrgsFilter{Type: "gov", Query: "Test", Page: 1, PageSize: 10}
	s.mockUserUsecase.On("GetOrgs", mock.Anything, filter).Return(expectedOrgs, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orgs?type=gov&q=Test&page=1&pageSize=10", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]OrgsListPaginated
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Len(resp["data"].Orgs, 1)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleGetOrgById_Success() {
	// Arrange
	orgID := "org-123"
	expectedOrg := &domain.Account{ID: orgID, Name: "Specific Org", Role: domain.RoleOrg, OrganizationDetail: &domain.OrganizationDetail{}}
	s.mockUserUsecase.On("GetOrgById", mock.Anything, orgID).Return(expectedOrg, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/orgs/%s", orgID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]OrganizationResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedOrg.ID, resp["data"].ID)
	s.mockUserUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleUpdateOrgs_Success() {
	// Arrange
	orgID := "org-123"
	name := "Updated Name"
	reqBody := UpdateOrgRequest{Name: &name}
	updateMap := map[string]interface{}{"name": "Updated Name"}
	s.mockUserUsecase.On("UpdateOrgFields", mock.Anything, orgID, updateMap).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/orgs/%s", orgID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.mockUserUsecase.AssertExpectations(s.T())
}

// --- Search Tests (Sad Path) ---
func (s *UserControllerTestSuite) TestHandleSearch_InvalidQuery() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?q=passport&page=abc&limit=10", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "invalid request")
}

// --- Checklist Tests ---

func (s *UserControllerTestSuite) TestHandleGetChecklistById_Success() {
	// Arrange
	userProcedureID := "up-123"
	expectedResult := []*domain.Checklist{{ID: "cl-1", UserProcedureID: userProcedureID}}
	s.mockChecklistUsecase.On("GetChecklistByUserProcedureID", mock.Anything, userProcedureID).Return(expectedResult, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/checklists/%s", userProcedureID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string][]ChecklistResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Len(resp["message"], 1)
	s.mockChecklistUsecase.AssertExpectations(s.T())
}

func (s *UserControllerTestSuite) TestHandleUpdateChecklist_Success() {
	// Arrange
	checklistID := "cl-123"
	expectedResult := &domain.Checklist{ID: checklistID, IsChecked: true}
	s.mockChecklistUsecase.On("UpdateChecklist", mock.Anything, checklistID).Return(expectedResult, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/checklists/%s", checklistID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]ChecklistResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.True(resp["message"].IsChecked)
	s.mockChecklistUsecase.AssertExpectations(s.T())
}

// TestUserControllerTestSuite runs the entire suite
func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
