package usecase

import (
    "context"
    "ethio-guide/domain"
)

type ChatUsecase struct {
    embedService domain.EmbeddingService
    searchRepo   domain.SearchRepository // abstracts Mongo / vector DB
    llmService   domain.LLMService       // abstracts Gemini / OpenAI
}

func NewChatUsecase(e domain.EmbeddingService, s domain.SearchRepository, l domain.LLMService) *ChatUsecase {
    return &ChatUsecase{embedService: e, searchRepo: s, llmService: l}
}

func (u *ChatUsecase) Answer(ctx context.Context, query string) (string, error) {
    // 1. embed query
    vec, err := u.embedService.Generate(ctx, query)
    if err != nil {
        return "", err
    }

    // 2. vector search
    docs, err := u.searchRepo.SearchByEmbedding(ctx, vec, 3)
    if err != nil {
        return "", err
    }

    // 3. call LLM
    answer, err := u.llmService.Answer(ctx, query, docs)
    if err != nil {
        return "", err
    }

    return answer, nil
}
