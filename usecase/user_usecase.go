package usecase

import (
	"EthioGuide/domain"
	"context"
	"net/mail"
	"regexp"
)

type userUsecase struct {
	userRepo domain.IAccountRepository
	passwordService domain.IPasswordService
	jwtService    domain.IJWTService
}

func NewUserUsecase(userRepo domain.IAccountRepository, passwordService domain.IPasswordService, jwtService domain.IJWTService) *userUsecase {
	return &userUsecase{
		userRepo: userRepo,
		passwordService: passwordService,
		jwtService: jwtService,
	}
}

func (uc *userUsecase) Login(ctx context.Context, identifier, password string) (*domain.Account, string, string, error) {
	var account *domain.Account
	var err error
	if _, mailErr := mail.ParseAddress(identifier); mailErr == nil {
		account, err = uc.userRepo.GetByEmail(ctx, identifier)
	} else if result, _ := regexp.MatchString(`^+?[0-9]{8,14}$`, identifier); result {
		account, err = uc.userRepo.GetByPhone(ctx, identifier)
	} else {
		account, err = uc.userRepo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		return nil, "", "", err
	}
	if account == nil {
		return nil, "", "", domain.ErrAuthenticationFailed
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
		return nil, "", "", err
	}

	refreshToken, _, err := uc.jwtService.GenerateRefreshToken(account.ID)
	if err != nil {
		return nil, "", "", err
	}

	return account, accessToken, refreshToken, nil
}