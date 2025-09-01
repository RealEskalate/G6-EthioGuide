package infrastructure

import (
	"EthioGuide/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"time"
)

type HuggingFaceEmbedding struct {
    APIKeys []string
    Model   string
	EmbeddingUrl string
    counter uint64
}

// NewHuggingFaceEmbedding creates an embedding client with a list of API keys and a model name
func NewHuggingFaceEmbedding(apiKeys []string, model, embeddingUrl string) domain.IEmbeddingService {
    if len(apiKeys) == 0 {
        panic("at least one vector Embedding API key must be provided")
    }
    return &HuggingFaceEmbedding{
        APIKeys: apiKeys,
        Model:   model,
		EmbeddingUrl: embeddingUrl,
    }
}

// getNextKey rotates through the available API keys in a thread-safe way
func (h *HuggingFaceEmbedding) getNextKey() string {
    idx := atomic.AddUint64(&h.counter, 1)
    return h.APIKeys[int(idx-1)%len(h.APIKeys)]
}

func (h *HuggingFaceEmbedding) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
    url := h.EmbeddingUrl + h.Model

    reqBody, _ := json.Marshal(map[string]string{"inputs": text})
    req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", "Bearer "+h.getNextKey())
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, errors.New("huggingface embedding request failed")
    }

    var embedding [][]float64
    if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
        return nil, err
    }

    return embedding[0], nil
}
