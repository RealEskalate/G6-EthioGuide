package domain

import "context"

type IUserUsecase interface {
	Register(ctx context.Context, user *Account) error
	Login(ctx context.Context, identifier, password string) (*Account, string, string, error)
	// VerifyAccount(ctx context.Context, activationTokenValue string) error
	RefreshTokenForWeb(ctx context.Context, refreshToken string) (string, error)
	RefreshTokenForMobile(ctx context.Context, refreshToken string) (string, string, error)

	// // Password Management
	// ForgetPassword(ctx context.Context, email string) error
	// ResetPassword(ctx context.Context, resetToken, newPassword string) error
}

type IGeminiUseCase interface {
	TranslateContent(ctx context.Context, content, targetLang string) (string, error)
}

type ICategoryUsecase interface {
	GetCategories(ctx context.Context, options *CategorySearchAndFilter) ([]*Category, int64, error)
}