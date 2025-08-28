package usecase_test

import (
	"EthioGuide/domain"
	. "EthioGuide/usecase" // Use dot import for convenience in test files
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// --- Mocks & Placeholders ---

// MockAccountRepository mocks IAccountRepository
type MockAccountRepository struct{ mock.Mock }

func (m *MockAccountRepository) Create(ctx context.Context, user *domain.Account) error {
	return m.Called(ctx, user).Error(0)
}
func (m *MockAccountRepository) GetById(ctx context.Context, id string) (*domain.Account, error) {
	args := m.Called(ctx, id)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.Error(1)
}
func (m *MockAccountRepository) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	args := m.Called(ctx, email)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.Error(1)
}
func (m *MockAccountRepository) GetByUsername(ctx context.Context, username string) (*domain.Account, error) {
	args := m.Called(ctx, username)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.Error(1)
}
func (m *MockAccountRepository) GetByPhoneNumber(ctx context.Context, phone string) (*domain.Account, error) {
	args := m.Called(ctx, phone)
	var acc *domain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*domain.Account)
	}
	return acc, args.Error(1)
}
func (m *MockAccountRepository) UpdatePassword(ctx context.Context, accountID, newPassword string) error {
	return m.Called(ctx, accountID, newPassword).Error(0)
}

// MockTokenRepository mocks ITokenRepository
type MockTokenRepository struct{ mock.Mock }

func (m *MockTokenRepository) CreateToken(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	args := m.Called(ctx, token)
	var tok *domain.Token
	if args.Get(0) != nil {
		tok = args.Get(0).(*domain.Token)
	}
	return tok, args.Error(1)
}
func (m *MockTokenRepository) GetToken(ctx context.Context, tokentype, tokenID string) (string, error) {
	args := m.Called(ctx, tokentype, tokenID)
	return args.String(0), args.Error(1)
}
func (m *MockTokenRepository) DeleteToken(ctx context.Context, tokentype, tokenID string) error {
	return m.Called(ctx, tokentype, tokenID).Error(0)
}

// MockPasswordService mocks IPasswordService
type MockPasswordService struct{ mock.Mock }

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
func (m *MockPasswordService) ComparePassword(hashedPassword, password string) error {
	return m.Called(hashedPassword, password).Error(0)
}

// MockJWTService mocks IJWTService
type MockJWTService struct{ mock.Mock }

func (m *MockJWTService) GenerateAccessToken(userID string, role domain.Role) (string, *domain.JWTClaims, error) {
	args := m.Called(userID, role)
	var claims *domain.JWTClaims
	if args.Get(1) != nil {
		claims = args.Get(1).(*domain.JWTClaims)
	}
	return args.String(0), claims, args.Error(2)
}
func (m *MockJWTService) GenerateRefreshToken(userID string) (string, *domain.JWTClaims, error) {
	args := m.Called(userID)
	var claims *domain.JWTClaims
	if args.Get(1) != nil {
		claims = args.Get(1).(*domain.JWTClaims)
	}
	return args.String(0), claims, args.Error(2)
}
func (m *MockJWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)
	var claims *domain.JWTClaims
	if args.Get(0) != nil {
		claims = args.Get(0).(*domain.JWTClaims)
	}
	return claims, args.Error(1)
}

func (m *MockJWTService) ParseExpiredToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)
	var claims *domain.JWTClaims
	if args.Get(0) != nil {
		claims = args.Get(0).(*domain.JWTClaims)
	}
	return claims, args.Error(1)
}

func (m *MockJWTService) GetRefreshTokenExpiry() time.Duration {
	args := m.Called()
	// The return value is retrieved by its index.
	// We need to perform a type assertion to the expected type.
	return args.Get(0).(time.Duration)
}

// Setup domain errors for testing
func setupDomainErrors() {
	domain.ErrNotFound = errors.New("not found")
	domain.ErrEmailExists = errors.New("email already exists")
	domain.ErrUsernameExists = errors.New("username already exists")
	domain.ErrAuthenticationFailed = errors.New("authentication failed")
	domain.ErrAccountNotActive = errors.New("account not active")
	domain.ErrUserNotFound = errors.New("user not found")
	domain.ErrPasswordTooShort = errors.New("password too short")
}

