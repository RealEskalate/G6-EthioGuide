package domain

import (
	"context"
	"mime/multipart"
	"time"

	"golang.org/x/oauth2"
)

type IAIService interface {
	GenerateCompletion(ctx context.Context, prompt string) (string, error)
}

type GoogleUserInfo struct {
	ID                string
	Email             string
	Name              string
	ProfilePictureURL string
}

type IGoogleOAuthService interface {
	ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error)
}

type ImageUploaderService interface {
	UploadProfilePicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type ICacheService interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	AddToSet(ctx context.Context, key string, members ...any) error
	GetSetMembers(ctx context.Context, key string) ([]string, error)
	DeleteKeys(ctx context.Context, keys []string) error
}
