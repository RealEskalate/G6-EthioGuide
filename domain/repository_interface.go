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
	UpdateUserFields(ctx context.Context, userIDstr string, update map[string]interface{}) error
}

type ITokenRepository interface {
	CreateToken(ctx context.Context, token *Token) (*Token, error)
	GetToken(ctx context.Context, tokentype, token string) (string, error)
	DeleteToken(ctx context.Context, tokentype, token string) error
}

type ISearchRepository interface {
	Search(ctx context.Context, filter SearchFilterRequest) (*SearchResult, error)
	FindProcedures(ctx context.Context, filter SearchFilterRequest) ([]*Procedure, error)
	FindOrganizations(ctx context.Context, filter SearchFilterRequest) ([]*AccountOrgSearch, error)
}
