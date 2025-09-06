package usecase

import (
	"EthioGuide/domain"
	"context"
	"errors"
	"fmt"
	"log"
)

type AIChatUsecase struct {
	EmbedService  domain.IEmbeddingService
	ProcedureRepo domain.IProcedureRepository // abstracts Mongo / vector DB
	AiChatRepo    domain.IAIChatRepository
	LLMService    domain.IAIService // abstracts Gemini / OpenAI
}

func NewChatUsecase(e domain.IEmbeddingService, s domain.IProcedureRepository, aiChatRepo domain.IAIChatRepository, l domain.IAIService) domain.IAIChatUsecase {
	return &AIChatUsecase{EmbedService: e, ProcedureRepo: s, AiChatRepo: aiChatRepo, LLMService: l}
}

func (u *AIChatUsecase) AIchat(ctx context.Context, query string) (string, error) {
	classifierPrompt := fmt.Sprintf(`
	Classify the following user query into one of these categories:
	- procedure   (government services, documents, licenses, permits, taxes, etc.)
	- irrelevant  (math, politics, jokes, casual talk, health, personal advice, etc.)
	- offensive   (insults, hate speech, profanity, harassment, NSFW, etc.)

	Query: "%s"
	Return only one word: procedure, irrelevant, or offensive.`, query)

	category, err := u.LLMService.GenerateCompletion(ctx, classifierPrompt)
	if err != nil {
		return "", err
	}
	if category == "offensive"{
		return "", errors.New("your query contains offensive content and cannot be processed")

	}else if category == "irrelevant"{
		return "Sorry, I can only answer questions related to Ethiopian government procedures.", nil
	}

	// detect the language
	prompt := fmt.Sprintf("I want you to identify the language of this promt %s and i want to give me the only the language in small later like if it is Amharic give me amharic. and if you do not know the language just give me only  a word 'unkown'.", query)
	orglang, err := u.LLMService.GenerateCompletion(ctx, prompt)
	if err != nil {
		return "", err
	}
	if orglang != "english" {
		prompt := fmt.Sprintf("translate this query %s in to English lanuage. And i do not want you to add another thing by yourself", query)
		query, err = u.LLMService.GenerateCompletion(ctx, prompt)
		if err != nil {
			return "", nil
		}
	} else if orglang != "unkown" {
		return "", errors.New("unknown language")
	}

	// 1. embed query
	vec, err := u.EmbedService.GenerateEmbedding(ctx, query)
	if err != nil {
		return "", err
	}

	// 2. vector search
	docs, err := u.ProcedureRepo.SearchByEmbedding(ctx, vec, 3)
	if err != nil {
		return "", err
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

	prompt = fmt.Sprintf(`You are an assistant helping Ethiopians navigate bureaucracy.
    Here are the most relevant procedures:
    %s
    Now answer the user query: %s`, contextText, query)

	answer, err := u.LLMService.GenerateCompletion(ctx, prompt)
	if err != nil {
		return "", err
	}
	source := "unofficial"
	if len(docs) > 0 {
		source = "official"
	}
	// After: return answer, nil

	// Example: Save chat history (pseudo, adjust as needed)
	userID, _ := ctx.Value("userID").(string)
	chat := &domain.AIChat{
		UserID:   userID, // You need to get this from context or as a parameter
		Source:   source,
		Request:  query,
		Response: answer,
		// Timestamp will be set in the repository
	}
	err = u.AiChatRepo.Save(ctx, chat)
	if err != nil {
		// Optionally log or handle the error, but don't block the user
		log.Println("the chat is not saved")
	}
	if orglang != "english" {
		prompt := fmt.Sprintf(`I want you to translate procedure into this %s language by keeping its format as it is.
		here is the procedure
		%s
		`, orglang, answer)
		answer, err = u.LLMService.GenerateCompletion(ctx, prompt)
		if err != nil {
			return "", err
		}
	}

	return answer, nil
}
