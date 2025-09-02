package usecase_test

import (
	"EthioGuide/domain"
	"EthioGuide/usecase"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mock for IProcedureRepository ---
type MockProcedureRepository struct {
	mock.Mock
}

func (m *MockProcedureRepository) Create(ctx context.Context, procedure *domain.Procedure) error {
	args := m.Called(ctx, procedure)
	return args.Error(0)
}

// --- Test Suite ---
type ProcedureUsecaseTestSuite struct {
	suite.Suite
	mockRepo      *MockProcedureRepository
	usecase       domain.IProcedureUsecase
	contextTimeout time.Duration
}

func (s *ProcedureUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockProcedureRepository)
	s.contextTimeout = 2 * time.Second
	s.usecase = usecase.NewProcedureUsecase(s.mockRepo, s.contextTimeout)
}

func TestProcedureUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureUsecaseTestSuite))
}

// --- Tests ---
func (s *ProcedureUsecaseTestSuite) TestCreateProcedure() {
	validProc := &domain.Procedure{
		Name:           "Test Procedure",
		GroupID:        "group123",
		OrganizationID: "org456",
		Content: domain.Content{
			Prerequisites: []string{"A", "B"},
			Steps:         []string{"Step1", "Step2"},
			Result:        []string{"Result1"},
		},
		Fees: domain.Fees{
			Label:    "FeeLabel",
			Currency: "ETB",
			Amount:   100.0,
		},
		ProcessingTime: domain.ProcessingTime{
			MinDays: 1,
			MaxDays: 5,
		},
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Procedure")).Return(nil).Once()

		err := s.usecase.CreateProcedure(context.Background(), validProc)
		s.NoError(err)
		s.mockRepo.AssertExpectations(s.T())
	})
}