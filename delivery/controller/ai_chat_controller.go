package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIChatController struct {
	usecase domain.IAIChatUsecase
}

func NewAIChatController(usecase domain.IAIChatUsecase) *AIChatController {
	return &AIChatController{usecase: usecase}
}

// @Summary      AI Chat
// @Description  Interact with the chatbot.
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body AIChatRequest  true "Prompt"
// @Success      200 {object} AIChatResponse  "Response"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /ai/guide [post]
func (c *AIChatController) AIChatController(ctx *gin.Context) {
	var req AIChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	answer, err := c.usecase.AIchat(ctx.Request.Context(), req.Query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, AIChatResponse{Answer: answer})
}
