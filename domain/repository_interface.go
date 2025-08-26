package domain

import (
	"context"
)

type IAuthRepository interface {
	CreateToken(ctx context.Context, token *TokenModel) (*TokenModel, error)
	GetToken(ctx context.Context, tokentype, token string) (string, error)
	DeleteToken(ctx context.Context, tokentype, token string) error
}
