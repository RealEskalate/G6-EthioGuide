package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedbackController struct {
	feedbackUsecase domain.IFeedbackUsecase
}

func NewFeedbackController(fu domain.IFeedbackUsecase) *FeedbackController {
	return &FeedbackController{
		feedbackUsecase: fu,
	}
}

func (ctrl *FeedbackController) SubmitFeedback(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		return
	}

	procedureID := c.Param("id")

	var req FeedbackCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if procedureID == "" || req.Content == "" || req.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID, ProcedureID, Content, and Type are required fields"})
		return
	}

	feedback := &domain.Feedback{
		UserID: userID,
		ProcedureID: procedureID,
		Content: req.Content,
		Type: domain.FeedbackType(req.Type),
		Tags: req.Tags,
	}

	err := ctrl.feedbackUsecase.SubmitFeedback(c.Request.Context(), feedback)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback submitted successfully", "feedback": fromDomainFeedback(feedback)})
}

func (ctrl *FeedbackController) GetAllFeedbacksForProcedure(c *gin.Context) {
	filter := domain.FeedbackFilter{}
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'page' parameter"})
		return
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'limit' parameter"})
		return
	}
	filter.Page = page
	filter.Limit = limit

	if status := c.Query("title"); status != "" {
		filter.Status = &status
	}

	procedureID := c.Param("id")

	feedbacks, total, err := ctrl.feedbackUsecase.GetAllFeedbacksForProcedure(c.Request.Context(), procedureID, &filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": toFeedbackListResponse(feedbacks, total, page, limit)})
}

func (ctrl *FeedbackController) UpdateFeedbackStatus(c *gin.Context) {
	var req FeedbackStatePatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		return
	}

	feedbackID := c.Param("id")
	err := ctrl.feedbackUsecase.UpdateFeedbackStatus(c.Request.Context(), feedbackID, userID, domain.FeedbackStatus(req.Status), &req.AdminResponse)

	if err != nil {
		HandleError(c, err)
		return 
	}
	c.JSON(http.StatusOK, gin.H{"message": "Feedback status updated successfully"})
}