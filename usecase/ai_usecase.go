package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
)

type geminiUseCase struct {
	geminiServices domain.IAIService
}


func NewGeminiUsecase(geminiServices domain.IAIService) domain.IGeminiUseCase {
	return &geminiUseCase{
		geminiServices: geminiServices,
	}
}

func (g *geminiUseCase) TranslateContent(ctx context.Context, content string, targetLang string) (string, error) {

	prompt := fmt.Sprintf("Translate the following text to %s. Ensure the translation is accurate, maintains the original meaning, and uses natural, fluent language appropriate for the context. Do not add or omit any information. %s", targetLang, content)
	supportedLangs := map[string]bool{
		"en": true, "am": true,
	}

	if !supportedLangs[targetLang]{
		return "", fmt.Errorf("unsupported target language: %s", targetLang)
	}
	
	translated, err := g.geminiServices.GenerateCompletion(ctx, prompt)
	if err != nil{
		return "", err
	}

	return translated, nil

}