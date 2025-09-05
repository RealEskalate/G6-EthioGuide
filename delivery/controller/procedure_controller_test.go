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

// --- Mock for IProcedureUsecase ---
type MockProcedureUsecase struct {
	mock.Mock
}

func (m *MockProcedureUsecase) CreateProcedure(ctx context.Context, procedure *domain.Procedure) error {
	args := m.Called(ctx, procedure)
	if procedure != nil && args.Error(0) == nil {
		procedure.ID = "new-proc-id-123"
	}
	return args.Error(0)
}

func (m *MockProcedureUsecase) GetProcedureByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Procedure), args.Error(1)
}

func (m *MockProcedureUsecase) UpdateProcedure(ctx context.Context, id string, procedure *domain.Procedure) error {
	args := m.Called(ctx, id, procedure)
	return args.Error(0)
}

func (m *MockProcedureUsecase) DeleteProcedure(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Test Suite ---
type ProcedureControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockProcedureUsecase
	controller  *controller.ProcedureController
	recorder    *httptest.ResponseRecorder
}

func (s *ProcedureControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *ProcedureControllerTestSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUsecase = new(MockProcedureUsecase)
	s.controller = controller.NewProcedureController(s.mockUsecase)
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "org456")
		c.Set("userRole", domain.RoleOrg)
		c.Next()
	}
	procRoutes := s.router.Group("/procedures")
	{
		procRoutes.POST("", authMiddleware, s.controller.CreateProcedure)
		procRoutes.GET("/:id", s.controller.GetProcedureByID)
		procRoutes.PUT("/:id", authMiddleware, s.controller.UpdateProcedure)
		procRoutes.DELETE("/:id", authMiddleware, s.controller.DeleteProcedure)
	}
}

func TestProcedureControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureControllerTestSuite))
}

// --- Tests ---
func (s *ProcedureControllerTestSuite) TestCreateProcedure() {
	validReq := controller.ProcedureCreateRequest{
		Name:          "Test Procedure",
		GroupID:       "group123",
		Prerequisites: []string{"A", "B"},
		Steps:         map[int]string{1: "Step1", 2: "Step2"},
		Result:        "Result1",
		Label:         "FeeLabel",
		Currency:      "ETB",
		Amount:        100.0,
		MinDays:       1,
		MaxDays:       5,
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUsecase.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*domain.Procedure")).Return(nil).Once()

		body, _ := json.Marshal(validReq)
		req, _ := http.NewRequest(http.MethodPost, "/procedures", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusCreated, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Procedure created successfully")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid Body", func() {
		s.SetupTest()
		req, _ := http.NewRequest(http.MethodPost, "/procedures", bytes.NewBuffer([]byte(`{invalid json`)))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusBadRequest, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Invalid request body")
	})

	s.Run("Failure - Usecase Error", func() {
		s.SetupTest()
		s.mockUsecase.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*domain.Procedure")).Return(domain.ErrConflict).Once()

		body, _ := json.Marshal(validReq)
		req, _ := http.NewRequest(http.MethodPost, "/procedures", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusConflict, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), domain.ErrConflict.Error())
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *ProcedureControllerTestSuite) TestGetProcedureByID() {
	procID := "proc-123"
	mockProc := &domain.Procedure{ID: procID, Name: "Found Procedure"}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUsecase.On("GetProcedureByID", mock.Anything, procID).Return(mockProc, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/procedures/"+procID, nil)
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusOK, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Found Procedure")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Not Found", func() {
		s.SetupTest()
		s.mockUsecase.On("GetProcedureByID", mock.Anything, procID).Return(nil, domain.ErrNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/procedures/"+procID, nil)
		s.router.ServeHTTP(s.recorder, req)

		// Assuming your HandleError correctly maps domain.ErrNotFound to 404
		s.Equal(http.StatusNotFound, s.recorder.Code)
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *ProcedureControllerTestSuite) TestUpdateProcedure() {
	procID := "proc-123"
	updateData := gin.H{"name": "Updated Name"} // Use gin.H for simple body creation

	s.Run("Success", func() {
		s.SetupTest()
		// Expect a call to the usecase and return no error
		s.mockUsecase.On("UpdateProcedure", mock.Anything, procID, mock.AnythingOfType("*domain.Procedure")).Return(nil).Once()

		body, _ := json.Marshal(updateData)
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusOK, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Updated Procedure Successfully")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Permission Denied", func() {
		s.SetupTest()
		// Mock the usecase to return a permission error
		s.mockUsecase.On("UpdateProcedure", mock.Anything, procID, mock.AnythingOfType("*domain.Procedure")).Return(domain.ErrPermissionDenied).Once()

		body, _ := json.Marshal(updateData)
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		s.router.ServeHTTP(s.recorder, req)

		// Assuming HandleError maps this to 403 Forbidden
		s.Equal(http.StatusForbidden, s.recorder.Code)
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *ProcedureControllerTestSuite) TestDeleteProcedure() {
	procID := "proc-123"

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUsecase.On("DeleteProcedure", mock.Anything, procID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/procedures/"+procID, nil)
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusNoContent, s.recorder.Code)
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Not Found", func() {
		s.SetupTest()
		s.mockUsecase.On("DeleteProcedure", mock.Anything, procID).Return(domain.ErrNotFound).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/procedures/"+procID, nil)
		s.router.ServeHTTP(s.recorder, req)

		s.Equal(http.StatusNotFound, s.recorder.Code)
		s.mockUsecase.AssertExpectations(s.T())
	})
}
