package domain

import (
	"context"
)

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


type IPostUseCase interface {
	CreatePost(ctx context.Context, discussion *Post) (*Post, error)
	GetPosts(ctx context.Context, opts PostFilters) ([]*Post, int64, error)
	GetPostByID(ctx context.Context, id string) (*Post, error)
	UpdatePost(ctx context.Context, Post *Post) (*Post, error)
	DeletePost(ctx context.Context, id, userID, role string) error
}