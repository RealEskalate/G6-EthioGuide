package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock IJWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID string, role domain.Role) (string, *domain.JWTClaims, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Get(1).(*domain.JWTClaims), args.Error(2)
}

func (m *MockJWTService) GenerateRefreshToken(userID string) (string, *domain.JWTClaims, error) {
	args := m.Called(userID)
	return args.String(0), args.Get(1).(*domain.JWTClaims), args.Error(2)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*domain.JWTClaims), args.Error(1)
}

func (m *MockJWTService) ParseExpiredToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*domain.JWTClaims), args.Error(1)
}

func (m *MockJWTService) GetRefreshTokenExpiry() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

// Mock IAuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateToken(ctx context.Context, token *domain.TokenModel) (*domain.TokenModel, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*domain.TokenModel), args.Error(1)
}

func (m *MockAuthRepository) GetToken(ctx context.Context, tokenType, token string) (string, error) {
	args := m.Called(ctx, tokenType, token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRepository) DeleteToken(ctx context.Context, tokenType, token string) error {
	args := m.Called(ctx, tokenType, token)
	return args.Error(0)
}

type AuthUsecaseTestSuite struct {
	suite.Suite
	mockJwtService *MockJWTService
	mockAuthRepo   *MockAuthRepository
	authUsecase    *AuthUsecase
}

func (suite *AuthUsecaseTestSuite) SetupTest() {
	suite.mockJwtService = new(MockJWTService)
	suite.mockAuthRepo = new(MockAuthRepository)
	suite.authUsecase = NewAuthUsecase(suite.mockJwtService, suite.mockAuthRepo)
}

func TestAuthUsecase(t *testing.T) {
	suite.Run(t, new(AuthUsecaseTestSuite))
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_Success() {

	userID := uuid.New().String()
	refreshToken := "test_refresh_token"

	claims := &domain.JWTClaims{
		UserID: userID,
		Role:   domain.RoleUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			ID:        uuid.New().String(),
		},
	}
	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"
	newRefreshClaims := &domain.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			ID:        uuid.New().String(),
		},
	}

	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(claims, nil).Once()
	suite.mockAuthRepo.On("GetToken", mock.Anything, string(RefreshToken), refreshToken).Return("", nil).Once()
	suite.mockAuthRepo.On("DeleteToken", mock.Anything, string(RefreshToken), refreshToken).Return(nil).Once()
	suite.mockJwtService.On("GenerateAccessToken", claims.UserID, claims.Role).Return(newAccessToken, &domain.JWTClaims{}, nil).Once()
	suite.mockJwtService.On("GenerateRefreshToken", claims.UserID).Return(newRefreshToken, newRefreshClaims, nil).Once()
	suite.mockAuthRepo.On("CreateToken", mock.Anything, mock.AnythingOfType("*domain.TokenModel")).Return(&domain.TokenModel{}, nil).Once()

	accessToken, refreshToken, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)

	suite.NoError(err)
	suite.Equal(newAccessToken, accessToken)
	suite.Equal(newRefreshToken, refreshToken)

	suite.mockJwtService.AssertExpectations(suite.T())
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_EmptyToken() {

	_, _, err := suite.authUsecase.RefreshToken(context.Background(), "")
	suite.Error(err)
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_ParseExpiredTokenError() {

	refreshToken := "test_refresh_token"

	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(&domain.JWTClaims{}, fmt.Errorf("some error")).Once()
	_, _, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)
	suite.Error(err)
	suite.mockJwtService.AssertExpectations(suite.T())
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_GetTokenError() {

	userID := uuid.New().String()
	refreshToken := "test_refresh_token"

	claims := &domain.JWTClaims{UserID: userID}
	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(claims, nil).Once()
	suite.mockAuthRepo.On("GetToken", mock.Anything, string(RefreshToken), refreshToken).Return("", fmt.Errorf("some error")).Once()
	_, _, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)
	suite.Error(err)
	suite.mockJwtService.AssertExpectations(suite.T())
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_DeleteTokenError() {

	userID := uuid.New().String()
	refreshToken := "test_refresh_token"

	claims := &domain.JWTClaims{UserID: userID}
	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(claims, nil).Once()
	suite.mockAuthRepo.On("GetToken", mock.Anything, string(RefreshToken), refreshToken).Return("", nil).Once()
	suite.mockAuthRepo.On("DeleteToken", mock.Anything, string(RefreshToken), refreshToken).Return(fmt.Errorf("some error")).Once()
	_, _, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)
	suite.Error(err)
	suite.mockJwtService.AssertExpectations(suite.T())
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_GenerateAccessTokenError() {

	userID := uuid.New().String()
	refreshToken := "test_refresh_token"

	claims := &domain.JWTClaims{UserID: userID}
	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(claims, nil).Once()
	suite.mockAuthRepo.On("GetToken", mock.Anything, string(RefreshToken), refreshToken).Return("", nil).Once()
	suite.mockAuthRepo.On("DeleteToken", mock.Anything, string(RefreshToken), refreshToken).Return(nil).Once()
	suite.mockJwtService.On("GenerateAccessToken", claims.UserID, claims.Role).Return("", &domain.JWTClaims{}, fmt.Errorf("some error")).Once()
	_, _, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)
	suite.Error(err)
	suite.mockJwtService.AssertExpectations(suite.T())
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_GenerateRefreshTokenError() {

	userID := uuid.New().String()
	refreshToken := "test_refresh_token"

	claims := &domain.JWTClaims{UserID: userID}
	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(claims, nil).Once()
	suite.mockAuthRepo.On("GetToken", mock.Anything, string(RefreshToken), refreshToken).Return("", nil).Once()
	suite.mockAuthRepo.On("DeleteToken", mock.Anything, string(RefreshToken), refreshToken).Return(nil).Once()
	suite.mockJwtService.On("GenerateAccessToken", claims.UserID, claims.Role).Return("new_access_token", &domain.JWTClaims{}, nil).Once()
	suite.mockJwtService.On("GenerateRefreshToken", claims.UserID).Return("", &domain.JWTClaims{}, fmt.Errorf("some error")).Once()

	_, _, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)

	suite.Error(err)
	suite.mockJwtService.AssertExpectations(suite.T())
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUsecaseTestSuite) TestRefreshToken_CreateTokenError() {

	userID := uuid.New().String()
	refreshToken := "test_refresh_token"

	claims := &domain.JWTClaims{UserID: userID}
	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"
	newRefreshClaims := &domain.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	suite.mockJwtService.On("ParseExpiredToken", refreshToken).Return(claims, nil).Once()
	suite.mockAuthRepo.On("GetToken", mock.Anything, string(RefreshToken), refreshToken).Return("", nil).Once()
	suite.mockAuthRepo.On("DeleteToken", mock.Anything, string(RefreshToken), refreshToken).Return(nil).Once()
	suite.mockJwtService.On("GenerateAccessToken", claims.UserID, claims.Role).Return(newAccessToken, &domain.JWTClaims{}, nil).Once()
	suite.mockJwtService.On("GenerateRefreshToken", claims.UserID).Return(newRefreshToken, newRefreshClaims, nil).Once()
	suite.mockAuthRepo.On("CreateToken", mock.Anything, mock.AnythingOfType("*domain.TokenModel")).Return(&domain.TokenModel{}, fmt.Errorf("some error")).Once()

	_, _, err := suite.authUsecase.RefreshToken(context.Background(), refreshToken)

	suite.Error(err)
	suite.mockJwtService.AssertExpectations(suite.T())
	suite.mockAuthRepo.AssertExpectations(suite.T())
}