// --- Test Suite Definition ---
type UserUsecaseTestSuite struct {
	suite.Suite
	mockUserRepo  *MockAccountRepository
	mockTokenRepo *MockTokenRepository
	mockPassSvc   *MockPasswordService
	mockJwtSvc    *MockJWTService
	usecase       domain.IUserUsecase
}

func (s *UserUsecaseTestSuite) SetupSuite() {
	setupDomainErrors()
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.mockUserRepo = new(MockAccountRepository)
	s.mockTokenRepo = new(MockTokenRepository)
	s.mockPassSvc = new(MockPasswordService)
	s.mockJwtSvc = new(MockJWTService)

	s.usecase = NewUserUsecase(
		s.mockUserRepo,
		s.mockTokenRepo,
		s.mockPassSvc,
		s.mockJwtSvc,
		5*time.Second, // Default timeout for tests
	)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

// --- Register Tests ---
func (s *UserUsecaseTestSuite) TestRegister() {
	user := &domain.Account{
		Email:        "test@example.com",
		PasswordHash: "password123",
		UserDetail:   &domain.UserDetail{Username: "testuser"},
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(nil, domain.ErrNotFound).Once()
		s.mockUserRepo.On("GetByUsername", mock.Anything, user.UserDetail.Username).Return(nil, domain.ErrNotFound).Once()
		s.mockPassSvc.On("HashPassword", user.PasswordHash).Return("hashed_password", nil).Once()
		s.mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		err := s.usecase.Register(context.Background(), user)

		s.NoError(err)
		s.mockUserRepo.AssertExpectations(s.T())
		s.mockPassSvc.AssertExpectations(s.T())
	})

	s.Run("Failure - Email Exists", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(&domain.Account{}, nil).Once()

		err := s.usecase.Register(context.Background(), user)

		s.ErrorIs(err, domain.ErrEmailExists)
		s.mockUserRepo.AssertCalled(s.T(), "GetByEmail", mock.Anything, user.Email)
		s.mockUserRepo.AssertNotCalled(s.T(), "GetByUsername") // Should fail fast
	})

	s.Run("Failure - Username Exists", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetByEmail", mock.Anything, user.Email).Return(nil, domain.ErrNotFound).Once()
		s.mockUserRepo.On("GetByUsername", mock.Anything, user.UserDetail.Username).Return(&domain.Account{}, nil).Once()

		err := s.usecase.Register(context.Background(), user)

		s.ErrorIs(err, domain.ErrUsernameExists)
		s.mockUserRepo.AssertExpectations(s.T())
		s.mockPassSvc.AssertNotCalled(s.T(), "HashPassword")
	})
}

