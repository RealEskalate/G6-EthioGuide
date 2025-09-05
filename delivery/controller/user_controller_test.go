package controller_test

// import (
// 	. "EthioGuide/delivery/controller"
// 	"EthioGuide/domain"
// 	"bytes"
// 	"context"
// 	"encoding/json"

// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// // --- Mocks & Placeholders ---

// // MockUserUsecase is a mock implementation of the IUserUsecase interface.
// type MockUserUsecase struct {
// 	mock.Mock
// }

// var _ domain.IUserUsecase = (*MockUserUsecase)(nil)

// func (m *MockUserUsecase) Register(ctx context.Context, account *domain.Account) error {
// 	args := m.Called(ctx, account)
// 	return args.Error(0)
// }

// func (m *MockUserUsecase) Login(ctx context.Context, identifier, password string) (*domain.Account, string, string, error) {
// 	args := m.Called(ctx, identifier, password)
// 	var acc *domain.Account
// 	if args.Get(0) != nil {
// 		acc = args.Get(0).(*domain.Account)
// 	}
// 	return acc, args.String(1), args.String(2), args.Error(3)
// }

// func (m *MockUserUsecase) RefreshTokenForWeb(ctx context.Context, refreshToken string) (string, error) {
// 	args := m.Called(ctx, refreshToken)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockUserUsecase) RefreshTokenForMobile(ctx context.Context, refreshToken string) (string, string, error) {
// 	args := m.Called(ctx, refreshToken)
// 	return args.String(0), args.String(1), args.Error(2)
// }

// // MockChecklistUsecase is a mock implementation of the IChecklistUsecase interface.
// type MockChecklistUsecase struct {
// 	mock.Mock
// }

// var _ domain.IChecklistUsecase = (*MockChecklistUsecase)(nil)

// func (m *MockChecklistUsecase) CreateChecklist(ctx context.Context, userID, procedureID string) (*domain.UserProcedure, error) {
// 	args := m.Called(ctx, userID, procedureID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*domain.UserProcedure), args.Error(1)
// }

// func (m *MockChecklistUsecase) GetProcedures(ctx context.Context, userID string) ([]*domain.UserProcedure, error) {
// 	args := m.Called(ctx, userID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]*domain.UserProcedure), args.Error(1)
// }

// func (m *MockChecklistUsecase) GetChecklistByUserProcedureID(ctx context.Context, userProcedureID string) ([]*domain.Checklist, error) {
// 	args := m.Called(ctx, userProcedureID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]*domain.Checklist), args.Error(1)
// }

// func (m *MockChecklistUsecase) UpdateChecklist(ctx context.Context, checklistID string) (*domain.Checklist, error) {
// 	args := m.Called(ctx, checklistID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*domain.Checklist), args.Error(1)
// }

// // --- Test Suite Definition ---

// type UserControllerTestSuite struct {
// 	suite.Suite
// 	router               *gin.Engine
// 	mockUserUsecase      *MockUserUsecase
// 	mockChecklistUsecase *MockChecklistUsecase
// 	controller           *UserController
// 	recorder             *httptest.ResponseRecorder
// 	refreshTokenTTL      time.Duration
// }

// func (s *UserControllerTestSuite) SetupSuite() {
// 	gin.SetMode(gin.TestMode)
// }

// func (s *UserControllerTestSuite) SetupTest() {
// 	s.recorder = httptest.NewRecorder()
// 	s.router = gin.Default()
// 	s.mockUserUsecase = new(MockUserUsecase)
// 	s.mockChecklistUsecase = new(MockChecklistUsecase)
// 	s.refreshTokenTTL = 15 * time.Minute

// 	s.controller = NewUserController(s.mockUserUsecase, s.mockChecklistUsecase, s.refreshTokenTTL)

// 	// Setup routes
// 	s.router.POST("/register", s.controller.Register)
// 	s.router.POST("/login", s.controller.Login)
// 	s.router.POST("/refresh", s.controller.HandleRefreshToken)

