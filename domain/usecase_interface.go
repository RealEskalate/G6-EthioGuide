package domain

import "context"

type IGeminiUseCase interface {
    TranslateContent(ctx context.Context, content, targetLang string) (string, error)
}