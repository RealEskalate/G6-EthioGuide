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

	ctx.JSON(http.StatusOK, ToAIChatResponse(answer))
}

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

	conversations, total, err := c.usecase.AIHistory(userID, page, limit)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toPaginatedAIHisory(conversations, total, page, limit))
}
