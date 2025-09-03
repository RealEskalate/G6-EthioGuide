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