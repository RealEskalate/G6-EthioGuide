package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strconv"

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
	res, err := dc.useCase.CreatePost(c.Request.Context(), Post)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": res})
}

func (dc *PostController) GetPosts(c *gin.Context) {
    title := c.Query("title")
	
    procedureIDs := c.QueryArray("procedure_ids")
    tags := c.QueryArray("tags")
    page, _ := strconv.ParseInt(c.DefaultQuery("page", "0"), 10, 64)
    limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
    sortBy := c.DefaultQuery("sort_by", "created_at")
    sortOrder := c.DefaultQuery("sort_order", "desc")
	
	opts := domain.PostFilters{
		Title: &title,
		ProcedureID: procedureIDs,
		Tags: tags,
		Page:  page,
		Limit: limit,
		SortBy: sortBy,
		SortOrder: domain.SortOrder(sortOrder),
    }
    

    posts, total, err := dc.useCase.GetPosts(c, opts)
    if err != nil {
		HandleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "posts":      posts,
        "totalCount": total,
        "page":       page,
        "limit":      limit,
    })
}


func (dc *PostController) GetPostByID(c *gin.Context){
	id := c.Param("id")
	if id == ""{
		HandleError(c, domain.ErrEmptyParamField)
		return
	}
	res , err := dc.useCase.GetPostByID(c, id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": res})
}

func (dc *PostController) UpdatePost(c *gin.Context){
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists{
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	if id == ""{
		HandleError(c, domain.ErrEmptyParamField)
		return
	}
	var dto UpdatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil{
		HandleError(c, domain.ErrInvalidBody)
		return
	}
	Post := &domain.Post{
		ID:			id,
		UserID:   userID.(string),
		Title:    dto.Title,
		Content:  dto.Content,
		Procedures: dto.Procedures,
		Tags: dto.Tags,
	}
	res, err := dc.useCase.UpdatePost(c, Post)

	if err != nil{
		HandleError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"post": res})

}

func (dc *PostController) DeletePost(c *gin.Context){
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists{
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	userRole, exists := c.Get("userRole")
	if !exists{
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	if id == ""{
		HandleError(c, domain.ErrEmptyParamField)
		return
	}
	err := dc.useCase.DeletePost(c, id, userID.(string),string(userRole.(domain.Role))) 
	if err != nil{
		HandleError(c, err)
		return
	}
	c.Status(204)
}
