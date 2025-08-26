package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"time"
)

type AuthUsecase struct {
	Jwtservice domain.IJWTService
	Tokenrepo  domain.IAuthRepository
}

func NewAuthUsecase(jwtser domain.IJWTService, repo domain.IAuthRepository) *AuthUsecase {
	return &AuthUsecase{
		Jwtservice: jwtser,
		Tokenrepo:  repo,
	}
}

type TokenTypes string

const (
	RefreshToken       TokenTypes = "refreshToken"
	PasswordResetToken TokenTypes = "passwordResetToken"
	ActivationToken    TokenTypes = "activationToken"
)

func (auc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken == "" {
		return "", "", fmt.Errorf("empty token string")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	claims, err := auc.Jwtservice.ParseExpiredToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if _, err := auc.Tokenrepo.GetToken(ctx, string(RefreshToken), refreshToken); err != nil {
		return "", "", err
	}

	if err := auc.Tokenrepo.DeleteToken(ctx, string(RefreshToken), refreshToken); err != nil {
		return "", "", err
	}

	newAccessToken, _, errAccess := auc.Jwtservice.GenerateAccessToken(claims.UserID, claims.Role)
	if errAccess != nil {
		return "", "", errAccess
	}

	newRefreshToken, refreshClaim, errRefresh := auc.Jwtservice.GenerateRefreshToken(claims.UserID)
	if errRefresh != nil {
		return "", "", errRefresh
	}

	token := &domain.TokenModel{
		Token:     newRefreshToken,
		TokenType: string(RefreshToken),
		ExpiresAt: refreshClaim.ExpiresAt.Time,
	}
	if _, err := auc.Tokenrepo.CreateToken(ctx, token); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
