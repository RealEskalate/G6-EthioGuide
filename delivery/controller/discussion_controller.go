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

// @Summary      Create Post
// @Description  Create a new post.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body CreatePostDTO true "Post Detail"
// @Success      200 {object} domain.Post "Post Created"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /discussions [post]
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
		UserID:     userID.(string),
		Title:      dto.Title,
		Content:    dto.Content,
		Procedures: dto.Procedures,
		Tags:       dto.Tags,
	}
	res, err := dc.useCase.CreatePost(c.Request.Context(), Post)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"post": res})
}

// @Summary      Fetch Posts
// @Description  Fetch list of posts
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        title query string false "title"
// @Param        userId query string fase "user id"
// @Param        procedure_ids query array false "procedure ids"
// @Param        tags query array false "tags"
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        sort_by query string false "sort_by"
// @Param        sort_order query string false "sort_order"
// @Success      200 {object} PaginatedPostsResponse "Posts"
// @Failure      400 {string}  "Bad Request"
// @Failure      500 {string}  "Internal"
// @Router       /discussions [get]
func (dc *PostController) GetPosts(c *gin.Context) {
	title := c.Query("title")
	userId := c.Query("userId")
	procedureIDs := c.QueryArray("procedure_ids")
	tags := c.QueryArray("tags")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	opts := domain.PostFilters{
		Title:       &title,
		UserId:      &userId,
		ProcedureID: procedureIDs,
		Tags:        tags,
		Page:        page,
		Limit:       limit,
		SortBy:      sortBy,
		SortOrder:   domain.SortOrder(sortOrder),
	}

	posts, total, err := dc.useCase.GetPosts(c, opts)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Posts": PaginatedPostsResponse{
			Posts: posts,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

// @Summary      Fetch Post
// @Description  Fetch a post
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        id path string true "Post ID"
// @Success      200 {object} domain.Post "Post"
// @Failure      400 {string}  "Bad Request"
// @Failure      500 {string}  "Internal"
// @Router       /discussions/{id} [get]
func (dc *PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		HandleError(c, domain.ErrEmptyParamField)
		return
	}
	res, err := dc.useCase.GetPostByID(c, id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": res})
}

// @Summary      Update Post
// @Description  Update a post
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Post ID"
// @Param        request body UpdatePostDTO  true "Post update details"
// @Success      200 {object} domain.Post "Post"
// @Failure      400 {string}  "Bad Request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      403 {string}  "Forbidden"
// @Failure      500 {string}  "Internal"
// @Router       /discussions/{id} [patch]
func (dc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	if id == "" {
		HandleError(c, domain.ErrEmptyParamField)
		return
	}
	var dto UpdatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		HandleError(c, domain.ErrInvalidBody)
		return
	}
	Post := &domain.Post{
		ID:         id,
		UserID:     userID.(string),
		Title:      dto.Title,
		Content:    dto.Content,
		Procedures: dto.Procedures,
		Tags:       dto.Tags,
	}
	res, err := dc.useCase.UpdatePost(c, Post)

	if err != nil {
		HandleError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"post": res})

}

// @Summary      Delete Post
// @Description  Delete a post
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Post ID"
// @Success      204 "No Content"
// @Failure      400 {string}  "Bad Request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      403 {string}  "Forbidden"
// @Failure      500 {string}  "Internal"
// @Router       /discussions/{id} [delete]
func (dc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	userRole, exists := c.Get("userRole")
	if !exists {
		HandleError(c, domain.ErrAuthenticationFailed)
		return
	}
	if id == "" {
		HandleError(c, domain.ErrEmptyParamField)
		return
	}
	err := dc.useCase.DeletePost(c, id, userID.(string), string(userRole.(domain.Role)))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.Status(204)
}
