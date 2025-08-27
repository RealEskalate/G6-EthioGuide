package domain

import "context"

type IAccountRepository interface {
	GetByEmail(ctx context.Context, email string) (*Account, error)
	GetByPhone(ctx context.Context, phone string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
}