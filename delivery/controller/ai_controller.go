package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeminiController struct {
	geminiUseCase domain.IGeminiUseCase
	aiChatUsecase domain.IAIChatUsecase
}

func NewGeminiController(geminiUseCase domain.IGeminiUseCase) *GeminiController {
	return &GeminiController{
		geminiUseCase: geminiUseCase,
	}
}

// @Summary      Translate Content
// @Description  Get translation of content in different languages.
// @Tags         AI
// @Param        lang  header  string  false  "Preferred language (default: en)"
// @Param        Authorization header string true "Bearer token"
// @Param        request body TranslateDTO true "Content to be translated"
// @Success      200  {string}  "Translated output"
// @Failure      401  {string}  Unauthorized
// @Failure      403  {string}  Permission Denied
// @Failure      404  {string}  Procedure not found
// @Router       /ai/translate [post]
func (gc *GeminiController) Translate(c *gin.Context) {

	var request TranslateDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	preferredLang := c.GetHeader("lang")
	if preferredLang == "" {
		preferredLang = "en"
	}

	translated, err := gc.geminiUseCase.TranslateContent(c.Request.Context(), request.Content, preferredLang)

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": translated})

}

// @Summary      AI Chat
// @Description  Interact with the chatbot.
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body ChatRequest true "Prompt"
// @Success      200 {string}  "Response"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /ai/guide [post]
func (gc *GeminiController) AIChat(c *gin.Context) {
	var request ChatRequest 
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	answer, err := gc.aiChatUsecase.AIchat(c.Request.Context(), request.Content)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
