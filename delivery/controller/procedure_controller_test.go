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
	s.router.POST("/procedures", func(ctx *gin.Context) {
		ctx.Set("userID", "org456")
	}, s.controller.CreateProcedure)
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
		Steps:         []string{"Step1", "Step2"},
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
