package domain

import "context"

type IUserUsecase interface {
	Register(ctx context.Context, user *Account) error
	// ActivateAccount(ctx context.Context, activationTokenValue string) error
	// Login(ctx context.Context, identifier, password string) (accessToken, refreshToken string, err error)
	// Logout(ctx context.Context, refreshToken string) error
	// RefreshAccessToken(ctx context.Context, refreshToken, accessToken string) (newAccessToken, newRefreshToken string, err error)

	// // Password Management
	// ForgetPassword(ctx context.Context, email string) error
	// ResetPassword(ctx context.Context, resetToken, newPassword string) error

	// //Profile Management
	// UpdateProfile(c context.Context, userID, bio string, profilePicFile multipart.File, profilePicHeader *multipart.FileHeader) (*domain.User, error)
	// GetProfile(c context.Context, userID string) (*domain.User, error)

	// // User Management
	// SearchAndFilter(ctx context.Context, options domain.UserSearchFilterOptions) ([]*domain.User, int64, error)
	// SetUserRole(ctx context.Context, actorUserID string, actorRole domain.Role, targetUserID string, newRole domain.Role) (*domain.User, error)
}

type IGeminiUseCase interface {
	TranslateContent(ctx context.Context, content, targetLang string) (string, error)
}
