package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NoticeController struct {
	noticeUsecase domain.INoticeUseCase
}

func NewNoticeController(noticeUsecase domain.INoticeUseCase) *NoticeController {
	return &NoticeController{
		noticeUsecase: noticeUsecase,
	}
}

// @Summary      Create Notice
// @Description  Create new notice.
// @Tags         Notice
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body NoticeDTO true "Procedure Detail"
// @Success      200 {string}  "Notice Created successfully."
// @Failure      400 {string}  "Invalid request"
// @Failure      404 {string}  "Not Found"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /notices [post]
func (nc *NoticeController) CreateNotice(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var noticeDTO NoticeDTO
	if err := ctx.ShouldBindJSON(&noticeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notice := noticeDTO.ToDomainNotice()
	notice.OrganizationID = userID
	if err := nc.noticeUsecase.CreateNotice(ctx, notice); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "Notice Created successfully.")
}

// @Summary      Get List of Notices
// @Description  Search and filter Notices with pagination, sorting, and various filters.
// @Tags         Notice
// @Accept       json
// @Produce      json
// @Param        page              query     int     false  "Page number (default 1)"
// @Param        limit             query     int     false  "Results per page (default 10)"
// @Param        organizationId    query     string  false  "Filter by organization ID"
// @Param        sortBy            query     string  false  "Sort by field (e.g. createdAt, fee, processingTime)"
// @Param        sortOrder         query     string  false  "Sort order: ASC or DESC (default DESC)"
// @Success      200  {object}  NoticeListResponse "List of Notices"
// @Failure      400  {string}   "Bad Request"
// @Failure      500  {string}   "Server error"
// @Router       /notices [get]
func (nc *NoticeController) GetNoticesByFilter(ctx *gin.Context) {
	filter := &domain.NoticeFilter{}

	// Optional: organization filter
	if orgID := ctx.Query("organizationId"); orgID != "" {
		filter.OrganizationID = orgID
	}

	// Tags AND logic: provide multiple via comma-separated list
	if tagsStr := ctx.Query("tags"); tagsStr != "" {
		filter.Tags = strings.Split(tagsStr, ",")
	} else if single := ctx.Query("tag"); single != "" {
		// Backward compat for single tag param
		filter.Tags = []string{single}
	}

	// Pagination
	if p := ctx.DefaultQuery("page", "1"); p != "" {
		if page, err := strconv.ParseInt(p, 10, 64); err == nil {
			filter.Page = page
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter"})
			return
		}
	}
	if l := ctx.DefaultQuery("limit", "10"); l != "" {
		if limit, err := strconv.ParseInt(l, 10, 64); err == nil {
			filter.Limit = limit
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' parameter"})
			return
		}
	}

	// Sorting (only created_At supported)
	filter.SortBy = ctx.Query("sortBy")
	filter.SortOrder = domain.SortOrder(strings.ToLower(ctx.DefaultQuery("sortOrder", string(domain.SortDesc))))

	notices, total, err := nc.noticeUsecase.GetNoticesByFilter(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]NoticeDTO, len(notices))
	for i, n := range notices {
		out[i] = *FromDomainNotice(n)
	}
	resp := NoticeListResponse{
		Data:  out,
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary      Update Notice
// @Description  Update a notice.
// @Tags         Notice
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Notice ID"
// @param        request body NoticeDTO true "Details of updated notice"
// @Success      200  {string}  "Notice Updated successfully."
// @Failure      400  {string}   "Bad Request"
// @Failure      500  {string}   "Server error"
// @Router       /notices/{id} [patch]
func (nc *NoticeController) UpdateNotice(ctx *gin.Context) {
	id := ctx.Param("id")
	var noticeDTO NoticeDTO
	if err := ctx.ShouldBindJSON(&noticeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notice := noticeDTO.ToDomainNotice()
	if err := nc.noticeUsecase.UpdateNotice(ctx, id, notice); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Notice Updated successfully.")
}

// @Summary      Delete Notice
// @Description  Deletes an existing Notice. Requires admin or organization ownership.
// @Tags         Notice
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Notice ID"
// @Success      200  "Notice Deleted successfully."
// @Failure      401  {string} Unauthorized
// @Failure      403  {string} Permission Denied
// @Failure      404  {string} Procedure not found
// @Router       /notices/{id} [delete]
func (nc *NoticeController) DeleteNotice(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := nc.noticeUsecase.DeleteNotice(ctx, id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Notice Deleted successfully.")
}
