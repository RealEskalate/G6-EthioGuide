package domain

import (
	"context"
)

type IUserRepository interface {
	Create(ctx context.Context, user *Account) error
	GetByEmail(ctx context.Context, email string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
	// GetByID(ctx context.Context, id string) (*domain.User, error)
	// Update(ctx context.Context, user *domain.User) error
	// FindUserIDsByName(ctx context.Context, authorName string) ([]string, error)
	// FindByProviderID(ctx context.Context, provider domain.AuthProvider, providerID string) (*domain.User, error)
	// SearchAndFilter(ctx context.Context, options domain.UserSearchFilterOptions) ([]*domain.User, int64, error)
}