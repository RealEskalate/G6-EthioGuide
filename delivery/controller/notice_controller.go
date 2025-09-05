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

func (nc *NoticeController) CreateNotice(ctx *gin.Context) {
	var noticeDTO NoticeDTO
	if err := ctx.ShouldBindJSON(&noticeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notice := noticeDTO.ToDomainNotice()
	if err := nc.noticeUsecase.CreateNotice(ctx, notice); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "Notice Created successfully.")
}

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
	filter.SortOrder = domain.SortOrder(strings.ToUpper(ctx.DefaultQuery("sortOrder", string(domain.SortDesc))))

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

func (nc *NoticeController) UpdateNotice(ctx *gin.Context, id string) {
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

func (nc *NoticeController) DeleteNotice(ctx *gin.Context, id string) {
	if err := nc.noticeUsecase.DeleteNotice(ctx, id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Notice Deleted successfully.")
}
