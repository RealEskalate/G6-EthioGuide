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

	GetProfile(ctx context.Context, userID string) (*Account, error)
	UpdatePassword(ctx context.Context, userID, currentPassword, newPassword string) error
	LoginWithSocial(ctx context.Context, provider AuthProvider, code string) (*Account, string, string, error)
	UpdateProfile(ctx context.Context, userID string, updates map[string]interface{}) (*Account, error)
}

type IGeminiUseCase interface {
	TranslateContent(ctx context.Context, content, targetLang string) (string, error)
}

type ICategoryUsecase interface {
	CreateCategory(ctx context.Context, category *Category) error
	GetCategories(ctx context.Context, options *CategorySearchAndFilter) ([]*Category, int64, error)
}

type IProcedureUsecase interface {
	CreateProcedure(ctx context.Context, procedure *Procedure) error
}

type IFeedbackUsecase interface {
	SubmitFeedback(ctx context.Context, feedback *Feedback) error
	GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *FeedbackFilter) ([]*Feedback, int64, error)
	UpdateFeedbackStatus(ctx context.Context, feedbackID, userID string, status FeedbackStatus, adminResponse *string) error
}
