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

// @Summary      Submit Feedback
// @Description  Submit a feedback for a procedure
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Procedure ID"
// @Param        request body FeedbackCreateRequest true "Feedback Detail"
// @Success      201 {object} domain.Feedback "Feedback Submitted Successfully"
// @Failure      400 {string}  "Bad Request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Internal"
// @Router       /procedures/{id}/feedback [post]
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
		UserID:      userID,
		ProcedureID: procedureID,
		Content:     req.Content,
		Type:        domain.FeedbackType(req.Type),
		Tags:        req.Tags,
	}

	err := ctrl.feedbackUsecase.SubmitFeedback(c.Request.Context(), feedback)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback submitted successfully", "feedback": fromDomainFeedback(feedback)})
}

// @Summary      Fetch Feedbacks
// @Description  Fetch list of feedbacks for a procedure
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        status query string false "status"
// @Param        id path string true "Procedure ID"
// @Success      200 {object} FeedbackListResponse "Feedbacks"
// @Failure      400 {string}  "Bad Request"
// @Failure      500 {string}  "Internal"
// @Router       /procedures/{id}/feedback [get]
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

	if status := c.Query("status"); status != "" {
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

// @Summary      Update Feedback
// @Description  Update the status of a feedback
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Feedback ID"
// @Param        request body FeedbackStatePatchRequest true "Update Status of Feedback"
// @Success      201 {string}  "Feedback Updated Successfully"
// @Failure      400 {string}  "Bad Request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      404 {string}  "Not Found"
// @Failure      500 {string}  "Internal"
// @Router       /feedback/{id} [patch]
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

// @Summary      Fetch Feedbacks
// @Description  Fetch list of feedbacks for admin
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        status query string false "status"
// @Param        procedure_id query string false "Procedure ID"
// @Success      200 {object} FeedbackListResponse "Feedbacks"
// @Failure      400 {string}  "Bad Request"
// @Failure      500 {string}  "Internal"
// @Router       /feedback [get]
func (ctrl *FeedbackController) GetAllFeedbacks(c *gin.Context) {
	filter := domain.FeedbackFilter{}
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		HandleError(c, domain.ErrInvalidQueryParam)
		return
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		HandleError(c, domain.ErrInvalidQueryParam)
		return
	}
	filter.Page = page
	filter.Limit = limit
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if procedureID := c.Query("procedure_id"); procedureID != "" {
		filter.ProcedureID = &procedureID
	}
	feedbacks, total, err := ctrl.feedbackUsecase.GetAllFeedbacks(c.Request.Context(), &filter)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"feedbacks": toFeedbackListResponse(feedbacks, total, page, limit)})
}
