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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockProcedureUsecase is a mock for the IProcedureUsecase interface
type MockProcedureUsecase struct {
	mock.Mock
}

func (m *MockProcedureUsecase) CreateProcedure(ctx context.Context, procedure *domain.Procedure, userId string, userRole domain.Role) error {
	args := m.Called(ctx, procedure, userId, userRole)
	return args.Error(0)
}

func (m *MockProcedureUsecase) GetProcedureByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Procedure), args.Error(1)
}

func (m *MockProcedureUsecase) UpdateProcedure(ctx context.Context, id string, procedure *domain.Procedure, userId string, userRole domain.Role) error {
	args := m.Called(ctx, id, procedure, userId, userRole)
	return args.Error(0)
}

func (m *MockProcedureUsecase) DeleteProcedure(ctx context.Context, id string, userId string, userRole domain.Role) error {
	args := m.Called(ctx, id, userId, userRole)
	return args.Error(0)
}

func (m *MockProcedureUsecase) SearchAndFilter(ctx context.Context, options domain.ProcedureSearchFilterOptions) ([]*domain.Procedure, int64, error) {
	args := m.Called(ctx, options)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Procedure), args.Get(1).(int64), args.Error(2)
}

// ProcedureControllerTestSuite is the test suite for ProcedureController
type ProcedureControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockProcedureUsecase
	controller  *ProcedureController
}

// SetupTest is run before each test in the suite
func (s *ProcedureControllerTestSuite) SetupTest() {
	s.router = gin.Default()
	s.mockUsecase = new(MockProcedureUsecase)
	s.controller = NewProcedureController(s.mockUsecase)

	authMiddleware := func(c *gin.Context) {
		c.Set("userID", "test-org-id")
		c.Set("userRole", domain.RoleOrg)
		c.Next()
	}

	procedures := s.router.Group("/procedures")
	{
		procedures.POST("", authMiddleware, s.controller.CreateProcedure)
		procedures.GET("", s.controller.SearchAndFilter)
		procedures.GET("/:id", s.controller.GetProcedureByID)
		procedures.PATCH("/:id", authMiddleware, s.controller.UpdateProcedure)
		procedures.DELETE("/:id", authMiddleware, s.controller.DeleteProcedure)
	}
}

func (s *ProcedureControllerTestSuite) TestCreateProcedure_Success() {
	// Arrange
	reqBody := ProcedureCreateRequest{Name: "New Procedure", GroupID: "group-1"}
	s.mockUsecase.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*domain.Procedure"), "test-org-id", domain.RoleOrg).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/procedures", toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Procedure created successfully")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestGetProcedureByID_Success() {
	// Arrange
	procID := "proc-123"
	expectedProc := &domain.Procedure{ID: procID, Name: "Test Procedure"}
	s.mockUsecase.On("GetProcedureByID", mock.Anything, procID).Return(expectedProc, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/procedures/%s", procID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp domain.Procedure
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(expectedProc.ID, resp.ID)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestGetProcedureByID_NotFound() {
	// Arrange
	procID := "proc-404"
	s.mockUsecase.On("GetProcedureByID", mock.Anything, procID).Return(nil, domain.ErrProcedureNotFound).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/procedures/%s", procID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), domain.ErrProcedureNotFound.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestUpdateProcedure_Success() {
	// Arrange
	procID := "proc-to-update"
	reqBody := domain.Procedure{Name: "Updated Name"}
	s.mockUsecase.On("UpdateProcedure", mock.Anything, procID, &reqBody, "test-org-id", domain.RoleOrg).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/procedures/%s", procID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Updated Procedure Successfully.")
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestUpdateProcedure_PermissionDenied() {
	// Arrange
	procID := "proc-to-update"
	reqBody := domain.Procedure{Name: "Updated Name"}
	s.mockUsecase.On("UpdateProcedure", mock.Anything, procID, &reqBody, "test-org-id", domain.RoleOrg).Return(domain.ErrPermissionDenied).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/procedures/%s", procID), toJSON(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusForbidden, w.Code)
	s.Contains(w.Body.String(), domain.ErrPermissionDenied.Error())
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestDeleteProcedure_Success() {
	// Arrange
	procID := "proc-to-delete"
	s.mockUsecase.On("DeleteProcedure", mock.Anything, procID, "test-org-id", domain.RoleOrg).Return(nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/procedures/%s", procID), nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusNoContent, w.Code)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestSearchAndFilter_Success() {
	// Arrange
	expectedProcs := []*domain.Procedure{{ID: "proc-1", Name: "Search Result"}}
	var total int64 = 1
	name := "Test"
	options := domain.ProcedureSearchFilterOptions{
		Page:        1,
		Limit:       10,
		Name:        &name,
		SortOrder:   domain.SortDesc,
		GlobalLogic: domain.GlobalLogicAND,
	}

	s.mockUsecase.On("SearchAndFilter", mock.Anything, options).Return(expectedProcs, total, nil).Once()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/procedures?name=Test", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	var resp PaginatedProcedureResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal(total, resp.Pagination.Total)
	s.Len(resp.Data, 1)
	s.mockUsecase.AssertExpectations(s.T())
}

func (s *ProcedureControllerTestSuite) TestSearchAndFilter_InvalidParam() {
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/procedures?minFee=abc", nil)
	s.router.ServeHTTP(w, req)

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid 'minFee' parameter")
}

// TestProcedureControllerTestSuite runs the entire suite
func TestProcedureControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureControllerTestSuite))
}
