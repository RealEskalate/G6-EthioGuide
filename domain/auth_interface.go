package domain

import "context"

type IAuthInterface interface {
	RefreshTokenForWeb(ctx context.Context, refreshToken string) (string, error)
	RefreshTokenForMobile(ctx context.Context, refreshToken string) (string, string, error)
}
