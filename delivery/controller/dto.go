package controller

import (
	"EthioGuide/domain"
	"time"
)

type RegisterRequest struct {
	Name          string `json:"name" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone"`
	Password      string `json:"password" binding:"required"`
	PreferredLang string `json:"preferredLang"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`

	ProfilePicture string      `json:"profile_picture,omitempty"`
	Role           domain.Role `json:"role"`
	IsVerified     bool        `json:"is_verified"`
	CreatedAt      time.Time   `json:"created_at"`
}

func toUserResponse(a *domain.Account) UserResponse {
	return UserResponse{
		ID:       a.ID,
		Name:     a.Name,
		Username: a.UserDetail.Username,
		Email:    a.Email,

		ProfilePicture: a.ProfilePicURL,
		Role:           a.Role,
		IsVerified:     a.UserDetail.IsVerified,
		CreatedAt:      a.CreatedAt,
	}
}

type TranslateDTO struct {
	Content string `json:"content" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         domain.Account `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token,omitempty"`
}

type CreateCategoryRequest struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id" binding:"required"`
	ParentID       string `json:"parent_id"`
	Title          string `json:"title" binding:"required"`
}

type CategoryResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	ParentID       string `json:"parent_id,omitempty"`
	Title          string `json:"title"`
}

type PaginatedCategoryResponse struct {
	Data    []*CategoryResponse `json:"data"`
	Total   int64               `json:"total"`
	Page    int64               `json:"page"`
	Limit   int64               `json:"limit"`
}