package domain

import (
	"context"
)

type IAccountRepository interface {
	Create(ctx context.Context, user *Account) error
	GetById(ctx context.Context, id string) (*Account, error)
	GetByEmail(ctx context.Context, email string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
	// GetByPhoneNumber(ctx context.Context, phone string) (*Account, error)
	UpdatePassword(ctx context.Context, accountID, newPassword string) error
	UpdateProfile(ctx context.Context, account Account) error
}

type ITokenRepository interface {
	CreateToken(ctx context.Context, token *Token) (*Token, error)
	GetToken(ctx context.Context, tokentype, token string) (string, error)
	DeleteToken(ctx context.Context, tokentype, token string) error
}

type ICategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetCategories(ctx context.Context, options *CategorySearchAndFilter) ([]*Category, int64, error)
}

type IProcedureRepository interface {
	Create(ctx context.Context, procedure *Procedure) error
}