// --- Login Tests ---
func (s *UserUsecaseTestSuite) TestLogin() {
	identifier := "test@example.com"
	password := "password123"
	account := &domain.Account{
		ID:           "user-123",
		PasswordHash: "hashed_password",
		Role:         domain.RoleUser,
		UserDetail:   &domain.UserDetail{IsVerified: true},
	}
	refreshClaims := &domain.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "token-abc",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetByEmail", mock.Anything, identifier).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", account.PasswordHash, password).Return(nil).Once()
		s.mockJwtSvc.On("GenerateAccessToken", account.ID, account.Role).Return("access_token", &domain.JWTClaims{}, nil).Once()
		s.mockJwtSvc.On("GenerateRefreshToken", account.ID).Return("refresh_token", refreshClaims, nil).Once()
		s.mockTokenRepo.On("CreateToken", mock.Anything, mock.AnythingOfType("*domain.Token")).Return(nil, nil).Once()

		_, _, _, err := s.usecase.Login(context.Background(), identifier, password)

		s.NoError(err)
		s.mockUserRepo.AssertExpectations(s.T())
		s.mockPassSvc.AssertExpectations(s.T())
		s.mockJwtSvc.AssertExpectations(s.T())
		s.mockTokenRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - User Not Found", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetByEmail", mock.Anything, identifier).Return(nil, domain.ErrNotFound).Once()

		_, _, _, err := s.usecase.Login(context.Background(), identifier, password)

		s.ErrorIs(err, domain.ErrAuthenticationFailed)
		s.mockPassSvc.AssertNotCalled(s.T(), "ComparePassword")
	})

	s.Run("Failure - Incorrect Password", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetByEmail", mock.Anything, identifier).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", account.PasswordHash, password).Return(errors.New("password mismatch")).Once()

		_, _, _, err := s.usecase.Login(context.Background(), identifier, password)

		s.ErrorIs(err, domain.ErrAuthenticationFailed)
		s.mockJwtSvc.AssertNotCalled(s.T(), "GenerateAccessToken")
	})

	s.Run("Failure - Account Not Verified", func() {
		s.SetupTest()
		unverifiedAccount := *account
		unverifiedAccount.UserDetail.IsVerified = false
		s.mockUserRepo.On("GetByEmail", mock.Anything, identifier).Return(&unverifiedAccount, nil).Once()

		_, _, _, err := s.usecase.Login(context.Background(), identifier, password)

		s.ErrorIs(err, domain.ErrAccountNotActive)
		s.mockPassSvc.AssertNotCalled(s.T(), "ComparePassword")
	})
}

// --- Refresh Token Tests ---
func (s *UserUsecaseTestSuite) TestRefreshToken() {
	refreshToken := "valid_refresh_token"
	claims := &domain.JWTClaims{
		UserID:           "user-123",
		Role:             domain.RoleUser,
		RegisteredClaims: jwt.RegisteredClaims{ID: "token-abc"},
	}

	s.Run("Success - Web", func() {
		s.SetupTest()
		s.mockJwtSvc.On("ValidateToken", refreshToken).Return(claims, nil).Once()
		s.mockTokenRepo.On("GetToken", mock.Anything, string(domain.RefreshToken), claims.ID).Return("some_token", nil).Once()
		s.mockJwtSvc.On("GenerateAccessToken", claims.UserID, claims.Role).Return("new_access_token", &domain.JWTClaims{}, nil).Once()

		_, err := s.usecase.RefreshTokenForWeb(context.Background(), refreshToken)

		s.NoError(err)
		s.mockJwtSvc.AssertExpectations(s.T())
		s.mockTokenRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Web Token Not In DB", func() {
		s.SetupTest()
		s.mockJwtSvc.On("ValidateToken", refreshToken).Return(claims, nil).Once()
		s.mockTokenRepo.On("GetToken", mock.Anything, string(domain.RefreshToken), claims.ID).Return("", domain.ErrNotFound).Once()

		_, err := s.usecase.RefreshTokenForWeb(context.Background(), refreshToken)

		s.ErrorIs(err, domain.ErrAuthenticationFailed)
		s.mockJwtSvc.AssertNotCalled(s.T(), "GenerateAccessToken")
	})

	s.Run("Success - Mobile (Rotation)", func() {
		s.SetupTest()
		newRefreshTokenClaims := &domain.JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        "token-xyz",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}

		s.mockJwtSvc.On("ValidateToken", refreshToken).Return(claims, nil).Once()
		s.mockTokenRepo.On("DeleteToken", mock.Anything, string(domain.RefreshToken), claims.ID).Return(nil).Once()
		s.mockJwtSvc.On("GenerateAccessToken", claims.UserID, claims.Role).Return("new_access_token", &domain.JWTClaims{}, nil).Once()
		s.mockJwtSvc.On("GenerateRefreshToken", claims.UserID).Return("new_refresh_token", newRefreshTokenClaims, nil).Once()
		s.mockTokenRepo.On("CreateToken", mock.Anything, mock.AnythingOfType("*domain.Token")).Return(nil, nil).Once()

		_, _, err := s.usecase.RefreshTokenForMobile(context.Background(), refreshToken)

		s.NoError(err)
		s.mockJwtSvc.AssertExpectations(s.T())
		s.mockTokenRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Mobile Token Not In DB (Delete Fails)", func() {
		s.SetupTest()
		s.mockJwtSvc.On("ValidateToken", refreshToken).Return(claims, nil).Once()
		s.mockTokenRepo.On("DeleteToken", mock.Anything, string(domain.RefreshToken), claims.ID).Return(domain.ErrNotFound).Once()

		_, _, err := s.usecase.RefreshTokenForMobile(context.Background(), refreshToken)

		s.ErrorIs(err, domain.ErrAuthenticationFailed)
		s.mockJwtSvc.AssertNotCalled(s.T(), "GenerateAccessToken")
	})
}

