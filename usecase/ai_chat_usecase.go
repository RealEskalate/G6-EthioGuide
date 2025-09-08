package usecase

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type AIChatUsecase struct {
	EmbedService   domain.IEmbeddingService
	ProcedureRepo  domain.IProcedureRepository // abstracts Mongo / vector DB
	AiChatRepo     domain.IAIChatRepository
	LLMService     domain.IAIService // abstracts Gemini / OpenAI
	contextTimeout time.Duration
}

func NewChatUsecase(e domain.IEmbeddingService, s domain.IProcedureRepository, aiChatRepo domain.IAIChatRepository, l domain.IAIService, timeOut time.Duration) domain.IAIChatUsecase {
	return &AIChatUsecase{
		EmbedService:   e,
		ProcedureRepo:  s,
		AiChatRepo:     aiChatRepo,
		LLMService:     l,
		contextTimeout: timeOut,
	}
}

func (u *AIChatUsecase) AIchat(ctx context.Context, userId, query string) (*domain.AIChat, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userFriendlyErrorResponse := "I'm sorry, but I encountered a technical issue and can't answer your question right now. Please try again in a few moments."
	orgQuery := query
	classifierPrompt := fmt.Sprintf(`
	Classify the following user query into one of these categories:
	- procedure   (government services, documents, licenses, permits, taxes, etc.)
	- irrelevant  (math, politics, jokes, casual talk, health, personal advice, etc.)
	- offensive   (insults, hate speech, profanity, harassment, NSFW, etc.)

	Query: "%s"
	Return only one word: procedure, irrelevant, or offensive.`, query)

	category, err := u.LLMService.GenerateCompletion(ctx, classifierPrompt)
	if err != nil {
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  userFriendlyErrorResponse,
		}, err
	}
	switch category {
	case "offensive":
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  "I am unable to process this request as it appears to violate our content policy. Please rephrase your question to continue.",
		}, errors.New("your query contains offensive content and cannot be processed")

	case "irrelevant":
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  "I can only answer questions about Ethiopian government procedures. Could you please ask a question related to that topic?",
		}, nil
	}

	// detect the language
	prompt := fmt.Sprintf("I want you to identify the language of this promt %s and i want to give me the only the language in small later like if it is Amharic give me amharic. and if you do not know the language just give me only  a word 'unknown'.", query)
	orglang, err := u.LLMService.GenerateCompletion(ctx, prompt)
	if err != nil {
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  userFriendlyErrorResponse,
		}, err
	}
	if orglang == "unknown" {
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  "I'm sorry, I could not determine the language of your request. I currently support Amharic and English. Please try again in one of these languages.",
		}, domain.ErrUnsupportedLanguage
	} else if orglang != "english" {
		prompt := fmt.Sprintf("translate this query %s in to English lanuage. And i do not want you to add another thing by yourself", query)
		query, err = u.LLMService.GenerateCompletion(ctx, prompt)
		if err != nil {
			return &domain.AIChat{
				ID:        time.Now().String(),
				Timestamp: time.Now(),
				UserID:    userId,
				Source:    "Couldn't translate the request",
				Request:   orgQuery,
				Response:  "I'm sorry, I had trouble understanding your request. Could you please try rephrasing it?",
			}, fmt.Errorf("%w %v", domain.ErrUnsupportedLanguage, err)
		}
	}

	// 1. embed query
	vec, err := u.EmbedService.GenerateEmbedding(ctx, query)
	if err != nil {
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  userFriendlyErrorResponse,
		}, err
	}

	// 2. vector search
	docs, err := u.ProcedureRepo.SearchByEmbedding(ctx, vec, 3)
	if err != nil {
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  userFriendlyErrorResponse,
		}, err
	}

	// 3. call LLM
	var contextText string
	for _, d := range docs {
		contextText += fmt.Sprintf(
			"Name: %s\nPrerequisites: %v\nSteps: %v\nResult: %v\nFees: %s %.2f (%s)\nProcessing Time: %d-%d days\n\n",
			d.Name,
			d.Content.Prerequisites,
			d.Content.Steps,
			d.Content.Result,
			d.Fees.Currency, d.Fees.Amount, d.Fees.Label,
			d.ProcessingTime.MinDays, d.ProcessingTime.MaxDays,
		)
	}

	prompt = fmt.Sprintf(`You are an assistant for EthioGuide, helping people with Ethiopian bureaucracy.
Your tone should be clear, helpful, and encouraging.
Use the following information to answer the user's query.

*Relevant Procedures:*
---
%s
---

*User's Query:* "%s"

*Instructions:*
1.  Provide a direct and concise answer to the user's query.
2.  Format your answer for easy reading. Use asterisks for bolding (e.g., *Key Requirement*) and numbered lists for steps.
3.  If the provided procedures are a good match, present the most relevant one clearly.
4.  If the procedures don't seem to be a good match, politely state that you couldn't find a specific procedure but can offer general advice.`, contextText, query)

	answer, err := u.LLMService.GenerateCompletion(ctx, prompt)
	if err != nil {
		return &domain.AIChat{
			ID:        time.Now().String(),
			Timestamp: time.Now(),
			UserID:    userId,
			Source:    "unofficial",
			Request:   orgQuery,
			Response:  userFriendlyErrorResponse,
		}, err
	}
	source := "unofficial"
	if len(docs) > 0 {
		source = "official"
	}
	// After: return answer, nil

	if orglang != "english" {
		prompt := fmt.Sprintf(`Translate the following answer into the language '%s'.
    IMPORTANT: Preserve the formatting exactly, including any asterisks for bolding and any numbered lists. Do not add any extra commentary.

    *Answer to Translate:*
    %s
    `, orglang, answer)
		answer, err = u.LLMService.GenerateCompletion(ctx, prompt)
		if err != nil {
			return &domain.AIChat{
				ID:        time.Now().String(),
				Timestamp: time.Now(),
				UserID:    userId,
				Source:    "unofficial",
				Request:   orgQuery,
				Response:  userFriendlyErrorResponse,
			}, err
		}
	}

	related_procedures := make([]*domain.AIProcedure, len(docs))
	for i, proc := range docs {
		related_procedures[i] = &domain.AIProcedure{
			Id:   proc.ID,
			Name: proc.Name,
		}
	}

	// Example: Save chat history (pseudo, adjust as needed)
	chat := &domain.AIChat{
		UserID:            userId,
		Source:            source,
		Request:           orgQuery,
		Response:          answer,
		RelatedProcedures: related_procedures,
	}

	err = u.AiChatRepo.Save(ctx, chat)
	if err != nil {
		// Optionally log or handle the error, but don't block the user
		log.Println("the chat is not saved")
	}

	return chat, nil
}

func (u *AIChatUsecase) AIHistory(ctx context.Context, userId string, page, limit int64) ([]*domain.AIChat, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.AiChatRepo.GetByUser(ctx, userId, page, limit)
}
