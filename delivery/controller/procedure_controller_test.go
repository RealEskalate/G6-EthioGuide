package controller_test

import (
	. "EthioGuide/delivery/controller"
	"EthioGuide/domain"
	"bytes"
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

// --- Mocks & Placeholders ---

// MockProcedureUsecase is a mock implementation of the IProcedureUseCase interface.
type MockProcedureUsecase struct {
	mock.Mock
}

// Ensure MockProcedureUsecase implements IProcedureUseCase at compile time.
var _ domain.IProcedureUseCase = (*MockProcedureUsecase)(nil)

func (m *MockProcedureUsecase) GetProcedureByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	var procedure *domain.Procedure
	if args.Get(0) != nil {
		procedure = args.Get(0).(*domain.Procedure)
	}
	return procedure, args.Error(1)
}

func (m *MockProcedureUsecase) UpdateProcedure(ctx context.Context, id string, procedure *domain.Procedure) error {
	args := m.Called(ctx, id, procedure)
	return args.Error(0)
}

func (m *MockProcedureUsecase) DeleteProcedure(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Test Suite Definition ---

type ProcedureControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockProcedureUsecase
	controller  *ProcedureController
	recorder    *httptest.ResponseRecorder
}

func (s *ProcedureControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *ProcedureControllerTestSuite) SetupTest() {
	s.recorder = httptest.NewRecorder()
	s.router = gin.Default()
	s.mockUsecase = new(MockProcedureUsecase)

	// Create a new controller for each test, injecting the mock.
	s.controller = NewProcedureController(s.mockUsecase)

	// Setup routes
	s.router.GET("/procedures/:id", s.controller.GetProcedureByID)
	s.router.PUT("/procedures/:id", s.controller.UpdateProcedure)
}

func TestProcedureControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureControllerTestSuite))
}

// --- Test Cases ---

func (s *ProcedureControllerTestSuite) TestGetProcedureByID() {
	s.Run("Success", func() {
		// Arrange
		procedureID := "123"
		mockProcedure := &domain.Procedure{
			ID:   procedureID,
			Name: "Test Procedure",
		}

		s.mockUsecase.On("GetProcedureByID", mock.Anything, procedureID).
			Return(mockProcedure, nil).Once()

		// Act
		req, _ := http.NewRequest(http.MethodGet, "/procedures/"+procedureID, nil)
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)

		var response ProcedureDTO
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		s.Equal(procedureID, response.ID)
		s.Equal("Test Procedure", response.Name)

		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Procedure Not Found", func() {
		// Arrange
		s.SetupTest()
		procedureID := "non-existent-id"
		notFoundError := errors.New("procedure not found")

		s.mockUsecase.On("GetProcedureByID", mock.Anything, procedureID).
			Return(nil, notFoundError).Once()

		// Act
		req, _ := http.NewRequest(http.MethodGet, "/procedures/"+procedureID, nil)
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusNotFound, s.recorder.Code)
		
		var response gin.H
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		s.Equal(notFoundError.Error(), response["error"])
		
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Database Error", func() {
		// Arrange
		s.SetupTest()
		procedureID := "123"
		dbError := errors.New("database connection failed")

		s.mockUsecase.On("GetProcedureByID", mock.Anything, procedureID).
			Return(nil, dbError).Once()

		// Act
		req, _ := http.NewRequest(http.MethodGet, "/procedures/"+procedureID, nil)
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusNotFound, s.recorder.Code)
		
		var response gin.H
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		s.Equal(dbError.Error(), response["error"])
		
		s.mockUsecase.AssertExpectations(s.T())
	})
}

func (s *ProcedureControllerTestSuite) TestUpdateProcedure() {
	s.Run("Success", func() {
		// Arrange
		procedureID := "123"
		updateDTO := ProcedureDTO{
			ID:   procedureID,
			Name: "Updated Procedure Name",
		}

		s.mockUsecase.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*domain.Procedure")).
			Return(nil).Once()

		jsonBody, _ := json.Marshal(updateDTO)

		// Act
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procedureID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusOK, s.recorder.Code)
		s.Contains(s.recorder.Body.String(), "Updated Procedure Successfully")
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Invalid JSON", func() {
		// Arrange
		s.SetupTest()
		procedureID := "123"
		invalidJSON := `{"name": "test", "description":}`

		// Act
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procedureID, bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusBadRequest, s.recorder.Code)
		
		var response gin.H
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		// Check that we get an error field in the response
		s.Contains(response, "error")
		// The specific error message might vary, so we just check that it exists
		s.NotEmpty(response["error"])
		
		s.mockUsecase.AssertNotCalled(s.T(), "UpdateProcedure")
	})

	s.Run("Failure - Procedure Not Found", func() {
		// Arrange
		s.SetupTest()
		procedureID := "non-existent-id"
		updateDTO := ProcedureDTO{
			Name: "Updated Name",
		}
		notFoundError := errors.New("procedure not found")

		s.mockUsecase.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*domain.Procedure")).
			Return(notFoundError).Once()

		jsonBody, _ := json.Marshal(updateDTO)

		// Act
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procedureID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusNotFound, s.recorder.Code)
		
		var response gin.H
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		s.Equal(notFoundError.Error(), response["error"])
		
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Validation Error", func() {
		// Arrange
		s.SetupTest()
		procedureID := "123"
		updateDTO := ProcedureDTO{
			Name: "", // Empty name might trigger validation error
		}
		validationError := errors.New("name is required")

		s.mockUsecase.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*domain.Procedure")).
			Return(validationError).Once()

		jsonBody, _ := json.Marshal(updateDTO)

		// Act
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procedureID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusNotFound, s.recorder.Code)
		
		var response gin.H
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		s.Equal(validationError.Error(), response["error"])
		
		s.mockUsecase.AssertExpectations(s.T())
	})

	s.Run("Failure - Database Error", func() {
		// Arrange
		s.SetupTest()
		procedureID := "123"
		updateDTO := ProcedureDTO{
			Name: "Valid Name",
		}
		dbError := errors.New("database update failed")

		s.mockUsecase.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*domain.Procedure")).
			Return(dbError).Once()

		jsonBody, _ := json.Marshal(updateDTO)

		// Act
		req, _ := http.NewRequest(http.MethodPut, "/procedures/"+procedureID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(s.recorder, req)

		// Assert
		s.Equal(http.StatusNotFound, s.recorder.Code)
		
		var response gin.H
		err := json.Unmarshal(s.recorder.Body.Bytes(), &response)
		s.NoError(err)
		s.Equal(dbError.Error(), response["error"])
		
		s.mockUsecase.AssertExpectations(s.T())
	})
}