// --- Update Password Tests ---
func (s *UserUsecaseTestSuite) TestUpdatePassword() {
	userID := "user-123"
	currentPassword := "oldpass"
	newPassword := "newpass123"
	hashedCurrent := "hashed_oldpass"
	hashedNew := "hashed_newpass"

	account := &domain.Account{
		ID:           userID,
		PasswordHash: hashedCurrent,
		UserDetail:   &domain.UserDetail{},
	}

	s.Run("Success", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetById", mock.Anything, userID).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", hashedCurrent, currentPassword).Return(nil).Once()
		s.mockPassSvc.On("HashPassword", newPassword).Return(hashedNew, nil).Once()
		s.mockUserRepo.On("UpdatePassword", mock.Anything, userID, hashedNew).Return(nil).Once()

		err := s.usecase.UpdatePassword(context.Background(), userID, currentPassword, newPassword)

		s.NoError(err)
		s.mockUserRepo.AssertExpectations(s.T())
		s.mockPassSvc.AssertExpectations(s.T())
	})

	s.Run("Failure - User Not Found", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetById", mock.Anything, userID).Return(nil, domain.ErrUserNotFound).Once()

		err := s.usecase.UpdatePassword(context.Background(), userID, currentPassword, newPassword)

		s.ErrorIs(err, domain.ErrUserNotFound)
		s.mockUserRepo.AssertExpectations(s.T())
	})

	s.Run("Failure - Incorrect Current Password", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetById", mock.Anything, userID).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", hashedCurrent, currentPassword).Return(errors.New("wrong password")).Once()

		err := s.usecase.UpdatePassword(context.Background(), userID, currentPassword, newPassword)

		s.ErrorIs(err, domain.ErrAuthenticationFailed)
		s.mockPassSvc.AssertExpectations(s.T())
	})

	s.Run("Failure - New Password Too Short", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetById", mock.Anything, userID).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", hashedCurrent, currentPassword).Return(nil).Once()

		err := s.usecase.UpdatePassword(context.Background(), userID, currentPassword, "short")

		s.ErrorIs(err, domain.ErrPasswordTooShort)
	})

	s.Run("Failure - Hashing Error", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetById", mock.Anything, userID).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", hashedCurrent, currentPassword).Return(nil).Once()
		s.mockPassSvc.On("HashPassword", newPassword).Return("", errors.New("hash error")).Once()

		err := s.usecase.UpdatePassword(context.Background(), userID, currentPassword, newPassword)

		s.Error(err)
	})

	s.Run("Failure - Repo Update Error", func() {
		s.SetupTest()
		s.mockUserRepo.On("GetById", mock.Anything, userID).Return(account, nil).Once()
		s.mockPassSvc.On("ComparePassword", hashedCurrent, currentPassword).Return(nil).Once()
		s.mockPassSvc.On("HashPassword", newPassword).Return(hashedNew, nil).Once()
		s.mockUserRepo.On("UpdatePassword", mock.Anything, userID, hashedNew).Return(errors.New("db error")).Once()

		err := s.usecase.UpdatePassword(context.Background(), userID, currentPassword, newPassword)

		s.Error(err)
	})
}
