package usecase

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type UserUsecase struct {
	userRepo        domain.IAccountRepository
	tokenRepo       domain.ITokenRepository
	passwordService domain.IPasswordService
	jwtService      domain.IJWTService
	contextTimeout  time.Duration
}

func NewUserUsecase(
	ur domain.IAccountRepository,
	tr domain.ITokenRepository,
	ps domain.IPasswordService,
	js domain.IJWTService,
	timeout time.Duration,
) domain.IUserUsecase {
	return &UserUsecase{
		userRepo:        ur,
		tokenRepo:       tr,
		passwordService: ps,
		jwtService:      js,
		contextTimeout:  timeout,
	}
}

func (uc *UserUsecase) Register(c context.Context, user *domain.Account) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return domain.ErrInvalidEmailFormat
	}
	if len(user.PasswordHash) < 8 {
		return domain.ErrPasswordTooShort
	}
	if strings.TrimSpace(user.UserDetail.Username) == "" {
		return domain.ErrUsernameEmpty
	}

	_, err := uc.userRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		return domain.ErrEmailExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("error checking email existence: %w", err)
	}

	_, err = uc.userRepo.GetByUsername(ctx, user.UserDetail.Username)
	if err == nil {
		return domain.ErrUsernameExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("error checking username existence: %w", err)
	}

	hashedPassword, err := uc.passwordService.HashPassword(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashedPassword
	user.Role = domain.RoleUser
	user.UserDetail.IsVerified = false // Enforce business rules

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user in repository: %w", err)
	}

	return nil
}

// Login method is already quite good. Minimal changes for consistency.
func (uc *UserUsecase) Login(c context.Context, identifier, password string) (*domain.Account, string, string, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	var account *domain.Account
	var err error

	if _, mailErr := mail.ParseAddress(identifier); mailErr == nil {
		account, err = uc.userRepo.GetByEmail(ctx, identifier)
		// } else if result, _ := regexp.MatchString(`^\+?[0-9]{8,14}$`, identifier); result {
		// 	account, err = uc.userRepo.GetByPhoneNumber(ctx, identifier)
	} else {
		account, err = uc.userRepo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, "", "", domain.ErrAuthenticationFailed
		}
		return nil, "", "", fmt.Errorf("repository error during login: %w", err)
	}

	if !account.UserDetail.IsVerified {
		return nil, "", "", domain.ErrAccountNotActive
	}

	err = uc.passwordService.ComparePassword(account.PasswordHash, password)
	if err != nil {
		return nil, "", "", domain.ErrAuthenticationFailed
	}

	accessToken, _, err := uc.jwtService.GenerateAccessToken(account.ID, account.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshClaims, err := uc.jwtService.GenerateRefreshToken(account.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenToSave := &domain.Token{
		Id:        refreshClaims.ID,
		Token:     refreshToken,
		TokenType: domain.RefreshToken,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}

	// Save the token to the repository
	if _, err := uc.tokenRepo.CreateToken(ctx, tokenToSave); err != nil {
		// This is a critical error, as login succeeds but refresh will fail.
		return nil, "", "", fmt.Errorf("CRITICAL: failed to store refresh token after login: %w", err)
	}

	return account, accessToken, refreshToken, nil
}

func (uc *UserUsecase) RefreshTokenForWeb(ctx context.Context, refreshToken string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	claims, err := uc.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if the token exists in our database (it hasn't been revoked/logged out).
	// We use the token's unique ID (JTI) to look it up.
	_, err = uc.tokenRepo.GetToken(ctx, string(domain.RefreshToken), claims.ID)
	if err != nil {
		// If the token is not found in the repo, it's invalid, even if the signature is okay.
		return "", domain.ErrAuthenticationFailed
	}

	newAccessToken, _, err := uc.jwtService.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}

func (uc *UserUsecase) RefreshTokenForMobile(ctx context.Context, refreshToken string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	claims, err := uc.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Atomically find and delete the token. This prevents race conditions.
	// If your ITokenRepository can't do this atomically, a check-then-delete is the next best.
	// We assume DeleteToken fails if the token doesn't exist.
	err = uc.tokenRepo.DeleteToken(ctx, string(domain.RefreshToken), claims.ID)
	if err != nil {
		// This error means the token was not found or a DB error occurred.
		// In either case, it's an authentication failure.
		return "", "", domain.ErrAuthenticationFailed
	}

	newAccessToken, _, err := uc.jwtService.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken, refreshClaims, err := uc.jwtService.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	tokenToSave := &domain.Token{
		Id:        refreshClaims.ID, // Use the new token's ID
		Token:     newRefreshToken,
		TokenType: domain.RefreshToken,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}
	if _, err := uc.tokenRepo.CreateToken(ctx, tokenToSave); err != nil {
		return "", "", fmt.Errorf("CRITICAL: failed to store new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (uc *UserUsecase) GetProfile(c context.Context, userID string) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	account, err := uc.userRepo.GetById(ctx, userID)
	if err != nil || account == nil{
		return nil, domain.ErrUserNotFound

	}

	return account, nil
}
