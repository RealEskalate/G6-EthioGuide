package usecase

import (
	"EthioGuide/domain"
	"context"
	"time"
)

type userUsecase struct {
	userRepo             domain.IUserRepository
	// tokenRepo            TokenRepository
	passwordService      domain.IPasswordService
	// jwtService           infrastructure.JWTService
	// emailService         infrastructure.EmailService
	contextTimeout       time.Duration
}

func NewUserUsecase(ur domain.IUserRepository, ps domain.IPasswordService, timeout time.Duration) domain.IUserUsecase {
	return &userUsecase{
		userRepo:             ur,
		passwordService:      ps,
		// jwtService:           js,
		// emailService:         es,
		// imageUploaderService: ius,
		contextTimeout:       timeout,
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

	hashedPassword, err := uc.passwordService.HashPassword(*(&user.PasswordHash))
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
