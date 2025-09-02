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
type ChatRequest struct{
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

type ProcedureCreateRequest struct {
	Name           string   `json:"name"`
	GroupID        string   `json:"groupId"`
	OrganizationID string   `json:"organizationId"`

	// content
	Prerequisites  []string `json:"prerequisites"`
	Steps          []string `json:"steps"`
	Result         []string `json:"result"`

	// Fees
	Label          string   `json:"label"`
	Currency       string   `json:"currency"`
	Amount         float64  `json:"amount"` 

	// ProcessingTime
	MinDays        int      `json:"minDays"`
	MaxDays        int      `json:"maxDays"`
}

func toDomainProcedure(p *ProcedureCreateRequest) *domain.Procedure {
	content := domain.Content{
		Prerequisites: p.Prerequisites,
		Steps: p.Steps,
		Result: p.Result,
	}
	fees := domain.Fees {
		Label: p.Label,
		Currency: p.Currency,
		Amount: p.Amount,
	}
	processingTime := domain.ProcessingTime {
		MinDays: p.MinDays,
		MaxDays: p.MaxDays,
	}

	return &domain.Procedure{
		GroupID:        p.GroupID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
		Content:        content,
		Fees:           fees,
		ProcessingTime: processingTime,
	}
}