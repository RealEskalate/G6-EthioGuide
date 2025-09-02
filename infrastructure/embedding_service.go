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
	APIKeys      []string
	EmbeddingUrl string
	counter      *uint64
}

// NewHuggingFaceEmbedding creates an embedding client with a list of API keys and a model name
func NewEmbeddingService(apiKeys []string, embeddingUrl string) (domain.IEmbeddingService, error) {
	if len(apiKeys) == 0 {
		panic("at least one vector Embedding API key must be provided")
	}
	var counter uint64 = 0
	return &HuggingFaceEmbedding{
		APIKeys:      apiKeys,
		EmbeddingUrl: embeddingUrl,
		counter:      &counter,
	}, nil
}

// getNextKey rotates through the available API keys in a thread-safe way
func (h *HuggingFaceEmbedding) getNextKey() string {
	idx := atomic.AddUint64(h.counter, 1)
	return h.APIKeys[int(idx-1)%len(h.APIKeys)]
}

func (h *HuggingFaceEmbedding) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	if text == "" {
		return nil, errors.New("text input cannot be empty")
	}
	url := h.EmbeddingUrl

	reqBody, _ := json.Marshal(map[string][]string{"inputs": {text}})
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
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
