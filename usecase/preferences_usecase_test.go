package usecase_test

import (
    "context"
    "testing"

    "EthioGuide/domain"
    "EthioGuide/usecase"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

// --- Mock for IPreferencesRepository ---
type MockPreferencesRepository struct {
    mock.Mock
}

func (m *MockPreferencesRepository) Create(ctx context.Context, preferences *domain.Preferences) error {
    args := m.Called(ctx, preferences)
    return args.Error(0)
}
func (m *MockPreferencesRepository) GetByUserID(ctx context.Context, userID string) (*domain.Preferences, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) != nil {
        return args.Get(0).(*domain.Preferences), args.Error(1)
    }
    return nil, args.Error(1)
}
func (m *MockPreferencesRepository) UpdateByUserID(ctx context.Context, userID string, preferences *domain.Preferences) error {
    args := m.Called(ctx, userID, preferences)
    return args.Error(0)
}

// --- Test Suite ---
type PreferencesUsecaseSuite struct {
    suite.Suite
    mockRepo *MockPreferencesRepository
    usecase  domain.IPreferencesUsecase
}

func (s *PreferencesUsecaseSuite) SetupTest() {
    s.mockRepo = new(MockPreferencesRepository)
    s.usecase = usecase.NewPreferencesUsecase(s.mockRepo)
}

func TestPreferencesUsecaseSuite(t *testing.T) {
    suite.Run(t, new(PreferencesUsecaseSuite))
}

func (s *PreferencesUsecaseSuite) TestCreateUserPreferences() {
    ctx := context.Background()
    expected := &domain.Preferences{
        UserID:            "user123",
        PreferredLang:     domain.English,
        PushNotification:  false,
        EmailNotification: false,
    }
    s.mockRepo.On("Create", ctx, mock.MatchedBy(func(p *domain.Preferences) bool {
        return p.UserID == expected.UserID &&
            p.PreferredLang == expected.PreferredLang &&
            !p.PushNotification && !p.EmailNotification
    })).Return(nil).Once()

    err := s.usecase.CreateUserPreferences(ctx, "user123")
    s.NoError(err)
    s.mockRepo.AssertExpectations(s.T())
}

func (s *PreferencesUsecaseSuite) TestGetUserPreferences() {
    ctx := context.Background()
    expected := &domain.Preferences{
        UserID:            "user123",
        PreferredLang:     domain.English,
        PushNotification:  true,
        EmailNotification: true,
    }
    s.mockRepo.On("GetByUserID", ctx, "user123").Return(expected, nil).Once()

    got, err := s.usecase.GetUserPreferences(ctx, "user123")
    s.NoError(err)
    s.Equal(expected, got)
    s.mockRepo.AssertExpectations(s.T())
}

func (s *PreferencesUsecaseSuite) TestUpdateUserPreferences() {
    ctx := context.Background()
    pref := &domain.Preferences{
        UserID:            "user123",
        PreferredLang:     domain.English,
        PushNotification:  true,
        EmailNotification: false,
    }
    s.mockRepo.On("UpdateByUserID", ctx, "user123", pref).Return(nil).Once()

    err := s.usecase.UpdateUserPreferences(ctx, pref)
    s.NoError(err)
    s.mockRepo.AssertExpectations(s.T())
}