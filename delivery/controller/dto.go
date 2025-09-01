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

type ForgotDTO struct {
	Email string `json:"email"`
}

type ResetDTO struct {
	ResetToken  string `json:"resetToken"`
	NewPassword string `json:"new_password"`
}

type ActivateDTO struct {
	ActivateToken string `json:"activatationToken"`
}

type Procedure struct {
	ID             int                    `json:"id"`
	GroupID        *int                   `json:"group_id,omitempty"`
	OrganizationID int                    `json:"organization_id" `
	Name           string                 `json:"name"`
	Content        map[string]interface{} `json:"content,omitempty"`
	Fees           map[string]interface{} `json:"fees,omitempty"`
	ProcessingTime map[string]interface{} `json:"processing_time,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}
