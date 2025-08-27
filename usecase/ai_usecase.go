package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

var supportedLangs = map[string]bool{
	"en": true,
	"am": true,
}

type geminiUseCase struct {
	geminiServices domain.IAIService
	contextTimeout time.Duration
}

func NewGeminiUsecase(geminiServices domain.IAIService, timeOut time.Duration) domain.IGeminiUseCase {
	return &geminiUseCase{
		geminiServices: geminiServices,
		contextTimeout: timeOut,
	}
}

func (g *geminiUseCase) TranslateContent(ctx context.Context, content string, targetLang string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, g.contextTimeout)
	defer cancel()

	content = strings.TrimSpace(content)
	if content == "" {
		return "", nil // Return early for empty content.
	}

	if !supportedLangs[targetLang] {
		// Fail fast for unsupported languages.
		return "", fmt.Errorf("%w: %s", domain.ErrUnsupportedLanguage, targetLang)
	}

	prompt := fmt.Sprintf(`
	Translate the following content to %s.
	Ensure the translation is accurate, maintains the original meaning, and uses natural, fluent language appropriate for the context. Do not add or omit any information. 
	Do not include any other text or markdown formatting.
	If you don't know the language I requested, send only the text "unknown language".
	Here is the content:
	%s
	`, targetLang, content)

	translated, err := g.geminiServices.GenerateCompletion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("gemini service failed to generate completion: %w", err)
	}

	// This check remains a valuable safety net.
	if translated == "unknown language" {
		return "", domain.ErrUnsupportedLanguage
	}

	return strings.TrimSpace(translated), nil
}
