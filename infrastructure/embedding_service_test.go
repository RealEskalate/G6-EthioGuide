package infrastructure_test

import (
	"EthioGuide/config"
	"EthioGuide/domain"
	. "EthioGuide/infrastructure"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EmbeddingServiceSuite struct {
	suite.Suite
	apiURL  string
	apiKey  []string
	service domain.IEmbeddingService
}

func (s *EmbeddingServiceSuite) SetupSuite() {
	cfg := config.LoadForTest()
	s.apiURL = cfg.EmbeddingUrl
	if s.apiURL == "" {
		s.T().Skip("Skipping EmbeddingService integration tests: EMBEDDING_URL not set.")
	}
	s.apiKey = append(s.apiKey, cfg.EmbeddingApiKey)
	if len(s.apiKey) == 0 || s.apiKey[0] == "" {
		s.T().Skip("Skipping EmbeddingService integration tests: api key not set.")

	}
}

func (s *EmbeddingServiceSuite) SetupTest() {
	service, err := NewEmbeddingService(s.apiKey, s.apiURL)
	s.Require().NoError(err, "Failed to create a new Embedding Service instance")
	s.service = service
}

func TestEmbeddingServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}
	suite.Run(t, new(EmbeddingServiceSuite))
}

// --- Tests ---

func (s *EmbeddingServiceSuite) TestGenerateEmbedding() {
	ctx := context.Background()

	s.Run("Success - Simple Text", func() {
		text := "Hello, world!"
		embedding, err := s.service.GenerateEmbedding(ctx, text)
		s.NoError(err, "The API call should not return an error")
		s.NotNil(embedding, "The embedding should not be nil")
		s.Greater(len(embedding), 0, "The embedding should have a non-zero length")
	})

	s.Run("Failure - Empty Text", func() {
		embedding, err := s.service.GenerateEmbedding(ctx, "")
		s.Error(err, "Should return an error for empty input")
		s.Nil(embedding, "Embedding should be nil for empty input")
	})
}
