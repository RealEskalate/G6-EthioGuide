package usecase

import (
    "EthioGuide/domain"
    "context"
)

type PreferencesUsecase struct {
    repo domain.IPreferencesRepository
}

func NewPreferencesUsecase(repo domain.IPreferencesRepository) domain.IPreferencesUsecase{
    return &PreferencesUsecase{repo: repo}
}

func (u *PreferencesUsecase) CreateUserPreferences(ctx context.Context, userID string) error {
	preferences := &domain.Preferences{
		UserID: userID,
		PreferredLang: domain.English,
		PushNotification: false,
		EmailNotification: false,

	}
    return u.repo.Create(ctx, preferences)
}

func (u *PreferencesUsecase) GetUserPreferences(ctx context.Context, userId string) (*domain.Preferences, error) {
    return u.repo.GetByUserID(ctx, userId)
}

func (u *PreferencesUsecase) UpdateUserPreferences(ctx context.Context, preferences *domain.Preferences) error {
    return  u.repo.UpdateByUserID(ctx, preferences.UserID, preferences)
}