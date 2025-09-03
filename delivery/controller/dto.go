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


// notice related dto
type NoticeDTO struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Tags 		   []string	`json:"tags"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (dto *NoticeDTO) ToDomainNotice() *domain.Notice {
	return &domain.Notice{
		ID:             dto.ID,
		OrganizationID: dto.OrganizationID,
		Title:          dto.Title,
		Content:        dto.Content,
		Tags:          dto.Tags,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}

func FromDomainNotice(n *domain.Notice) *NoticeDTO {
	return &NoticeDTO{
		ID:             n.ID,
		OrganizationID: n.OrganizationID,
		Title:          n.Title,
		Content:        n.Content,
		Tags:          n.Tags,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
	}
}

// Paginated response for notices
type NoticeListResponse struct {
	Data  []NoticeDTO `json:"data"`
	Total int64       `json:"total"`
	Page  int64       `json:"page"`
	Limit int64       `json:"limit"`
}