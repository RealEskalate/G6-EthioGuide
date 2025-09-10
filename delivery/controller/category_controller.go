package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryUsecase domain.ICategoryUsecase
}

func NewCategoryController(cc domain.ICategoryUsecase) *CategoryController {
	return &CategoryController{
		categoryUsecase: cc,
	}
}

// @Summary      Create Category
// @Description  Create a category for procedures
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body CreateCategoryRequest true "Category Details"
// @Success      200 {object} domain.Category "New Category"
// @Failure      400 {string}  "invalid
// @Failure      500 {string}  "invalid"
// @Router       /categories [post]
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	// Implementation for creating a category will go here
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID := c.GetString("userID")

	category := &domain.Category{
		OrganizationID: orgID,
		ParentID:       req.ParentID,
		Title:          req.Title,
	}

	err := ctrl.categoryUsecase.CreateCategory(c.Request.Context(), category)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "category": category})
}

// @Summary      Get Categories
// @Description  Get list of categories
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        sortOrder query string false "asc | desc"
// @Param        parentID query string false "parentID"
// @Param        organizationID query string false "organizationID"
// @Param        title query string false "title"
// @Success      200 {object} PaginatedCategoryResponse "List of Categories"
// @Failure      400 {string}  "invalid
// @Failure      500 {string}  "invalid"
// @Router       /categories [get]
func (ctrl *CategoryController) GetCategory(c *gin.Context) {
	// Implementation for creating a category will go here
	options := domain.CategorySearchAndFilter{}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'page' parameter"})
		return
	}
	options.Page = page

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'limit' parameter"})
		return
	}
	options.Limit = limit

	options.SortBy = c.Query("sortBy") // e.g., "date", "popularity", "title"
	if strings.ToLower(c.Query("sortOrder")) == string(domain.SortAsc) {
		options.SortOrder = domain.SortAsc
	} else {
		options.SortOrder = domain.SortDesc
	}

	options.ParentID = c.DefaultQuery("parentID", "")
	options.OrganizationID = c.DefaultQuery("organizationID", "")
	options.Title = c.DefaultQuery("title", "")

	categories, total, err := ctrl.categoryUsecase.GetCategories(c.Request.Context(), &options)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPaginatedCategoryResponse(categories, total, options.Page, options.Limit))
}

// --- helper functions ---
func toCategoryResponse(c *domain.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:             c.ID,
		OrganizationID: c.OrganizationID,
		ParentID:       c.ParentID,
		Title:          c.Title,
	}
}

func toPaginatedCategoryResponse(categories []*domain.Category, total, page, limit int64) PaginatedCategoryResponse {
	catResponse := make([]*CategoryResponse, len(categories))
	for i, c := range categories {
		catResponse[i] = toCategoryResponse(c)
	}

	return PaginatedCategoryResponse{
		Data:  catResponse,
		Total: total,
		Page:  page,
		Limit: limit,
	}
}
