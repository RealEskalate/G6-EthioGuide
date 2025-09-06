package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strconv"

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
// @Success      200 {object} AIConversationResponse  "Response"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /ai/guide [post]
func (c *AIChatController) AIChatController(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var req AIChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	answer, err := c.usecase.AIchat(ctx.Request.Context(), userID.(string), req.Query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ToAIChatResponse(answer))
}

// @Summary      Get Paginated ai chat history
// @Description  Get Paginated ai chat history
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        page  query     int     false  "Page number (default 1)"
// @Param        limit query     int     false  "Results per page (default 10)"
// @Success      200   {object}  PaginatedAIHisoryResponse
// @Failure      400   {object}  map[string]string "Invalid parameter"
// @Failure      500   {object}  map[string]string "Server error"
// @Router       /ai/history [get]
func (c *AIChatController) AIChatHistoryController(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'page' parameter"})
		return
	}

	limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'limit' parameter"})
		return
	}

	conversations, total, err := c.usecase.AIHistory(ctx.Request.Context(), userID.(string), page, limit)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toPaginatedAIHisory(conversations, total, page, limit))
}
