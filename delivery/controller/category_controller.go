package controller

import (
	"EthioGuide/domain"
	"net/http"

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

func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	// Implementation for creating a category will go here
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &domain.Category{
		OrganizationID: req.OrganizationID,
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