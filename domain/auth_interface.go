package domain

import "context"

type IAuthInterface interface {
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}
