package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
)

type ChatUsecase struct {
    embedService domain.IEmbeddingService
    procedureRepo    domain.IProcedureRepository // abstracts Mongo / vector DB
    llmService   domain.IAIService       // abstracts Gemini / OpenAI
}

func NewChatUsecase(e domain.IEmbeddingService, s domain.IProcedureRepository, l domain.IAIService) *ChatUsecase {
    return &ChatUsecase{embedService: e, procedureRepo: s, llmService: l}
}

func (u *ChatUsecase) AIchat(ctx context.Context, query string) (string, error) {
    // 1. embed query
    vec, err := u.embedService.GenerateEmbedding(ctx, query)
    if err != nil {
        return "", err
    }

    // 2. vector search
    docs, err := u.procedureRepo.SearchByEmbedding(ctx, vec, 3)
    if err != nil {
        return "", err
    }

    // 3. call LLM
	    var contextText string
    for _, d := range docs {
        contextText += fmt.Sprintf("Title: %s\nRequirements: %s\nSteps: %s\nFees: %s\n\n",
            d.Title, d.Requirements, d.Steps, d.Fees)
    }

    prompt := fmt.Sprintf(`You are an assistant helping Ethiopians navigate bureaucracy.
    Here are the most relevant procedures:
    %s
    Now answer the user query: %s`, contextText, query)

    answer, err := u.llmService.GenerateCompletion(ctx, prompt)
    if err != nil {
        return "", err
    }

    return answer, nil
}



////


///