// 	checklistRoutes := s.router.Group("/checklists")
// 	checklistRoutes.Use(func(c *gin.Context) { // Mock middleware to set user_id
// 		c.Set("user_id", uuid.New().String())
// 		c.Next()
// 	})
// 	{
// 		checklistRoutes.POST("", s.controller.HandleCreateChecklist)
// 		checklistRoutes.GET("", s.controller.HandleGetProcedures)
// 		checklistRoutes.GET("/by-procedure", s.controller.HandleGetChecklistById)
// 		checklistRoutes.PATCH("/:id", s.controller.HandleUpdateChecklist)
// 	}
// }

// func TestUserControllerTestSuite(t *testing.T) {
// 	suite.Run(t, new(UserControllerTestSuite))
// }

// // --- Test Cases ---

// func (s *UserControllerTestSuite) TestRegister() {
// 	s.Run("Success", func() {
// 		s.SetupTest()
// 		reqBody := RegisterRequest{Name: "Test User", Username: "testuser", Email: "test@example.com", Password: "password123"}
// 		jsonBody, _ := json.Marshal(reqBody)

// 		s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusCreated, s.recorder.Code)
// 		s.mockUserUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Email Exists", func() {
// 		s.SetupTest()
// 		reqBody := RegisterRequest{Name: "Test User", Username: "testuser", Email: "test@example.com", Password: "password123"}
// 		jsonBody, _ := json.Marshal(reqBody)

// 		s.mockUserUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(domain.ErrEmailExists).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusConflict, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), domain.ErrEmailExists.Error())
// 		s.mockUserUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestLogin() {
// 	reqBody := LoginRequest{Identifier: "test@example.com", Password: "password"}
// 	jsonBody, _ := json.Marshal(reqBody)
// 	mockAccount := &domain.Account{Email: "test@example.com"}

// 	s.Run("Success - Web Client", func() {
// 		s.SetupTest()
// 		s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
// 			Return(mockAccount, "access_token_123", "refresh_token_abc", nil).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)
// 		cookie := s.recorder.Result().Cookies()[0]
// 		s.Equal("refresh_token", cookie.Name)
// 		s.Equal("refresh_token_abc", cookie.Value)

// 		var resp LoginResponse
// 		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
// 		s.Equal("access_token_123", resp.AccessToken)
// 		s.Empty(resp.RefreshToken)
// 		s.mockUserUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestHandleCreateChecklist() {
// 	procedureID := uuid.New().String()
// 	reqBody := CreateChecklistRequest{ProcedureID: procedureID}
// 	jsonBody, _ := json.Marshal(reqBody)
// 	expectedResponse := &domain.UserProcedure{ID: uuid.New().String()}

// 	s.Run("Success", func() {
// 		s.SetupTest()
// 		s.mockChecklistUsecase.On("CreateChecklist", mock.Anything, mock.AnythingOfType("string"), procedureID).
// 			Return(expectedResponse, nil).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/checklists", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusCreated, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Usecase Error", func() {
// 		s.SetupTest()
// 		usecaseErr := errors.New("failed to create")
// 		s.mockChecklistUsecase.On("CreateChecklist", mock.Anything, mock.AnythingOfType("string"), procedureID).
// 			Return(nil, usecaseErr).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/checklists", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusInternalServerError, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestHandleGetProcedures() {
// 	expectedResponse := []*domain.UserProcedure{{ID: uuid.New().String()}}

