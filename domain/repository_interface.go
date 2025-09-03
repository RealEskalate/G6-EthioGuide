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

type IChecklistRepository interface {
	CreateChecklist(ctx context.Context, userid, procdureID string) (*UserProcedure, error)
	GetProcedures(ctx context.Context, userid string) ([]*UserProcedure, error)
	GetChecklistByUserProcedureID(ctx context.Context, userprocedureID string) ([]*Checklist, error)
	ToggleCheck(ctx context.Context, checklistID string) error
	FindCheck(ctx context.Context, checklistID string) (*Checklist, error)
	CountDocumentsChecklist(ctx context.Context, filter interface{}) (int64, error)
	UpdateUserProcedure(ctx context.Context, filter interface{}, update map[string]interface{}) error
}
