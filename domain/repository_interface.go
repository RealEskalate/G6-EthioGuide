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
}

type ITokenRepository interface {
	CreateToken(ctx context.Context, token *Token) (*Token, error)
	GetToken(ctx context.Context, tokentype, token string) (string, error)
	DeleteToken(ctx context.Context, tokentype, token string) error
}

type IProcedureRepository interface {
	Create(ctx context.Context, procedure *Procedure) error
	SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*Procedure, error)
}

type IAIChatRepository interface {
	Save(ctx context.Context, chat *AIChat) error
	GetByUser(ctx context.Context, userID string, limit int) ([]*AIChat, error)
	DeleteByUser(ctx context.Context, userID string) error
}