// 	s.Run("Success", func() {
// 		s.SetupTest()
// 		s.mockChecklistUsecase.On("GetProcedures", mock.Anything, mock.AnythingOfType("string")).
// 			Return(expectedResponse, nil).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/checklists", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Usecase Error", func() {
// 		s.SetupTest()
// 		usecaseErr := errors.New("failed to get procedures")
// 		s.mockChecklistUsecase.On("GetProcedures", mock.Anything, mock.AnythingOfType("string")).
// 			Return(nil, usecaseErr).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/checklists", nil)
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusInternalServerError, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestHandleGetChecklistById() {
// 	userProcedureID := uuid.New().String()
// 	reqBody := GetChecklistByID{UserProcedureID: userProcedureID}
// 	jsonBody, _ := json.Marshal(reqBody)
// 	expectedResponse := []*domain.Checklist{{ID: uuid.New().String()}}

// 	s.Run("Success", func() {
// 		s.SetupTest()
// 		s.mockChecklistUsecase.On("GetChecklistByUserProcedureID", mock.Anything, userProcedureID).
// 			Return(expectedResponse, nil).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/checklists/by-procedure", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Usecase Error", func() {
// 		s.SetupTest()
// 		usecaseErr := errors.New("failed to get checklists")
// 		s.mockChecklistUsecase.On("GetChecklistByUserProcedureID", mock.Anything, userProcedureID).
// 			Return(nil, usecaseErr).Once()

// 		req, _ := http.NewRequest(http.MethodGet, "/checklists/by-procedure", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusInternalServerError, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestLoginMobile() {
// 	reqBody := LoginRequest{Identifier: "test@example.com", Password: "password"}
// 	jsonBody, _ := json.Marshal(reqBody)
// 	mockAccount := &domain.Account{Email: "test@example.com"}

// 	s.Run("Success - Mobile Client", func() {
// 		s.SetupTest()
// 		s.mockUserUsecase.On("Login", mock.Anything, reqBody.Identifier, reqBody.Password).
// 			Return(mockAccount, "access_token_123", "refresh_token_abc", nil).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("X-Client-Type", "mobile")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)

// 		var resp LoginResponse
// 		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
// 		s.Equal("access_token_123", resp.AccessToken)
// 		s.Equal("refresh_token_abc", resp.RefreshToken)
// 		s.mockUserUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestHandleRefreshTokenMobile() {
// 	s.Run("Success - Mobile Client", func() {
// 		s.SetupTest()
// 		s.mockUserUsecase.On("RefreshTokenForMobile", mock.Anything, "refresh_token_abc").
// 			Return("new_access_token", "new_refresh_token", nil).Once()

// 		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
// 		req.Header.Set("Authorization", "Bearer refresh_token_abc")
// 		req.Header.Set("X-Client-Type", "mobile")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)
// 		var resp map[string]string
// 		json.Unmarshal(s.recorder.Body.Bytes(), &resp)
// 		s.Equal("new_access_token", resp["access_token"])
// 		s.Equal("new_refresh_token", resp["refresh_token"])
// 		s.mockUserUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestHandleUpdateChecklist() {
// 	checklistID := uuid.New().String()
// 	reqBody := UpdateChecklistRequest{ChecklistID: checklistID}
// 	jsonBody, _ := json.Marshal(reqBody)
// 	expectedResponse := &domain.Checklist{ID: checklistID, IsChecked: true}

// 	s.Run("Success", func() {
// 		s.SetupTest()

// 		// The controller now gets the ID from the URL param
// 		s.mockChecklistUsecase.On("UpdateChecklist", mock.Anything, checklistID).
// 			Return(expectedResponse, nil).Once()

// 		// The request body is no longer needed to carry the ID
// 		req, _ := http.NewRequest(http.MethodPatch, "/checklists/"+checklistID, bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusOK, s.recorder.Code)
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})

// 	s.Run("Failure - Usecase Error", func() {
// 		s.SetupTest()
// 		usecaseErr := errors.New("failed to update")
// 		s.mockChecklistUsecase.On("UpdateChecklist", mock.Anything, checklistID).
// 			Return(nil, usecaseErr).Once()

// 		req, _ := http.NewRequest(http.MethodPatch, "/checklists/"+checklistID, bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusInternalServerError, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), "An unexpected internal error occurred. Please try again later.")
// 		s.mockChecklistUsecase.AssertExpectations(s.T())
// 	})
// }

// func (s *UserControllerTestSuite) TestHandleUpdateChecklist_BadRequest() {
// 	s.Run("Failure - Invalid Body", func() {
// 		s.SetupTest()
// 		checklistID := uuid.New().String()

// 		// The controller's UpdateChecklist now expects a JSON body,
// 		// so sending a malformed one should result in a bad request.
// 		req, _ := http.NewRequest(http.MethodPatch, "/checklists/"+checklistID, bytes.NewBufferString("not-json"))
// 		req.Header.Set("Content-Type", "application/json")
// 		s.router.ServeHTTP(s.recorder, req)

// 		s.Equal(http.StatusBadRequest, s.recorder.Code)
// 		s.Contains(s.recorder.Body.String(), "Invalid request body")
// 	})
// }
