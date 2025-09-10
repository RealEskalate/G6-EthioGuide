package infrastructure_test

import (
	"EthioGuide/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewCohereEmbedding validates the constructor logic.
func TestNewCohereEmbedding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service, err := infrastructure.NewCohereEmbedding([]string{"test-key"}, "http://example.com")
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("Error_NoAPIKeys", func(t *testing.T) {
		service, err := infrastructure.NewCohereEmbedding([]string{}, "http://example.com")
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.EqualError(t, err, "at least one Cohere API key must be provided")
	})
}

// TestGenerateEmbedding_IntegrationWithMockServer tests the full flow against a mock HTTP server.
func TestGenerateEmbedding_IntegrationWithMockServer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange: Create a mock server that simulates a successful Cohere API response.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Assert that the request from our service is correct
			assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Send a successful response
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"embeddings": [[0.1, 0.2, 0.3]]}`)
		}))
		defer mockServer.Close()

		service, _ := infrastructure.NewCohereEmbedding([]string{"test-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, embedding)
		assert.Equal(t, []float64{0.1, 0.2, 0.3}, embedding)
	})

	t.Run("Error_API_Returns_Non_200", func(t *testing.T) {
		// Arrange: Create a mock server that returns an error status code.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{"message": "Invalid API Key"}`)
		}))
		defer mockServer.Close()

		service, _ := infrastructure.NewCohereEmbedding([]string{"invalid-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, embedding)
		assert.Contains(t, err.Error(), "cohere embedding request failed with status 401 Unauthorized")
	})

	t.Run("Error_No_Embeddings_In_Response", func(t *testing.T) {
		// Arrange: Server returns 200 OK but with an empty embeddings array.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"embeddings": []}`)
		}))
		defer mockServer.Close()

		service, _ := infrastructure.NewCohereEmbedding([]string{"test-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, embedding)
		assert.EqualError(t, err, "no embeddings returned from Cohere")
	})

	t.Run("Error_Invalid_JSON_Response", func(t *testing.T) {
		// Arrange: Server returns malformed JSON.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"embeddings": [[0.1, 0.2]`) // Malformed
		}))
		defer mockServer.Close()

		service, _ := infrastructure.NewCohereEmbedding([]string{"test-key"}, mockServer.URL)

		// Act
		embedding, err := service.GenerateEmbedding(context.Background(), "hello world")

		// Assert
		assert.Error(t, err, "An error should be returned for malformed JSON")
		assert.Nil(t, embedding, "Embedding should be nil on parsing error")

		// Check that the error is indeed a JSON decoding error.
		_, isJsonError := err.(*json.SyntaxError)
		_, isJsonUnmarshalTypeError := err.(*json.UnmarshalTypeError)
		isUnexpectedEof := err.Error() == "unexpected EOF"

		assert.True(t, isJsonError || isJsonUnmarshalTypeError || isUnexpectedEof, "Error should be a JSON parsing error")
	})

	t.Run("Error_Context_Timeout", func(t *testing.T) {
		// Arrange: Server takes too long to respond.
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		service, _ := infrastructure.NewCohereEmbedding([]string{"test-key"}, mockServer.URL)

		// Create a context with a very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		// Act
		embedding, err := service.GenerateEmbedding(ctx, "hello world")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, embedding)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

// TestKeyRotation verifies the round-robin key rotation logic.
func TestKeyRotation(t *testing.T) {
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
		fmt.Fprintln(w, `{"embeddings": [[1.0]]}`)
	}))
	defer mockServer.Close()

	service, _ := infrastructure.NewCohereEmbedding(apiKeys, mockServer.URL)

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
	assert.Equal(t, expectedKeySequence, receivedKeys)
}
