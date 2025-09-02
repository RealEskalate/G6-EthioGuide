package controller

import (
	"EthioGuide/domain"

	"github.com/gin-gonic/gin"
)



type PostController struct {
    useCase domain.IPostUseCase  
}

func NewPostController(uc domain.IPostUseCase) *PostController {
	return &PostController{
		useCase: uc,
	}
}

func (dc *PostController) CreatePost(c *gin.Context) {
	var dto CreatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, domain.ErrInvalidBody)
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	Post := &domain.Post{
		UserID:   userID.(string),
		Title:    dto.Title,
		Content:  dto.Content,
		Procedures: dto.Procedures,
		Tags: dto.Tags,
	}
	err := dc.useCase.CreatePost(c.Request.Context(), Post)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, gin.H{"message": "Post created successfully"})
}

