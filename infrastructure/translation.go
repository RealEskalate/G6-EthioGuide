package infrastructure

import (
	"EthioGuide/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings" // Ensure strings is imported

	"github.com/gin-gonic/gin"
)

// bodyLogWriter now ONLY captures the body. It does not write to the original writer.
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write ONLY captures the body. The original response writer is NOT called here.
func (w bodyLogWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// TranslationMiddleware creates a Gin middleware that intercepts JSON responses
// and translates them based on the 'lang' header.
func NewTranslationMiddleware(geminiUseCase domain.IGeminiUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetLang := c.GetHeader("lang")
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		// If no translation is needed, we don't need the custom writer.
		if targetLang == "" || targetLang == "en" {
			c.Next()
			return
		}

		// Replace the writer
		c.Writer = blw
		c.Next()

		// After the handler, get the captured original body
		originalBody := blw.body.Bytes()
		statusCode := blw.Status() // Use the captured status
		contentType := blw.Header().Get("Content-Type")

		// Important: Set the original writer back before we write anything.
		c.Writer = blw.ResponseWriter

		// --- Translation Logic ---
		if statusCode != http.StatusOK || !strings.Contains(contentType, "application/json") || strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			// Write the original body and exit if we are not translating
			c.Writer.Write(originalBody)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(originalBody, &data); err != nil {
			// Not valid JSON, write the original body back and exit
			c.Writer.Write(originalBody)
			return
		}

		translatedData, err := geminiUseCase.TranslateJSON(c.Request.Context(), data, targetLang)
		if err != nil {
			fmt.Printf("Translation to '%s' failed: %v\n", targetLang, err)
			// On failure, write the original body back
			c.Writer.Write(originalBody)
			return
		}

		translatedBody, err := json.Marshal(translatedData)
		if err != nil {
			fmt.Printf("Failed to marshal translated data: %v\n", err)
			// On failure, write the original body back
			c.Writer.Write(originalBody)
			return
		}

		// On success, write the NEW translated body
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(translatedBody)))
		c.Writer.Write(translatedBody)
	}
}
