package usecase

import (
	"EthioGuide/domain"
	"context"
	"net/mail"
	"regexp"
	"time"
)

type userUsecase struct {
	userRepo domain.IUserRepository
	// tokenRepo            TokenRepository
	passwordService domain.IPasswordService
	jwtService      domain.IJWTService
	// emailService         infrastructure.EmailService
	contextTimeout time.Duration
}

func NewUserUsecase(ur domain.IUserRepository, ps domain.IPasswordService, js domain.IJWTService, timeout time.Duration) domain.IUserUsecase {
	return &userUsecase{
		userRepo:        ur,
		passwordService: ps,
		jwtService:      js,
		// emailService:         es,
		// imageUploaderService: ius,
		contextTimeout: timeout,
	}
}

func (uc *userUsecase) Register(c context.Context, user *domain.Account) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// if err := user.Validate(); err != nil {
	// 	return err
	// }

	existingUser, _ := uc.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return domain.ErrEmailExists
	}

	existingUser, _ = uc.userRepo.GetByUsername(ctx, user.UserDetail.Username)
	if existingUser != nil {
		return domain.ErrUsernameExists
	}

	hashedPassword, err := uc.passwordService.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword
	user.Role = domain.RoleUser
	user.UserDetail.IsVerified = false

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
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
