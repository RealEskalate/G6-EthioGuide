package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
)

type AIChatUsecase struct {
	embedService  domain.IEmbeddingService
	procedureRepo domain.IProcedureRepository // abstracts Mongo / vector DB
    aiChatRepo  domain.AIChatRepository
	llmService    domain.IAIService           // abstracts Gemini / OpenAI
}

func NewChatUsecase(e domain.IEmbeddingService, s domain.IProcedureRepository, l domain.IAIService) domain.IAIChatUsecase {
	return &AIChatUsecase{embedService: e, procedureRepo: s, llmService: l}
}

func (u *AIChatUsecase) AIchat(ctx context.Context, query string) (string, error) {
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

	prompt := fmt.Sprintf(`You are an assistant helping Ethiopians navigate bureaucracy.
    Here are the most relevant procedures:
    %s
    Now answer the user query: %s`, contextText, query)

	answer, err := u.llmService.GenerateCompletion(ctx, prompt)
	if err != nil {
		return "", err
	}
    // After: return answer, nil

    // Example: Save chat history (pseudo, adjust as needed)
    userID, _ := ctx.Value("userID").(string)
    chat := &domain.AIChat{
        UserID:   userID, // You need to get this from context or as a parameter
        Source:   "ai_chat",
        Request:  query,
        Response: answer,
        // Timestamp will be set in the repository
    }
    err = u.aiChatRepo.Save(ctx, chat)
    if err != nil {
        // Optionally log or handle the error, but don't block the user
}

	return answer, nil
}

////
// func (r *ProcedureRepositoryMongo) SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*domain.Procedure, error) {
//     filter := bson.M{}
//     opts := options.Find().
//         SetLimit(int64(limit)).
//         SetSort(bson.D{{"embedding", bson.D{{"$meta", "vectorSearch"}}}})

//     // vector search stage
//     searchStage := bson.D{
//         {"$vectorSearch", bson.D{
//             {"queryVector", queryVec},
//             {"path", "embedding"},
//             {"numCandidates", 100},  // wider search, filtered down to limit
//             {"limit", limit},
//             {"index", "embedding_index"},
//         }},
//     }

//     cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{searchStage})
//     if err != nil {
//         return nil, err
//     }
//     defer cursor.Close(ctx)

//     var results []*domain.Procedure
//     if err := cursor.All(ctx, &results); err != nil {
//         return nil, err
//     }
//     return results, nil
// }
// 🔸 Before this works:

// You must create a vector index on embedding:

// js
// Copy code
// db.procedures.createIndex(
//   { embedding: "vector" },
//   { vectorOptions: { dims: 384, similarity: "cosine" } }
// )
///
