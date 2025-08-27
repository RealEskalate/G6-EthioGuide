package domain

import "context"

type IUserUsecase interface {
	Login(ctx context.Context, identifier, password string) (*Account, string, string, error)
}