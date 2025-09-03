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

type MockProcedureRepo struct {
	mock.Mock
}

func (m *MockProcedureRepo) GetByID(ctx context.Context, id string) (*domain.Procedure, error) {
	args := m.Called(ctx, id)
	var proc *domain.Procedure
	if args.Get(0) != nil {
		proc = args.Get(0).(*domain.Procedure)
	}
	return proc, args.Error(1)
}
func (m *MockProcedureRepo) Update(ctx context.Context, id string, procedure *domain.Procedure) error {
	args := m.Called(ctx, id, procedure)
	return args.Error(0)
}
func (m *MockProcedureRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProcedureRepo) Create(ctx context.Context, procedure *domain.Procedure) error {
	args := m.Called(ctx, procedure)
	return args.Error(0)
}

// --- Test Suite ---
type ProcedureUsecaseTestSuite struct {
	suite.Suite
	mockRepo       *MockProcedureRepo
	usecase        domain.IProcedureUsecase
	contextTimeout time.Duration
}

func (s *ProcedureUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockProcedureRepo)
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
		OrganizationID: "org456",
		Content: domain.ProcedureContent{
			Prerequisites: []string{"A", "B"},
			Steps:         []string{"Step1", "Step2"},
			Result:        "Result1",
		},
		Fees: domain.ProcedureFee{
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
