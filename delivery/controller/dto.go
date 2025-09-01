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
	RefreshToken string         `json:"refresh_token"`
}


// procedure-related DTOs
type ProcedureDTO struct {
	ID             string   				`json:"id"`
	GroupID        *string  				`json:"group_id"`
	OrganizationID string   				 `json:"organization_id"`
	Name           string   				`json:"name"`
	Content        domain.ProcedureContent `json:"content"`
	Fees           domain.ProcedureFee     `json:"fees"`
	ProcessingTime domain.ProcessingTime   `json:"processing_time"`
	CreatedAt      time.Time               `json:"created_at"`
	NoticeIDs      []string                 `json:"notice_ids"`
}

func (dto *ProcedureDTO) FromDTOToDomainProcedure() *domain.Procedure {

	return &domain.Procedure{
		ID:             dto.ID,
		GroupID:        dto.GroupID,
		OrganizationID: dto.OrganizationID,
		Name:           dto.Name,
		Content:        dto.Content,
		Fees:           dto.Fees,
		ProcessingTime: dto.ProcessingTime,
		CreatedAt:      dto.CreatedAt,
		NoticeIDs:      dto.NoticeIDs,
	}
}

func FromDomainProcedureToDTO(p *domain.Procedure) *ProcedureDTO {
	return &ProcedureDTO{
		ID:             p.ID,
		GroupID:        p.GroupID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
		Content:        p.Content,
		Fees:           p.Fees,
		ProcessingTime: p.ProcessingTime,
		CreatedAt:      p.CreatedAt,
		NoticeIDs:      p.NoticeIDs,
	}
}