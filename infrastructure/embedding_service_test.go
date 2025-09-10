package infrastructure_test

import (
	. "EthioGuide/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// EmbeddingServiceTestSuite is the test suite for HuggingFaceEmbedding
type EmbeddingServiceTestSuite struct {
	suite.Suite
}

// TestNewEmbeddingService validates the constructor logic.
func (s *EmbeddingServiceTestSuite) TestNewEmbeddingService() {
	s.Run("Success", func() {
		service, err := NewEmbeddingService([]string{"test-key"}, "http://example.com")
		s.NoError(err)
		s.NotNil(service)
	})

	s.Run("Panic_NoAPIKeys", func() {
		// The function is expected to panic, so we assert that.
		s.Panics(func() {
			_, _ = NewEmbeddingService([]string{}, "http://example.com")
		}, "Should panic when no API keys are provided")
	})
}

// TestGenerateEmbedding_IntegrationWithMockServer tests the full flow against a mock HTTP server.
func (s *EmbeddingServiceTestSuite) TestGenerateEmbedding_IntegrationWithMockServer() {
	s.Run("Success", func() {
		// Arrange: Create a mock server that simulates a successful HF API response.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Assert that the request from our service is correct
			s.Equal("Bearer test-key", r.Header.Get("Authorization"))
			s.Equal("application/json", r.Header.Get("Content-Type"))

			var body map[string][]string
			err := json.NewDecoder(r.Body).Decode(&body)
			s.NoError(err)
			s.Equal([]string{"hello world"}, body["inputs"])

			// Send a successful response
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `[[0.1, 0.2, 0.3]]`)
		}))
		defer mockServer.Close()

		service, _ := NewEmbeddingService([]string{"test-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		s.NoError(err)
		s.NotNil(embedding)
		s.Equal([]float64{0.1, 0.2, 0.3}, embedding)
	})

	s.Run("Error_EmptyInputText", func() {
		// Arrange
		service, _ := NewEmbeddingService([]string{"test-key"}, "http://dummy.url")

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "")

		// Assert
		s.Error(err)
		s.Nil(embedding)
		s.EqualError(err, "text input cannot be empty")
	})

	s.Run("Error_API_Returns_Non_200", func() {
		// Arrange: Create a mock server that returns an error status code.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer mockServer.Close()

		service, _ := NewEmbeddingService([]string{"invalid-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		s.Error(err)
		s.Nil(embedding)
		s.EqualError(err, "huggingface embedding request failed")
	})

	s.Run("Error_Invalid_JSON_Response", func() {
		// Arrange: Server returns malformed JSON.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `[[0.1, 0.2]`) // Malformed
		}))
		defer mockServer.Close()

		service, _ := NewEmbeddingService([]string{"test-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		s.Error(err)
		s.Nil(embedding)
		_, isJsonError := err.(*json.SyntaxError)
		_, isJsonUnmarshalTypeError := err.(*json.UnmarshalTypeError)
		isUnexpectedEof := err.Error() == "unexpected EOF"

		s.True(isJsonError || isJsonUnmarshalTypeError || isUnexpectedEof, "Error should be a JSON parsing error")
	})

	s.Run("Error_Context_Timeout", func() {
		// Arrange: Server takes too long to respond.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		service, _ := NewEmbeddingService([]string{"test-key"}, mockServer.URL)

		// Create a context with a very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		// Act
		embedding, err := service.GenerateEmbedding(ctx, "hello world")

		// Assert
		s.Error(err)
		s.Nil(embedding)
		s.ErrorIs(err, context.DeadlineExceeded)
	})
}

// TestKeyRotation verifies the round-robin key rotation logic.
func (s *EmbeddingServiceTestSuite) TestKeyRotation() {
	apiKeys := []string{"key-1", "key-2", "key-3"}
	var receivedKeys []string
	var mu sync.Mutex

	// Arrange: Create a mock server that records which API key it received.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		mu.Lock()
		receivedKeys = append(receivedKeys, authHeader)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[[1.0]]`)
	}))
	defer mockServer.Close()

	service, _ := NewEmbeddingService(apiKeys, mockServer.URL)

	// Act: Call the service multiple times to trigger key rotation.
	for i := 0; i < 5; i++ {
		_, _ = service.GenerateEmbedding(context.Background(), "test")
	}

	// Assert
	expectedKeySequence := []string{
		"Bearer key-1",
		"Bearer key-2",
		"Bearer key-3",
		"Bearer key-1",
		"Bearer key-2",
	}
	s.Equal(expectedKeySequence, receivedKeys)
}

// TestEmbeddingServiceTestSuite runs the entire suite
func TestEmbeddingServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmbeddingServiceTestSuite))
}
