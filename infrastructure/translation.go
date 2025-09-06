package infrastructure

import (
	"EthioGuide/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// bodyLogWriter is a custom ResponseWriter that captures the response body.
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the body while also writing it to the original ResponseWriter.
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// TranslationMiddleware creates a Gin middleware that intercepts JSON responses
// and translates them based on the 'lang' header.
func NewTranslationMiddleware(geminiUseCase domain.IGeminiUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get target language from header
		targetLang := c.GetHeader("lang")

		// If 'lang' is not set, is 'en', or is empty, we don't need to do anything.
		if targetLang == "" || targetLang == "en" {
			c.Next()
			return
		}

		// 2. Create a custom response writer to capture the body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 3. Process the request through the next handlers in the chain
		c.Next()

		// 4. After the handler, check the response status and content type
		statusCode := c.Writer.Status()
		contentType := c.Writer.Header().Get("Content-Type")

		// We only want to translate successful JSON responses.
		// Exclude swagger docs from translation.
		if statusCode != http.StatusOK || !strings.Contains(contentType, "application/json") || strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			return
		}

		// 5. Unmarshal the captured body
		originalBody := blw.body.Bytes()
		var data map[string]interface{}
		if err := json.Unmarshal(originalBody, &data); err != nil {
			// If it's not valid JSON (e.g., empty body or just a string "OK"), we can't translate it.
			// This is not an error, just something we can't process.
			return
		}

		// 6. Perform the translation
		translatedData, err := geminiUseCase.TranslateJSON(c.Request.Context(), data, targetLang)
		if err != nil {
			// Log the translation error but don't fail the original request.
			// The client will receive the untranslated (English) version.
			// In a production system, you'd use a structured logger.
			fmt.Printf("Translation to '%s' failed: %v\n", targetLang, err)
			return
		}

		// 7. Marshal the new translated data back to JSON
		translatedBody, err := json.Marshal(translatedData)
		if err != nil {
			fmt.Printf("Failed to marshal translated data: %v\n", err)
			return
		}

		// 8. Overwrite the response
		// First, clear the original response body that was already written.
		// We can't actually clear it, but we can replace it by writing a new one.
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(translatedBody)))

		// Reset the writer to the original one to avoid recursion
		c.Writer = blw.ResponseWriter
		c.Writer.Write(translatedBody)
	}
}
