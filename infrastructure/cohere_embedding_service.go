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

// CohereEmbedding is the infrastructure service that calls Cohere Embedding API
type CohereEmbedding struct {
	APIKeys []string
	Url     string
	counter uint64
}

// NewCohereEmbedding creates a new embedding service for Cohere (English model only)
func NewCohereEmbedding(apiKeys []string, url string) (domain.IEmbeddingService, error) {
	if len(apiKeys) == 0 {
		return nil, errors.New("at least one Cohere API key must be provided")

	}
	return &CohereEmbedding{
		APIKeys: apiKeys,
		Url:     url,
	}, nil
}

// getNextKey rotates through API keys in a thread-safe way
func (c *CohereEmbedding) getNextKey() string {
	idx := atomic.AddUint64(&c.counter, 1)
	return c.APIKeys[int(idx-1)%len(c.APIKeys)]
}

// GenerateEmbedding calls Cohere API and returns the embedding vector for the given text
func (c *CohereEmbedding) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	// url := "https://api.cohere.ai/v1/embed"

	// Cohere request body
	body := map[string]interface{}{
		"model":      "embed-english-v3.0", // English-only model
		"texts":      []string{text},
		"input_type": "search_query", // Good for semantic search
	}
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.Url, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+c.getNextKey())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("cohere embedding request failed with status " + resp.Status)
	}

	// Parse response
	var res struct {
		Embeddings [][]float64 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Embeddings) == 0 {
		return nil, errors.New("no embeddings returned from Cohere")
	}

	return res.Embeddings[0], nil
}
